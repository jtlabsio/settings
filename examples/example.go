package main

import (
	"fmt"
	"log"

	"go.jtlabs.io/settings"
)

type config struct {
	Data struct {
		Name string `json:"name" yaml:"name"`
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"data" yaml:"data"`
	Logging struct {
		Level string `json:"level"`
	} `json:"logging"`
	Name   string `json:"name"`
	Server struct {
		Address string `json:"address"`
	} `json:"server"`
	Version string `json:"version"`
}

func main() {
	var c config
	options := settings.Options().
		SetBasePath("./examples/defaults.yaml").
		SetSearchPaths("./", "./config", "./settings", "./examples").
		SetDefaultsMap(map[string]interface{}{
			"Server.Address": ":8080",
		}).
		SetArgsFileOverride("--config-file", "-cf").
		SetArgsMap(map[string]string{
			"--data-name": "Data.Name",
			"--data-host": "Data.Host",
			"--data-port": "Data.Port",
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
