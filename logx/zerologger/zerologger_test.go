package zerologger

import (
	"bytes"
	"testing"

	"github.com/BYT0723/go-tools/logx/logcore"
	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"
)

func TestLevelWriter(t *testing.T) {
	t.Run("LevelWriter 测试", func(t *testing.T) {
		t.Run("filter返回true时写入数据", func(t *testing.T) {
			var buf bytes.Buffer
			w := NewLevelWriter(&buf, func(l zerolog.Level) bool {
				return l >= zerolog.WarnLevel
			})

			n, err := w.WriteLevel(zerolog.WarnLevel, []byte("hello"))
			assert.Nil(t, err)
			assert.Equal(t, 5, n)
			assert.Equal(t, "hello", buf.String())
		})

		t.Run("filter返回false时不写入数据", func(t *testing.T) {
			var buf bytes.Buffer
			w := NewLevelWriter(&buf, func(l zerolog.Level) bool {
				return l >= zerolog.WarnLevel
			})

			n, err := w.WriteLevel(zerolog.DebugLevel, []byte("debug msg"))
			assert.Nil(t, err)
			assert.Equal(t, len("debug msg"), n)
			assert.Empty(t, buf.String())
		})
	})
}

func TestNewInstance(t *testing.T) {
	t.Run("NewInstance 测试", func(t *testing.T) {
		t.Run("有效配置创建实例", func(t *testing.T) {
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

		t.Run("无效level返回错误", func(t *testing.T) {
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
			assert.NotNil(t, err)
		})
	})
}

func TestZeroLoggerMethods(t *testing.T) {
	t.Run("zeroLogger 方法测试", func(t *testing.T) {
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

		t.Run("With typed fields", func(t *testing.T) {
			f := logcore.Field{Key: "num", Value: int(42)}
			l2 := zl.With(f)
			assert.NotNil(t, l2)
		})

		t.Run("AddCallerSkip", func(t *testing.T) {
			l2 := zl.AddCallerSkip(1)
			assert.NotNil(t, l2)
		})

		t.Run("Sync", func(t *testing.T) {
			err := zl.Sync()
			assert.Nil(t, err)
		})
	})
}

func TestZeroLoggerInterface(t *testing.T) {
	t.Run("zeroLogger 实现 Logger 接口", func(t *testing.T) {
		cfg := logcore.DefaultLoggerConf()
		zl, err := NewInstance(cfg)
		assert.Nil(t, err)

		var l logcore.Logger = zl
		assert.NotNil(t, l)
	})
}
