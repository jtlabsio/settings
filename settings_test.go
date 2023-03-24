package settings

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

type verboseConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Nested  struct {
		Bool        bool
		Name        string
		Number      int
		NestedAgain struct {
			Desc string
		}
	}
	Numbers struct {
		I   int
		I8  int8
		I16 int16
		I32 int32
		I64 int64
		U   uint
		U8  uint8
		U16 uint16
		U32 uint32
		U64 uint64
		F32 float32
		F64 float64
	}
	Lists struct {
		B   []bool
		I   []int
		I8  []int8
		I16 []int16
		I32 []int32
		I64 []int64
		U   []uint
		U8  []uint8
		U16 []uint16
		U32 []uint32
		U64 []uint64
		F32 []float32
		F64 []float64
		S   []string
		T   []struct{}
	}
}

func TestGather(t *testing.T) {
	type testConfig struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	type args struct {
		opts ReadOptions
		out  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"should error with a nil out value",
			args{
				opts: Options(),
				out:  nil,
			},
			nil,
			true,
		},
		{
			"should error with a non struct type out value",
			args{
				opts: Options(),
				out:  map[string]interface{}{},
			},
			map[string]interface{}{},
			true,
		},
		{
			"should read base settings",
			args{
				opts: Options().SetBasePath("./tests/simple.yaml"),
				out:  &testConfig{},
			},
			&testConfig{
				Name:    "example",
				Version: "1.1",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Gather(tt.args.opts, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("Gather() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.out, tt.want) {
				t.Errorf("Gather() = %v, want %v", tt.args.out, tt.want)
			}
		})
	}
}

