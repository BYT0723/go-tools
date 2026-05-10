package logx

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/logx/logcore"

	"github.com/stretchr/testify/assert"
)

func TestFieldBuilders(t *testing.T) {
	t.Run("Field 构建函数测试", func(t *testing.T) {
		t.Run("Any", func(t *testing.T) {
			f := Any("key", 42)
			assert.Equal(t, "key", f.Key)
			assert.Equal(t, 42, f.Value)
		})

		t.Run("Bool", func(t *testing.T) {
			f := Bool("active", true)
			assert.Equal(t, "active", f.Key)
			assert.Equal(t, reflect.Bool, f.Kind)
			assert.Equal(t, true, f.Value)
		})

		t.Run("String", func(t *testing.T) {
			f := String("name", "test")
			assert.Equal(t, "name", f.Key)
			assert.Equal(t, reflect.String, f.Kind)
			assert.Equal(t, "test", f.Value)
		})

		t.Run("Int", func(t *testing.T) {
			f := Int("count", 42)
			assert.Equal(t, "count", f.Key)
			assert.Equal(t, reflect.Int, f.Kind)
			assert.Equal(t, 42, f.Value)
		})

		t.Run("Int8", func(t *testing.T) {
			f := Int8("val", 8)
			assert.Equal(t, "val", f.Key)
			assert.Equal(t, reflect.Int8, f.Kind)
			assert.Equal(t, int8(8), f.Value)
		})

		t.Run("Int16", func(t *testing.T) {
			f := Int16("val", 16)
			assert.Equal(t, reflect.Int16, f.Kind)
			assert.Equal(t, int16(16), f.Value)
		})

		t.Run("Int32", func(t *testing.T) {
			f := Int32("val", 32)
			assert.Equal(t, reflect.Int32, f.Kind)
			assert.Equal(t, int32(32), f.Value)
		})

		t.Run("Int64", func(t *testing.T) {
			f := Int64("val", 64)
			assert.Equal(t, reflect.Int64, f.Kind)
			assert.Equal(t, int64(64), f.Value)
		})

		t.Run("Uint", func(t *testing.T) {
			f := Uint("val", 100)
			assert.Equal(t, reflect.Uint, f.Kind)
			assert.Equal(t, uint(100), f.Value)
		})

		t.Run("Uint8", func(t *testing.T) {
			f := Uint8("val", 8)
			assert.Equal(t, reflect.Uint8, f.Kind)
			assert.Equal(t, uint8(8), f.Value)
		})

		t.Run("Uint16", func(t *testing.T) {
			f := Uint16("val", 16)
			assert.Equal(t, reflect.Uint16, f.Kind)
		})

		t.Run("Uint32", func(t *testing.T) {
			f := Uint32("val", 32)
			assert.Equal(t, reflect.Uint32, f.Kind)
		})

		t.Run("Uint64", func(t *testing.T) {
			f := Uint64("val", 64)
			assert.Equal(t, reflect.Uint64, f.Kind)
		})

		t.Run("Float32", func(t *testing.T) {
			f := Float32("val", 3.14)
			assert.Equal(t, reflect.Float32, f.Kind)
			assert.Equal(t, float32(3.14), f.Value)
		})

		t.Run("Float64", func(t *testing.T) {
			f := Float64("val", 3.14159)
			assert.Equal(t, reflect.Float64, f.Kind)
			assert.Equal(t, 3.14159, f.Value)
		})

		t.Run("Err", func(t *testing.T) {
			e := errors.New("test error")
			f := Err(e)
			assert.Equal(t, "error", f.Key)
			assert.Equal(t, "test error", f.Value)
		})

		t.Run("Duration", func(t *testing.T) {
			f := Duration("elapsed", time.Second)
			assert.Equal(t, "elapsed", f.Key)
			assert.Equal(t, time.Second, f.Value)
		})
	})
}

func TestLoggerType(t *testing.T) {
	t.Run("LoggerType 测试", func(t *testing.T) {
		assert.Equal(t, LoggerType(0), TypeZap)
		assert.Equal(t, LoggerType(1), TypeZeroLog)
		assert.Equal(t, LoggerType(15), TypeInvalid)
	})
}

func TestInitConf(t *testing.T) {
	t.Run("InitConf 测试", func(t *testing.T) {
		cfg := &InitConf{}
		assert.NotNil(t, cfg)
	})
}

func TestOptionFuncs(t *testing.T) {
	t.Run("Option 函数测试", func(t *testing.T) {
		t.Run("WithLoggerType 无效类型默认TypeZap", func(t *testing.T) {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithLoggerType(TypeInvalid)(cfg)
			assert.Equal(t, TypeZap, cfg.Type)
		})

		t.Run("WithLoggerType TypeZeroLog", func(t *testing.T) {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithLoggerType(TypeZeroLog)(cfg)
			assert.Equal(t, TypeZeroLog, cfg.Type)
		})

		t.Run("WithLevel", func(t *testing.T) {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithLevel("error")(cfg)
			assert.Equal(t, "error", cfg.LogCfg.Level)
		})

		t.Run("WithName", func(t *testing.T) {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithName("myapp")(cfg)
			assert.Equal(t, "myapp", cfg.LogCfg.Name)
		})

		t.Run("WithPath", func(t *testing.T) {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithPath("/var/log")(cfg)
			assert.Equal(t, "/var/log", cfg.LogCfg.Dir)
		})

		t.Run("WithMaxBackups", func(t *testing.T) {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithMaxBackups(100)(cfg)
			assert.Equal(t, 100, cfg.LogCfg.MaxBackups)
		})

		t.Run("WithMaxSize", func(t *testing.T) {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithMaxSize(500)(cfg)
			assert.Equal(t, 500, cfg.LogCfg.MaxSize)
		})

		t.Run("WithMaxAge", func(t *testing.T) {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithMaxAge(30)(cfg)
			assert.Equal(t, 30, cfg.LogCfg.MaxAge)
		})

		t.Run("WithConf 合并配置并清理Ext", func(t *testing.T) {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithConf(&logcore.LoggerConf{Name: "merged", Ext: ".log"})(cfg)
			assert.Equal(t, "merged", cfg.LogCfg.Name)
			assert.Equal(t, ".log", cfg.LogCfg.Ext)
		})
	})
}

func TestFieldTypeAlias(t *testing.T) {
	t.Run("Field 类型别名测试", func(t *testing.T) {
		var f Field = logcore.Field{Key: "k", Value: "v"}
		assert.Equal(t, "k", f.Key)
	})
}
