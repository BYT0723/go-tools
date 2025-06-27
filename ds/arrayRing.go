package ds

import "iter"

var _ Ring[int] = (*ArrayRing[int])(nil)

// ArrayRing Non-concurrency-safe Ring
type ArrayRing[T any] struct {
	data []T
	size int
	next int
	full bool
}

func NewArrayRing[T any]() *ArrayRing[T] {
	return NewArrayRingWithSize[T](DefaultRingSize)
}

func NewArrayRingWithSize[T any](size int) *ArrayRing[T] {
	if size <= 0 {
		size = DefaultRingSize
	}
	return &ArrayRing[T]{
		data: make([]T, size),
		size: size,
	}
}

func (r *ArrayRing[T]) Push(value T) {
	r.data[r.next] = value
	r.next = (r.next + 1) % r.size
	if r.next == 0 {
		r.full = true
	}
}

func (r *ArrayRing[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		if r.full {
			for i := r.next; i < r.size; i++ {
				if !yield(r.data[i]) {
					return
				}
			}
			for i := 0; i < r.next; i++ {
				if !yield(r.data[i]) {
					return
				}
			}
		} else {
			for i := 0; i < r.next; i++ {
				if !yield(r.data[i]) {
					return
				}
			}
		}
	}
}

func (r *ArrayRing[T]) Values() []T {
	var result []T
	if r.full {
		result = make([]T, r.size)
		for i := 0; i < r.size; i++ {
			index := (r.next + i) % r.size
			result[i] = r.data[index]
		}
	} else {
		result = make([]T, r.next)
		copy(result, r.data)
	}
	return result
}

func (r *ArrayRing[T]) Len() int {
	if r.full {
		return r.size
	} else {
		return r.next
	}
}

func (r *ArrayRing[T]) Cap() int {
	return r.size
}
