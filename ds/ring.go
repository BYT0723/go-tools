package ds

const DefaultRingSize = 1024

type Ring[T any] interface {
	Push(elem T)
	Peek() T
	Values() []T
	Len() int
}
