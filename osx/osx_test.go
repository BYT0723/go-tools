package osx

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetTermSize(t *testing.T) {
	Convey("GetTermSize 测试", t, func() {
		Convey("调用GetTermSize不panic", func() {
			So(func() { GetTermSize() }, ShouldNotPanic)
		})
	})
}

func TestCharmapDecode(t *testing.T) {
	Convey("CharmapDecode 测试", t, func() {
		Convey("UTF-8 code page 65001 返回原始数据", func() {
			input := []byte("hello")
			output, err := CharmapDecode(65001, input)
			So(err, ShouldBeNil)
			So(string(output), ShouldEqual, "hello")
		})

		Convey("未知 code page 返回错误", func() {
			_, err := CharmapDecode(99999, []byte("hello"))
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unknown OEM Code Page")
		})

		Convey("GBK code page 936", func() {
			_, err := CharmapDecode(936, []byte("hello"))
			So(err, ShouldBeNil)
		})

		Convey("Shift-JIS code page 932", func() {
			_, err := CharmapDecode(932, []byte("hello"))
			So(err, ShouldBeNil)
		})
	})
}
