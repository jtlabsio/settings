package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"

	"gopkg.in/yaml.v2"
)

var (
	dotRE   = regexp.MustCompile(`\.`)
	equalRE = regexp.MustCompile(`=`)
)

type settings struct {
	baseSettings []byte
	fieldTypeMap map[string]reflect.Kind
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
		baseSettings: []byte{},
		fieldTypeMap: map[string]reflect.Kind{},
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

	// apply command line arguments

	// apply environment variables

	return nil
}

func (s *settings) applyDefaultsMap(d map[string]interface{}) error {
	// only apply defaults where applicable
	if len(d) == 0 {
		return nil
	}

	// iterate the defaults and apply them (as appropriate)
	for fieldName, defaultValue := range d {
		// ensure the field exists in the out object
		if t, ok := s.fieldTypeMap[fieldName]; ok {
			// we found a match... ensure the type matches
			if t != reflect.ValueOf(defaultValue).Kind() {
				// type mismatch error
				return fmt.Errorf(
					"type mismatch for field %s: expected %v and default value is %v",
					fieldName,
					t,
					reflect.ValueOf(defaultValue).Kind())
			}

			// find the field within the out struct and set it (if we can)
			v := s.findOutFieldValue(fieldName)
			if v.CanSet() {
				dv := reflect.ValueOf(defaultValue)
				v.Set(dv)
				continue
			}

			// unable to set the value
			return SettingsFieldSetError(fieldName, t)
		}

		// default field is not in the out struct
		return SettingsFieldDoesNotExist("DefaultsMap", fieldName)
	}

	return nil
}

func (s *settings) determineFieldTypes() error {
	ct := reflect.TypeOf(s.out)

	// when a pointer, find the type that it is pointing to
	for ct.Kind() == reflect.Ptr {
		ct = ct.Elem()
	}

	// the target for settings must be a struct of some sort
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

		fmt.Printf("what is this %v \n", v)
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
		// TODO: do not know how to handle a Ptr in this scenario...
		s.fieldTypeMap[fieldName] = field.Type.Kind()
		return
	}

	fields := field.Type.NumField()
	for i := 0; i < fields; i++ {
		f := field.Type.FieldByIndex([]int{i})
		s.iterateFields(fieldName, f)
	}
}

func (s *settings) readBaseSettings(path string) error {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// base path doesn't exist
			return err
		}

		// unable to stat the file for other reasons...
		return SettingsFileReadError(path, err.Error())
	}

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		// unable to read the file
		return SettingsFileReadError(path, err.Error())
	}

	s.baseSettings = bs

	if err := s.unmarshalFile(path, s.baseSettings, &s.out); err != nil {
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

	b, err := ioutil.ReadFile(path)
	if err != nil {
		// unable to read the file
		return SettingsFileReadError(path, err.Error())
	}

	oo := map[interface{}]interface{}{}
	if err := s.unmarshalFile(path, b, &oo); err != nil {
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
		totalArgs := len(os.Args)

		for i, oa := range os.Args {
			// check for `--cli-arg=` scenario (where value is specified after =)
			if equalRE.MatchString(oa) {
				al := len(a)
				if len(oa) > al && oa[0:al] == a {
					// we have a match...
					path = oa[al+1:] // grab everything after the =
					break
				}
			}

			// check for direct arg match
			if oa == a && i < totalArgs-1 {
				// path should be the next argument specified
				path = os.Args[i+1]
				break
			}
		}

		// we found a path...
		if path != "" {
			fmt.Printf("path found with command line arg (%s): %s\n", a, path)
			if err := s.readOverrideFile(path); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *settings) unmarshalFile(path string, in []byte, out interface{}) error {
	t, err := s.determineFileType(path)
	if err != nil {
		// unable to determine base settings file type
		return err
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
