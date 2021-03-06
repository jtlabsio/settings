// package main in the example can be used to demo the settings
// parsing from command line, environment variables, specified
// override files and environment overrides. A few example commands
// to try:
//
// For an environment override (partial settings override)
// $ GO_ENV=example go run examples/example.go
//
// For a argument provided file override (partial settings override)
// $ go run examples/example.go --config-file examples/config-override.yml
//
// For an environment override (partial settings override) with additional
// env var settings
// $ GO_ENV=example LISTS_ANIMALS=cat,dog,bear,hare go run examples/example.go
//
// Try mixing and matching ENV VAR and command line arguments to observe the
// order of priority
// $ GO_ENV=example DATA_PORT=27019 go run examples/example.go --data-port 27018
package main

import (
	"fmt"
	"log"

	"go.jtlabs.io/settings"
)

// this is a silly struct for storing config / settings
// with example fields for demonstration purposes
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
	Lists   struct {
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
		}).
		SetVarsMap(map[string]string{
			"DATA_NAME":      "Data.Name",
			"DATA_HOST":      "Data.Host",
			"DATA_PORT":      "Data.Port",
			"SERVER_ADDRESS": "Server.Address",
			"LISTS_ANIMALS":  "Lists.Animals",
		})

	// read in configuration from all sources
	if err := settings.Gather(options, &c); err != nil {
		log.Fatal(err)
	}

	fmt.Println(c)
}
