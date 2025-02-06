package ds

var _ Stack[int] = (*LinkStack[int])(nil)

type (
	LinkStack[T any] struct {
		top  *element[T]
		size int
	}
	element[T any] struct {
		value T
		next  *element[T]
	}
)

func NewLinkStack[T any]() *LinkStack[T] {
	return &LinkStack[T]{}
}

func (m *LinkStack[T]) Pop() (value T, exist bool) {
	value, exist = m.Peek()
	m.top = m.top.next
	return
}

func (m *LinkStack[T]) Peek() (value T, exist bool) {
	if m.top == nil {
		return
	}
	value = m.top.value
	exist = true
	return
}

func (m *LinkStack[T]) Push(value T) {
	m.top = &element[T]{value: value, next: m.top}
	m.size++
}

func (m *LinkStack[T]) Empty() bool {
	return m.Size() == 0
}

func (m *LinkStack[T]) Size() int {
	return m.size
}
