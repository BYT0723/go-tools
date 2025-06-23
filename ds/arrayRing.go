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

func (r *ArrayRing[T]) Peek() (value T) {
	if r.full {
		return r.data[r.next]
	} else {
		return r.data[0]
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
		for i, v := range r.data {
			if i < r.next {
				result[i-r.next+r.size] = v
			} else {
				result[i-r.next] = v
			}
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
