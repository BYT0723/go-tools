package arraystack

import (
	"fmt"
	"strings"

	"github.com/BYT0723/go-tools/datastruct/stack"
)

var _ stack.Stack[int] = (*Stack[int])(nil)

type Stack[T any] struct {
	items []T
}

func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (m *Stack[T]) Pop() (value T, exist bool) {
	value, exist = m.Peek()
	m.items = m.items[:len(m.items)-1]
	return
}

func (m *Stack[T]) Peek() (value T, exist bool) {
	if m.Empty() {
		return
	}
	value = m.items[len(m.items)-1]
	exist = true
	return
}

func (m *Stack[T]) Push(value T) {
	m.items = append(m.items, value)
}

func (m *Stack[T]) Empty() bool {
	return m.Size() == 0
}

func (m *Stack[T]) Size() int {
	return len(m.items)
}

func (m *Stack[T]) String() string {
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
