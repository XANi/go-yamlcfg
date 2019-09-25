package yamlcfg

import (
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
	"testing"
	"path/filepath"
	"os"
)

type testCfg1 struct {
	ConfigSource_ string
	Test1 string `yaml:"test1"`
	Test2 map[string]string `yaml:"test2"`
	Test3 string
}

func (c *testCfg1)SetConfigPath(s string) {
	c.ConfigSource_ = s
}

type testCfg2 struct {
	Test1 string `yaml:"test1"`
}

type testCfg3 struct {
	Test3 string `yaml:"test3"`
	path string
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
	err := LoadConfig([]string{"./t-data/t1.cfg"}, &c)
	Convey("LoadConfig", t, func() {
		So(err, ShouldEqual, nil)
		So(c.Test1,ShouldEqual,"aaa")
		So(c.Test2["test3"],ShouldEqual,"123")
		So(c.Test3,ShouldEqual,"bbb")
	})
	Convey("Show Config source", t, func() {
		So(c.ConfigSource_,ShouldContainSubstring,"t-data/t1.cfg")
	})
	err = LoadConfig([]string{"/nonexisting/file"}, &c)
	Convey("Error on no file found",t,func() {
		So(err,ShouldNotEqual,nil)
	})
	var partial testCfg2
	err = LoadConfig([]string{"./t-data/t1.cfg"}, &partial)
	Convey("Load partial config", t, func() {
		So(err, ShouldEqual, nil)
		So(c.Test1,ShouldEqual,"aaa")
	})
	err = LoadConfig([]string{}, &c)
	Convey("No config file in slice",t,func() {
		So(err,ShouldNotBeNil)
	})
}

func TestDefaultGeneration(t *testing.T) {
	var c testCfg3
	tPath, err := filepath.Abs(os.ExpandEnv(`./t-data/t3.cfg.local`))
	if err != nil {
		Convey("test setup error - can't get path",t,func() {
			So(err,ShouldEqual,nil)
		})
	}
	// ensure file does not exist
	if _, err := os.Stat(tPath); err == nil {
		err := os.Remove(tPath)
		if err != nil {
			Convey("test setup error - can't remove file",t,func() {
				So(err,ShouldEqual,nil)
			})
		}
	}
	err = LoadConfig([]string{"./t-data/t3.cfg.local"}, &c)
	Convey("Create default config",t,func() {
		So(err,ShouldEqual,nil)
		So(c.path,ShouldContainSubstring,`/t-data/t3.cfg.local`)
	})


}

func TestRandomString(t *testing.T) {
	str1 := RandomString("abcd1",100)
	abcdRe := regexp.MustCompile(`^[abcd1]+$`)
	matched := abcdRe.MatchString(str1)
	Convey("string should contain only abcd1: " + str1 , t, func () {
		So(len(str1), ShouldEqual, 100)
		So(matched,ShouldBeTrue)
	})
}
