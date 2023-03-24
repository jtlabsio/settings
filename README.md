# 12-Factor Compliant Application Configuration

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/go.jtlabs.io/settings) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/jtlabsio/settings/main/LICENSE) [![codecov](https://codecov.io/gh/jtlabsio/settings/branch/main/graph/badge.svg?token=TO0CS5E9UH)](https://codecov.io/gh/jtlabsio/settings) [![GoReportCard example](https://goreportcard.com/badge/github.com/jtlabsio/settings)](https://goreportcard.com/report/github.com/jtlabsio/settings)

This package gathers values (typically used for application settings and configuration) from various sources outside of the application, and layers them together in a single output struct (supplied as an argument) for use by the application. This package supports [12-factor](https://12factor.net) configuration use cases to facilitate cloud-native API and application development.

The package will first attempt to load settings from the following sources in the order arranged below:
1. a base file (in `yaml` or `json` format)
2. from the default values map (if provided in `ReadOptions`)
3. from any command line provided override files (if `ArgsFileOverride` switches are defined in `ReadOptions`)
4. from any environment override files (if `EnvOverride` and `EnvSearchPaths` are provided in `ReadOptions`)
5. from all command line argument field specific value overrides (if `ArgsMap` is provided in `ReadOptions`)
6. from all environment variable field specific value overrides (if `VarsMap` is provided in `ReadOptions`)

## Installation

```bash
go get -u go.jtlabs.io/settings
```

## Usage

The following snippet demonstrates creating `settings.ReadOptions` with instructions that the `settings.Gather` function uses to populate the supplied config struct:

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

For a more verbose example along with execution instructions, see [examples/example.go](examples/example.go).

### ReadOptions

ReadOptions are used to instruct the package where to find override values from a base file, a command line override file, an environment override file, command line arguments, or from environment variables.

#### EnvDefault

This function exists so that default environment variable driven overrides, similar to those defined in [settings-lib](https://github.com/brozeph/settings-lib), can be provided easily to the Gather function.

```go
options := settings.Options().EnvDefault()
settings.Gather(options, &config)
```

#### SetArg

Similar to [`SetArgsMap`](#SetArgsMap), this can be used to attach command line arguments, individually, to fields for settings.

```go
options := settings.Options().
  SetArg("--a-flag", "Field.Name")
settings.Gather(options, &config)
```

#### SetArgsFileOverride

When providing a value to this method, one can override the underlying settings via one or more specific files that are provided via command line arguments.

```go
options := settings.
  Options().
  SetArgsFileOverride("./path/to/first-file.yml", "./path/to/another/second-file.json")
settings.Gather(options, &config)
```

The `first-file.yml` will be read and applied, and then the `second-file.json` will be read and applied over the top of the first. These files can be partial files with a subset of the fields from the out struct defined as `&config` in the example above if desired.

#### SetArgsMap

The arguments map is used by `Gather` to determine how command line switches can be applied to specific out struct fields.

```go
options := settings.
  Options().
  SetArgsMap(map[string]string{
    "--switch-to-look-for": "CaseSensitive.Field.Where.Hiearchy.Is.Noted.By.Dot",
    "--logging-level":      "Logging.Level",
    "-l":                   "Logging.Level",
  })
settings.Gather(options, &config)
```

Any switches that are provided in the map that do not appear in the list of `os.Args` for the application are effectively ignored. If the desired outcome is to have an alias for a command line argument (i.e. `--logging-level` and `-l` both capable of overriding `Logging.Level`), each value can be independently added to the map. When processing arguments, `--some-switch=value` (notice the `=` character) is processed the same as `--some-switch value` so that the value will properly read and applied in either scenario.

#### SetBasePath

The base path for settings is the initial (yaml or json) file that is loaded to populate the out argument to the gather method. As with the command line override file and with the environment override file, this base settings file is not required to be a complete serialization of the out struct... it can be partially defined if desired. If a file is specified, and the file can't be found or read, the `Gather` method will return a file doesn't exist (i.e. `os.ErrNotExist`) or a `SettingsFileReadError` in the event there is some other read problem.

```go
options := settings.
  Options().
  SetBasePath("./settings.yml")
settings.Gather(options, &config)
```

#### SetDefaultsMap

The defaults map is used by settings to apply default values to fields in the out struct. These defaults are applied immediately after the base settings (if provided) are applied.

```go
options := settings.
  Options().
  SetArgsMap(map[string]interface{}{
    "CaseSensitive.Field.Where.Hiearchy.Is.Noted.By.Dot": true,
    "Data.Port":                                          27017,
    "Server.Address":                                     ":8080",
    "Logging.Level":                                      "trace",
    "Name":                                               "cool name",
  })
settings.Gather(options, &config)
```

The string value of the map is the field path where hierarchy / depth is noted by the `.` character.

#### SetEnvOverride and SetEnvSearchPaths and SetEnvSearchPattern

Environment override and search paths can be provided to the package to enable virtually named environment level overrides at a partial or complete configuration level.

```go
options := settings.
  Options().
  SetEnvSearchPaths("./", "./settings"). // look for files in "./" and "./settings
  SetEnvOverride("GO_ENV", "GO_ENVIRONMENT")
settings.Gather(options, &config)
```

In the above example, if a value is set in the `GO_ENV` or `GO_ENVIRONMENT` variables for the application, the value will be used in a search for matching `yaml` or `json` files that exist in the paths provided as search paths (in the above example, `./` and `./settings`). To illustrate:

```bash
GO_ENV=testing go run cmd/app.go
```

The `GO_ENV` value is `testing`. Combined with the code snippet above, the app would search for the following files:

* `./testing.yml`
* `./testing.yaml`
* `./testing.json`
* `./settings/testing.yml`
* `./settings/testing.yaml`
* `./settings/testing.json`

Upon finding a file that matches (the first match), that file is read and the fields defined therein are applied to the out struct.

If `SetEnvSearchPattern` is used to defined a file name pattern, in addition to the above steps, files are searched using the file name pattern provided...

```go
options := settings.
  Options().
  SetEnvSearchPaths("./", "./settings"). // look for files in "./" and "./settings
  SetEnvSearchPattern("config.%s").
  SetEnvOverride("GO_ENV", "GO_ENVIRONMENT")
settings.Gather(options, &config)
```

Using the above example, when the following is used to start the application:

```bash
GO_ENV=testing go run cmd/app.go
```

The settings package will look for the following files:

* `./config.testing.yml`
* `./config.testing.yaml`
* `./config.testing.json`
* `./settings/config.testing.yml`
* `./settings/config.testing.yaml`
* `./settings/config.testing.json`

In this scenario, if both `./testing.yml` and `./config.testing.yml` are found, only the `./testing.yml` will be loaded.

#### SetVar

Similar to [`SetVarsMap`](#SetVarMap), this can be used to associate environment variables, individually, to fields for settings.

```go
options := settings.Options().
  SetVar("FIELD_NAME", "Field.Name")
settings.Gather(options, &config)
```

#### SetVarsMap

Similar to the Args map, the Vars map can be used to override individual fields with values defined as environment variables.

## Q & A

### Why build this?

For our use case, we desired a simple and extensible configuration mechanism that spiritually adheres to 12-factor principles and facilitates the layered specification of nn settings / configuration. An initial base configuration file to be the start, with an ability to specify override files as a layer of additional settings (either full or partial) on top of the base file (the locations for which specified via command line arguments or environment variables). Finally, having an ability to override individual keys within the configuration with specific environment variables or command line arguments.

This settings library was built in the spirit of <https://github.com/brozeph/settings-lib>, which in many ways could be considered this package's Node.js flavored older sibling. This package is not a direct port of [settings-lib](https://github.com/brozeph/settings-lib), though, and makes use of Go specific idioms.

### How is this different than <https://github.com/spf13/viper>?

Viper is an incredible and feature rich configuration utility that also aligns, philosophically, with 12-factor principles. [Viper](https://github.com/spf13/viper) supports several features that this package does not:

* loading configuration from external sources (i.e. Consul, etcd, and k/v stores, etc.)
* reading configuration from more sources (i.e. HCL, INI, TOML, dotenv files, etc.)
* saving configuration back out to a destination

Where [Viper](https://github.com/spf13/viper) differs is in the order in which configuration is loaded. Additionally, to load additional full or partial files specified through command line arguments or environment variables, custom code is required.

Ultimately, [Viper](https://github.com/spf13/viper) is a great choice for configuration as well. This package provides a subset of the functionality of [Viper](https://github.com/spf13/viper), and approaches the loading of configuration layers in a different order and with a completely different underlying approach.

### Viper doesn't support case-senstive keys... does this library?

TL;DR: yes, but this isn't really applicable in this package.

The approach within Viper involves loading configuration from various sources and each source is a potential origin of the configuration value. It would be a non-trivial matter for Viper to support case-sensitve key lookup.

In this package, there is a single source of truth for the final configuration, which is the `out interface{}` struct provided to the `Gather` method. All other sources for configuration are mapped to the desired output struct.