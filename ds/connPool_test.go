package ds

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockConn struct {
	id     int
	closed atomic.Bool
}

func (m *mockConn) Close() error {
	m.closed.Store(true)
	return nil
}

func (m *mockConn) isClosed() bool {
	return m.closed.Load()
}

func TestConnPoolNew(t *testing.T) {
	t.Run("ConnPool 创建测试", func(t *testing.T) {
		t.Run("正常创建", func(t *testing.T) {
			var createCount atomic.Int32
			p := NewConnPool(func() (*mockConn, error) {
				createCount.Add(1)
				return &mockConn{id: int(createCount.Load())}, nil
			}, 5)
			assert.NotNil(t, p)
		})

		t.Run("factory为nil panic", func(t *testing.T) {
			assert.Panics(t, func() {
				NewConnPool[*mockConn](nil, 5)
			})
		})

		t.Run("residence <= 0 使用默认值", func(t *testing.T) {
			p := NewConnPool(func() (*mockConn, error) {
				return &mockConn{}, nil
			}, 0)
			assert.NotNil(t, p)
		})
	})
}

func TestConnPoolGet(t *testing.T) {
	t.Run("ConnPool Get 测试", func(t *testing.T) {
		t.Run("获取连接-首次创建新连接", func(t *testing.T) {
			var createCount atomic.Int32
			p := NewConnPool(func() (*mockConn, error) {
				createCount.Add(1)
				return &mockConn{id: int(createCount.Load())}, nil
			}, 5)

			conn, err := p.Get()
			assert.Nil(t, err)
			assert.NotNil(t, conn)
			assert.Equal(t, int32(1), createCount.Load())
		})

		t.Run("归还后获取会复用连接", func(t *testing.T) {
			var createCount atomic.Int32
			p := NewConnPool(func() (*mockConn, error) {
				createCount.Add(1)
				return &mockConn{id: int(createCount.Load())}, nil
			}, 5)

			conn1, _ := p.Get()
			p.Put(conn1)

			conn2, _ := p.Get()
			assert.Equal(t, conn1.id, conn2.id)
		})

		t.Run("连接池满时Put会关闭连接", func(t *testing.T) {
			var createCount atomic.Int32
			p := NewConnPool(func() (*mockConn, error) {
				createCount.Add(1)
				return &mockConn{id: int(createCount.Load())}, nil
			}, 1)

			conn1, _ := p.Get()
			p.Put(conn1)

			conn2, _ := p.Get()
			p.Put(conn2)

			conn3, _ := p.Get()
			p.Put(conn3)
		})

		t.Run("关闭的池Get返回错误", func(t *testing.T) {
			p := NewConnPool(func() (*mockConn, error) {
				return &mockConn{}, nil
			}, 5)
			p.Close()

			_, err := p.Get()
			assert.Equal(t, ErrConnPoolClosed, err)
		})
	})
}

func TestConnPoolPut(t *testing.T) {
	t.Run("ConnPool Put 测试", func(t *testing.T) {
		t.Run("Put到已关闭的池返回错误", func(t *testing.T) {
			p := NewConnPool(func() (*mockConn, error) {
				return &mockConn{}, nil
			}, 5)

			conn, _ := p.Get()
			p.Close()

			err := p.Put(conn)
			assert.Equal(t, ErrConnPoolClosed, err)
		})
	})
}

func TestConnPoolClose(t *testing.T) {
	t.Run("ConnPool Close 测试", func(t *testing.T) {
		t.Run("Close 关闭池中所有连接", func(t *testing.T) {
			p := NewConnPool(func() (*mockConn, error) {
				return &mockConn{}, nil
			}, 5)

			conn, _ := p.Get()
			p.Put(conn)
			p.Close()

			assert.True(t, conn.isClosed())
		})

		t.Run("Close 是幂等的", func(t *testing.T) {
			p := NewConnPool(func() (*mockConn, error) {
				return &mockConn{}, nil
			}, 5)
			p.Close()
			assert.NotPanics(t, func() { p.Close() })
		})
	})
}

func TestConnPoolWithCheck(t *testing.T) {
	t.Run("ConnPool WithCheck 测试", func(t *testing.T) {
		t.Run("check返回false时创建新连接", func(t *testing.T) {
			var createCount atomic.Int32
			p := NewConnPool(func() (*mockConn, error) {
				createCount.Add(1)
				return &mockConn{id: int(createCount.Load())}, nil
			}, 5, WithCheck(func(c *mockConn) bool {
				return false
			}))

			conn1, _ := p.Get()
			p.Put(conn1)
			conn2, _ := p.Get()

			assert.Equal(t, int32(2), createCount.Load())
			assert.NotEqual(t, conn1.id, conn2.id)
		})
	})
}
