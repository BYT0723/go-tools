package logx

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/logx/logcore"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFieldBuilders(t *testing.T) {
	Convey("Field 构建函数测试", t, func() {
		Convey("Any", func() {
			f := Any("key", 42)
			So(f.Key, ShouldEqual, "key")
			So(f.Value, ShouldEqual, 42)
		})

		Convey("Bool", func() {
			f := Bool("active", true)
			So(f.Key, ShouldEqual, "active")
			So(f.Kind, ShouldEqual, reflect.Bool)
			So(f.Value, ShouldEqual, true)
		})

		Convey("String", func() {
			f := String("name", "test")
			So(f.Key, ShouldEqual, "name")
			So(f.Kind, ShouldEqual, reflect.String)
			So(f.Value, ShouldEqual, "test")
		})

		Convey("Int", func() {
			f := Int("count", 42)
			So(f.Key, ShouldEqual, "count")
			So(f.Kind, ShouldEqual, reflect.Int)
			So(f.Value, ShouldEqual, 42)
		})

		Convey("Int8", func() {
			f := Int8("val", 8)
			So(f.Key, ShouldEqual, "val")
			So(f.Kind, ShouldEqual, reflect.Int8)
			So(f.Value, ShouldEqual, int8(8))
		})

		Convey("Int16", func() {
			f := Int16("val", 16)
			So(f.Kind, ShouldEqual, reflect.Int16)
			So(f.Value, ShouldEqual, int16(16))
		})

		Convey("Int32", func() {
			f := Int32("val", 32)
			So(f.Kind, ShouldEqual, reflect.Int32)
			So(f.Value, ShouldEqual, int32(32))
		})

		Convey("Int64", func() {
			f := Int64("val", 64)
			So(f.Kind, ShouldEqual, reflect.Int64)
			So(f.Value, ShouldEqual, int64(64))
		})

		Convey("Uint", func() {
			f := Uint("val", 100)
			So(f.Kind, ShouldEqual, reflect.Uint)
			So(f.Value, ShouldEqual, uint(100))
		})

		Convey("Uint8", func() {
			f := Uint8("val", 8)
			So(f.Kind, ShouldEqual, reflect.Uint8)
			So(f.Value, ShouldEqual, uint8(8))
		})

		Convey("Uint16", func() {
			f := Uint16("val", 16)
			So(f.Kind, ShouldEqual, reflect.Uint16)
		})

		Convey("Uint32", func() {
			f := Uint32("val", 32)
			So(f.Kind, ShouldEqual, reflect.Uint32)
		})

		Convey("Uint64", func() {
			f := Uint64("val", 64)
			So(f.Kind, ShouldEqual, reflect.Uint64)
		})

		Convey("Float32", func() {
			f := Float32("val", 3.14)
			So(f.Kind, ShouldEqual, reflect.Float32)
			So(f.Value, ShouldEqual, float32(3.14))
		})

		Convey("Float64", func() {
			f := Float64("val", 3.14159)
			So(f.Kind, ShouldEqual, reflect.Float64)
			So(f.Value, ShouldEqual, 3.14159)
		})

		Convey("Err", func() {
			e := errors.New("test error")
			f := Err(e)
			So(f.Key, ShouldEqual, "error")
			So(f.Value, ShouldEqual, "test error")
		})

		Convey("Duration", func() {
			f := Duration("elapsed", time.Second)
			So(f.Key, ShouldEqual, "elapsed")
			So(f.Value, ShouldEqual, time.Second)
		})
	})
}

func TestLoggerType(t *testing.T) {
	Convey("LoggerType 测试", t, func() {
		So(TypeZap, ShouldEqual, LoggerType(0))
		So(TypeZeroLog, ShouldEqual, LoggerType(1))
		So(TypeInvalid, ShouldEqual, LoggerType(15))
	})
}

func TestInitConf(t *testing.T) {
	Convey("InitConf 测试", t, func() {
		cfg := &InitConf{}
		So(cfg, ShouldNotBeNil)
	})
}

func TestOptionFuncs(t *testing.T) {
	Convey("Option 函数测试", t, func() {
		Convey("WithLoggerType 无效类型默认TypeZap", func() {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithLoggerType(TypeInvalid)(cfg)
			So(cfg.Type, ShouldEqual, TypeZap)
		})

		Convey("WithLoggerType TypeZeroLog", func() {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithLoggerType(TypeZeroLog)(cfg)
			So(cfg.Type, ShouldEqual, TypeZeroLog)
		})

		Convey("WithLevel", func() {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithLevel("error")(cfg)
			So(cfg.LogCfg.Level, ShouldEqual, "error")
		})

		Convey("WithName", func() {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithName("myapp")(cfg)
			So(cfg.LogCfg.Name, ShouldEqual, "myapp")
		})

		Convey("WithPath", func() {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithPath("/var/log")(cfg)
			So(cfg.LogCfg.Dir, ShouldEqual, "/var/log")
		})

		Convey("WithMaxBackups", func() {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithMaxBackups(100)(cfg)
			So(cfg.LogCfg.MaxBackups, ShouldEqual, 100)
		})

		Convey("WithMaxSize", func() {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithMaxSize(500)(cfg)
			So(cfg.LogCfg.MaxSize, ShouldEqual, 500)
		})

		Convey("WithMaxAge", func() {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithMaxAge(30)(cfg)
			So(cfg.LogCfg.MaxAge, ShouldEqual, 30)
		})

		Convey("WithConf 合并配置并清理Ext", func() {
			cfg := &InitConf{LogCfg: logcore.DefaultLoggerConf()}
			WithConf(&logcore.LoggerConf{Name: "merged", Ext: ".log"})(cfg)
			So(cfg.LogCfg.Name, ShouldEqual, "merged")
			So(cfg.LogCfg.Ext, ShouldEqual, ".log")
		})
	})
}

func TestFieldTypeAlias(t *testing.T) {
	Convey("Field 类型别名测试", t, func() {
		var f Field = logcore.Field{Key: "k", Value: "v"}
		So(f.Key, ShouldEqual, "k")
	})
}
