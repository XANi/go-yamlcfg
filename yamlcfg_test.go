package yamlcfg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

type testCfg1 struct {
	ConfigSource_ string
	Test1         string            `yaml:"test1"`
	Test2         map[string]string `yaml:"test2"`
	Test3         string
	Bool1         bool `yaml:"bool1"`
	Bool2         bool `yaml:"bool2"`
}

func (c *testCfg1) SetConfigPath(s string) {
	c.ConfigSource_ = s
}

type testCfg2 struct {
	Test1 string `yaml:"test1"`
}

type testCfg3 struct {
	Test3 string `yaml:"test3"`
	path  string
}

var testCfg3Default = `---
test3: testing
`

func (c *testCfg3) GetDefaultConfig() string {
	return testCfg3Default
}
func (c *testCfg3) SetConfigPath(s string) {
	c.path = s
}

func TestService(t *testing.T) {
	var c testCfg1
	y := testCfg1{
		ConfigSource_: "t-data/t1.cfg",
		Test1:         "aaa",
		Test2:         map[string]string{"test3": "123"},
		Test3:         "bbb",
		Bool1:         true,
		Bool2:         true,
	}
	t.Run("load yaml", func(t *testing.T) {
		err := LoadConfig([]string{"./t-data/t1.cfg"}, &c)
		require.NoError(t, err)
		assert.Contains(t, y.ConfigSource_, "t-data/t1.cfg")
		y.ConfigSource_ = c.ConfigSource_
		assert.Equal(t, y, c)
	})
	t.Run("nonexistent error", func(t *testing.T) {
		err := LoadConfig([]string{"/nonexisting/file"}, &c)
		assert.Error(t, err)
	})
	t.Run("partial config", func(t *testing.T) {
		var partial testCfg2
		err := LoadConfig([]string{"./t-data/t1.cfg"}, &partial)
		require.NoError(t, err)
		assert.Equal(t, "aaa", partial.Test1)
	})
	t.Run("no config", func(t *testing.T) {
		err := LoadConfig([]string{}, &c)
		require.Error(t, err)
	})
}

func TestDefaultGeneration(t *testing.T) {
	var c testCfg3
	tPath, err := filepath.Abs(os.ExpandEnv(`./t-data/t3.cfg.local`))
	require.NoError(t, err)
	// ensure file does not exist
	if _, err := os.Stat(tPath); err == nil {
		err := os.Remove(tPath)
		if err != nil {
			t.Fatal("test setup error - can't remove file")
		}
	}
	err = LoadConfig([]string{"./t-data/t3.cfg.local"}, &c)
	assert.NoError(t, err)
	assert.Contains(t, c.path, `/t-data/t3.cfg.local`)
}

func TestRandomString(t *testing.T) {
	str1 := RandomString("abcd1", 100)
	abcdRe := regexp.MustCompile(`^[abcd1]+$`)
	matched := abcdRe.MatchString(str1)
	assert.Len(t, str1, 100)
	assert.True(t, matched)
}

type testCfgValidate struct {
	Test1 string `yaml:"test1"`
}

func (c *testCfgValidate) Validate() error {
	return fmt.Errorf("err out on validate")
}

func TestValidate(t *testing.T) {
	var cfg testCfgValidate
	err := LoadConfig([]string{"./t-data/t1.cfg"}, &cfg)
	assert.EqualError(t, err, "err out on validate")
}
