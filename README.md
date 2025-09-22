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

## Templating parsed config (env, secrets etc.)

Adding method `GetSecret(string) string` to a struct enables template parsing via `text/template` before parsing the YAML

Aside from standard `text/template` functions additional ones are available:

* `{{ secret "secretname"}}` will call `GetSecret("secretname") string` method on the config struct. 
  Connect any serial retrieval there. The method should report errors separately as there is not really sensible way to push errors up
* `{{ env "USER"}}` will call `os.Getenv`

Both outputs are string only so lack of key should be signalled with empty string if you want to base config logic on it. 
But do try to avoid making config into an application, this is not a library for that and this function is designed so
loading vars from k8s or environment is easier.

**Be warned**, the value is inserted into raw text of YAML so it is entirely possible to have undesirable config injected via ENV (if say it contains newlines),
so it should only be used when inputs are secure.

## Partial config parsing

If you need to have more flexible config format, say a plugin list with each plugin having its own separate config definition, 
you might want to use ast.Node` (from github.com/goccy/go-yaml/ast module) to specify a part of config as to be parsed later, like

```go
    type Partial struct {
	    Config map[string]ast.Node `yaml:"config"`
    }
	type SubPartial1 struct {
	    Option1 string `yaml:"option1"`
	    Option2 string `yaml:"option2"`
    }
	type SubPartial2 struct {
	    Option1 int `yaml:"option1"`
	    Option2 int `yaml:"option2"`
    }
    ...
	c := Partial{}
	err := LoadConfig([]string{"./t-data/t4.cfg"}, &c)
	... // somewhere in plugin1 code
	o1 := SubPartial1{}
    err = yaml.Unmarshal([]byte(c.Config["plugin1"].String()), &o1)
	... // somewhere in plugin2 code
	err = yaml.Unmarshal([]byte(c.Config["plugin2"].String()), &o2)
	require.NoError(t, err)
```


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

## validation / adding defaults

Add `Validate() error` method. It will be called at the end and err returned will be returned from `LoadConfig`. 
It is also a good place to put any default value handling.


## helpers

Default config (and any sub-dirs leading to it) will be created at first entry of cfgFiles, then loaded

# TODO

* basic include support
* validation and default values (altho preinializing config struct kinda does that now )
