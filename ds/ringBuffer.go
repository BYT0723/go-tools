package ds

import (
	"iter"
	"sync"
	"sync/atomic"
)

// RingBuffer Non-concurrency-safe Ring
type RingBuffer[T any] struct {
	m    sync.Mutex
	data []T
	size int64
	next int64
	full bool
}

func NewRingBuffer[T any]() *RingBuffer[T] {
	return NewRingBufferWithSize[T](DefaultRingSize)
}

const DefaultRingSize = 1024

func NewRingBufferWithSize[T any](size int) *RingBuffer[T] {
	if size <= 0 {
		size = DefaultRingSize
	}
	return &RingBuffer[T]{
		data: make([]T, size),
		size: int64(size),
	}
}

func (r *RingBuffer[T]) Push(values ...T) {
	r.m.Lock()
	defer r.m.Unlock()
	for _, value := range values {
		r.data[r.next] = value
		atomic.StoreInt64(&r.next, (r.next+1)%r.size)

		if !r.full && r.next == 0 {
			r.full = true
		}
	}
}

func (r *RingBuffer[T]) Iterator() iter.Seq[T] {
	return func(yield func(T) bool) {
		r.m.Lock()
		defer r.m.Unlock()
		if r.full {
			for _, v := range r.data[r.next:] {
				if !yield(v) {
					return
				}
			}
			for _, v := range r.data[:r.next] {
				if !yield(v) {
					return
				}
			}
		} else {
			for _, v := range r.data[:r.next] {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func (r *RingBuffer[T]) Values() (result []T) {
	r.m.Lock()
	defer r.m.Unlock()

	if !r.full {
		result = make([]T, r.next)
		copy(result, r.data)
		return result
	}

	result = make([]T, r.size)
	copy(result, r.data[r.next:])
	copy(result[r.size-r.next:], r.data[:r.next])
	return result
}

func (r *RingBuffer[T]) Len() int {
	if !r.full {
		return int(r.next)
	}
	return int(r.size)
}

func (r *RingBuffer[T]) Cap() int {
	return int(r.size)
}
