package ds

import (
	"fmt"
	"strings"
)

var _ Stack[int] = (*ArrayStack[int])(nil)

// ArrayStack is a stack implementation using a slice as the underlying storage.
// It provides O(1) average time complexity for push and pop operations.
//
// Type parameters:
//   - T: The element type stored in the stack
type ArrayStack[T any] struct {
	items []T
}

// NewArrayStack creates a new empty ArrayStack.
func NewArrayStack[T any]() *ArrayStack[T] {
	return &ArrayStack[T]{}
}

// Pop removes and returns the top element from the stack.
// Returns the element and a boolean indicating if the stack was not empty.
func (m *ArrayStack[T]) Pop() (value T, exist bool) {
	value, exist = m.Peek()
	m.items = m.items[:len(m.items)-1]
	return
}

// Peek returns the top element without removing it from the stack.
// Returns the element and a boolean indicating if the stack was not empty.
func (m *ArrayStack[T]) Peek() (value T, exist bool) {
	if m.Empty() {
		return
	}
	value = m.items[len(m.items)-1]
	exist = true
	return
}

// Push adds an element to the top of the stack.
func (m *ArrayStack[T]) Push(value T) {
	m.items = append(m.items, value)
}

// Empty returns true if the stack contains no elements.
func (m *ArrayStack[T]) Empty() bool {
	return m.Size() == 0
}

// Size returns the number of elements in the stack.
func (m *ArrayStack[T]) Size() int {
	return len(m.items)
}

// String returns a string representation of the stack with elements from top to bottom.
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
