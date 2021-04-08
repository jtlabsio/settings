package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v2"
)

type settings struct {
	baseSettings     *[]byte
	baseSettingsType string
	fieldTypeMap     map[string]string
	target           *interface{}
}

func Read(opts ReadOptions, target interface{}) error {
	s := settings{
		target: &target,
	}

	// create an internal map for each field and its type
	if err := s.determineFieldTypes(); err != nil {
		return err
	}

	/*

		// read in base path (should be the base config file)
		if err := s.readBaseSettings(opts.BasePath); err != nil {
			// TODO: consider wrapping error
			return err
		}

		//*/

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

func (s settings) applyDefaultsMap(d map[string]interface{}) error {
	return nil
}

func (s settings) determineFieldTypes() error {
	s.fieldTypeMap = map[string]string{}

	t := *s.target
	ct := reflect.TypeOf(t)

	if ct.Kind() == reflect.Ptr {
	}

	if ct.Kind() != reflect.Struct {
		fmt.Println(ct.Kind())
		// target is not suitable to populate
		return errors.New("unable to read settings into unsupported type")
	}

	fields := ct.NumField()
	for i := 0; i < fields; i++ {
		field := ct.FieldByIndex([]int{i})
		fmt.Println(field)
	}

	return nil
}

func (s settings) determineFileType(path string) error {
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

func (s settings) readBaseSettings(path string) error {
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

	s.baseSettings = &bs

	if err := s.determineFileType(path); err != nil {
		// unable to determine base settings file type
		return err
	}

	// unmarshal base YAML
	if s.baseSettingsType == "yaml" {
		if err := yaml.Unmarshal(*s.baseSettings, s.target); err != nil {
			// unable to unmarshal as YAML
			return err
		}
	}

	// unmarshal base JSON
	if err := json.Unmarshal(*s.baseSettings, s.target); err != nil {
		// unable to unmarshal as JSON
		return err
	}

	return nil
}
