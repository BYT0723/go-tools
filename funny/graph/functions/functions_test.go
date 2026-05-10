package functions

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBinaryFunctionType(t *testing.T) {
	Convey("BinaryFunction 类型测试", t, func() {
		var f BinaryFunction = func(x, y float64) bool { return x == y }
		So(f, ShouldNotBeNil)
	})
}

func TestLove(t *testing.T) {
	Convey("Love 函数测试", t, func() {
		f := Love()
		So(f, ShouldNotBeNil)

		// 爱心中心点应在函数范围内
		So(f(0, 0), ShouldBeTrue)
		// 远离中心的点不应在范围内
		So(f(10, 10), ShouldBeFalse)
	})
}

func TestCircularLove(t *testing.T) {
	Convey("CircularLove 函数测试", t, func() {
		f := CircularLove()
		So(f, ShouldNotBeNil)
	})
}

func TestRoseLine(t *testing.T) {
	Convey("RoseLine 函数测试", t, func() {
		f := RoseLine(4, 1.0)
		So(f, ShouldNotBeNil)
	})
}
