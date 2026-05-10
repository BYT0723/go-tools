package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayStackPush(t *testing.T) {
	t.Run("ArrayStack Push 测试", func(t *testing.T) {
		t.Run("Push 空栈", func(t *testing.T) {
			s := NewArrayStack[int]()
			s.Push(1)
			assert.Equal(t, 1, s.Size())
			assert.False(t, s.Empty())
		})

		t.Run("Push 多个元素", func(t *testing.T) {
			s := NewArrayStack[int]()
			for i := 0; i < 10; i++ {
				s.Push(i)
			}
			assert.Equal(t, 10, s.Size())
		})
	})
}

func TestArrayStackPop(t *testing.T) {
	t.Run("ArrayStack Pop 测试", func(t *testing.T) {
		t.Run("空栈 Pop panic", func(t *testing.T) {
			s := NewArrayStack[int]()
			assert.Panics(t, func() { s.Pop() })
		})

		t.Run("Pop 返回最后Push的元素", func(t *testing.T) {
			s := NewArrayStack[int]()
			s.Push(1)
			s.Push(2)
			s.Push(3)
			v, ok := s.Pop()
			assert.True(t, ok)
			assert.Equal(t, 3, v)
			assert.Equal(t, 2, s.Size())
		})

		t.Run("Pop 直到栈空", func(t *testing.T) {
			s := NewArrayStack[int]()
			s.Push(1)
			s.Push(2)
			s.Pop()
			s.Pop()
			assert.True(t, s.Empty())
			assert.Equal(t, 0, s.Size())
		})
	})
}

func TestArrayStackPeek(t *testing.T) {
	t.Run("ArrayStack Peek 测试", func(t *testing.T) {
		t.Run("空栈 Peek 返回零值和false", func(t *testing.T) {
			s := NewArrayStack[int]()
			v, ok := s.Peek()
			assert.False(t, ok)
			assert.Equal(t, 0, v)
		})

		t.Run("Peek 返回栈顶元素但不移除", func(t *testing.T) {
			s := NewArrayStack[int]()
			s.Push(1)
			s.Push(2)
			v, ok := s.Peek()
			assert.True(t, ok)
			assert.Equal(t, 2, v)
			assert.Equal(t, 2, s.Size())
		})
	})
}

func TestArrayStackEmpty(t *testing.T) {
	t.Run("ArrayStack Empty 测试", func(t *testing.T) {
		t.Run("新栈为空", func(t *testing.T) {
			s := NewArrayStack[int]()
			assert.True(t, s.Empty())
		})

		t.Run("Push后不为空", func(t *testing.T) {
			s := NewArrayStack[int]()
			s.Push(1)
			assert.False(t, s.Empty())
		})
	})
}

func TestArrayStackSize(t *testing.T) {
	t.Run("ArrayStack Size 测试", func(t *testing.T) {
		t.Run("新栈大小为0", func(t *testing.T) {
			s := NewArrayStack[int]()
			assert.Equal(t, 0, s.Size())
		})

		t.Run("Push/Pop 后大小正确", func(t *testing.T) {
			s := NewArrayStack[int]()
			s.Push(1)
			s.Push(2)
			assert.Equal(t, 2, s.Size())
			s.Pop()
			assert.Equal(t, 1, s.Size())
		})
	})
}

func TestArrayStackString(t *testing.T) {
	t.Run("ArrayStack String 测试", func(t *testing.T) {
		t.Run("空栈输出", func(t *testing.T) {
			s := NewArrayStack[int]()
			assert.Equal(t, "ArrayStack[]", s.String())
		})

		t.Run("单元素输出", func(t *testing.T) {
			s := NewArrayStack[int]()
			s.Push(1)
			assert.Equal(t, "ArrayStack[1]", s.String())
		})

		t.Run("多元素输出（从顶到底）", func(t *testing.T) {
			s := NewArrayStack[int]()
			s.Push(1)
			s.Push(2)
			s.Push(3)
			assert.Equal(t, "ArrayStack[3 2 1]", s.String())
		})
	})
}

func TestArrayStackInterface(t *testing.T) {
	t.Run("ArrayStack 实现 Stack 接口", func(t *testing.T) {
		var s Stack[int] = NewArrayStack[int]()
		s.Push(1)
		v, ok := s.Pop()
		assert.True(t, ok)
		assert.Equal(t, 1, v)
	})
}
