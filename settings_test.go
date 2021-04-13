package settings

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

type config struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Nested  struct {
		Name        string `yaml:"name"`
		Number      int    `yaml:"number"`
		NestedAgain struct {
			Desc string `yaml:"desc"`
		} `yaml:"nestedAgain"`
	} `yaml:"nested"`
	Numbers struct {
		V8  int8    `yaml:"v8"`
		V16 int16   `yaml:"v16"`
		V32 int32   `yaml:"v32"`
		V64 int64   `yaml:"v64"`
		U8  uint8   `yaml:"u8"`
		U16 uint16  `yaml:"u16"`
		U32 uint32  `yaml:"u32"`
		U64 uint64  `yaml:"u64"`
		F32 float32 `yaml:"f32"`
		F64 float64 `yaml:"f64"`
	} `yaml:"numbers"`
	Lists struct {
		LuckyNumbers []int     `yaml:"luckyNumbers"`
		Animals      []string  `yaml:"animals"`
		Dollars      []float32 `yaml:"dollars"`
	} `yaml:"lists"`
}

func TestGather(t *testing.T) {
	type args struct {
		opts ReadOptions
		out  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"nil interface",
			args{
				opts: Options(),
				out:  nil,
			},
			true,
		},
		{
			"blank interface",
			args{
				opts: Options(),
				out:  map[string]interface{}{},
			},
			true,
		},
		{
			"string interface",
			args{
				opts: Options(),
				out:  map[string]string{},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Gather(tt.args.opts, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("Gather() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_settings_applyArgs(t *testing.T) {
	type testConfig struct {
		Name string
	}
	type fields struct {
		fieldTypeMap map[string]reflect.Type
		out          interface{}
	}
	type args struct {
		a map[string]string
	}
	tests := []struct {
		name    string
		osArgs  []string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"should properly apply command line arguments",
			[]string{"--name", "test name"},
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Name": reflect.TypeOf(""),
				},
				out: &testConfig{},
			},
			args{
				a: map[string]string{
					"--name": "Name",
				},
			},
			&testConfig{
				Name: "test name",
			},
			false,
		},
		{
			"should properly apply command line arguments (one arg with =)",
			[]string{"--name=test name"},
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Name": reflect.TypeOf(""),
				},
				out: &testConfig{},
			},
			args{
				a: map[string]string{
					"--name": "Name",
				},
			},
			&testConfig{
				Name: "test name",
			},
			false,
		},
	}
	for _, tt := range tests {
		os.Args = tt.osArgs
		t.Run(tt.name, func(t *testing.T) {
			s := &settings{
				fieldTypeMap: tt.fields.fieldTypeMap,
				out:          tt.fields.out,
			}
			if err := s.applyArgs(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("settings.applyArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(s.out, tt.want) {
				t.Errorf("settings.applyArgs() = %v, want %v", s.out, tt.want)
			}
		})

		// reset the args
		os.Args = []string{}
	}
}

func Test_settings_applyDefaultsMap(t *testing.T) {
	type testConfig struct {
		Name   string
		Nested struct {
			Count int
		}
	}
	type fields struct {
		fieldTypeMap map[string]reflect.Type
		out          interface{}
	}
	type args struct {
		dm map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"should apply defaults map overrides to struct",
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Name":         reflect.TypeOf(""),
					"Nested.Count": reflect.TypeOf(1),
				},
				out: &testConfig{},
			},
			args{
				dm: map[string]interface{}{
					"Name":         "default name",
					"Nested.Count": 10,
				},
			},
			&testConfig{
				Name:   "default name",
				Nested: struct{ Count int }{10},
			},
			false,
		},
		{
			"should error if the default value type does not match",
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Name":         reflect.TypeOf(""),
					"Nested.Count": reflect.TypeOf(1),
				},
				out: &testConfig{},
			},
			args{
				dm: map[string]interface{}{
					"Name":         "default name",
					"Nested.Count": "10",
				},
			},
			&testConfig{},
			true,
		},
		{
			"should error if the default value is not in the destination struct",
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Name":         reflect.TypeOf(""),
					"Nested.Count": reflect.TypeOf(1),
				},
				out: &testConfig{},
			},
			args{
				dm: map[string]interface{}{
					"Name":         "default name",
					"Nested.Count": 10,
					"Random.Field": []string{"this", "is", "random"},
				},
			},
			&testConfig{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &settings{
				fieldTypeMap: tt.fields.fieldTypeMap,
				out:          tt.fields.out,
			}
			if err := s.applyDefaultsMap(tt.args.dm); (err != nil) != tt.wantErr {
				t.Errorf("settings.applyDefaultsMap() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(s.out, tt.want) {
				t.Errorf("settings.applyDefaultsMap() = %v, want %v", s.out, tt.want)
			}
		})
	}
}

