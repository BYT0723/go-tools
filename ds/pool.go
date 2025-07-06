package ds

import (
	"context"
	"sync/atomic"
)

type (
	Pool[K, V any] struct {
		New        poolNewFunc[K, V]
		Identifier poolIDFunc[K]
		Destroy    poolDestroyFunc[V]
		entries    SyncMap[string, *poolItem[V]]
	}
	poolItem[V any] struct {
		value  V
		borrow atomic.Int32
	}
	poolNewFunc[K, V any]  func(ctx context.Context, key K) (V, error)
	poolIDFunc[K any]      func(key K) string
	poolDestroyFunc[V any] func(ctx context.Context, value V) error
)

func (p *Pool[K, V]) Get(ctx context.Context, key K) (value V, err error) {
	return p.GetWithCtx(context.Background(), key)
}

func (p *Pool[K, V]) GetWithCtx(ctx context.Context, key K) (value V, err error) {
	k := p.Identifier(key)

	item, ok := p.entries.Load(k)
	if ok {
		item.borrow.Add(1)
		value = item.value
		return
	}

	v, err := p.New(ctx, key)
	if err != nil {
		return
	}
	item = &poolItem[V]{value: v}
	item.borrow.Add(1)
	p.entries.Store(k, item)

	value = v
	return
}

func (p *Pool[K, V]) Put(key K) (err error) {
	return p.PutWithCtx(context.Background(), key)
}

func (p *Pool[K, V]) PutWithCtx(ctx context.Context, key K) (err error) {
	k := p.Identifier(key)
	item, ok := p.entries.Load(k)
	if !ok {
		return
	}
	if item.borrow.Add(-1) == 0 {
		p.entries.Delete(k)
	}
	return p.Destroy(ctx, item.value)
}
