package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	commaRE     = regexp.MustCompile(`\,\s?`)
	dotRE       = regexp.MustCompile(`\.`)
	settingsExt = []string{".yml", ".yaml", ".json", ""}
	timeType    = reflect.TypeOf(time.Now())
)

type settings struct {
	fieldTypeMap map[string]reflect.Type
	out          interface{}
}

// Gather compiles configuration from various sources and
// iteratively builds up the out object with the values
// that are retrieved successively from the following sources:
// 1. base settings file
// 2. defaults as configured in options (*diverges from github.com/brozeph/settings-lib)
// 3. override files (from command line)
// 4. override files (from environment)
// 5. command line arguments
// 6. environment variables
func Gather(opts ReadOptions, out interface{}) error {
	s := settings{
		fieldTypeMap: map[string]reflect.Type{},
		out:          out,
	}

	// create an internal map for each field and its type
	if err := s.determineFieldTypes(); err != nil {
		return err
	}

	// read in base path (should be the base config file)
	if err := s.readBaseSettings(opts.BasePath); err != nil {
		return err
	}

	// apply default mapped values
	// iterate through the options.DefaultsMap and
	// apply the values that match the field names in the
	// inbound pointer argument that is an interface{} with
	// variable name "s"
	if err := s.applyDefaultsMap(opts.DefaultsMap); err != nil {
		return err
	}

	// iterate each arg file override
	if err := s.searchForArgOverrides(opts.ArgsFileOverride); err != nil {
		return err
	}

	// read any applicable environment override files
	if err := s.searchForEnvOverrides(opts.EnvOverride, opts.EnvSearchPaths, opts.EnvSearchPattern); err != nil {
		return err
	}

	// apply command line arguments
	if err := s.applyArgs(opts.ArgsMap); err != nil {
		return err
	}

	// apply environment variables
	if err := s.applyVars(opts.VarsMap); err != nil {
		return err
	}

	return nil
}

func (s *settings) applyArgs(a map[string]string) error {
	eq := []byte(`=`)
	totalArgs := len(os.Args)

	// iterate each element in args map
	for arg, field := range a {
		// iterate each arg provided to the application
		for i, oa := range os.Args {
			// check for `--cli-arg=` scenario (where value is specified after =)
			al := len(arg)
			if len(oa) > al && oa[0:al] == arg && oa[al] == eq[0] {
				// we have a match...
				if err := s.setFieldValue(
					field,
					s.cleanArgValue(oa[al:]),
					"Args"); err != nil {
					return err
				}

				break
			}

			// check for direct arg match
			if oa == arg && i < totalArgs-1 {
				if err := s.setFieldValue(
					field,
					s.cleanArgValue(os.Args[i+1]),
					"Args"); err != nil {
					return err
				}

				// next os.Arg is the value, skip trying to match it
				break
			}
		}
	}

	return nil
}

func (s *settings) applyVars(v map[string]string) error {
	// options.SetVarsMap shouldn't ever pass a nil map, so this is defensively safe
	if v == nil {
		return nil
	}

	// iterate the vars map
	for evar, fieldPath := range v {
		// lookup the var from the environment
		v := os.Getenv(evar)

		// if there is no value, continue on
		if v == "" {
			continue
		}

		// set the value
		if err := s.setFieldValue(fieldPath, v, "Vars"); err != nil {
			return err
		}
	}

	return nil
}

func (s *settings) applyDefaultsMap(d map[string]interface{}) error {
	// only apply defaults where applicable
	if len(d) == 0 {
		return nil
	}

	a := []struct {
		defVal    interface{}
		fieldName string
		fieldVal  reflect.Value
	}{}

	// validate each default value type before setting
	for fieldName, defVal := range d {
		if t, ok := s.fieldTypeMap[fieldName]; ok {
			if t.Kind() != reflect.ValueOf(defVal).Kind() {
				// type mismatch error
				return SettingsFieldTypeMismatch(
					fieldName,
					t.Kind(),
					reflect.ValueOf(defVal).Kind())
			}

			fieldVal := s.findOutFieldValue(fieldName)

			if !fieldVal.CanSet() {
				// unable to set the value
				return SettingsFieldSetError(fieldName, t.Kind())
			}

			a = append(
				a,
				struct {
					defVal    interface{}
					fieldName string
					fieldVal  reflect.Value
				}{
					defVal,
					fieldName,
					fieldVal,
				})

			continue
		}

		// default field is not in the out struct
		return SettingsFieldDoesNotExist("DefaultsMap", fieldName)
	}

	// iterate the default to apply and apply them
	for _, aa := range a {
		dv := reflect.ValueOf(aa.defVal)
		aa.fieldVal.Set(dv)
	}

	return nil
}

