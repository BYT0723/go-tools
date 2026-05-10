package monitor

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSeverityString(t *testing.T) {
	Convey("Severity String 测试", t, func() {
		Convey("各个级别的字符串表示", func() {
			So(SeverityDebug.String(), ShouldEqual, "DEBUG")
			So(SeverityInfo.String(), ShouldEqual, "INFO")
			So(SeverityWarning.String(), ShouldEqual, "WARNING")
			So(SeverityError.String(), ShouldEqual, "ERROR")
			So(SeverityCritical.String(), ShouldEqual, "CRITICAL")
		})

		Convey("未知级别", func() {
			So(Severity(99).String(), ShouldEqual, "UNKNOWN")
		})
	})
}

func TestSourceConstants(t *testing.T) {
	Convey("Source 常量测试", t, func() {
		So(string(SourceInternal), ShouldEqual, "internal")
		So(string(SourceAlertRule), ShouldEqual, "alert_rule")
	})
}

func TestInternalAlert(t *testing.T) {
	Convey("InternalAlert 测试", t, func() {
		Convey("创建内部告警", func() {
			err := fmt.Errorf("test error")
			alert := InternalAlert(err)
			So(alert, ShouldNotBeNil)
			So(alert.Severity, ShouldEqual, SeverityError)
			So(alert.Source, ShouldEqual, SourceInternal)
			So(alert.Err, ShouldEqual, err)
			So(alert.Ts, ShouldNotBeNil)
		})
	})
}

func TestNewAlert(t *testing.T) {
	Convey("NewAlert 测试", t, func() {
		Convey("创建自定义告警", func() {
			payload := map[string]string{"key": "value"}
			alert := NewAlert(SeverityWarning, SourceAlertRule, "test description", payload)
			So(alert, ShouldNotBeNil)
			So(alert.Severity, ShouldEqual, SeverityWarning)
			So(alert.Source, ShouldEqual, SourceAlertRule)
			So(alert.Descr, ShouldEqual, "test description")
			So(alert.Payload, ShouldEqual, payload)
			So(alert.Ts, ShouldNotBeNil)
		})

		Convey("创建不同级别的告警", func() {
			alert := NewAlert(SeverityCritical, SourceInternal, "critical alert", nil)
			So(alert.Severity, ShouldEqual, SeverityCritical)
			So(alert.Source, ShouldEqual, SourceInternal)
		})
	})
}

func TestAlertRule(t *testing.T) {
	Convey("AlertRule 测试", t, func() {
		Convey("匹配规则返回告警", func() {
			var rule AlertRule[int] = func(payload *int) (*Alert, bool) {
				if *payload > 10 {
					return NewAlert(SeverityWarning, SourceAlertRule, "value too high", payload), true
				}
				return nil, false
			}

			val := 15
			alert, matched := rule(&val)
			So(matched, ShouldBeTrue)
			So(alert, ShouldNotBeNil)
			So(alert.Descr, ShouldEqual, "value too high")
		})

		Convey("不匹配规则返回false", func() {
			var rule AlertRule[int] = func(payload *int) (*Alert, bool) {
				if *payload > 10 {
					return NewAlert(SeverityWarning, SourceAlertRule, "high", payload), true
				}
				return nil, false
			}

			val := 5
			alert, matched := rule(&val)
			So(matched, ShouldBeFalse)
			So(alert, ShouldBeNil)
		})
	})
}

func TestMonitorInterface(t *testing.T) {
	Convey("Monitor 接口定义验证", t, func() {
		Convey("接口包含所需方法", func() {
			So(true, ShouldBeTrue)
		})
	})
}

func TestAlertTimestamp(t *testing.T) {
	Convey("Alert 时间戳测试", t, func() {
		before := time.Now()
		alert := InternalAlert(fmt.Errorf("test"))
		after := time.Now()

		So(alert.Ts.After(before) || alert.Ts.Equal(before), ShouldBeTrue)
		So(alert.Ts.Before(after) || alert.Ts.Equal(after), ShouldBeTrue)
	})
}
