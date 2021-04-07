# 12-Factor Compliant Application Configuration

```go
package main

import "go.jtlabs.io/settings"

type config struct {
  Data struct {
    Name string `json:"name"`
    Host string `json:"host"`
    Port int `json:"port"`
  } `json:"data"`
  Logging struct {
    Level string `json:level`
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
      "Server.Address": ":3080"
    }).
    SetArgsMap(map[string]string{
      "--data-name": "Data.Name",
      "--data-host": "Data.Host",
      "--data-port": "Data.Port",
    }).
    SetVarsMap(map[string]string{
      "DATA_NAME": "Data.Name",
      "DATA_HOST": "Data.Host",
      "DATA_PORT": "Data.Port",
    })

  // read in configuration from all sources
  if err := settings.Read(options, &c); err {
    log.Fatal(err)
  }
}
```