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
			want: ReadOptions{},
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
			"should properly set the base path when provided",
			args{
				[]string{".", "./config", "./settings"},
			},
			ReadOptions{
				SearchPaths: []string{".", "./config", "./settings"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ro := Options()

			if got := ro.SetSearchPaths(tt.args.paths...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadOptions.SetSearchPaths() = %v, want %v", got, tt.want)
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
