package main

import (
	"fmt"
	"log"

	"go.jtlabs.io/settings"
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

func main() {
	var c config
	options := settings.Options().
		SetBasePath("./examples/defaults.yaml").
		SetDefaultsMap(map[string]interface{}{
			"Server.Address": ":8080",
		}).
		SetEnvOverride("GO_ENV").
		SetEnvSearchPaths("./examples", "./", "./config", "./settings").
		SetArgsFileOverride("--config-file", "-cf").
		SetArgsMap(map[string]string{
			"--data-name":           "Data.Name",
			"--data-host":           "Data.Host",
			"--data-port":           "Data.Port",
			"--lists-animals":       "Lists.Animals",
			"--lists-lucky-numbers": "Lists.LuckyNumbers",
			"--logging-verbose":     "Logging.Verbose",
			"--num-v8":              "Numbers.V8",
			"--num-v16":             "Numbers.V16",
			"--num-v32":             "Numbers.V32",
			"--num-v64":             "Numbers.V64",
			"--num-u8":              "Numbers.U8",
			"--num-u16":             "Numbers.U16",
			"--num-u32":             "Numbers.U32",
			"--num-u64":             "Numbers.U64",
			"--num-f32":             "Numbers.F32",
			"--num-f64":             "Numbers.F64",
		}).
		SetVarsMap(map[string]string{
			"DATA_NAME":      "Data.Name",
			"DATA_HOST":      "Data.Host",
			"DATA_PORT":      "Data.Port",
			"SERVER_ADDRESS": "Server.Address",
		})

	// read in configuration from all sources
	if err := settings.Gather(options, &c); err != nil {
		log.Fatal(err)
	}

	fmt.Println(c)
}
