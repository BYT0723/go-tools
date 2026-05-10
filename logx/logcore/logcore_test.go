package logcore

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultLoggerConf(t *testing.T) {
	Convey("DefaultLoggerConf 测试", t, func() {
		cfg := DefaultLoggerConf()
		So(cfg, ShouldNotBeNil)
		So(cfg.Dir, ShouldEqual, "logs")
		So(cfg.Name, ShouldEqual, "app")
		So(cfg.Ext, ShouldEqual, ".log")
		So(cfg.Level, ShouldEqual, "debug")
		So(cfg.Multi, ShouldBeFalse)
		So(cfg.MaxBackups, ShouldEqual, 20)
		So(cfg.MaxSize, ShouldEqual, 20)
		So(cfg.MaxAge, ShouldEqual, 7)
		So(cfg.Console, ShouldBeTrue)
	})
}

func TestLoggerConfMerge(t *testing.T) {
	Convey("LoggerConf Merge 测试", t, func() {
		Convey("Merge 覆盖非零值字段", func() {
			base := DefaultLoggerConf()
			other := &LoggerConf{
				Dir:        "custom_logs",
				Name:       "custom_app",
				Level:      "info",
				MaxBackups: 10,
				MaxSize:    50,
				Console:    false,
			}
			base.Merge(other)
			So(base.Dir, ShouldEqual, "custom_logs")
			So(base.Name, ShouldEqual, "custom_app")
			So(base.Level, ShouldEqual, "info")
			So(base.MaxBackups, ShouldEqual, 10)
			So(base.MaxSize, ShouldEqual, 50)
			So(base.Console, ShouldBeFalse)
		})

		Convey("Merge 空字段不覆盖", func() {
			base := DefaultLoggerConf()
			other := &LoggerConf{}
			base.Merge(other)
			So(base.Dir, ShouldEqual, "logs")
			So(base.Name, ShouldEqual, "app")
			So(base.Level, ShouldEqual, "debug")
		})

		Convey("Merge Multi 和 Console 直接赋值", func() {
			base := DefaultLoggerConf()
			base.Multi = false
			base.Console = true
			other := &LoggerConf{Multi: true, Console: false}
			base.Merge(other)
			So(base.Multi, ShouldBeTrue)
			So(base.Console, ShouldBeFalse)
		})
	})
}

func TestField(t *testing.T) {
	Convey("Field 结构测试", t, func() {
		f := Field{Key: "test", Value: 42}
		So(f.Key, ShouldEqual, "test")
		So(f.Value, ShouldEqual, 42)
	})
}
