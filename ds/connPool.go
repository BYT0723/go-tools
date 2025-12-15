package ds

import (
	"fmt"
	"io"
	"sync/atomic"
)

var ErrConnPoolClosed = fmt.Errorf("conn pool is closed")

const defaultConnPoolResidence = 10

type (
	connPool[T io.Closer] struct {
		factory func() (T, error) // required, 新建连接function
		check   func(T) bool      // 检查连接是否可用
		conns   chan T            // 连接池
		closed  atomic.Bool
	}

	ConnPoolOption[T io.Closer] func(*connPool[T])
)

func WithCheck[T io.Closer](check func(T) bool) ConnPoolOption[T] {
	return func(p *connPool[T]) {
		p.check = check
	}
}

// NewConnPool returns a new connection pool.
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

// Close closes the connection pool.
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