func Test_settings_applyVars(t *testing.T) {
	type testConfig struct {
		Name   string
		Nested struct {
			Count int
		}
	}
	type fields struct {
		fieldTypeMap map[string]reflect.Type
		out          interface{}
	}
	type args struct {
		v map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		evars   map[string]string
		want    interface{}
		wantErr bool
	}{
		{
			"should apply environment override to struct",
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Name":         reflect.TypeOf(""),
					"Nested.Count": reflect.TypeOf(1),
				},
				out: &testConfig{},
			},
			args{
				v: map[string]string{
					"NAME":         "Name",
					"NESTED_COUNT": "Nested.Count",
				},
			},
			map[string]string{
				"NAME":         "testing name assignment",
				"NESTED_COUNT": "10",
			},
			&testConfig{
				Name:   "testing name assignment",
				Nested: struct{ Count int }{10},
			},
			false,
		},
	}
	for _, tt := range tests {
		// set the environment test values
		for ev, val := range tt.evars {
			os.Setenv(ev, val)
		}

		t.Run(tt.name, func(t *testing.T) {
			s := &settings{
				fieldTypeMap: tt.fields.fieldTypeMap,
				out:          tt.fields.out,
			}
			if err := s.applyVars(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("settings.applyVars() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(s.out, tt.want) {
				t.Errorf("settings.applyVars() = %v, want %v", s.out, tt.want)
			}
		})

		// clear the environment
		os.Clearenv()
	}
}

