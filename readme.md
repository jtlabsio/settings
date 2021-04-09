# 12-Factor Compliant Application Configuration



## Installation

```bash
go get -u go.jtlabs.io/settings
```

## Usage

```go
package main

import (
	"log"

	"go.jtlabs.io/settings"
)

type config struct {
	Data struct {
		Name string `json:"name"`
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"data"`
	Logging struct {
		Level string `json:"level"`
	} `json:"logging"`
	Server struct {
		Address string `json:"address"`
	} `json:"server"`
}

func main() {
	var c config
	options := settings.Options().
		SetBasePath("./defaults.yaml").
		SetSearchPaths("./", "./config", "./settings").
		SetDefaultsMap(map[string]interface{}{
			"Server.Address": ":3080",
		}).
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
}
```

## Q & A

### Why build this?

For our use case, we desired a simple and extensible configuration mechanism that spiritually adheres to 12-factor principles and facilitates the layered specification of application settings / configuration. An initial base configuration file to be the start, with an ability to specify override files as a layer of additional settings (either full or partial) on top of the base file (the locations for which specified via command line arguments or environment variables). Finally, having an ability to override individual keys within the configuration with specific environment variables or command line arguments. 

This settings library was built in the spirit of <https://github.com/brozeph/settings-lib>, which in many ways could be considered this package's Node.js flavored older sibling. This package is not a direct port of [settings-lib](https://github.com/brozeph/settings-lib), though, and makes use of Go specific idioms.

### How is this different than <https://github.com/spf13/viper>?

Viper is an incredible and feature rich configuration utility that also aligns, philosophically, with 12-factor principles. [Viper](https://github.com/spf13/viper) supports several features that this package does not:

* loading configuration from external sources (i.e. Consul, etcd, and k/v stores, etc.)
* reading configuration from more sources (i.e. HCL, INI and dotenv files, etc.)
* saving configuration back out to a destination

Where [Viper](https://github.com/spf13/viper) differs is in the order in which configuration is loaded. Additionally, to load additional full or partial files specified through command line arguments or environment variables, custom code is required. 

Ultimately, [Viper](https://github.com/spf13/viper) is a great choice for configuration as well. This package provides a subset of the functionality of [Viper](https://github.com/spf13/viper), and approaches the loading of configuration layers in a different order and with some nuance.

### Viper doesn't support case-senstive keys... does this library?

TL;DR: yes

The approach within Viper involves loading configuration from various sources and each source is a potential source of truth for the configuration value. As such, it would be a non-trivial matter for Viper to support case-sensitve key lookup. In this package, the source of truth for the configuration value is considered the `out interface{}` struct provided to the `Gather` method. All other sources for configuration are mapped to the desired output struct.