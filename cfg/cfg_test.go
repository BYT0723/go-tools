package cfg

import (
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func TestOptionFuncs(t *testing.T) {
	t.Run("Option 函数测试", func(t *testing.T) {
		t.Run("WithConfigName", func(t *testing.T) {
			c := &_config{viper: viper.New()}
			WithConfigName("test")(c)
			assert.NotNil(t, c.viper.GetString("test"))
		})

		t.Run("WithConfigType", func(t *testing.T) {
			c := &_config{viper: viper.New()}
			WithConfigType("yaml")(c)
		})

		t.Run("WithConfigFile", func(t *testing.T) {
			c := &_config{viper: viper.New()}
			WithConfigFile("/tmp/test.yaml")(c)
		})

		t.Run("WithConfigPath", func(t *testing.T) {
			c := &_config{viper: viper.New()}
			WithConfigPath("/tmp", "/etc")(c)
		})

		t.Run("OnConfigChange", func(t *testing.T) {
			c := &_config{viper: viper.New()}
			h := func(e fsnotify.Event) {}
			OnConfigChange(h)(c)
			assert.NotNil(t, c.onConfigChange)
		})

		t.Run("WithConfigTag", func(t *testing.T) {
			c := &_config{viper: viper.New()}
			WithConfigTag("json")(c)
		})

		t.Run("WithCustomDeocodeOpt", func(t *testing.T) {
			c := &_config{viper: viper.New()}
			WithCustomDeocodeOpt(func(dc *mapstructure.DecoderConfig) {})(c)
		})

		t.Run("WithDefaultUnMarshal", func(t *testing.T) {
			c := &_config{viper: viper.New()}
			var payload map[string]any
			WithDefaultUnMarshal(&payload)(c)
			assert.NotNil(t, c.unmarshaler)
		})

		t.Run("WithCustomUnMarshal", func(t *testing.T) {
			c := &_config{viper: viper.New()}
			um := func(v *viper.Viper) error { return nil }
			WithCustomUnMarshal(um)(c)
			assert.NotNil(t, c.unmarshaler)
		})

		t.Run("WithRemoteConfig returns option", func(t *testing.T) {
			opt := WithRemoteConfig("localhost", "/path")
			assert.NotNil(t, opt)
		})
	})
}

func TestChangeHandler(t *testing.T) {
	t.Run("ChangeHandler 测试", func(t *testing.T) {
		t.Run("Restart 生成 ChangeHandler", func(t *testing.T) {
			h := Restart()
			assert.NotNil(t, h)
		})

		t.Run("Reload 生成 ChangeHandler", func(t *testing.T) {
			var target map[string]string
			h := Reload(&target)
			assert.NotNil(t, h)
		})

		t.Run("ReloadKey 生成 ChangeHandler", func(t *testing.T) {
			var target map[string]string
			h := ReloadKey("key", &target)
			assert.NotNil(t, h)
		})
	})
}

func TestChangeMatcher(t *testing.T) {
	t.Run("ChangeMatcher 测试", func(t *testing.T) {
		t.Run("matcher返回false时不触发handler操作", func(t *testing.T) {
			alwaysFalse := func(e fsnotify.Event) bool { return false }
			h := Restart(alwaysFalse)
			assert.NotPanics(t, func() {
				h(fsnotify.Event{Name: "test", Op: fsnotify.Write})
			})
		})

		t.Run("多个matcher全部返回true时触发handler", func(t *testing.T) {
			alwaysTrue1 := func(e fsnotify.Event) bool { return true }
			alwaysTrue2 := func(e fsnotify.Event) bool { return true }
			h := Reload(new(int), alwaysTrue1, alwaysTrue2)
			// 由于 Unmarshal 在 nil config 上会 panic
			assert.Panics(t, func() {
				h(fsnotify.Event{Name: "test", Op: fsnotify.Write})
			})
		})

		t.Run("多个matcher中有一个返回false则不触发", func(t *testing.T) {
			h := Reload(new(int),
				func(e fsnotify.Event) bool { return false },
				func(e fsnotify.Event) bool { return true },
			)
			assert.NotPanics(t, func() {
				h(fsnotify.Event{Name: "test", Op: fsnotify.Write})
			})
		})
	})
}

func TestUnmarshaler(t *testing.T) {
	t.Run("Unmarshaler 类型测试", func(t *testing.T) {
		var u Unmarshaler = func(v *viper.Viper) error {
			return nil
		}
		assert.NotNil(t, u)
	})
}

func TestRestartChangeHandler(t *testing.T) {
	t.Run("Restart ChangeHandler", func(t *testing.T) {
		t.Run("Restart 构造 ChangeHandler 类型正确", func(t *testing.T) {
			h := Restart(func(e fsnotify.Event) bool { return false })
			assert.NotNil(t, h)
		})

		t.Run("matcher返回false时不触发重启", func(t *testing.T) {
			h := Restart(func(e fsnotify.Event) bool { return false })
			assert.NotPanics(t, func() {
				h(fsnotify.Event{Name: "test", Op: fsnotify.Write})
			})
		})
	})
}

func TestReloadChangeHandler(t *testing.T) {
	t.Run("Reload ChangeHandler", func(t *testing.T) {
		t.Run("Reload 构造 ChangeHandler", func(t *testing.T) {
			var target struct {
				Name string `cfg:"name"`
			}
			h := Reload(&target)
			assert.NotNil(t, h)
		})
	})
}

func TestReloadKeyChangeHandler(t *testing.T) {
	t.Run("ReloadKey ChangeHandler", func(t *testing.T) {
		t.Run("ReloadKey 构造 ChangeHandler", func(t *testing.T) {
			var target struct {
				Name string `cfg:"name"`
			}
			h := ReloadKey("section", &target)
			assert.NotNil(t, h)
		})
	})
}
