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
	if err := readBaseSettings(opts.BasePath); err != nil {
		log.Fatal(err)
	}
	// apply default mapped values
	defOptions, err := readDefaultsMap(opts.BasePath)
	if err != nil {
		log.Fatal(err)
	}
	opts.SetDefaultsMap(defOptions, true)
	// apply environment override files

	// apply environment variables

	// apply command line arguments

	return nil
}

func readBaseSettings(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			// base path doesn't exist
			return err
		}

		// unable to stat the file for other reasons...
		return err
	}

	Options().SetBasePath(path)

	return nil
}

func readDefaultsMap(path string) (map[string]interface{}, error) {
	yamlConfig, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	c := make(map[string]interface{})
	err = yaml.Unmarshal(yamlConfig, &c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", path, err)
	}

	return c, nil
}
