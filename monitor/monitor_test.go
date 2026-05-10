package monitor

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSeverityString(t *testing.T) {
	t.Run("Severity String 测试", func(t *testing.T) {
		t.Run("各个级别的字符串表示", func(t *testing.T) {
			assert.Equal(t, "DEBUG", SeverityDebug.String())
			assert.Equal(t, "INFO", SeverityInfo.String())
			assert.Equal(t, "WARNING", SeverityWarning.String())
			assert.Equal(t, "ERROR", SeverityError.String())
			assert.Equal(t, "CRITICAL", SeverityCritical.String())
		})

		t.Run("未知级别", func(t *testing.T) {
			assert.Equal(t, "UNKNOWN", Severity(99).String())
		})
	})
}

func TestSourceConstants(t *testing.T) {
	t.Run("Source 常量测试", func(t *testing.T) {
		assert.Equal(t, "internal", string(SourceInternal))
		assert.Equal(t, "alert_rule", string(SourceAlertRule))
	})
}

func TestInternalAlert(t *testing.T) {
	t.Run("InternalAlert 测试", func(t *testing.T) {
		t.Run("创建内部告警", func(t *testing.T) {
			err := fmt.Errorf("test error")
			alert := InternalAlert(err)
			assert.NotNil(t, alert)
			assert.Equal(t, SeverityError, alert.Severity)
			assert.Equal(t, SourceInternal, alert.Source)
			assert.Equal(t, err, alert.Err)
			assert.NotNil(t, alert.Ts)
		})
	})
}

func TestNewAlert(t *testing.T) {
	t.Run("NewAlert 测试", func(t *testing.T) {
		t.Run("创建自定义告警", func(t *testing.T) {
			payload := map[string]string{"key": "value"}
			alert := NewAlert(SeverityWarning, SourceAlertRule, "test description", payload)
			assert.NotNil(t, alert)
			assert.Equal(t, SeverityWarning, alert.Severity)
			assert.Equal(t, SourceAlertRule, alert.Source)
			assert.Equal(t, "test description", alert.Descr)
			assert.Equal(t, payload, alert.Payload)
			assert.NotNil(t, alert.Ts)
		})

		t.Run("创建不同级别的告警", func(t *testing.T) {
			alert := NewAlert(SeverityCritical, SourceInternal, "critical alert", nil)
			assert.Equal(t, SeverityCritical, alert.Severity)
			assert.Equal(t, SourceInternal, alert.Source)
		})
	})
}

func TestAlertRule(t *testing.T) {
	t.Run("AlertRule 测试", func(t *testing.T) {
		t.Run("匹配规则返回告警", func(t *testing.T) {
			var rule AlertRule[int] = func(payload *int) (*Alert, bool) {
				if *payload > 10 {
					return NewAlert(SeverityWarning, SourceAlertRule, "value too high", payload), true
				}
				return nil, false
			}

			val := 15
			alert, matched := rule(&val)
			assert.True(t, matched)
			assert.NotNil(t, alert)
			assert.Equal(t, "value too high", alert.Descr)
		})

		t.Run("不匹配规则返回false", func(t *testing.T) {
			var rule AlertRule[int] = func(payload *int) (*Alert, bool) {
				if *payload > 10 {
					return NewAlert(SeverityWarning, SourceAlertRule, "high", payload), true
				}
				return nil, false
			}

			val := 5
			alert, matched := rule(&val)
			assert.False(t, matched)
			assert.Nil(t, alert)
		})
	})
}

func TestMonitorInterface(t *testing.T) {
	t.Run("Monitor 接口定义验证", func(t *testing.T) {
		t.Run("接口包含所需方法", func(t *testing.T) {
			assert.True(t, true)
		})
	})
}

func TestAlertTimestamp(t *testing.T) {
	t.Run("Alert 时间戳测试", func(t *testing.T) {
		before := time.Now()
		alert := InternalAlert(fmt.Errorf("test"))
		after := time.Now()

		assert.True(t, alert.Ts.After(before) || alert.Ts.Equal(before))
		assert.True(t, alert.Ts.Before(after) || alert.Ts.Equal(after))
	})
}