func Test_settings_cleanArgValue(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"should strip leading =",
			args{
				"=value",
			},
			"value",
		},
		{
			"should strip surrounding \"",
			args{
				"\"value\"",
			},
			"value",
		},
		{
			"should strip surrounding '",
			args{
				"'value'",
			},
			"value",
		},
		{
			"should strip both equal and quotes",
			args{
				"=\"value\"",
			},
			"value",
		},
		{
			"should not remove equal in the middle...",
			args{
				"value = value",
			},
			"value = value",
		},
		{
			"should not remove quote in the middle...",
			args{
				"value \"value\" value",
			},
			"value \"value\" value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := settings{}
			if got := s.cleanArgValue(tt.args.v); got != tt.want {
				t.Errorf("settings.cleanArgValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_settings_determineFieldTypes(t *testing.T) {
	tests := []struct {
		name    string
		s       interface{}
		want    map[string]reflect.Type
		wantErr bool
	}{
		{
			"should error when not a struct",
			"not a struct",
			map[string]reflect.Type{},
			true,
		},
		{
			"should properly parse fields",
			struct {
				Name   string
				Truthy bool
				Age    int
				Nested struct {
					Birthday string
				}
			}{},
			map[string]reflect.Type{
				"Age":             reflect.TypeOf(1),
				"Name":            reflect.TypeOf(""),
				"Nested.Birthday": reflect.TypeOf(""),
				"Truthy":          reflect.TypeOf(true),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &settings{
				fieldTypeMap: map[string]reflect.Type{},
				out:          tt.s,
			}
			err := s.determineFieldTypes()
			if (err != nil) != tt.wantErr {
				t.Errorf("settings.determineFileType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(s.fieldTypeMap, tt.want) {
				t.Errorf("settings.determineFieldTypes() = %v, want %v", s.fieldTypeMap, tt.want)
			}
		})
	}
}

func Test_settings_determineFileType(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"should properly detect json",
			args{
				path: "./config.json",
			},
			"json",
			false,
		},
		{
			"should properly detect yml",
			args{
				path: "./config.yml",
			},
			"yaml",
			false,
		},
		{
			"should properly detect yaml",
			args{
				path: "./config.yaml",
			},
			"yaml",
			false,
		},
		{
			"should error when unsupported",
			args{
				path: "./config.toml",
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &settings{
				out: &config{},
			}
			got, err := s.determineFileType(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("settings.determineFileType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("settings.determineFileType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_settings_readBaseSettings(t *testing.T) {
	type args struct {
		path string
	}
	// Options().SetBasePath("./tests/test.yaml")
	tests := []struct {
		name         string
		args         args
		want         *config
		wantErr      bool
		errorMessage string
	}{
		{
			"should set name and version",
			args{
				path: "./tests/simple.yaml",
			},
			&config{
				Name:    "example",
				Version: "1.1",
			},
			false,
			"",
		},
		{
			"should return an error = <nil> if path is blank",
			args{path: ""},
			&config{},
			false,
			"",
		},
		{
			// change wanterr to check match with SettingsFileReadError for line 323 coverage
			"should return os.ErrNotexist if bad path",
			args{
				path: "./not/found.yml",
			},
			&config{},
			true,
			"no such file",
		},
		{
			// add unmarshal file return error check for line 327
			"should return invalid character on badly formatted file",
			args{
				path: "./tests/broken.json",
			},
			&config{},
			true,
			"invalid character",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := settings{
				out: &config{},
			}
			if err := s.readBaseSettings(tt.args.path); err != nil {
				if !tt.wantErr {
					t.Errorf("settings.readBaseSettings() error = %v, wantErr %v", err, tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.errorMessage) {
					t.Errorf("settings.readBaseSettings() = %v, want %v", err.Error(), tt.errorMessage)
				}
			}
			if !reflect.DeepEqual(s.out, tt.want) {
				t.Errorf("settings.readBaseSettings() = %v, want %v", s.out, tt.want)
			}
		})
	}
}

func Test_settings_searchForArgOverrides(t *testing.T) {
	type testConfig struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	type fields struct {
		fieldTypeMap map[string]reflect.Type
		out          interface{}
	}
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		osArgs  []string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"should not apply override from command line when no args are set to be read",
			[]string{"--config-file", "./tests/simple.yaml"},
			fields{
				map[string]reflect.Type{
					"Name":    reflect.TypeOf(""),
					"Version": reflect.TypeOf(""),
				},
				&testConfig{},
			},
			args{
				[]string{},
			},
			&testConfig{},
			false,
		},
		{
			"should apply override from command line args",
			[]string{"--config-file", "./tests/simple.yaml"},
			fields{
				map[string]reflect.Type{
					"Name":    reflect.TypeOf(""),
					"Version": reflect.TypeOf(""),
				},
				&testConfig{},
			},
			args{
				[]string{"--config-file"},
			},
			&testConfig{
				Name:    "example",
				Version: "1.1",
			},
			false,
		},
		{
			"should apply override from command line arg (one arg with =)",
			[]string{"--config-file=./tests/simple.yaml"},
			fields{
				map[string]reflect.Type{
					"Name":    reflect.TypeOf(""),
					"Version": reflect.TypeOf(""),
				},
				&testConfig{},
			},
			args{
				[]string{"--config-file"},
			},
			&testConfig{
				Name:    "example",
				Version: "1.1",
			},
			false,
		},
		{
			"should error when there is a problem reading the file",
			[]string{"--config-file", "./tests/broken.json"},
			fields{
				map[string]reflect.Type{
					"Name":    reflect.TypeOf(""),
					"Version": reflect.TypeOf(""),
				},
				&testConfig{},
			},
			args{
				[]string{"--config-file"},
			},
			&testConfig{},
			true,
		},
		{
			"should error when file does not exist",
			[]string{"--config-file", "./not/exists.yml"},
			fields{
				map[string]reflect.Type{
					"Name":    reflect.TypeOf(""),
					"Version": reflect.TypeOf(""),
				},
				&testConfig{},
			},
			args{
				[]string{"--config-file"},
			},
			&testConfig{},
			true,
		},
	}
	for _, tt := range tests {
		os.Args = tt.osArgs

		t.Run(tt.name, func(t *testing.T) {
			s := &settings{
				fieldTypeMap: tt.fields.fieldTypeMap,
				out:          tt.fields.out,
			}
			if err := s.searchForArgOverrides(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("settings.searchForArgOverrides() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(s.out, tt.want) {
				t.Errorf("settings.searchForArgOverrides() = %v, want %v", s.out, tt.want)
			}
		})

		os.Args = []string{}
	}
}

func Test_settings_searchForEnvOverrides(t *testing.T) {
	type testConfig struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	type fields struct {
		fieldTypeMap map[string]reflect.Type
		out          interface{}
	}
	type args struct {
		searchPaths []string
		vars        []string
	}
	tests := []struct {
		name    string
		env     map[string]string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"should not apply override from command line when no vars are set to be read",
			map[string]string{
				"GO_ENV": "simple",
			},
			fields{
				map[string]reflect.Type{
					"Name":    reflect.TypeOf(""),
					"Version": reflect.TypeOf(""),
				},
				&testConfig{},
			},
			args{
				[]string{"./tests"},
				[]string{},
			},
			&testConfig{},
			false,
		},
		{
			"should not apply override from command line when no search paths are defined",
			map[string]string{
				"GO_ENV": "simple",
			},
			fields{
				map[string]reflect.Type{
					"Name":    reflect.TypeOf(""),
					"Version": reflect.TypeOf(""),
				},
				&testConfig{},
			},
			args{
				[]string{},
				[]string{"GO_ENV"},
			},
			&testConfig{},
			false,
		},
		{
			"should apply override from environment override",
			map[string]string{
				"GO_ENV": "simple",
			},
			fields{
				map[string]reflect.Type{
					"Name":    reflect.TypeOf(""),
					"Version": reflect.TypeOf(""),
				},
				&testConfig{},
			},
			args{
				[]string{"./tests"},
				[]string{"GO_ENV"},
			},
			&testConfig{
				Name:    "example",
				Version: "1.1",
			},
			false,
		},
		{
			"should error when there is a problem reading the file",
			map[string]string{
				"GO_ENV": "broken",
			},
			fields{
				map[string]reflect.Type{
					"Name":    reflect.TypeOf(""),
					"Version": reflect.TypeOf(""),
				},
				&testConfig{},
			},
			args{
				[]string{"./tests"},
				[]string{"GO_ENV"},
			},
			&testConfig{},
			true,
		},
		{
			"should not error when file does not exist",
			map[string]string{
				"GO_ENV": "no-environment",
			},
			fields{
				map[string]reflect.Type{
					"Name":    reflect.TypeOf(""),
					"Version": reflect.TypeOf(""),
				},
				&testConfig{},
			},
			args{
				[]string{"./tests"},
				[]string{"GO_ENV"},
			},
			&testConfig{},
			false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			// set environment variables
			for ev, val := range tt.env {
				os.Setenv(ev, val)
			}

			s := &settings{
				fieldTypeMap: tt.fields.fieldTypeMap,
				out:          tt.fields.out,
			}
			if err := s.searchForEnvOverrides(tt.args.vars, tt.args.searchPaths); (err != nil) != tt.wantErr {
				t.Errorf("settings.searchForEnvOverrides() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(s.out, tt.want) {
				t.Errorf("settings.searchForEnvOverrides() = %v, want %v", s.out, tt.want)
			}
		})

		os.Clearenv()
	}
}
