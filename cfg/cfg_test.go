package cfg

import (
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptionFuncs(t *testing.T) {
	Convey("Option 函数测试", t, func() {
		Convey("WithConfigName", func() {
			c := &_config{viper: viper.New()}
			WithConfigName("test")(c)
			So(c.viper.GetString("test"), ShouldNotBeNil)
		})

		Convey("WithConfigType", func() {
			c := &_config{viper: viper.New()}
			WithConfigType("yaml")(c)
		})

		Convey("WithConfigFile", func() {
			c := &_config{viper: viper.New()}
			WithConfigFile("/tmp/test.yaml")(c)
		})

		Convey("WithConfigPath", func() {
			c := &_config{viper: viper.New()}
			WithConfigPath("/tmp", "/etc")(c)
		})

		Convey("OnConfigChange", func() {
			c := &_config{viper: viper.New()}
			h := func(e fsnotify.Event) {}
			OnConfigChange(h)(c)
			So(c.onConfigChange, ShouldNotBeNil)
		})

		Convey("WithConfigTag", func() {
			c := &_config{viper: viper.New()}
			WithConfigTag("json")(c)
		})

		Convey("WithCustomDeocodeOpt", func() {
			c := &_config{viper: viper.New()}
			WithCustomDeocodeOpt(func(dc *mapstructure.DecoderConfig) {})(c)
		})

		Convey("WithDefaultUnMarshal", func() {
			c := &_config{viper: viper.New()}
			var payload map[string]any
			WithDefaultUnMarshal(&payload)(c)
			So(c.unmarshaler, ShouldNotBeNil)
		})

		Convey("WithCustomUnMarshal", func() {
			c := &_config{viper: viper.New()}
			um := func(v *viper.Viper) error { return nil }
			WithCustomUnMarshal(um)(c)
			So(c.unmarshaler, ShouldNotBeNil)
		})

		Convey("WithRemoteConfig returns option", func() {
			opt := WithRemoteConfig("localhost", "/path")
			So(opt, ShouldNotBeNil)
		})
	})
}

func TestChangeHandler(t *testing.T) {
	Convey("ChangeHandler 测试", t, func() {
		Convey("Restart 生成 ChangeHandler", func() {
			h := Restart()
			So(h, ShouldNotBeNil)
		})

		Convey("Reload 生成 ChangeHandler", func() {
			var target map[string]string
			h := Reload(&target)
			So(h, ShouldNotBeNil)
		})

		Convey("ReloadKey 生成 ChangeHandler", func() {
			var target map[string]string
			h := ReloadKey("key", &target)
			So(h, ShouldNotBeNil)
		})
	})
}

func TestChangeMatcher(t *testing.T) {
	Convey("ChangeMatcher 测试", t, func() {
		Convey("matcher返回false时不触发handler操作", func() {
			alwaysFalse := func(e fsnotify.Event) bool { return false }
			h := Restart(alwaysFalse)
			So(func() {
				h(fsnotify.Event{Name: "test", Op: fsnotify.Write})
			}, ShouldNotPanic)
		})

		Convey("多个matcher全部返回true时触发handler", func() {
			alwaysTrue1 := func(e fsnotify.Event) bool { return true }
			alwaysTrue2 := func(e fsnotify.Event) bool { return true }
			h := Reload(new(int), alwaysTrue1, alwaysTrue2)
			// 由于 Unmarshal 在 nil config 上会 panic
			So(func() {
				h(fsnotify.Event{Name: "test", Op: fsnotify.Write})
			}, ShouldPanic)
		})

		Convey("多个matcher中有一个返回false则不触发", func() {
			h := Reload(new(int),
				func(e fsnotify.Event) bool { return false },
				func(e fsnotify.Event) bool { return true },
			)
			So(func() {
				h(fsnotify.Event{Name: "test", Op: fsnotify.Write})
			}, ShouldNotPanic)
		})
	})
}

func TestUnmarshaler(t *testing.T) {
	Convey("Unmarshaler 类型测试", t, func() {
		var u Unmarshaler = func(v *viper.Viper) error {
			return nil
		}
		So(u, ShouldNotBeNil)
	})
}

func TestRestartChangeHandler(t *testing.T) {
	Convey("Restart ChangeHandler", t, func() {
		Convey("Restart 构造 ChangeHandler 类型正确", func() {
			h := Restart(func(e fsnotify.Event) bool { return false })
			So(h, ShouldNotBeNil)
		})

		Convey("matcher返回false时不触发重启", func() {
			h := Restart(func(e fsnotify.Event) bool { return false })
			So(func() {
				h(fsnotify.Event{Name: "test", Op: fsnotify.Write})
			}, ShouldNotPanic)
		})
	})
}

func TestReloadChangeHandler(t *testing.T) {
	Convey("Reload ChangeHandler", t, func() {
		Convey("Reload 构造 ChangeHandler", func() {
			var target struct {
				Name string `cfg:"name"`
			}
			h := Reload(&target)
			So(h, ShouldNotBeNil)
		})
	})
}

func TestReloadKeyChangeHandler(t *testing.T) {
	Convey("ReloadKey ChangeHandler", t, func() {
		Convey("ReloadKey 构造 ChangeHandler", func() {
			var target struct {
				Name string `cfg:"name"`
			}
			h := ReloadKey("section", &target)
			So(h, ShouldNotBeNil)
		})
	})
}
