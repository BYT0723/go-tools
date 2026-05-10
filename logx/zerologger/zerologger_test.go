package zerologger

import (
	"bytes"
	"testing"

	"github.com/BYT0723/go-tools/logx/logcore"
	"github.com/rs/zerolog"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLevelWriter(t *testing.T) {
	Convey("LevelWriter 测试", t, func() {
		Convey("filter返回true时写入数据", func() {
			var buf bytes.Buffer
			w := NewLevelWriter(&buf, func(l zerolog.Level) bool {
				return l >= zerolog.WarnLevel
			})

			// filter returns true for warn level
			n, err := w.WriteLevel(zerolog.WarnLevel, []byte("hello"))
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 5)
			So(buf.String(), ShouldEqual, "hello")
		})

		Convey("filter返回false时不写入数据", func() {
			var buf bytes.Buffer
			w := NewLevelWriter(&buf, func(l zerolog.Level) bool {
				return l >= zerolog.WarnLevel
			})

			// filter returns false for debug level
			n, err := w.WriteLevel(zerolog.DebugLevel, []byte("debug msg"))
			So(err, ShouldBeNil)
			So(n, ShouldEqual, len("debug msg"))
			So(buf.String(), ShouldBeEmpty)
		})
	})
}

func TestNewInstance(t *testing.T) {
	Convey("NewInstance 测试", t, func() {
		Convey("有效配置创建实例", func() {
			cfg := &logcore.LoggerConf{
				Dir:        ".",
				Name:       "test",
				Ext:        ".log",
				Level:      "debug",
				MaxBackups: 1,
				MaxSize:    1,
				MaxAge:     1,
				Console:    false,
			}
			ins, err := NewInstance(cfg)
			So(err, ShouldBeNil)
			So(ins, ShouldNotBeNil)
		})

		Convey("无效level返回错误", func() {
			cfg := &logcore.LoggerConf{
				Level:      "invalid-level",
				Dir:        ".",
				Name:       "test",
				Ext:        ".log",
				MaxBackups: 1,
				MaxSize:    1,
				MaxAge:     1,
			}
			_, err := NewInstance(cfg)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestZeroLoggerMethods(t *testing.T) {
	Convey("zeroLogger 方法测试", t, func() {
		cfg := &logcore.LoggerConf{
			Dir:        ".",
			Name:       "test",
			Ext:        ".log",
			Level:      "debug",
			MaxBackups: 1,
			MaxSize:    1,
			MaxAge:     1,
			Console:    false,
		}
		zl, err := NewInstance(cfg)
		So(err, ShouldBeNil)

		Convey("Debug/Info/Warn/Error 不panic", func() {
			So(func() { zl.Debug("test") }, ShouldNotPanic)
			So(func() { zl.Info("test") }, ShouldNotPanic)
			So(func() { zl.Warn("test") }, ShouldNotPanic)
			So(func() { zl.Error("test") }, ShouldNotPanic)
		})

		Convey("Debugf/Infof/Warnf/Errorf 不panic", func() {
			So(func() { zl.Debugf("test %d", 1) }, ShouldNotPanic)
			So(func() { zl.Infof("test %d", 1) }, ShouldNotPanic)
			So(func() { zl.Warnf("test %d", 1) }, ShouldNotPanic)
			So(func() { zl.Errorf("test %d", 1) }, ShouldNotPanic)
		})

		Convey("Log/Logf 不panic", func() {
			So(func() { zl.Log("info", "test") }, ShouldNotPanic)
			So(func() { zl.Logf("info", "test %d", 1) }, ShouldNotPanic)
		})

		Convey("With 返回新Logger", func() {
			l2 := zl.With(logcore.Field{Key: "key", Value: "value"})
			So(l2, ShouldNotBeNil)
		})

		Convey("With typed fields", func() {
			f := logcore.Field{Key: "num", Value: int(42)}
			l2 := zl.With(f)
			So(l2, ShouldNotBeNil)
		})

		Convey("AddCallerSkip", func() {
			l2 := zl.AddCallerSkip(1)
			So(l2, ShouldNotBeNil)
		})

		Convey("Sync", func() {
			err := zl.Sync()
			So(err, ShouldBeNil)
		})
	})
}

func TestZeroLoggerInterface(t *testing.T) {
	Convey("zeroLogger 实现 Logger 接口", t, func() {
		cfg := logcore.DefaultLoggerConf()
		zl, err := NewInstance(cfg)
		So(err, ShouldBeNil)

		var l logcore.Logger = zl
		So(l, ShouldNotBeNil)
	})
}
