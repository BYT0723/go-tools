package ds

import (
	"fmt"
	"io"
	"sync/atomic"
)

var ErrConnPoolClosed = fmt.Errorf("conn pool is closed")

const defaultConnPoolResidence = 10

type (
	// connPool is a generic connection pool for managing reusable connections.
	// It maintains a fixed-size pool of connections that can be borrowed and returned.
	// Connections must implement the io.Closer interface.
	//
	// Type parameters:
	//   - T: The connection type, must implement io.Closer
	connPool[T io.Closer] struct {
		factory func() (T, error) // required, creates new connections
		check   func(T) bool      // optional, checks if a connection is still usable
		conns   chan T            // channel storing available connections
		closed  atomic.Bool       // atomic flag indicating if pool is closed
	}

	// ConnPoolOption defines optional configuration for the connection pool.
	ConnPoolOption[T io.Closer] func(*connPool[T])
)

// WithCheck returns a ConnPoolOption that sets a connection validation function.
// The check function is called when retrieving a connection from the pool
// to verify it's still usable. If check returns false, the connection is closed
// and a new one is created.
func WithCheck[T io.Closer](check func(T) bool) ConnPoolOption[T] {
	return func(p *connPool[T]) {
		p.check = check
	}
}

// NewConnPool returns a new connection pool.
//
// Parameters:
//   - factory: Required function that creates new connections
//   - residence: Maximum number of connections to keep in the pool (capacity)
//   - opts: Optional configuration options
//
// Returns:
//   - *connPool[T]: A new connection pool instance
//
// Panics:
//   - If factory is nil
func NewConnPool[T io.Closer](
	factory func() (T, error),
	residence int,
	opts ...ConnPoolOption[T],
) *connPool[T] {
	if factory == nil {
		panic("conn pool: nil factory")
	}
	if residence <= 0 {
		residence = defaultConnPoolResidence
	}

	p := &connPool[T]{
		factory: factory,
		conns:   make(chan T, residence),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Get returns a connection from the pool.
// If a connection is available in the pool, it's returned immediately.
// If the pool is empty, a new connection is created using the factory.
// If the pool is closed, returns ErrConnPoolClosed.
//
// Returns:
//   - v: The connection
//   - err: Error if pool is closed or connection creation fails
func (c *connPool[T]) Get() (v T, err error) {
	if c.closed.Load() {
		return v, ErrConnPoolClosed
	}

	select {
	case conn, ok := <-c.conns:
		if !ok {
			return v, ErrConnPoolClosed
		}
		if c.check != nil && !c.check(conn) {
			conn.Close()
			if c.closed.Load() {
				return v, ErrConnPoolClosed
			}
			return c.factory()
		}
		return conn, nil
	default:
		return c.factory()
	}
}

// Put returns a connection to the pool.
// If the pool is full, the connection is closed.
// If the pool is closed, returns ErrConnPoolClosed.
//
// Parameters:
//   - conn: The connection to return to the pool
//
// Returns:
//   - error: ErrConnPoolClosed if pool is closed, nil otherwise
func (c *connPool[T]) Put(conn T) error {
	if c.closed.Load() {
		return ErrConnPoolClosed
	}
	select {
	case c.conns <- conn:
		return nil
	default:
		return conn.Close()
	}
}

// Close closes the connection pool and all connections in it.
// After closing, Get and Put will return ErrConnPoolClosed.
// This method is idempotent and can be called multiple times.
func (c *connPool[T]) Close() {
	if !c.closed.CompareAndSwap(false, true) {
		return
	}

	for {
		select {
		case c := <-c.conns:
			c.Close()
		default:
			return
		}
	}
}
