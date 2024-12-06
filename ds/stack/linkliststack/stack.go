package linkliststack

import (
	"github.com/BYT0723/go-tools/ds"
)

var _ ds.Stack[int] = (*Stack[int])(nil)

type Stack[T any] struct {
	top  *element[T]
	size int
}

type element[T any] struct {
	value T
	next  *element[T]
}

func (m *Stack[T]) Pop() (value T, exist bool) {
	value, exist = m.Peek()
	m.top = m.top.next
	return
}

func (m *Stack[T]) Peek() (value T, exist bool) {
	if m.top == nil {
		return
	}
	value = m.top.value
	exist = true
	return
}

func (m *Stack[T]) Push(value T) {
	m.top = &element[T]{value: value, next: m.top}
	m.size++
}

func (m *Stack[T]) Empty() bool {
	return m.Size() == 0
}

func (m *Stack[T]) Size() int {
	return m.size
}
