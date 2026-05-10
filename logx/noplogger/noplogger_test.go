package noplogger

import (
	"testing"

	"github.com/BYT0723/go-tools/logx/logcore"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNopLogger(t *testing.T) {
	Convey("NopLogger 测试", t, func() {
		l := NopLogger{}

		Convey("所有方法不panic", func() {
			So(func() { l.Debug("test") }, ShouldNotPanic)
			So(func() { l.Debugf("test %s", "arg") }, ShouldNotPanic)
			So(func() { l.Info("test") }, ShouldNotPanic)
			So(func() { l.Infof("test %s", "arg") }, ShouldNotPanic)
			So(func() { l.Warn("test") }, ShouldNotPanic)
			So(func() { l.Warnf("test %s", "arg") }, ShouldNotPanic)
			So(func() { l.Error("test") }, ShouldNotPanic)
			So(func() { l.Errorf("test %s", "arg") }, ShouldNotPanic)
			So(func() { l.Log("info", "test") }, ShouldNotPanic)
			So(func() { l.Logf("info", "test %s", "arg") }, ShouldNotPanic)
		})

		Convey("Sync 返回nil", func() {
			So(l.Sync(), ShouldBeNil)
		})

		Convey("With 返回自身", func() {
			result := l.With(logcore.Field{Key: "k", Value: "v"})
			So(result, ShouldResemble, l)
		})

		Convey("AddCallerSkip 返回自身", func() {
			result := l.AddCallerSkip(2)
			So(result, ShouldResemble, l)
		})
	})
}

func TestNopLoggerInterface(t *testing.T) {
	Convey("NopLogger 实现 Logger 接口", t, func() {
		var l logcore.Logger = NopLogger{}
		So(l, ShouldNotBeNil)
	})
}
