package sync

import "sync"

type Pool[T any] struct {
	p sync.Pool
}

func NewPool[T any](newf func() T) *Pool[T] {
	var pool Pool[T]
	pool.p.New = func() any {
		return newf()
	}
	return &pool
}

func (p *Pool[T]) Get() T {
	return p.p.Get().(T)
}

func (p *Pool[T]) Put(x T) {
	p.p.Put(x)
}
