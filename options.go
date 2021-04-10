package settings

// ReadOptions define additional optional instructions for
// the Settings package when reading and compiling layers of
// configuration settings from various sources
type ReadOptions struct {
	ArgsFileOverride []string
	ArgsMap          map[string]string
	BasePath         string
	DefaultsMap      map[string]interface{}
	EnvOverride      []string
	EnvSearchPaths   []string
	VarsMap          map[string]string
}

// Options returns an empty ReadOptions for use with the
// Settings package
func Options() ReadOptions {
	return ReadOptions{}
}

// SetArgsFileOverride instructs the settings package on where to look
// for any potential override file locations that are provided as command
// line arguments
func (ro ReadOptions) SetArgsFileOverride(args ...string) ReadOptions {
	if len(ro.ArgsFileOverride) == 0 {
		ro.ArgsFileOverride = []string{}
	}

	ro.ArgsFileOverride = append(ro.ArgsFileOverride, args...)

	return ro
}

// SetArgsMap will either rewrite or, by default, augment the map
// that defines which configuration keys are related to the specified
// command line arguments
func (ro ReadOptions) SetArgsMap(argsMap map[string]string, rewrite ...bool) ReadOptions {
	// ensure it's not empty
	if ro.ArgsMap == nil {
		ro.ArgsMap = map[string]string{}
	}

	// add to the map when requested (or by default)
	if len(rewrite) == 0 || !rewrite[0] {
		populateMap(&ro.ArgsMap, argsMap)

		return ro
	}

	// rewrite the map entirely
	ro.ArgsMap = argsMap

	return ro
}

// SetBasePath can be used to define the path to the base settings
// file which is the first element loaded when the Settings package
// begins reading configuration
func (ro ReadOptions) SetBasePath(path string) ReadOptions {
	ro.BasePath = path
	return ro
}

// SetDefaultsMap can be used to define default values for config
// elements in the event that the value is not provided in one
// of the layered mechanisms used to read settings
func (ro ReadOptions) SetDefaultsMap(defMap map[string]interface{}, rewrite ...bool) ReadOptions {
	if ro.DefaultsMap == nil {
		ro.DefaultsMap = map[string]interface{}{}
	}

	// add to the map when requested (or by default)
	if len(rewrite) == 0 || !rewrite[0] {
		for k, v := range defMap {
			ro.DefaultsMap[k] = v
		}

		return ro
	}

	// rewrite the map entirely
	ro.DefaultsMap = defMap

	return ro
}

// SetEnvOverride instructs the settings package on where to look
// for any potential override file locations that are provided as environment
// variables to the application
func (ro ReadOptions) SetEnvOverride(vars ...string) ReadOptions {
	if len(ro.EnvOverride) == 0 {
		ro.EnvOverride = []string{}
	}

	ro.EnvOverride = append(ro.ArgsFileOverride, vars...)

	return ro
}

// SetEnvSearchPaths can be used to instruct the Settings package on
// where it might find additional configuration files for use when
// loading additional layers of configuration
func (ro ReadOptions) SetEnvSearchPaths(paths ...string) ReadOptions {
	if len(ro.EnvSearchPaths) == 0 {
		ro.EnvSearchPaths = []string{}
	}

	ro.EnvSearchPaths = append(ro.EnvSearchPaths, paths...)

	return ro
}

// SetVarsMap will either rewrite or, by default, augment the map
// that associates environment variables to various configuration keys
// specified in the base
func (ro ReadOptions) SetVarsMap(varsMap map[string]string, rewrite ...bool) ReadOptions {
	// ensure it's not empty
	if ro.VarsMap == nil {
		ro.VarsMap = map[string]string{}
	}

	// add to the map when requested (or by default)
	if len(rewrite) == 0 || !rewrite[0] {
		populateMap(&ro.VarsMap, varsMap)

		return ro
	}

	// rewrite the map entirely
	ro.VarsMap = varsMap

	return ro
}

func populateMap(tm *map[string]string, fm map[string]string) {
	for k, v := range fm {
		(*tm)[k] = v
	}
}
