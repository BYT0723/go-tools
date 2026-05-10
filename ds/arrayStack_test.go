package ds

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArrayStackPush(t *testing.T) {
	Convey("ArrayStack Push 测试", t, func() {
		Convey("Push 空栈", func() {
			s := NewArrayStack[int]()
			s.Push(1)
			So(s.Size(), ShouldEqual, 1)
			So(s.Empty(), ShouldBeFalse)
		})

		Convey("Push 多个元素", func() {
			s := NewArrayStack[int]()
			for i := 0; i < 10; i++ {
				s.Push(i)
			}
			So(s.Size(), ShouldEqual, 10)
		})
	})
}

func TestArrayStackPop(t *testing.T) {
	Convey("ArrayStack Pop 测试", t, func() {
		Convey("空栈 Pop panic", func() {
			s := NewArrayStack[int]()
			So(func() { s.Pop() }, ShouldPanic)
		})

		Convey("Pop 返回最后Push的元素", func() {
			s := NewArrayStack[int]()
			s.Push(1)
			s.Push(2)
			s.Push(3)
			v, ok := s.Pop()
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, 3)
			So(s.Size(), ShouldEqual, 2)
		})

		Convey("Pop 直到栈空", func() {
			s := NewArrayStack[int]()
			s.Push(1)
			s.Push(2)
			s.Pop()
			s.Pop()
			So(s.Empty(), ShouldBeTrue)
			So(s.Size(), ShouldEqual, 0)
		})
	})
}

func TestArrayStackPeek(t *testing.T) {
	Convey("ArrayStack Peek 测试", t, func() {
		Convey("空栈 Peek 返回零值和false", func() {
			s := NewArrayStack[int]()
			v, ok := s.Peek()
			So(ok, ShouldBeFalse)
			So(v, ShouldEqual, 0)
		})

		Convey("Peek 返回栈顶元素但不移除", func() {
			s := NewArrayStack[int]()
			s.Push(1)
			s.Push(2)
			v, ok := s.Peek()
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, 2)
			So(s.Size(), ShouldEqual, 2)
		})
	})
}

func TestArrayStackEmpty(t *testing.T) {
	Convey("ArrayStack Empty 测试", t, func() {
		Convey("新栈为空", func() {
			s := NewArrayStack[int]()
			So(s.Empty(), ShouldBeTrue)
		})

		Convey("Push后不为空", func() {
			s := NewArrayStack[int]()
			s.Push(1)
			So(s.Empty(), ShouldBeFalse)
		})
	})
}

func TestArrayStackSize(t *testing.T) {
	Convey("ArrayStack Size 测试", t, func() {
		Convey("新栈大小为0", func() {
			s := NewArrayStack[int]()
			So(s.Size(), ShouldEqual, 0)
		})

		Convey("Push/Pop 后大小正确", func() {
			s := NewArrayStack[int]()
			s.Push(1)
			s.Push(2)
			So(s.Size(), ShouldEqual, 2)
			s.Pop()
			So(s.Size(), ShouldEqual, 1)
		})
	})
}

func TestArrayStackString(t *testing.T) {
	Convey("ArrayStack String 测试", t, func() {
		Convey("空栈输出", func() {
			s := NewArrayStack[int]()
			So(s.String(), ShouldEqual, "ArrayStack[]")
		})

		Convey("单元素输出", func() {
			s := NewArrayStack[int]()
			s.Push(1)
			So(s.String(), ShouldEqual, "ArrayStack[1]")
		})

		Convey("多元素输出（从顶到底）", func() {
			s := NewArrayStack[int]()
			s.Push(1)
			s.Push(2)
			s.Push(3)
			So(s.String(), ShouldEqual, "ArrayStack[3 2 1]")
		})
	})
}

func TestArrayStackInterface(t *testing.T) {
	Convey("ArrayStack 实现 Stack 接口", t, func() {
		var s Stack[int] = NewArrayStack[int]()
		s.Push(1)
		v, ok := s.Pop()
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, 1)
	})
}
