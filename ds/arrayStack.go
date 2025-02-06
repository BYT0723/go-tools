package ds

import (
	"fmt"
	"strings"
)

var _ Stack[int] = (*ArrayStack[int])(nil)

type ArrayStack[T any] struct {
	items []T
}

func NewArrayStack[T any]() *ArrayStack[T] {
	return &ArrayStack[T]{}
}

func (m *ArrayStack[T]) Pop() (value T, exist bool) {
	value, exist = m.Peek()
	m.items = m.items[:len(m.items)-1]
	return
}

func (m *ArrayStack[T]) Peek() (value T, exist bool) {
	if m.Empty() {
		return
	}
	value = m.items[len(m.items)-1]
	exist = true
	return
}

func (m *ArrayStack[T]) Push(value T) {
	m.items = append(m.items, value)
}

func (m *ArrayStack[T]) Empty() bool {
	return m.Size() == 0
}

func (m *ArrayStack[T]) Size() int {
	return len(m.items)
}

func (m *ArrayStack[T]) String() string {
	var (
		str    = "ArrayStack"
		n      = m.Size()
		values = make([]string, m.Size())
	)
	for i, v := range m.items {
		values[n-i-1] = fmt.Sprint(v)
	}
	return str + "[" + strings.Join(values, " ") + "]"
}
