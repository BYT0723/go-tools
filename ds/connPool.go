package ds

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

var ErrConnPoolClosed = fmt.Errorf("conn pool is closed")

var defaultConnPoolResidence = 10

type ConnPool[T io.Closer] struct {
	Residence int               // 持久化连接数量
	New       func() (T, error) // required, 新建连接function
	Check     func(T) bool      // 检查连接是否可用
	conns     chan T            // 连接池
	once      sync.Once
	closed    atomic.Bool
}

func (c *ConnPool[T]) Get() (v T, err error) {
	c.once.Do(c.init)

	if c.closed.Load() {
		return v, ErrConnPoolClosed
	}

	select {
	case conn := <-c.conns:
		if c.Check != nil && !c.Check(conn) {
			conn.Close()
			return c.New()
		}
		return conn, nil
	default:
		return c.New()
	}
}

func (c *ConnPool[T]) Put(conn T) error {
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

func (c *ConnPool[T]) Close() {
	if !c.closed.CompareAndSwap(false, true) {
		return
	}

	close(c.conns)
	for c := range c.conns {
		c.Close()
	}
}

func (c *ConnPool[T]) init() {
	if c.New == nil {
		panic("ConnPool.New is nil")
	}
	if c.Residence <= 0 {
		c.Residence = defaultConnPoolResidence
	}
	c.conns = make(chan T, c.Residence)
}
