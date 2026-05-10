package ds

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLinkStackPush(t *testing.T) {
	Convey("LinkStack Push 测试", t, func() {
		Convey("Push 空栈", func() {
			s := NewLinkStack[int]()
			s.Push(1)
			So(s.Size(), ShouldEqual, 1)
			So(s.Empty(), ShouldBeFalse)
		})

		Convey("Push 多个元素", func() {
			s := NewLinkStack[int]()
			for i := 0; i < 10; i++ {
				s.Push(i)
			}
			So(s.Size(), ShouldEqual, 10)
		})
	})
}

func TestLinkStackPop(t *testing.T) {
	Convey("LinkStack Pop 测试", t, func() {
		Convey("空栈 Pop panic", func() {
			s := NewLinkStack[int]()
			So(func() { s.Pop() }, ShouldPanic)
		})

		Convey("Pop 返回最后Push的元素", func() {
			s := NewLinkStack[int]()
			s.Push(1)
			s.Push(2)
			s.Push(3)
			v, ok := s.Pop()
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, 3)
		})
	})
}

func TestLinkStackPeek(t *testing.T) {
	Convey("LinkStack Peek 测试", t, func() {
		Convey("空栈 Peek 返回零值和false", func() {
			s := NewLinkStack[int]()
			v, ok := s.Peek()
			So(ok, ShouldBeFalse)
			So(v, ShouldEqual, 0)
		})

		Convey("Peek 返回栈顶元素但不移除", func() {
			s := NewLinkStack[int]()
			s.Push(1)
			s.Push(2)
			v, ok := s.Peek()
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, 2)
			So(s.Size(), ShouldEqual, 2)
		})

		Convey("多次 Peek 返回相同结果", func() {
			s := NewLinkStack[int]()
			s.Push(42)
			for i := 0; i < 5; i++ {
				v, ok := s.Peek()
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, 42)
			}
			So(s.Size(), ShouldEqual, 1)
		})
	})
}

func TestLinkStackEmpty(t *testing.T) {
	Convey("LinkStack Empty 测试", t, func() {
		Convey("新栈为空", func() {
			s := NewLinkStack[int]()
			So(s.Empty(), ShouldBeTrue)
		})

		Convey("Push后不为空", func() {
			s := NewLinkStack[int]()
			s.Push(1)
			So(s.Empty(), ShouldBeFalse)
		})
	})
}

func TestLinkStackSize(t *testing.T) {
	Convey("LinkStack Size 测试", t, func() {
		Convey("新栈大小为0", func() {
			s := NewLinkStack[int]()
			So(s.Size(), ShouldEqual, 0)
		})

		Convey("Push 后大小正确", func() {
			s := NewLinkStack[int]()
			s.Push(1)
			s.Push(2)
			So(s.Size(), ShouldEqual, 2)
		})
	})
}

func TestLinkStackStringType(t *testing.T) {
	Convey("LinkStack 字符串类型测试", t, func() {
		s := NewLinkStack[string]()
		s.Push("hello")
		s.Push("world")
		v, ok := s.Pop()
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "world")
		v, ok = s.Pop()
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "hello")
	})
}

func TestLinkStackInterface(t *testing.T) {
	Convey("LinkStack 实现 Stack 接口", t, func() {
		var s Stack[int] = NewLinkStack[int]()
		s.Push(1)
		v, ok := s.Pop()
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, 1)
	})
}