func Test_settings_applyArgs(t *testing.T) {
	type testConfig struct {
		Name    string
		Created time.Time
	}
	type fields struct {
		fieldTypeMap map[string]reflect.Type
		out          interface{}
	}
	type args struct {
		a map[string]string
	}
	tests := []struct {
		name         string
		osArgs       []string
		fields       fields
		args         args
		want         interface{}
		errorMessage string
		wantErr      bool
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
			"",
			false,
		},
		{
			"should properly set time.Time from command line argument",
			[]string{"--created", "2021-02-16T00:00:00.000Z"},
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Created": reflect.TypeOf(time.Now()),
				},
				out: &testConfig{},
			},
			args{
				a: map[string]string{
					"--created": "Created",
				},
			},
			&testConfig{
				Created: time.Date(2021, time.February, 16, 0, 0, 0, 0, time.UTC),
			},
			"",
			false,
		},
		{
			"should gracefully handle when not enough arguments are provided",
			[]string{"--name"},
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
			&testConfig{},
			"",
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
			"",
			false,
		},
		{
			"should properly handle number string value conversion",
			[]string{
				"-i=-1",
				"-i8=-10",
				"-i16=-100",
				"-i32=-1000",
				"-i64=-10000",
				"-u", "1",
				"-u8=10",
				"-u16=100",
				"-u32=1000",
				"-u64=10000",
				"-f32=10.10",
				"-f64=100.100",
			},
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Numbers.I":   reflect.TypeOf(int(1)),
					"Numbers.I8":  reflect.TypeOf(int8(1)),
					"Numbers.I16": reflect.TypeOf(int16(1)),
					"Numbers.I32": reflect.TypeOf(int32(1)),
					"Numbers.I64": reflect.TypeOf(int64(1)),
					"Numbers.U":   reflect.TypeOf(uint(1)),
					"Numbers.U8":  reflect.TypeOf(uint8(1)),
					"Numbers.U16": reflect.TypeOf(uint16(1)),
					"Numbers.U32": reflect.TypeOf(uint32(1)),
					"Numbers.U64": reflect.TypeOf(uint64(1)),
					"Numbers.F32": reflect.TypeOf(float32(1.1)),
					"Numbers.F64": reflect.TypeOf(float64(1.1)),
				},
				out: &verboseConfig{},
			},
			args{
				map[string]string{
					"-i":   "Numbers.I",
					"-i8":  "Numbers.I8",
					"-i16": "Numbers.I16",
					"-i32": "Numbers.I32",
					"-i64": "Numbers.I64",
					"-u":   "Numbers.U",
					"-u8":  "Numbers.U8",
					"-u16": "Numbers.U16",
					"-u32": "Numbers.U32",
					"-u64": "Numbers.U64",
					"-f32": "Numbers.F32",
					"-f64": "Numbers.F64",
				},
			},
			&verboseConfig{
				Numbers: struct {
					I   int
					I8  int8
					I16 int16
					I32 int32
					I64 int64
					U   uint
					U8  uint8
					U16 uint16
					U32 uint32
					U64 uint64
					F32 float32
					F64 float64
				}{
					I:   -1,
					I8:  -10,
					I16: -100,
					I32: -1000,
					I64: -10000,
					U:   1,
					U8:  10,
					U16: 100,
					U32: 1000,
					U64: 10000,
					F32: 10.10,
					F64: 100.100,
				},
			},
			"",
			false,
		},
		{
			"should properly handle slice string value conversion",
			[]string{
				"-b=true,false",
				"-i=-1",
				"-i8=-10",
				"-i16=-100",
				"-i32=-1000",
				"-i64=-10000",
				"-u=1",
				"-u8=10",
				"-u16=100",
				"-u32=1000",
				"-u64=10000",
				"-f32=10.10",
				"-f64=100.100",
				"-s=testing,a,string,array",
			},
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Lists.B":   reflect.TypeOf([]bool{true}),
					"Lists.I":   reflect.TypeOf([]int{1}),
					"Lists.I8":  reflect.TypeOf([]int8{1}),
					"Lists.I16": reflect.TypeOf([]int16{1}),
					"Lists.I32": reflect.TypeOf([]int32{1}),
					"Lists.I64": reflect.TypeOf([]int64{1}),
					"Lists.U":   reflect.TypeOf([]uint{1}),
					"Lists.U8":  reflect.TypeOf([]uint8{1}),
					"Lists.U16": reflect.TypeOf([]uint16{1}),
					"Lists.U32": reflect.TypeOf([]uint32{1}),
					"Lists.U64": reflect.TypeOf([]uint64{1}),
					"Lists.F32": reflect.TypeOf([]float32{1.1}),
					"Lists.F64": reflect.TypeOf([]float64{1.1}),
					"Lists.S":   reflect.TypeOf([]string{""}),
				},
				out: &verboseConfig{},
			},
			args{
				map[string]string{
					"-b":   "Lists.B",
					"-i":   "Lists.I",
					"-i8":  "Lists.I8",
					"-i16": "Lists.I16",
					"-i32": "Lists.I32",
					"-i64": "Lists.I64",
					"-u":   "Lists.U",
					"-u8":  "Lists.U8",
					"-u16": "Lists.U16",
					"-u32": "Lists.U32",
					"-u64": "Lists.U64",
					"-f32": "Lists.F32",
					"-f64": "Lists.F64",
					"-s":   "Lists.S",
				},
			},
			&verboseConfig{
				Lists: struct {
					B   []bool
					I   []int
					I8  []int8
					I16 []int16
					I32 []int32
					I64 []int64
					U   []uint
					U8  []uint8
					U16 []uint16
					U32 []uint32
					U64 []uint64
					F32 []float32
					F64 []float64
					S   []string
					T   []struct{}
				}{
					B:   []bool{true, false},
					I:   []int{-1},
					I8:  []int8{-10},
					I16: []int16{-100},
					I32: []int32{-1000},
					I64: []int64{-10000},
					U:   []uint{1},
					U8:  []uint8{10},
					U16: []uint16{100},
					U32: []uint32{1000},
					U64: []uint64{10000},
					F32: []float32{10.10},
					F64: []float64{100.100},
					S:   []string{"testing", "a", "string", "array"},
				},
			},
			"",
			false,
		},
		{
			"should error when trying to convert an unsupported type",
			[]string{
				"-t={\"field\": \"val\"}",
			},
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Lists.T": reflect.TypeOf([]testConfig{}),
				},
				out: &verboseConfig{},
			},
			args{
				map[string]string{
					"-t": "Lists.T",
				},
			},
			&verboseConfig{
				Lists: struct {
					B   []bool
					I   []int
					I8  []int8
					I16 []int16
					I32 []int32
					I64 []int64
					U   []uint
					U8  []uint8
					U16 []uint16
					U32 []uint32
					U64 []uint64
					F32 []float32
					F64 []float64
					S   []string
					T   []struct{}
				}{},
			},
			"unsupported field type",
			true,
		},
		{
			"should error on setting a field with the wrong type",
			[]string{"--name", "test name"},
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Name": reflect.TypeOf(true),
				},
				out: &testConfig{},
			},
			args{
				a: map[string]string{
					"--name": "Name",
				},
			},
			&testConfig{},
			"unable to set",
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
			if err := s.applyArgs(tt.args.a); err != nil {
				if !tt.wantErr {
					t.Errorf("settings.applyArgs() error = %v, wantErr %v", err, tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.errorMessage) {
					t.Errorf("settings.applyArgs() = %v, want %v", err.Error(), tt.errorMessage)
				}
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
		Bool   bool
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
			"should apply environment variable overrides to struct",
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Bool":         reflect.TypeOf(true),
					"Name":         reflect.TypeOf(""),
					"Nested.Count": reflect.TypeOf(1),
				},
				out: &testConfig{},
			},
			args{
				v: map[string]string{
					"BOOL":         "Bool",
					"NAME":         "Name",
					"NESTED_COUNT": "Nested.Count",
				},
			},
			map[string]string{
				"BOOL":         "true",
				"NAME":         "testing name assignment",
				"NESTED_COUNT": "10",
			},
			&testConfig{
				Bool:   true,
				Name:   "testing name assignment",
				Nested: struct{ Count int }{10},
			},
			false,
		},
		{
			"should error when trying to set a non supported field type",
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Nested": reflect.TypeOf(testConfig{}),
				},
				out: &testConfig{},
			},
			args{
				v: map[string]string{
					"NESTED": "Nested",
				},
			},
			map[string]string{
				"NESTED": "{ \"Count\": 10 }",
			},
			&testConfig{},
			true,
		},
		{
			"should return nil if varsMap is nil",
			fields{
				fieldTypeMap: map[string]reflect.Type{
					"Name":         reflect.TypeOf(""),
					"Nested.Count": reflect.TypeOf(1),
				},
				out: &testConfig{},
			},
			args{
				v: nil,
			},
			map[string]string{
				"NAME":         "testing name assignment",
				"NESTED_COUNT": "10",
			},
			&testConfig{},
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
				out: &verboseConfig{},
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
	type testConfig struct {
		Name    string
		Created time.Time
		Version string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name         string
		args         args
		want         *testConfig
		wantErr      bool
		errorMessage string
	}{
		{
			"should set fields when unmarshalling",
			args{
				path: "./tests/simple.yaml",
			},
			&testConfig{
				Name:    "example",
				Created: time.Date(2021, time.February, 16, 0, 0, 0, 0, time.UTC),
				Version: "1.1",
			},
			false,
			"",
		},
		{
			"should raise an error when path is blank",
			args{path: ""},
			&testConfig{},
			false,
			"",
		},
		{
			"should raise os.ErrNotexist if bad path",
			args{
				path: "./not/found.yml",
			},
			&testConfig{},
			true,
			"no such file",
		},
		{
			"should return invalid character on badly formatted file",
			args{
				path: "./tests/broken.json",
			},
			&testConfig{},
			true,
			"invalid character",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := settings{
				out: &testConfig{},
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
		filePattern string
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
				"",
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
				"",
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
				"",
				[]string{"GO_ENV"},
			},
			&testConfig{
				Name:    "example",
				Version: "1.1",
			},
			false,
		},
		{
			"should only apply override from base file named by environment when there is also a match on file pattern",
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
				"config.%s",
				[]string{"GO_ENV"},
			},
			&testConfig{
				Name:    "example",
				Version: "1.1",
			},
			false,
		},
		{
			"should apply override from environment override with file pattern",
			map[string]string{
				"GO_ENV": "pattern",
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
				"config.%s",
				[]string{"GO_ENV"},
			},
			&testConfig{
				Name:    "example-config-file-pattern",
				Version: "1.1",
			},
			false,
		},
		{
			"should apply override from environment override with file pattern and extension",
			map[string]string{
				"GO_ENV": "pattern",
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
				"config.%s.yml",
				[]string{"GO_ENV"},
			},
			&testConfig{
				Name:    "example-config-file-pattern",
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
				"",
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
				"",
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
			if err := s.searchForEnvOverrides(tt.args.vars, tt.args.searchPaths, tt.args.filePattern); (err != nil) != tt.wantErr {
				t.Errorf("settings.searchForEnvOverrides() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(s.out, tt.want) {
				t.Errorf("settings.searchForEnvOverrides() = %v, want %v", s.out, tt.want)
			}
		})

		os.Clearenv()
	}
}
