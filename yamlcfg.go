package yamlcfg

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"os"
	"io/ioutil"
	"path/filepath"
)


func LoadConfig(cfgFiles []string, cfg interface{}) (err error) {
	var cfgFile string
	for _,element := range cfgFiles {
		filename, _ := filepath.Abs(os.ExpandEnv(element))
		if _, err := os.Stat(filename); err == nil {
			cfgFile = filename
			break
		}
	}
	if cfgFile == "" {
		return fmt.Errorf("could not find config file: %v", cfgFiles)
	}
	//fmt.Sprintf("Loading config file: %s", cfgFile))
	raw_cfg, err := ioutil.ReadFile(cfgFile)
	err = yaml.Unmarshal([]byte(raw_cfg), cfg)
	if cfg, ok := cfg.(interface{SetConfigPath(string)}); ok {
		cfg.SetConfigPath(cfgFile)
	}

	return err
}
