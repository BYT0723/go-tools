package component

import (
	"context"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/monitor"

	"github.com/stretchr/testify/assert"
)

func TestAlertComponentAddRule(t *testing.T) {
	t.Run("AlertComponent AddAlertRule 测试", func(t *testing.T) {
		t.Run("添加规则后Evaluate", func(t *testing.T) {
			var ac AlertComponent[int]
			ac.AddAlertRule(func(payload *int) (*monitor.Alert, bool) {
				if *payload > 10 {
					return monitor.NewAlert(monitor.SeverityWarning, monitor.SourceAlertRule, "high", payload), true
				}
				return nil, false
			})

			val := 15
			alerts := ac.Evaluate(&val)
			assert.Empty(t, alerts)
		})

		t.Run("添加多条规则", func(t *testing.T) {
			var ac AlertComponent[int]
			ac.AddAlertRule(
				func(payload *int) (*monitor.Alert, bool) {
					if *payload > 10 {
						return monitor.NewAlert(monitor.SeverityWarning, monitor.SourceAlertRule, "high", payload), true
					}
					return nil, false
				},
				func(payload *int) (*monitor.Alert, bool) {
					if *payload < 0 {
						return monitor.NewAlert(monitor.SeverityError, monitor.SourceAlertRule, "negative", payload), true
					}
					return nil, false
				},
			)

			val := 15
			alerts := ac.Evaluate(&val)
			assert.Empty(t, alerts)
		})
	})
}

func TestAlertComponentEvaluate(t *testing.T) {
	t.Run("AlertComponent Evaluate 测试", func(t *testing.T) {
		t.Run("不匹配时不产生告警", func(t *testing.T) {
			var ac AlertComponent[int]
			ac.AddAlertRule(func(payload *int) (*monitor.Alert, bool) {
				if *payload > 100 {
					return monitor.NewAlert(monitor.SeverityError, monitor.SourceAlertRule, "error", payload), true
				}
				return nil, false
			})

			val := 50
			alerts := ac.Evaluate(&val)
			assert.Empty(t, alerts)
		})

		t.Run("规则返回nil alert但matched为true时不产生告警", func(t *testing.T) {
			var ac AlertComponent[int]
			ac.AddAlertRule(func(payload *int) (*monitor.Alert, bool) {
				return nil, true
			})

			val := 1
			alerts := ac.Evaluate(&val)
			assert.Empty(t, alerts)
		})
	})
}

func TestMonitorComponentNew(t *testing.T) {
	t.Run("MonitorComponent 创建测试", func(t *testing.T) {
		t.Run("NewMonitorComponent 创建实例", func(t *testing.T) {
			mc := NewMonitorComponent()
			assert.NotNil(t, mc)
			assert.Equal(t, time.Minute, mc.Timeout)
		})
	})
}

func TestMonitorComponentSetCycle(t *testing.T) {
	t.Run("MonitorComponent SetCycle 测试", func(t *testing.T) {
		t.Run("设置为0不更新", func(t *testing.T) {
			mc := NewMonitorComponent()
			mc.SetCycle(0)
			assert.Equal(t, time.Minute, mc.Timeout)
		})

		t.Run("设置相同值不更新", func(t *testing.T) {
			mc := NewMonitorComponent()
			mc.SetCycle(time.Minute)
			assert.Equal(t, time.Minute, mc.Timeout)
		})

		t.Run("设置新值更新cycle", func(t *testing.T) {
			mc := NewMonitorComponent()
			mc.SetCycle(10 * time.Second)
			assert.Equal(t, 10*time.Second, mc.Timeout)
		})
	})
}

func TestMonitorComponentSetTimeout(t *testing.T) {
	t.Run("MonitorComponent SetTimeout 测试", func(t *testing.T) {
		t.Run("设置为0不更新", func(t *testing.T) {
			mc := NewMonitorComponent()
			mc.SetTimeout(0)
			assert.Equal(t, time.Minute, mc.Timeout)
		})

		t.Run("设置相同值不更新", func(t *testing.T) {
			mc := NewMonitorComponent()
			mc.SetTimeout(time.Minute)
			assert.Equal(t, time.Minute, mc.Timeout)
		})

		t.Run("设置新值更新timeout", func(t *testing.T) {
			mc := NewMonitorComponent()
			mc.SetTimeout(30 * time.Second)
			assert.Equal(t, 30*time.Second, mc.Timeout)
		})
	})
}

func TestMonitorComponentSubscribe(t *testing.T) {
	t.Run("MonitorComponent Subscribe 测试", func(t *testing.T) {
		t.Run("Subscribe 返回只读channel", func(t *testing.T) {
			mc := NewMonitorComponent()
			ch := mc.Subscribe()
			assert.NotNil(t, ch)
		})
	})
}

func TestMonitorComponentNotify(t *testing.T) {
	t.Run("MonitorComponent Notify 测试", func(t *testing.T) {
		t.Run("Notify 发送告警到channel", func(t *testing.T) {
			mc := NewMonitorComponent()
			alert := monitor.NewAlert(monitor.SeverityInfo, monitor.SourceInternal, "test", nil)
			mc.Notify(alert)

			select {
			case a := <-mc.Subscribe():
				assert.Equal(t, alert, a)
			default:
			}
		})
	})
}

func TestMonitorComponentStop(t *testing.T) {
	t.Run("MonitorComponent Stop 测试", func(t *testing.T) {
		t.Run("Stop 关闭channel", func(t *testing.T) {
			mc := NewMonitorComponent()
			mc.Stop(context.Background())
			select {
			case _, ok := <-mc.Subscribe():
				assert.False(t, ok)
			}
		})
	})
}

func TestMonitorComponentContext(t *testing.T) {
	t.Run("MonitorComponent Context 测试", func(t *testing.T) {
		t.Run("空context时自动初始化", func(t *testing.T) {
			mc := NewMonitorComponent()
			ctx := mc.Context()
			assert.NotNil(t, ctx)
		})

		t.Run("SetContext 设置context", func(t *testing.T) {
			mc := NewMonitorComponent()
			parentCtx := context.Background()
			mc.SetContext(parentCtx)
			ctx := mc.Context()
			assert.NotNil(t, ctx)
		})
	})
}

func TestMonitorComponentTicker(t *testing.T) {
	t.Run("MonitorComponent Ticker 测试", func(t *testing.T) {
		t.Run("Ticker 自动初始化", func(t *testing.T) {
			mc := NewMonitorComponent()
			ticker := mc.Ticker()
			assert.NotNil(t, ticker)
			ticker.Stop()
		})

		t.Run("多次调用返回相同ticker", func(t *testing.T) {
			mc := NewMonitorComponent()
			t1 := mc.Ticker()
			t2 := mc.Ticker()
			assert.Equal(t, t2, t1)
			t1.Stop()
		})
	})
}
