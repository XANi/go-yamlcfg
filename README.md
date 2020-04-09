[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/XANi/go-yamlcfg)


# go-yamlcfg

Dead simple loading YAMLs as config files

First, define your config struct:

```go
type MyConfig struct {
    Address string
    LoginName string `yaml:"username"`
}
```

then, load:

```go

import "github.com/XANi/go-yamlcfg"

cfgFiles := []string{
    "$HOME/.config/my/cnf.yaml",
    "./cfg/config.yaml",
    "/etc/my/cnf.yaml",
}
var cfg MyConfig
err := yamlcfg.LoadConfig(cfgFiles, &cfg)
```

It will err out on:

* no readable file in config file list
* first file found was unparseable

## getting loaded config name

Define method `SetConfigPath(string)` on config struct like that:

```go
func (c *testCfg1)SetConfigPath(s string) {
	log.Infof("Loaded config file from %s",s)
}
```

## creating default config file

Define method `GetDefaultConfig() string` that returns default config, like that:

```go
var testCfg3Default = `---
test3: testing
`
func (c *testCfg3) GetDefaultConfig() string {
    return testCfg3Default
}
```

or, just return yaml directly if you do not need comments: {

```go
func (c *testCfg) GetDefaultConfig() string {
	defaultCfg := testCfg{
		User: "root",
		Pass: yamlcfg.RandomString(yamlcfg.CharsetAlphanumeric, 32)
	}
	out, err := yaml.Marshal(&defaultCfg)
	if err != nil {panic(fmt.Errorf("can't marshal [%T- %+v] into YAML: %s",defaultCfg,defaultCfg,err))}
	return string(out)
}
```

if `GetDefaultConfig()` returns empty string, config will not be created



## helpers

Default config (and any sub-dirs leading to it) will be created at first entry of cfgFiles, then loaded

# TODO

* basic include support
* validation and default values (altho preinializing config struct kinda does that now )
