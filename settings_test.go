package settings

import (
	"reflect"
	"testing"
)

type config struct {
	Data struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"data"`
	Logging struct {
		Level   string `yaml:"level"`
		Verbose bool   `yaml:"verbose"`
	} `yaml:"logging"`
	Name   string `yaml:"name"`
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
	Version string `yaml:"version"`
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
		LuckyNumbers []int    `yaml:"luckyNumbers"`
		Animals      []string `yaml:"animals"`
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
	Options().SetBasePath("./tests/test.yaml")
	tests := []struct {
		name     string
		args     args
		expected config
		wantErr  bool
	}{
		{
			"should set name and version",
			args{
				path: "./tests/test.yaml",
			},
			config{
				Name:    "example",
				Version: "1.0",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &settings{
				out: &config{},
			}
			if err := s.readBaseSettings(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("settings.readBaseSettings() error = %v, wantErr %v", err, tt.wantErr)
			}
			o := s.out
			if !reflect.DeepEqual(tt.expected, s.out) {
				t.Errorf("setting.readBaseSettings() = %v, want %v", *&(o), tt.expected)
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
