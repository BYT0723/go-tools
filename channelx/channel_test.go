package channelx

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTryIn(t *testing.T) {
	Convey("TryIn 测试", t, func() {
		Convey("发送到缓冲channel成功", func() {
			ch := make(chan int, 1)
			So(TryIn(ch, 42), ShouldBeTrue)
			So(<-ch, ShouldEqual, 42)
		})

		Convey("发送到满的channel失败", func() {
			ch := make(chan int, 1)
			ch <- 1
			So(TryIn(ch, 42), ShouldBeFalse)
		})

		Convey("发送到无缓冲channel失败", func() {
			ch := make(chan int)
			So(TryIn(ch, 42), ShouldBeFalse)
		})
	})
}

func TestTryOut(t *testing.T) {
	Convey("TryOut 测试", t, func() {
		Convey("从有数据的channel读取成功", func() {
			ch := make(chan int, 1)
			ch <- 42
			v, ok := TryOut(ch)
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, 42)
		})

		Convey("从空的channel读取失败", func() {
			ch := make(chan int, 1)
			v, ok := TryOut(ch)
			So(ok, ShouldBeFalse)
			So(v, ShouldEqual, 0)
		})
	})
}

func TestIn(t *testing.T) {
	Convey("In 测试", t, func() {
		Convey("正常发送成功", func() {
			ch := make(chan int, 1)
			err := In(context.Background(), ch, 42)
			So(err, ShouldBeNil)
			So(<-ch, ShouldEqual, 42)
		})

		Convey("context取消时返回错误", func() {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			ch := make(chan int)
			err := In(ctx, ch, 42)
			So(err, ShouldEqual, context.Canceled)
		})
	})
}

func TestOut(t *testing.T) {
	Convey("Out 测试", t, func() {
		Convey("正常读取成功", func() {
			ch := make(chan int, 1)
			ch <- 42
			v, err := Out(context.Background(), ch)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 42)
		})

		Convey("context取消时返回错误", func() {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			ch := make(chan int)
			_, err := Out(ctx, ch)
			So(err, ShouldEqual, context.Canceled)
		})
	})
}

func TestInTimeout(t *testing.T) {
	Convey("InTimeout 测试", t, func() {
		Convey("超时内发送成功", func() {
			ch := make(chan int, 1)
			err := InTimeout(ch, 42, time.Second)
			So(err, ShouldBeNil)
			So(<-ch, ShouldEqual, 42)
		})

		Convey("超时发送失败", func() {
			ch := make(chan int)
			err := InTimeout(ch, 42, time.Millisecond)
			So(err, ShouldEqual, context.DeadlineExceeded)
		})
	})
}

func TestOutTimeout(t *testing.T) {
	Convey("OutTimeout 测试", t, func() {
		Convey("超时内读取成功", func() {
			ch := make(chan int, 1)
			ch <- 42
			v, err := OutTimeout(ch, time.Second)
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 42)
		})

		Convey("超时读取失败", func() {
			ch := make(chan int)
			_, err := OutTimeout(ch, time.Millisecond)
			So(err, ShouldEqual, context.DeadlineExceeded)
		})
	})
}

func TestInDeadline(t *testing.T) {
	Convey("InDeadline 测试", t, func() {
		Convey("截止时间前发送成功", func() {
			ch := make(chan int, 1)
			err := InDeadline(ch, 42, time.Now().Add(time.Second))
			So(err, ShouldBeNil)
			So(<-ch, ShouldEqual, 42)
		})

		Convey("截止时间后发送失败", func() {
			ch := make(chan int)
			err := InDeadline(ch, 42, time.Now().Add(-time.Second))
			So(err, ShouldEqual, context.DeadlineExceeded)
		})
	})
}

func TestOutDeadline(t *testing.T) {
	Convey("OutDeadline 测试", t, func() {
		Convey("截止时间前读取成功", func() {
			ch := make(chan int, 1)
			ch <- 42
			v, err := OutDeadline(ch, time.Now().Add(time.Second))
			So(err, ShouldBeNil)
			So(v, ShouldEqual, 42)
		})

		Convey("截止时间后读取失败", func() {
			ch := make(chan int)
			_, err := OutDeadline(ch, time.Now().Add(-time.Second))
			So(err, ShouldEqual, context.DeadlineExceeded)
		})
	})
}
