package noplogger

import (
	"testing"

	"github.com/BYT0723/go-tools/logx/logcore"

	"github.com/stretchr/testify/assert"
)

func TestNopLogger(t *testing.T) {
	t.Run("NopLogger 测试", func(t *testing.T) {
		l := NopLogger{}

		t.Run("所有方法不panic", func(t *testing.T) {
			assert.NotPanics(t, func() { l.Debug("test") })
			assert.NotPanics(t, func() { l.Debugf("test %s", "arg") })
			assert.NotPanics(t, func() { l.Info("test") })
			assert.NotPanics(t, func() { l.Infof("test %s", "arg") })
			assert.NotPanics(t, func() { l.Warn("test") })
			assert.NotPanics(t, func() { l.Warnf("test %s", "arg") })
			assert.NotPanics(t, func() { l.Error("test") })
			assert.NotPanics(t, func() { l.Errorf("test %s", "arg") })
			assert.NotPanics(t, func() { l.Log("info", "test") })
			assert.NotPanics(t, func() { l.Logf("info", "test %s", "arg") })
		})

		t.Run("Sync 返回nil", func(t *testing.T) {
			assert.Nil(t, l.Sync())
		})

		t.Run("With 返回自身", func(t *testing.T) {
			result := l.With(logcore.Field{Key: "k", Value: "v"})
			assert.EqualValues(t, l, result)
		})

		t.Run("AddCallerSkip 返回自身", func(t *testing.T) {
			result := l.AddCallerSkip(2)
			assert.EqualValues(t, l, result)
		})
	})
}

func TestNopLoggerInterface(t *testing.T) {
	t.Run("NopLogger 实现 Logger 接口", func(t *testing.T) {
		var l logcore.Logger = NopLogger{}
		assert.NotNil(t, l)
	})
}
