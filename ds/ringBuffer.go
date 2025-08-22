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
			for i := r.next; i < r.size; i++ {
				if !yield(r.data[i]) {
					return
				}
			}
			for i := int64(0); i < r.next; i++ {
				if !yield(r.data[i]) {
					return
				}
			}
		} else {
			for i := int64(0); i < r.next; i++ {
				if !yield(r.data[i]) {
					return
				}
			}
		}
	}
}

func (r *RingBuffer[T]) Values() []T {
	r.m.Lock()
	defer r.m.Unlock()

	var result []T
	if r.full {
		result = make([]T, r.size)
		for i := int64(0); i < r.size; i++ {
			index := (r.next + i) % r.size
			result[i] = r.data[index]
		}
	} else {
		result = make([]T, r.next)
		copy(result, r.data)
	}
	return result
}

func (r *RingBuffer[T]) Len() int {
	if r.full {
		return int(r.size)
	} else {
		return int(r.next)
	}
}

func (r *RingBuffer[T]) Cap() int {
	return int(r.size)
}
