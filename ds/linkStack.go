package ds

var _ Stack[int] = (*LinkStack[int])(nil)

type (
	// LinkStack is a stack implementation using a linked list as the underlying storage.
	// It provides O(1) time complexity for push and pop operations.
	//
	// Type parameters:
	//   - T: The element type stored in the stack
	LinkStack[T any] struct {
		top  *element[T]
		size int
	}

	// element represents a node in the linked list stack.
	element[T any] struct {
		value T
		next  *element[T]
	}
)

// NewLinkStack creates a new empty LinkStack.
func NewLinkStack[T any]() *LinkStack[T] {
	return &LinkStack[T]{}
}

// Pop removes and returns the top element from the stack.
// Returns the element and a boolean indicating if the stack was not empty.
func (m *LinkStack[T]) Pop() (value T, exist bool) {
	value, exist = m.Peek()
	m.top = m.top.next
	return
}

// Peek returns the top element without removing it from the stack.
// Returns the element and a boolean indicating if the stack was not empty.
func (m *LinkStack[T]) Peek() (value T, exist bool) {
	if m.top == nil {
		return
	}
	value = m.top.value
	exist = true
	return
}

// Push adds an element to the top of the stack.
func (m *LinkStack[T]) Push(value T) {
	m.top = &element[T]{value: value, next: m.top}
	m.size++
}

// Empty returns true if the stack contains no elements.
func (m *LinkStack[T]) Empty() bool {
	return m.Size() == 0
}

// Size returns the number of elements in the stack.
func (m *LinkStack[T]) Size() int {
	return m.size
}
