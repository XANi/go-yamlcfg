package yamlcfg

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"crypto/rand"
	"math/big"
	"os"
	"io/ioutil"
	"path/filepath"
)

var CharsetAlphanumeric = "1234567890ABCDEFGHIKLMNOPQRSTUVWXYZabcdefghiklmnopqrstuvwxyz"
// LoadConfig loads first YAML file on cfgFiles list into cfg interface
// if cfg have defined GetDefaultConfig() method it will also create a default config file in first location if it can't find any config
// see README ( https://github.com/XANi/go-yamlcfg/blob/master/README.md ) for details
func LoadConfig(cfgFiles []string, cfg interface{}) (err error) {
	var cfgFile string
	if len(cfgFiles) < 1 { return fmt.Errorf("cfgFiles slice needs at least one element") }
	for _,element := range cfgFiles {
		filename, _ := filepath.Abs(os.ExpandEnv(element))
		if _, err := os.Stat(filename); err == nil {
			cfgFile = filename
			break
		}
	}
	if cfgFile == "" {
		if cfg, ok := cfg.(interface{GetDefaultConfig() string}); ok {
			defaultCfg := cfg.GetDefaultConfig()
			defaultCfgPath, _ := filepath.Abs(os.ExpandEnv(cfgFiles[0]))
			err := os.MkdirAll(filepath.Dir(defaultCfgPath),os.FileMode(0755))
			if err != nil {
				return fmt.Errorf("could not create directory; tried to create default config in [%s]: %s",defaultCfg,  err)
			}
			err = ioutil.WriteFile(defaultCfgPath,[]byte(defaultCfg), os.FileMode(0600))
			if err != nil {
				return fmt.Errorf("could not find config file; tried to create default config in [%s]: %s",defaultCfg, err)
			} else {
				cfgFile = defaultCfgPath
			}
		} else {
			return fmt.Errorf("could not find config file: %v", cfgFiles)
		}
	}
	//fmt.Sprintf("Loading config file: %s", cfgFile))
	raw_cfg, err := ioutil.ReadFile(cfgFile)
	if err != nil { return err }
	if len(raw_cfg) < 1 { return fmt.Errorf("Something gone wrong, file %s is 0 size or can't be read", cfgFile) }
	err = yaml.Unmarshal([]byte(raw_cfg), cfg)
	if err != nil { return err }
	if cfg, ok := cfg.(interface{SetConfigPath(string)}); ok {
		cfg.SetConfigPath(cfgFile)
	}

	return err
}
// RandomString generates random string with the defined dictionary
func RandomString(chars string, length int) string {
	b := make([]byte, length)
	chLength := big.NewInt(int64(len(chars)))
	for i := range b {
		r, err := rand.Int(rand.Reader,chLength)
		if err != nil {panic("error getting random data: " + err.Error())}
		b[i] = chars[r.Int64()]
  }
  return string(b)
}