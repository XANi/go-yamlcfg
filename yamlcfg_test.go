package yamlcfg

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
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

func TestService(t *testing.T) {
	var c testCfg1
	err := LoadConfig([]string{"./test/t1.cfg"}, &c)
	Convey("LoadConfig", t, func() {
		So(err, ShouldEqual, nil)
		So(c.Test1,ShouldEqual,"aaa")
		So(c.Test2["test3"],ShouldEqual,"123")
		So(c.Test3,ShouldEqual,"bbb")
	})
	Convey("Show Config source", t, func() {
		So(c.ConfigSource_,ShouldContainSubstring,"test/t1.cfg")
	})
	err = LoadConfig([]string{"/nonexisting/file"}, &c)
	Convey("Error on no file found",t,func() {
		So(err,ShouldNotEqual,nil)

	})
	var partial testCfg2
	err = LoadConfig([]string{"./test/t1.cfg"}, &partial)
	Convey("Load partial config", t, func() {
		So(err, ShouldEqual, nil)
		So(c.Test1,ShouldEqual,"aaa")
		})

}
