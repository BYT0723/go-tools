package ds

import (
	"context"
	"sync/atomic"
)

// Pool is a generic, thread-safe object pool that manages reusable resources.
// It maintains a collection of objects that can be borrowed and returned,
// reducing the overhead of creating new objects repeatedly.
//
// Type parameters:
//   - K: The key type used to identify objects in the pool
//   - V: The value type of objects stored in the pool
//
// The pool tracks borrow counts for each object and automatically removes
// objects when they are no longer in use (borrow count reaches zero).
type (
	Pool[K, V any] struct {
		// New is a function that creates a new object when the pool doesn't have
		// an available instance for the given key.
		New poolNewFunc[K, V]

		// Identifier converts a key to a string identifier used for internal storage.
		// This allows the pool to use SyncMap for thread-safe operations.
		Identifier poolIDFunc[K]

		// Destroy is called when an object is removed from the pool (borrow count reaches zero).
		// Use this to clean up resources (e.g., close connections, free memory).
		Destroy poolDestroyFunc[V]

		// entries stores the pooled objects with their borrow counts.
		entries SyncMap[string, *poolItem[V]]
	}

	// poolItem wraps a pooled value with its borrow count.
	poolItem[V any] struct {
		value  V          // The actual pooled value
		borrow atomic.Int32 // Current number of borrowers
	}

	// poolNewFunc creates a new value for the given key.
	// Context can be used for cancellation or timeout during creation.
	poolNewFunc[K, V any] func(ctx context.Context, key K) (V, error)

	// poolIDFunc converts a key to a string identifier.
	// This must produce consistent identifiers for equal keys.
	poolIDFunc[K any] func(key K) string

	// poolDestroyFunc cleans up a value when it's removed from the pool.
	// Context can be used for cancellation or timeout during cleanup.
	poolDestroyFunc[V any] func(ctx context.Context, value V) error
)

// Get retrieves a value from the pool for the given key.
// If an object exists in the pool, its borrow count is incremented and the object is returned.
// If no object exists, a new one is created using the New function.
// This method uses context.Background() internally.
//
// Parameters:
//   - ctx: Context for cancellation/timeout (not used in this method, see GetWithCtx)
//   - key: The key identifying the object to retrieve
//
// Returns:
//   - value: The retrieved or newly created object
//   - err: Error if object creation fails
func (p *Pool[K, V]) Get(ctx context.Context, key K) (value V, err error) {
	return p.GetWithCtx(context.Background(), key)
}

// GetWithCtx retrieves a value from the pool for the given key with context support.
// If an object exists in the pool, its borrow count is incremented and the object is returned.
// If no object exists, a new one is created using the New function with the provided context.
//
// Parameters:
//   - ctx: Context for cancellation/timeout during object creation
//   - key: The key identifying the object to retrieve
//
// Returns:
//   - value: The retrieved or newly created object
//   - err: Error if object creation fails or context is cancelled
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

// Put returns a borrowed object to the pool, decrementing its borrow count.
// When the borrow count reaches zero, the object is removed from the pool
// and the Destroy function is called to clean up resources.
// This method uses context.Background() internally.
//
// Parameters:
//   - key: The key identifying the object to return
//
// Returns:
//   - err: Error if the Destroy function fails
func (p *Pool[K, V]) Put(key K) (err error) {
	return p.PutWithCtx(context.Background(), key)
}

// PutWithCtx returns a borrowed object to the pool with context support.
// When the borrow count reaches zero, the object is removed from the pool
// and the Destroy function is called with the provided context.
//
// Parameters:
//   - ctx: Context for cancellation/timeout during cleanup
//   - key: The key identifying the object to return
//
// Returns:
//   - err: Error if the Destroy function fails or context is cancelled
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
