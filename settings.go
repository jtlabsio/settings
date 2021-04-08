package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"

	"gopkg.in/yaml.v2"
)

var dotRE = regexp.MustCompile(`\.`)

type settings struct {
	baseSettings     []byte
	baseSettingsType string
	fieldTypeMap     map[string]reflect.Kind
	out              interface{}
}

// Read compiles configuration from various sources and
// iteratively builds up the out object with the values
// that are read successively from the following sources:
// 1. base settings file
// 2. defaults as configured in options (*diverges from github.com/brozeph/settings-lib)
// 3. subsequent override settings files
// 4. command line arguments
// 5. environment variables
func Read(opts ReadOptions, out interface{}) error {
	s := settings{
		baseSettings:     []byte{},
		baseSettingsType: "",
		fieldTypeMap:     map[string]reflect.Kind{},
		out:              out,
	}

	// create an internal map for each field and its type
	if err := s.determineFieldTypes(); err != nil {
		return err
	}

	// read in base path (should be the base config file)
	if err := s.readBaseSettings(opts.BasePath); err != nil {
		// TODO: consider wrapping error
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

	// read any applicable environment override files

	// apply environment variables

	// apply command line arguments

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
			v := s.findField(fieldName)
			if v.CanSet() {
				dv := reflect.ValueOf(defaultValue)
				v.Set(dv)
				continue
			}

			// unable to set the value
			// TODO: should we error?
			continue
		}

		// default field is not in the out struct
		// TODO: should we error?
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
		return fmt.Errorf("unable to read settings into unsupported type (%v)", ct.Kind())
	}

	fields := ct.NumField()
	for i := 0; i < fields; i++ {
		field := ct.FieldByIndex([]int{i})
		s.iterateFields("", field)
	}

	return nil
}

func (s *settings) determineFileType(path string) error {
	ext := filepath.Ext(path)
	switch ext {
	case ".yml", ".yaml":
		s.baseSettingsType = "yaml"
	case ".json":
		s.baseSettingsType = "json"
	default:
		return fmt.Errorf("unsupported file type for base settings: %s", path)
	}

	return nil
}

func (s *settings) findField(fieldPath string) reflect.Value {
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
		if os.IsNotExist(err) {
			// base path doesn't exist
			return err
		}

		// unable to stat the file for other reasons...
		return err
	}

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		// unable to read the file
		return err
	}

	s.baseSettings = bs

	if err := s.determineFileType(path); err != nil {
		// unable to determine base settings file type
		return err
	}

	// unmarshal base YAML
	if s.baseSettingsType == "yaml" {
		if err := yaml.Unmarshal(s.baseSettings, s.out); err != nil {
			// unable to unmarshal as YAML
			return err
		}

		return nil
	}

	// unmarshal base JSON
	if err := json.Unmarshal(s.baseSettings, s.out); err != nil {
		// unable to unmarshal as JSON
		return err
	}

	return nil
}
