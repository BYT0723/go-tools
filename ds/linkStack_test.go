package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkStackPush(t *testing.T) {
	t.Run("LinkStack Push 测试", func(t *testing.T) {
		t.Run("Push 空栈", func(t *testing.T) {
			s := NewLinkStack[int]()
			s.Push(1)
			assert.Equal(t, 1, s.Size())
			assert.False(t, s.Empty())
		})

		t.Run("Push 多个元素", func(t *testing.T) {
			s := NewLinkStack[int]()
			for i := 0; i < 10; i++ {
				s.Push(i)
			}
			assert.Equal(t, 10, s.Size())
		})
	})
}

func TestLinkStackPop(t *testing.T) {
	t.Run("LinkStack Pop 测试", func(t *testing.T) {
		t.Run("空栈 Pop panic", func(t *testing.T) {
			s := NewLinkStack[int]()
			assert.Panics(t, func() { s.Pop() })
		})

		t.Run("Pop 返回最后Push的元素", func(t *testing.T) {
			s := NewLinkStack[int]()
			s.Push(1)
			s.Push(2)
			s.Push(3)
			v, ok := s.Pop()
			assert.True(t, ok)
			assert.Equal(t, 3, v)
		})
	})
}

func TestLinkStackPeek(t *testing.T) {
	t.Run("LinkStack Peek 测试", func(t *testing.T) {
		t.Run("空栈 Peek 返回零值和false", func(t *testing.T) {
			s := NewLinkStack[int]()
			v, ok := s.Peek()
			assert.False(t, ok)
			assert.Equal(t, 0, v)
		})

		t.Run("Peek 返回栈顶元素但不移除", func(t *testing.T) {
			s := NewLinkStack[int]()
			s.Push(1)
			s.Push(2)
			v, ok := s.Peek()
			assert.True(t, ok)
			assert.Equal(t, 2, v)
			assert.Equal(t, 2, s.Size())
		})

		t.Run("多次 Peek 返回相同结果", func(t *testing.T) {
			s := NewLinkStack[int]()
			s.Push(42)
			for i := 0; i < 5; i++ {
				v, ok := s.Peek()
				assert.True(t, ok)
				assert.Equal(t, 42, v)
			}
			assert.Equal(t, 1, s.Size())
		})
	})
}

func TestLinkStackEmpty(t *testing.T) {
	t.Run("LinkStack Empty 测试", func(t *testing.T) {
		t.Run("新栈为空", func(t *testing.T) {
			s := NewLinkStack[int]()
			assert.True(t, s.Empty())
		})

		t.Run("Push后不为空", func(t *testing.T) {
			s := NewLinkStack[int]()
			s.Push(1)
			assert.False(t, s.Empty())
		})
	})
}

func TestLinkStackSize(t *testing.T) {
	t.Run("LinkStack Size 测试", func(t *testing.T) {
		t.Run("新栈大小为0", func(t *testing.T) {
			s := NewLinkStack[int]()
			assert.Equal(t, 0, s.Size())
		})

		t.Run("Push 后大小正确", func(t *testing.T) {
			s := NewLinkStack[int]()
			s.Push(1)
			s.Push(2)
			assert.Equal(t, 2, s.Size())
		})
	})
}

func TestLinkStackStringType(t *testing.T) {
	t.Run("LinkStack 字符串类型测试", func(t *testing.T) {
		s := NewLinkStack[string]()
		s.Push("hello")
		s.Push("world")
		v, ok := s.Pop()
		assert.True(t, ok)
		assert.Equal(t, "world", v)
		v, ok = s.Pop()
		assert.True(t, ok)
		assert.Equal(t, "hello", v)
	})
}

func TestLinkStackInterface(t *testing.T) {
	t.Run("LinkStack 实现 Stack 接口", func(t *testing.T) {
		var s Stack[int] = NewLinkStack[int]()
		s.Push(1)
		v, ok := s.Pop()
		assert.True(t, ok)
		assert.Equal(t, 1, v)
	})
}
