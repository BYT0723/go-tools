package logcore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLoggerConf(t *testing.T) {
	t.Run("DefaultLoggerConf 测试", func(t *testing.T) {
		cfg := DefaultLoggerConf()
		assert.NotNil(t, cfg)
		assert.Equal(t, "logs", cfg.Dir)
		assert.Equal(t, "app", cfg.Name)
		assert.Equal(t, ".log", cfg.Ext)
		assert.Equal(t, "debug", cfg.Level)
		assert.False(t, cfg.Multi)
		assert.Equal(t, 20, cfg.MaxBackups)
		assert.Equal(t, 20, cfg.MaxSize)
		assert.Equal(t, 7, cfg.MaxAge)
		assert.True(t, cfg.Console)
	})
}

func TestLoggerConfMerge(t *testing.T) {
	t.Run("LoggerConf Merge 测试", func(t *testing.T) {
		t.Run("Merge 覆盖非零值字段", func(t *testing.T) {
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
			assert.Equal(t, "custom_logs", base.Dir)
			assert.Equal(t, "custom_app", base.Name)
			assert.Equal(t, "info", base.Level)
			assert.Equal(t, 10, base.MaxBackups)
			assert.Equal(t, 50, base.MaxSize)
			assert.False(t, base.Console)
		})

		t.Run("Merge 空字段不覆盖", func(t *testing.T) {
			base := DefaultLoggerConf()
			other := &LoggerConf{}
			base.Merge(other)
			assert.Equal(t, "logs", base.Dir)
			assert.Equal(t, "app", base.Name)
			assert.Equal(t, "debug", base.Level)
		})

		t.Run("Merge Multi 和 Console 直接赋值", func(t *testing.T) {
			base := DefaultLoggerConf()
			base.Multi = false
			base.Console = true
			other := &LoggerConf{Multi: true, Console: false}
			base.Merge(other)
			assert.True(t, base.Multi)
			assert.False(t, base.Console)
		})
	})
}

func TestField(t *testing.T) {
	t.Run("Field 结构测试", func(t *testing.T) {
		f := Field{Key: "test", Value: 42}
		assert.Equal(t, "test", f.Key)
		assert.Equal(t, 42, f.Value)
	})
}
