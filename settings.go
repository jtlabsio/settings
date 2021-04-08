package settings

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func Read(opts ReadOptions, s *interface{}) error {
	// read in base path (should be the base config file)
	if err := readBaseSettings(opts.BasePath, s); err != nil {
		log.Fatal(err)
	}

	// apply default mapped values
	// iterate through the options.DefaultsMap and
	// apply the values that match the field names in the
	// inbound pointer argument that is an interface{} with
	// variable name "s"
	if err := applyDefaultsMap(opts.DefaultsMap, s); err != nil {
		log.Fatal(err)
	}

	// read any applicable environment override files

	// apply environment variables

	// apply command line arguments

	return nil
}

func readBaseSettings(path string, s *interface{}) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// base path doesn't exist
			return err
		}

		// unable to stat the file for other reasons...
		return err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		// unable to read the file
		return err
	}

	// determine if JSON or YAML
	if err := yaml.Unmarshal(b, s); err != nil {
		// unable to unmarshal
		return err
	}

	return nil
}

func applyDefaultsMap(d map[string]interface{}, s *interface{}) error {
	// only iterate if d contains values
	if len(d) != 0 {
		// create a map to hold the options and assign
		// to the passed in empty interface, making it a map
		map1 := map[string]interface{}{}
		for k, v := range d {
			map1[k] = v
		}
		*s = map1
		return nil
	}

	return fmt.Errorf("empty map %q", d)
}
