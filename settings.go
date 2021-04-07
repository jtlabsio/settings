package settings

import "os"

func Read(opts ReadOptions, s *interface{}) error {
	// read in base path (should be the base config file)

	// apply default mapped values

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

	return nil
}
