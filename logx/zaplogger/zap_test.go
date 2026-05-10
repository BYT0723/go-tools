package zaplogger

import (
	"reflect"
	"testing"

	"github.com/BYT0723/go-tools/logx/logcore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTransFields(t *testing.T) {
	Convey("transFields 测试", t, func() {
		Convey("Bool field", func() {
			fields := []logcore.Field{{Key: "active", Kind: reflect.Bool, Value: true}}
			result := transFields(fields)
			So(len(result), ShouldEqual, 1)
			So(result[0].Key, ShouldEqual, "active")
		})

		Convey("String field", func() {
			fields := []logcore.Field{{Key: "name", Kind: reflect.String, Value: "test"}}
			result := transFields(fields)
			So(len(result), ShouldEqual, 1)
		})

		Convey("Int field", func() {
			fields := []logcore.Field{{Key: "count", Kind: reflect.Int, Value: 42}}
			result := transFields(fields)
			So(len(result), ShouldEqual, 1)
		})

		Convey("Float64 field", func() {
			fields := []logcore.Field{{Key: "pi", Kind: reflect.Float64, Value: 3.14}}
			result := transFields(fields)
			So(len(result), ShouldEqual, 1)
		})

		Convey("Any field (default case)", func() {
			fields := []logcore.Field{{Key: "data", Value: map[string]int{"a": 1}}}
			result := transFields(fields)
			So(len(result), ShouldEqual, 1)
			So(result[0].Key, ShouldEqual, "data")
		})

		Convey("Error field from Value", func() {
			fields := []logcore.Field{{Key: "err", Value: "error message"}}
			result := transFields(fields)
			So(len(result), ShouldEqual, 1)
		})
	})
}

func TestNewInstance(t *testing.T) {
	Convey("NewInstance 测试", t, func() {
		Convey("有效配置创建实例", func() {
			cfg := logcore.DefaultLoggerConf()
			ins, err := NewInstance(cfg)
			So(err, ShouldBeNil)
			So(ins, ShouldNotBeNil)
		})

		Convey("无效level返回错误", func() {
			cfg := logcore.DefaultLoggerConf()
			cfg.Level = "invalid-level"
			_, err := NewInstance(cfg)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestZapLoggerMethods(t *testing.T) {
	Convey("zapLogger 方法测试", t, func() {
		cfg := logcore.DefaultLoggerConf()
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

		Convey("AddCallerSkip", func() {
			l2 := zl.AddCallerSkip(1)
			So(l2, ShouldNotBeNil)
		})

		Convey("Sync", func() {
			_ = zl.Sync()
		})
	})
}

func TestZapLoggerInterface(t *testing.T) {
	Convey("zapLogger 实现 Logger 接口", t, func() {
		cfg := logcore.DefaultLoggerConf()
		zl, err := NewInstance(cfg)
		So(err, ShouldBeNil)

		var l logcore.Logger = zl
		So(l, ShouldNotBeNil)
	})
}

func TestZapLoggerNullConfig(t *testing.T) {
	Convey("空配置创建实例", t, func() {
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
}

func TestNewConsoleCore(t *testing.T) {
	Convey("newConsoleCore 测试", t, func() {
		level := zap.NewAtomicLevel()
		core := newConsoleCore(level)
		So(core, ShouldNotBeNil)
	})
}

func TestNewCore(t *testing.T) {
	Convey("newCore 测试", t, func() {
		cfg := logcore.DefaultLoggerConf()
		level := zap.NewAtomicLevel()
		core := newCore(cfg, func(l zapcore.Level) bool { return l >= level.Level() }, "/tmp/test.log")
		So(core, ShouldNotBeNil)
	})
}
