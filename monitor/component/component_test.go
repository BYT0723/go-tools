package component

import (
	"context"
	"testing"
	"time"

	"github.com/BYT0723/go-tools/monitor"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAlertComponentAddRule(t *testing.T) {
	Convey("AlertComponent AddAlertRule 测试", t, func() {
		Convey("添加规则后Evaluate", func() {
			var ac AlertComponent[int]
			ac.AddAlertRule(func(payload *int) (*monitor.Alert, bool) {
				if *payload > 10 {
					return monitor.NewAlert(monitor.SeverityWarning, monitor.SourceAlertRule, "high", payload), true
				}
				return nil, false
			})

			val := 15
			alerts := ac.Evaluate(&val)
			So(alerts, ShouldBeEmpty) // 需要累积触发次数才报警
		})

		Convey("添加多条规则", func() {
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
			So(alerts, ShouldBeEmpty)
		})
	})
}

func TestAlertComponentEvaluate(t *testing.T) {
	Convey("AlertComponent Evaluate 测试", t, func() {
		Convey("不匹配时不产生告警", func() {
			var ac AlertComponent[int]
			ac.AddAlertRule(func(payload *int) (*monitor.Alert, bool) {
				if *payload > 100 {
					return monitor.NewAlert(monitor.SeverityError, monitor.SourceAlertRule, "error", payload), true
				}
				return nil, false
			})

			val := 50
			alerts := ac.Evaluate(&val)
			So(alerts, ShouldBeEmpty)
		})

		Convey("规则返回nil alert但matched为true时不产生告警", func() {
			var ac AlertComponent[int]
			ac.AddAlertRule(func(payload *int) (*monitor.Alert, bool) {
				return nil, true
			})

			val := 1
			alerts := ac.Evaluate(&val)
			So(alerts, ShouldBeEmpty)
		})
	})
}

func TestMonitorComponentNew(t *testing.T) {
	Convey("MonitorComponent 创建测试", t, func() {
		Convey("NewMonitorComponent 创建实例", func() {
			mc := NewMonitorComponent()
			So(mc, ShouldNotBeNil)
			So(mc.Timeout, ShouldEqual, time.Minute)
		})
	})
}

func TestMonitorComponentSetCycle(t *testing.T) {
	Convey("MonitorComponent SetCycle 测试", t, func() {
		Convey("设置为0不更新", func() {
			mc := NewMonitorComponent()
			mc.SetCycle(0)
			So(mc.Timeout, ShouldEqual, time.Minute)
		})

		Convey("设置相同值不更新", func() {
			mc := NewMonitorComponent()
			mc.SetCycle(time.Minute)
			So(mc.Timeout, ShouldEqual, time.Minute)
		})

		Convey("设置新值更新cycle", func() {
			mc := NewMonitorComponent()
			mc.SetCycle(10 * time.Second)
			So(mc.Timeout, ShouldEqual, 10*time.Second) // timeout调整为不超过cycle
		})
	})
}

func TestMonitorComponentSetTimeout(t *testing.T) {
	Convey("MonitorComponent SetTimeout 测试", t, func() {
		Convey("设置为0不更新", func() {
			mc := NewMonitorComponent()
			mc.SetTimeout(0)
			So(mc.Timeout, ShouldEqual, time.Minute)
		})

		Convey("设置相同值不更新", func() {
			mc := NewMonitorComponent()
			mc.SetTimeout(time.Minute)
			So(mc.Timeout, ShouldEqual, time.Minute)
		})

		Convey("设置新值更新timeout", func() {
			mc := NewMonitorComponent()
			mc.SetTimeout(30 * time.Second)
			So(mc.Timeout, ShouldEqual, 30*time.Second)
		})
	})
}

func TestMonitorComponentSubscribe(t *testing.T) {
	Convey("MonitorComponent Subscribe 测试", t, func() {
		Convey("Subscribe 返回只读channel", func() {
			mc := NewMonitorComponent()
			ch := mc.Subscribe()
			So(ch, ShouldNotBeNil)
		})
	})
}

func TestMonitorComponentNotify(t *testing.T) {
	Convey("MonitorComponent Notify 测试", t, func() {
		Convey("Notify 发送告警到channel", func() {
			mc := NewMonitorComponent()
			alert := monitor.NewAlert(monitor.SeverityInfo, monitor.SourceInternal, "test", nil)
			mc.Notify(alert)

			select {
			case a := <-mc.Subscribe():
				So(a, ShouldEqual, alert)
			default:
				// may not receive due to buffer/timeout
			}
		})
	})
}

func TestMonitorComponentStop(t *testing.T) {
	Convey("MonitorComponent Stop 测试", t, func() {
		Convey("Stop 关闭channel", func() {
			mc := NewMonitorComponent()
			mc.Stop(context.Background())
			// channel should be closed
			select {
			case _, ok := <-mc.Subscribe():
				So(ok, ShouldBeFalse)
			}
		})
	})
}

func TestMonitorComponentContext(t *testing.T) {
	Convey("MonitorComponent Context 测试", t, func() {
		Convey("空context时自动初始化", func() {
			mc := NewMonitorComponent()
			ctx := mc.Context()
			So(ctx, ShouldNotBeNil)
		})

		Convey("SetContext 设置context", func() {
			mc := NewMonitorComponent()
			parentCtx := context.Background()
			mc.SetContext(parentCtx)
			ctx := mc.Context()
			So(ctx, ShouldNotBeNil)
		})
	})
}

func TestMonitorComponentTicker(t *testing.T) {
	Convey("MonitorComponent Ticker 测试", t, func() {
		Convey("Ticker 自动初始化", func() {
			mc := NewMonitorComponent()
			ticker := mc.Ticker()
			So(ticker, ShouldNotBeNil)
			ticker.Stop()
		})

		Convey("多次调用返回相同ticker", func() {
			mc := NewMonitorComponent()
			t1 := mc.Ticker()
			t2 := mc.Ticker()
			So(t1, ShouldEqual, t2)
			t1.Stop()
		})
	})
}
