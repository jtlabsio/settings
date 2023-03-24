package settings

import (
	"reflect"
	"testing"
)

func TestOptions(t *testing.T) {
	tests := []struct {
		name string
		want ReadOptions
	}{
		{
			name: "should return empty ReadOptions",
			want: Options(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Options(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadOptions_EnvDefault(t *testing.T) {
	type fields struct {
		EnvOverride    []string
		EnvSearchPaths []string
	}
	tests := []struct {
		name   string
		fields fields
		want   ReadOptions
	}{
		{
			"should properly set default environment override settings",
			fields{},
			ReadOptions{
				EnvOverride:    []string{"GO_ENV"},
				EnvSearchPaths: []string{"./", "./config", "./settings"},
			},
		},
		{
			"should not override or clear any existing options",
			fields{
				EnvOverride:    []string{"SOME_OTHER_VAR"},
				EnvSearchPaths: []string{"./test"},
			},
			ReadOptions{
				EnvOverride:    []string{"SOME_OTHER_VAR", "GO_ENV"},
				EnvSearchPaths: []string{"./test", "./", "./config", "./settings"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ro := Options()

			if len(tt.fields.EnvOverride) > 0 {
				ro.EnvOverride = tt.fields.EnvOverride
			}

			if len(tt.fields.EnvSearchPaths) > 0 {
				ro.EnvSearchPaths = tt.fields.EnvSearchPaths
			}

			if got := ro.EnvDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOptions.EnvDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadOptions_SetArg(t *testing.T) {
	type args struct {
		arg       string
		fieldPath string
	}
	tests := []struct {
		name    string
		argsMap map[string]string
		args    args
		want    ReadOptions
	}{
		{
			"should properly add value to args map",
			nil,
			args{
				"--test-value",
				"Test.Value",
			},
			ReadOptions{
				ArgsMap: map[string]string{
					"--test-value": "Test.Value",
				},
			},
		},
		{
			"should not clear any existing values in args map",
			map[string]string{
				"--existing": "Existing",
			},
			args{
				"--test-value",
				"Test.Value",
			},
			ReadOptions{
				ArgsMap: map[string]string{
					"--existing":   "Existing",
					"--test-value": "Test.Value",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ro := Options()

			if tt.argsMap != nil {
				ro.ArgsMap = tt.argsMap
			}

			if got := ro.SetArg(tt.args.arg, tt.args.fieldPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOptions.SetArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadOptions_SetArgsFileOverride(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want ReadOptions
	}{
		{
			"should properly set the file override args when provided",
			args{
				[]string{"--config-file", "-cf"},
			},
			ReadOptions{
				ArgsFileOverride: []string{"--config-file", "-cf"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ro := Options()

			if got := ro.SetArgsFileOverride(tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOptions.SetArgsFileOverride() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadOptions_SetArgsMap(t *testing.T) {
	type args struct {
		argsMap map[string]string
		rewrite []bool
	}
	tests := []struct {
		name        string
		initArgsMap map[string]string
		args        args
		want        ReadOptions
	}{
		{
			"should create ArgsMap when not set (and no rewrite parameter)",
			nil,
			args{
				argsMap: map[string]string{},
			},
			ReadOptions{
				ArgsMap: map[string]string{},
			},
		},
		{
			"should create ArgsMap with supplied values when not set (and no rewrite parameter)",
			nil,
			args{
				argsMap: map[string]string{
					"test.test": "TEST_TEST",
				},
			},
			ReadOptions{
				ArgsMap: map[string]string{
					"test.test": "TEST_TEST",
				},
			},
		},
		{
			"should augment ArgsMap with supplied values when rewrite is not provided",
			map[string]string{
				"something.else": "SOMETHING_ELSE",
			},
			args{
				argsMap: map[string]string{
					"test.test": "TEST_TEST",
				},
			},
			ReadOptions{
				ArgsMap: map[string]string{
					"something.else": "SOMETHING_ELSE",
					"test.test":      "TEST_TEST",
				},
			},
		},
		{
			"should augment ArgsMap with supplied values when rewrite is false",
			map[string]string{
				"something.else": "SOMETHING_ELSE",
			},
			args{
				map[string]string{
					"test.test": "TEST_TEST",
				},
				[]bool{false},
			},
			ReadOptions{
				ArgsMap: map[string]string{
					"something.else": "SOMETHING_ELSE",
					"test.test":      "TEST_TEST",
				},
			},
		},
		{
			"should rewrite ArgsMap with supplied values when rewrite is true",
			map[string]string{
				"something.else": "SOMETHING_ELSE",
			},
			args{
				map[string]string{
					"test.test": "TEST_TEST",
				},
				[]bool{true},
			},
			ReadOptions{
				ArgsMap: map[string]string{
					"test.test": "TEST_TEST",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ro := Options()

			// init with args if desired
			if tt.initArgsMap != nil {
				ro.ArgsMap = tt.initArgsMap
			}

			var got ReadOptions

			if tt.args.rewrite == nil {
				got = ro.SetArgsMap(tt.args.argsMap)
			} else {
				got = ro.SetArgsMap(tt.args.argsMap, tt.args.rewrite...)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOptions.SetArgsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadOptions_SetBasePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want ReadOptions
	}{
		{
			"should properly set the base path when provided",
			args{
				"/usr/local/whatever.yml",
			},
			ReadOptions{
				BasePath: "/usr/local/whatever.yml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ro := Options()

			if got := ro.SetBasePath(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOptions.SetBasePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadOptions_SetDefaultsMap(t *testing.T) {
	type args struct {
		defMap  map[string]interface{}
		rewrite []bool
	}
	tests := []struct {
		name            string
		initDefaultsMap map[string]interface{}
		args            args
		want            ReadOptions
	}{
		{
			"should create DefaultsMap when not set (and no rewrite parameter)",
			nil,
			args{
				defMap: map[string]interface{}{},
			},
			ReadOptions{
				DefaultsMap: map[string]interface{}{},
			},
		},
		{
			"should create DefaultsMap with supplied values when not set (and no rewrite parameter)",
			nil,
			args{
				defMap: map[string]interface{}{
					"test.test": "testing 123",
				},
			},
			ReadOptions{
				DefaultsMap: map[string]interface{}{
					"test.test": "testing 123",
				},
			},
		},
		{
			"should augment DefaultsMap with supplied values with no rewrite parameter",
			map[string]interface{}{
				"something.else": 1234,
			},
			args{
				defMap: map[string]interface{}{
					"test.test": "testing 123",
				},
			},
			ReadOptions{
				DefaultsMap: map[string]interface{}{
					"something.else": 1234,
					"test.test":      "testing 123",
				},
			},
		},
		{
			"should augment DefaultsMap with supplied values with false rewrite parameter",
			map[string]interface{}{
				"something.else": 1234,
			},
			args{
				map[string]interface{}{
					"test.test": "testing 123",
				},
				[]bool{false},
			},
			ReadOptions{
				DefaultsMap: map[string]interface{}{
					"something.else": 1234,
					"test.test":      "testing 123",
				},
			},
		},
		{
			"should rewrite DefaultsMap with supplied values with true rewrite parameter",
			map[string]interface{}{
				"something.else": 1234,
			},
			args{
				map[string]interface{}{
					"test.test": "testing 123",
				},
				[]bool{true},
			},
			ReadOptions{
				DefaultsMap: map[string]interface{}{
					"test.test": "testing 123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ro := Options()

			// init with args if desired
			if tt.initDefaultsMap != nil {
				ro.DefaultsMap = tt.initDefaultsMap
			}

			var got ReadOptions

			if tt.args.rewrite == nil {
				got = ro.SetDefaultsMap(tt.args.defMap)
			} else {
				got = ro.SetDefaultsMap(tt.args.defMap, tt.args.rewrite...)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOptions.SetDefaultsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadOptions_SetSearchPaths(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name string
		args args
		want ReadOptions
	}{
		{
			"should properly set the search paths when provided",
			args{
				[]string{".", "./config", "./settings"},
			},
			ReadOptions{
				EnvSearchPaths: []string{".", "./config", "./settings"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ro := Options()

			if got := ro.SetEnvSearchPaths(tt.args.paths...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOptions.SetSearchPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadOptions_SetEnvSearchPattern(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want ReadOptions
	}{
		{
			"should properly set the search pattern when provided",
			args{
				"test.*",
			},
			ReadOptions{
				EnvSearchPattern: "test.*",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ro := Options()

			if got := ro.SetEnvSearchPattern(tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOptions.SetEnvSearchPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadOptions_SetVarsFileOverride(t *testing.T) {
	type args struct {
		vars []string
	}
	tests := []struct {
		name string
		args args
		want ReadOptions
	}{
		{
			"should properly set the file override args when provided",
			args{
				[]string{"GO_ENV", "CONFIG_FILE"},
			},
			ReadOptions{
				EnvOverride: []string{"GO_ENV", "CONFIG_FILE"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ro := Options()

			if got := ro.SetEnvOverride(tt.args.vars...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOptions.SetVarsFileOverride() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadOptions_SetVar(t *testing.T) {
	type args struct {
		v         string
		fieldPath string
	}
	tests := []struct {
		name    string
		varsMap map[string]string
		args    args
		want    ReadOptions
	}{
		{
			"should properly add value to vars map",
			nil,
			args{
				"TEST_VALUE",
				"Test.Value",
			},
			ReadOptions{
				VarsMap: map[string]string{
					"TEST_VALUE": "Test.Value",
				},
			},
		},
		{
			"should not clear any existing values in vars map",
			map[string]string{
				"EXISTING": "Existing",
			},
			args{
				"TEST_VALUE",
				"Test.Value",
			},
			ReadOptions{
				VarsMap: map[string]string{
					"EXISTING":   "Existing",
					"TEST_VALUE": "Test.Value",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ro := Options()

			if tt.varsMap != nil {
				ro.VarsMap = tt.varsMap
			}

			if got := ro.SetVar(tt.args.v, tt.args.fieldPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOptions.SetVar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadOptions_SetVarsMap(t *testing.T) {
	type args struct {
		varsMap map[string]string
		rewrite []bool
	}
	tests := []struct {
		name        string
		initVarsMap map[string]string
		args        args
		want        ReadOptions
	}{
		{
			"should create VarsMap when not set (and no rewrite parameter)",
			nil,
			args{
				varsMap: map[string]string{},
			},
			ReadOptions{
				VarsMap: map[string]string{},
			},
		},
		{
			"should create VarsMap with supplied values when not set (and no rewrite parameter)",
			nil,
			args{
				varsMap: map[string]string{
					"test.test": "--test-test",
				},
			},
			ReadOptions{
				VarsMap: map[string]string{
					"test.test": "--test-test",
				},
			},
		},
		{
			"should augment VarsMap with supplied values when rewrite is not provided",
			map[string]string{
				"something.else": "--something-else",
			},
			args{
				varsMap: map[string]string{
					"test.test": "--test-test",
				},
			},
			ReadOptions{
				VarsMap: map[string]string{
					"something.else": "--something-else",
					"test.test":      "--test-test",
				},
			},
		},
		{
			"should augment VarsMap with supplied values when rewrite is false",
			map[string]string{
				"something.else": "--something-else",
			},
			args{
				map[string]string{
					"test.test": "--test-test",
				},
				[]bool{false},
			},
			ReadOptions{
				VarsMap: map[string]string{
					"something.else": "--something-else",
					"test.test":      "--test-test",
				},
			},
		},
		{
			"should rewrite VarsMap with supplied values when rewrite is true",
			map[string]string{
				"something.else": "--something-else",
			},
			args{
				map[string]string{
					"test.test": "--test-test",
				},
				[]bool{true},
			},
			ReadOptions{
				VarsMap: map[string]string{
					"test.test": "--test-test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ro := Options()

			// init with args if desired
			if tt.initVarsMap != nil {
				ro.VarsMap = tt.initVarsMap
			}

			var got ReadOptions

			if tt.args.rewrite == nil {
				got = ro.SetVarsMap(tt.args.varsMap)
			} else {
				got = ro.SetVarsMap(tt.args.varsMap, tt.args.rewrite...)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOptions.SetVarsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