func (settings) cleanArgValue(v string) string {
	if len(v) == 0 {
		return v
	}

	charCheck := []byte(`='"`)

	for i, b := range charCheck {
		// look for = as first char and remove it
		if v[0] == b && i == 0 {
			v = v[1:]
			continue
		}

		// look for quotes (' or " surrounding the value)
		l := len(v)
		if v[0] == v[l-1] && v[0] == b {
			v = v[1 : l-1]
		}
	}

	return v
}

func (s *settings) determineFieldTypes() error {
	// if an out interface has somehow become nil, assign a new one
	if s.out == nil {
		return SettingsOutCannotBeNil()
	}
	ct := reflect.TypeOf(s.out)

	// when a pointer, find the type that it is pointing to
	for ct.Kind() == reflect.Ptr {
		ct = ct.Elem()
	}

	// check for a map target
	for ct.Kind() == reflect.Map {
		ct = ct.Elem()
	}

	// if the target isn't a map, then it must be a struct of some sort
	if ct.Kind() != reflect.Struct {
		// target is not suitable to populate
		return SettingsTypeDiscoveryError(ct.Kind())
	}

	fields := ct.NumField()
	for i := 0; i < fields; i++ {
		field := ct.FieldByIndex([]int{i})
		s.iterateFields("", field)
	}

	return nil
}

func (s *settings) determineFileType(path string) (string, error) {
	ext := filepath.Ext(path)
	var t string
	switch ext {
	case ".yml", ".yaml":
		t = "yaml"
	case ".json":
		t = "json"
	default:
		return t, SettingsFileTypeError(path, ext)
	}

	return t, nil
}

func (s *settings) findOutFieldValue(fieldPath string) reflect.Value {
	if fieldPath == "" {
		return reflect.Value{}
	}

	// create an array to iterate for the field hiearchy
	deepFields := dotRE.Split(fieldPath, -1)
	if len(deepFields) == 0 {
		deepFields = []string{fieldPath}
	}

	// find the value for the doc (which is the config)
	v := reflect.ValueOf(s.out)

	// iterate through each value until we get to the correct sub field
	for _, sf := range deepFields {
		// ensure we are working with the underlying value
		for v.Type().Kind() == reflect.Ptr {
			v = v.Elem()
		}

		v = v.FieldByName(sf)
	}

	return v
}

func (s *settings) iterateFields(parentPrefix string, field reflect.StructField) {
	fieldName := field.Name

	// make sure parent prefix is set for subsequent use...
	if parentPrefix != "" {
		fieldName = fmt.Sprintf("%s.%s", parentPrefix, fieldName)
	}

	// if field is not a struct, store the type
	if field.Type.Kind() != reflect.Struct {
		s.fieldTypeMap[fieldName] = field.Type
		return
	}

	fields := field.Type.NumField()
	for i := 0; i < fields; i++ {
		f := field.Type.FieldByIndex([]int{i})
		s.iterateFields(fieldName, f)
	}
}

func (s *settings) readBaseSettings(path string) error {
	// just return if path is empty
	if path == "" {
		return nil
	}

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// base path doesn't exist
			return err
		}

		// unable to stat the file for other reasons...
		return SettingsFileReadError(path, err.Error())
	}

	if err := s.unmarshalFile(path, s.out); err != nil {
		return err
	}

	return nil
}

func (s *settings) readOverrideFile(path string) error {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// base path doesn't exist
			return err
		}

		// unable to stat the file for other reasons...
		return SettingsFileReadError(path, err.Error())
	}

	// unmarshal over the top of the base...
	if err := s.unmarshalFile(path, s.out); err != nil {
		return err
	}

	return nil
}

