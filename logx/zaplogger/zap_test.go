package zaplogger

import (
	"reflect"
	"testing"

	"github.com/BYT0723/go-tools/logx/logcore"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
)

func TestTransFields(t *testing.T) {
	t.Run("transFields 测试", func(t *testing.T) {
		t.Run("Bool field", func(t *testing.T) {
			fields := []logcore.Field{{Key: "active", Kind: reflect.Bool, Value: true}}
			result := transFields(fields)
			assert.Equal(t, 1, len(result))
			assert.Equal(t, "active", result[0].Key)
		})

		t.Run("String field", func(t *testing.T) {
			fields := []logcore.Field{{Key: "name", Kind: reflect.String, Value: "test"}}
			result := transFields(fields)
			assert.Equal(t, 1, len(result))
		})

		t.Run("Int field", func(t *testing.T) {
			fields := []logcore.Field{{Key: "count", Kind: reflect.Int, Value: 42}}
			result := transFields(fields)
			assert.Equal(t, 1, len(result))
		})

		t.Run("Float64 field", func(t *testing.T) {
			fields := []logcore.Field{{Key: "pi", Kind: reflect.Float64, Value: 3.14}}
			result := transFields(fields)
			assert.Equal(t, 1, len(result))
		})

		t.Run("Any field (default case)", func(t *testing.T) {
			fields := []logcore.Field{{Key: "data", Value: map[string]int{"a": 1}}}
			result := transFields(fields)
			assert.Equal(t, 1, len(result))
			assert.Equal(t, "data", result[0].Key)
		})

		t.Run("Error field from Value", func(t *testing.T) {
			fields := []logcore.Field{{Key: "err", Value: "error message"}}
			result := transFields(fields)
			assert.Equal(t, 1, len(result))
		})
	})
}

func TestNewInstance(t *testing.T) {
	t.Run("NewInstance 测试", func(t *testing.T) {
		t.Run("有效配置创建实例", func(t *testing.T) {
			cfg := logcore.DefaultLoggerConf()
			ins, err := NewInstance(cfg)
			assert.Nil(t, err)
			assert.NotNil(t, ins)
		})

		t.Run("无效level返回错误", func(t *testing.T) {
			cfg := logcore.DefaultLoggerConf()
			cfg.Level = "invalid-level"
			_, err := NewInstance(cfg)
			assert.NotNil(t, err)
		})
	})
}

func TestZapLoggerMethods(t *testing.T) {
	t.Run("zapLogger 方法测试", func(t *testing.T) {
		cfg := logcore.DefaultLoggerConf()
		zl, err := NewInstance(cfg)
		assert.Nil(t, err)

		t.Run("Debug/Info/Warn/Error 不panic", func(t *testing.T) {
			assert.NotPanics(t, func() { zl.Debug("test") })
			assert.NotPanics(t, func() { zl.Info("test") })
			assert.NotPanics(t, func() { zl.Warn("test") })
			assert.NotPanics(t, func() { zl.Error("test") })
		})

		t.Run("Debugf/Infof/Warnf/Errorf 不panic", func(t *testing.T) {
			assert.NotPanics(t, func() { zl.Debugf("test %d", 1) })
			assert.NotPanics(t, func() { zl.Infof("test %d", 1) })
			assert.NotPanics(t, func() { zl.Warnf("test %d", 1) })
			assert.NotPanics(t, func() { zl.Errorf("test %d", 1) })
		})

		t.Run("Log/Logf 不panic", func(t *testing.T) {
			assert.NotPanics(t, func() { zl.Log("info", "test") })
			assert.NotPanics(t, func() { zl.Logf("info", "test %d", 1) })
		})

		t.Run("With 返回新Logger", func(t *testing.T) {
			l2 := zl.With(logcore.Field{Key: "key", Value: "value"})
			assert.NotNil(t, l2)
		})

		t.Run("AddCallerSkip", func(t *testing.T) {
			l2 := zl.AddCallerSkip(1)
			assert.NotNil(t, l2)
		})

		t.Run("Sync", func(t *testing.T) {
			_ = zl.Sync()
		})
	})
}

func TestZapLoggerInterface(t *testing.T) {
	t.Run("zapLogger 实现 Logger 接口", func(t *testing.T) {
		cfg := logcore.DefaultLoggerConf()
		zl, err := NewInstance(cfg)
		assert.Nil(t, err)

		var l logcore.Logger = zl
		assert.NotNil(t, l)
	})
}

func TestZapLoggerNullConfig(t *testing.T) {
	t.Run("空配置创建实例", func(t *testing.T) {
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
		assert.Nil(t, err)
		assert.NotNil(t, ins)
	})
}

func TestNewConsoleCore(t *testing.T) {
	t.Run("newConsoleCore 测试", func(t *testing.T) {
		level := zap.NewAtomicLevel()
		core := newConsoleCore(level)
		assert.NotNil(t, core)
	})
}

func TestNewCore(t *testing.T) {
	t.Run("newCore 测试", func(t *testing.T) {
		cfg := logcore.DefaultLoggerConf()
		level := zap.NewAtomicLevel()
		core := newCore(cfg, func(l zapcore.Level) bool { return l >= level.Level() }, "/tmp/test.log")
		assert.NotNil(t, core)
	})
}
