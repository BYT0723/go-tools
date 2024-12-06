package ds

type Stack[T any] interface {
	Push(T)
	Pop() (T, bool)
	Peek() (T, bool)
	Empty() bool
	Size() int
}