func (s *settings) searchForArgOverrides(args []string) error {
	if len(args) == 0 {
		return nil
	}

	for _, a := range args {
		var path string
		eq := []byte(`=`)
		totalArgs := len(os.Args)

		for i, oa := range os.Args {
			// check for `--cli-arg=` scenario (where value is specified after =)
			al := len(a)
			if len(oa) > al && oa[0:al] == a && oa[al] == eq[0] {
				// we have a match...
				path = s.cleanArgValue(oa[al:])

				break
			}

			// check for direct arg match
			if oa == a && i < totalArgs-1 {
				// path should be the next argument specified
				path = s.cleanArgValue(os.Args[i+1])
				break
			}
		}

		// we found a path...
		if path != "" {
			if err := s.readOverrideFile(path); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *settings) searchForEnvOverrides(vars []string, searchPaths []string, filePattern string) error {
	if len(vars) == 0 {
		return nil
	}

	if len(searchPaths) == 0 {
		return nil
	}

	var extensionSearch = func(sp string) (bool, error) {
		for _, ext := range settingsExt {
			spf := fmt.Sprintf("%s%s", sp, ext)

			// continue when the file can't be opened (presumably does not exist)
			if _, err := os.Stat(spf); err != nil {
				continue
			}

			// unmarshal the environment override over the base
			if err := s.readOverrideFile(spf); err != nil {
				return false, err
			}

			return true, nil
		}

		return false, nil
	}

	for _, v := range vars {
		envName := os.Getenv(v)

		// detected an environment name
		if envName != "" {
			// now iterate search paths
			for _, prefix := range searchPaths {
				var found bool

				// search file by env file pattern if provided
				if filePattern != "" {
					sp := path.Join(prefix, fmt.Sprintf(filePattern, envName))
					f, err := extensionSearch(sp)
					if err != nil {
						return err
					}
					found = f
				}

				if !found && filePattern == "" {
					// search file by environment name alone
					sp := path.Join(prefix, envName)
					f, err := extensionSearch(sp)
					if err != nil {
						return err
					}
					found = f
				}

				if found {
					break
				}
			}
		}
	}

	return nil
}

func (s *settings) setFieldValue(fieldPath string, sVal string, override string) error {
	// ensure the field exists in the out object
	if t, ok := s.fieldTypeMap[fieldPath]; ok {
		// we found a match... ensure the type matches
		var val interface{}

		switch t.Kind() {
		case reflect.Array, reflect.Slice:
			sVals := commaRE.Split(sVal, -1)
			ov := s.findOutFieldValue(fieldPath)
			st := ov.Type().Elem().Kind()
			pv := reflect.MakeSlice(reflect.Indirect(ov).Type(), len(sVals), cap(sVals))

			for i, sv := range sVals {
				switch st {
				case reflect.Bool:
					v, err := strconv.ParseBool(sv)
					if err != nil {
						return SettingsFieldSetError(fieldPath, t.Kind(), err)
					}
					pv.Index(i).Set(reflect.ValueOf(v))
				case reflect.Int:
					v, err := strconv.ParseInt(sv, 0, ov.Type().Elem().Bits())
					if err != nil {
						return SettingsFieldSetError(fieldPath, t.Kind(), err)
					}
					iv := int(v)
					pv.Index(i).Set(reflect.ValueOf(iv))
				case reflect.Int8:
					v, err := strconv.ParseInt(sv, 0, ov.Type().Elem().Bits())
					if err != nil {
						return SettingsFieldSetError(fieldPath, t.Kind(), err)
					}
					iv := int8(v)
					pv.Index(i).Set(reflect.ValueOf(iv))
				case reflect.Int16:
					v, err := strconv.ParseInt(sv, 0, ov.Type().Elem().Bits())
					if err != nil {
						return SettingsFieldSetError(fieldPath, t.Kind(), err)
					}
					iv := int16(v)
					pv.Index(i).Set(reflect.ValueOf(iv))
				case reflect.Int32:
					v, err := strconv.ParseInt(sv, 0, ov.Type().Elem().Bits())
					if err != nil {
						return SettingsFieldSetError(fieldPath, t.Kind(), err)
					}
					iv := int32(v)
					pv.Index(i).Set(reflect.ValueOf(iv))
				case reflect.Int64:
					v, err := strconv.ParseInt(sv, 0, ov.Type().Elem().Bits())
					if err != nil {
						return SettingsFieldSetError(fieldPath, t.Kind(), err)
					}
					pv.Index(i).Set(reflect.ValueOf(v))
				case reflect.Uint:
					v, err := strconv.ParseUint(sv, 0, ov.Type().Elem().Bits())
					if err != nil {
						return SettingsFieldSetError(fieldPath, t.Kind(), err)
					}
					iv := uint(v)
					pv.Index(i).Set(reflect.ValueOf(iv))
				case reflect.Uint8:
					v, err := strconv.ParseUint(sv, 0, ov.Type().Elem().Bits())
					if err != nil {
						return SettingsFieldSetError(fieldPath, t.Kind(), err)
					}
					iv := uint8(v)
					pv.Index(i).Set(reflect.ValueOf(iv))
				case reflect.Uint16:
					v, err := strconv.ParseUint(sv, 0, ov.Type().Elem().Bits())
					if err != nil {
						return SettingsFieldSetError(fieldPath, t.Kind(), err)
					}
					iv := uint16(v)
					pv.Index(i).Set(reflect.ValueOf(iv))
				case reflect.Uint32:
					v, err := strconv.ParseUint(sv, 0, ov.Type().Elem().Bits())
					if err != nil {
						return SettingsFieldSetError(fieldPath, t.Kind(), err)
					}
					iv := uint32(v)
					pv.Index(i).Set(reflect.ValueOf(iv))
				case reflect.Uint64:
					v, err := strconv.ParseUint(sv, 0, ov.Type().Elem().Bits())
					if err != nil {
						return SettingsFieldSetError(fieldPath, t.Kind(), err)
					}
					pv.Index(i).Set(reflect.ValueOf(v))
				case reflect.Float32:
					v, err := strconv.ParseFloat(sv, ov.Type().Elem().Bits())
					if err != nil {
						return SettingsFieldSetError(fieldPath, t.Kind(), err)
					}
					fv := float32(v)
					pv.Index(i).Set(reflect.ValueOf(fv))
				case reflect.Float64:
					v, err := strconv.ParseFloat(sv, ov.Type().Elem().Bits())
					if err != nil {
						return SettingsFieldSetError(fieldPath, t.Kind(), err)
					}
					pv.Index(i).Set(reflect.ValueOf(v))
				case reflect.String:
					pv.Index(i).Set(reflect.ValueOf(sv))
				default:
					// complex64, complex128, chan, func, interface, map, ptr, struct and unsafeptr
					return SettingsFieldSetError(
						fieldPath,
						t.Kind(),
						errors.New("unsupported field type"))
				}
			}

			val = pv.Interface()
		case reflect.Bool:
			pv, err := strconv.ParseBool(sVal)
			if err != nil {
				return SettingsFieldSetError(fieldPath, t.Kind(), err)
			}
			val = pv
		case reflect.Int:
			pv, err := strconv.ParseInt(sVal, 0, t.Bits())
			if err != nil {
				return SettingsFieldSetError(fieldPath, t.Kind(), err)
			}
			val = int(pv)
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			pv, err := strconv.ParseInt(sVal, 0, t.Bits())
			if err != nil {
				return SettingsFieldSetError(fieldPath, t.Kind(), err)
			}

			switch t.Bits() {
			case 8:
				val = int8(pv)
			case 16:
				val = int16(pv)
			case 32:
				val = int32(pv)
			default:
				val = pv
			}
		case reflect.Uint:
			pv, err := strconv.ParseInt(sVal, 0, t.Bits())
			if err != nil {
				return SettingsFieldSetError(fieldPath, t.Kind(), err)
			}
			val = uint(pv)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			pv, err := strconv.ParseUint(sVal, 0, t.Bits())
			if err != nil {
				return SettingsFieldSetError(fieldPath, t.Kind(), err)
			}

			switch t.Bits() {
			case 8:
				val = uint8(pv)
			case 16:
				val = uint16(pv)
			case 32:
				val = uint32(pv)
			default:
				val = pv
			}
		case reflect.Float32, reflect.Float64:
			pv, err := strconv.ParseFloat(sVal, t.Bits())
			if err != nil {
				return SettingsFieldSetError(fieldPath, t.Kind(), err)
			}

			switch t.Bits() {
			case 32:
				val = float32(pv)
			default:
				val = pv
			}
		case reflect.String:
			val = sVal
		default:
			if t == timeType {
				dv, err := time.Parse(time.RFC3339, sVal)
				if err != nil {
					return SettingsFieldSetError(fieldPath, t.Kind(), err)
				}
				val = dv
				break
			}

			// complex64, complex128, chan, func, interface, map, ptr, struct and unsafeptr
			return SettingsFieldSetError(
				fieldPath,
				t.Kind(),
				errors.New("unsupported field type"))
		}

		// don't try to set if there's no value to set
		if reflect.Zero(t) == val || val == nil {
			return nil
		}

		// find the field within the out struct and set it (if we can)
		v := s.findOutFieldValue(fieldPath)
		if v.CanSet() {
			dv := reflect.ValueOf(val)
			v.Set(dv)
			return nil
		}

		// unable to set the value
		return SettingsFieldSetError(fieldPath, t.Kind())
	}

	// default field is not in the out struct
	return SettingsFieldDoesNotExist(override, fieldPath)
}

func (s *settings) unmarshalFile(path string, out interface{}) error {
	t, err := s.determineFileType(path)
	if err != nil {
		// unable to determine base settings file type
		return err
	}

	in, err := ioutil.ReadFile(path)
	if err != nil {
		// unable to read the file
		return SettingsFileReadError(path, err.Error())
	}

	// unmarshal YAML
	if t == "yaml" {
		if err := yaml.Unmarshal(in, out); err != nil {
			// unable to unmarshal as YAML
			return SettingsFileParseError(path, err.Error())
		}

		return nil
	}

	// unmarshal JSON
	if t == "json" {
		if err := json.Unmarshal(in, out); err != nil {
			// unable to unmarshal as JSON
			return SettingsFileParseError(path, err.Error())
		}
	}

	return nil
}
