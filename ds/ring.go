package ds

import "iter"

const DefaultRingSize = 1024

type Ring[T any] interface {
	Push(elem T)
	Peek() T
	Iterator() iter.Seq[T]
	Values() []T
	Len() int
}
