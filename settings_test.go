package settings

import (
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

func Test_settings_readBaseSettings(t *testing.T) {
	type args struct {
		path string
	}
	// Options().SetBasePath("./tests/test.yaml")
	tests := []struct {
		name         string
		args         args
		expected     *config
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
			if err := s.readBaseSettings(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("settings.readBaseSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
			o := &s.out
			if !reflect.DeepEqual(tt.expected, s.out) {
				t.Errorf("settings.readBaseSettings() = %v, want %v", o, tt.expected)
			}
			err := s.readBaseSettings(tt.args.path)
			terr := tt.errorMessage
			tlen := len(terr)

			if tlen > 0 && !strings.Contains(err.Error(), terr) {
				t.Errorf("settings.readBaseSettings() = %v, want %v", err.Error(), tt.errorMessage)
			}
		})
	}
}

func Test_settings_determineFileType(t *testing.T) {
	type fields struct {
		fieldTypeMap map[string]reflect.Type
		out          interface{}
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &settings{
				fieldTypeMap: tt.fields.fieldTypeMap,
				out:          tt.fields.out,
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
