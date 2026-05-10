package ds

import (
	"sync/atomic"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
	Convey("ConnPool 创建测试", t, func() {
		Convey("正常创建", func() {
			var createCount atomic.Int32
			p := NewConnPool(func() (*mockConn, error) {
				createCount.Add(1)
				return &mockConn{id: int(createCount.Load())}, nil
			}, 5)
			So(p, ShouldNotBeNil)
		})

		Convey("factory为nil panic", func() {
			So(func() {
				NewConnPool[*mockConn](nil, 5)
			}, ShouldPanic)
		})

		Convey("residence <= 0 使用默认值", func() {
			p := NewConnPool(func() (*mockConn, error) {
				return &mockConn{}, nil
			}, 0)
			So(p, ShouldNotBeNil)
		})
	})
}

func TestConnPoolGet(t *testing.T) {
	Convey("ConnPool Get 测试", t, func() {
		Convey("获取连接-首次创建新连接", func() {
			var createCount atomic.Int32
			p := NewConnPool(func() (*mockConn, error) {
				createCount.Add(1)
				return &mockConn{id: int(createCount.Load())}, nil
			}, 5)

			conn, err := p.Get()
			So(err, ShouldBeNil)
			So(conn, ShouldNotBeNil)
			So(createCount.Load(), ShouldEqual, 1)
		})

		Convey("归还后获取会复用连接", func() {
			var createCount atomic.Int32
			p := NewConnPool(func() (*mockConn, error) {
				createCount.Add(1)
				return &mockConn{id: int(createCount.Load())}, nil
			}, 5)

			conn1, _ := p.Get()
			p.Put(conn1)

			conn2, _ := p.Get()
			So(conn2.id, ShouldEqual, conn1.id)
		})

		Convey("连接池满时Put会关闭连接", func() {
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

		Convey("关闭的池Get返回错误", func() {
			p := NewConnPool(func() (*mockConn, error) {
				return &mockConn{}, nil
			}, 5)
			p.Close()

			_, err := p.Get()
			So(err, ShouldEqual, ErrConnPoolClosed)
		})
	})
}

func TestConnPoolPut(t *testing.T) {
	Convey("ConnPool Put 测试", t, func() {
		Convey("Put到已关闭的池返回错误", func() {
			p := NewConnPool(func() (*mockConn, error) {
				return &mockConn{}, nil
			}, 5)

			conn, _ := p.Get()
			p.Close()

			err := p.Put(conn)
			So(err, ShouldEqual, ErrConnPoolClosed)
		})
	})
}

func TestConnPoolClose(t *testing.T) {
	Convey("ConnPool Close 测试", t, func() {
		Convey("Close 关闭池中所有连接", func() {
			p := NewConnPool(func() (*mockConn, error) {
				return &mockConn{}, nil
			}, 5)

			conn, _ := p.Get()
			p.Put(conn)
			p.Close()

			So(conn.isClosed(), ShouldBeTrue)
		})

		Convey("Close 是幂等的", func() {
			p := NewConnPool(func() (*mockConn, error) {
				return &mockConn{}, nil
			}, 5)
			p.Close()
			So(func() { p.Close() }, ShouldNotPanic)
		})
	})
}

func TestConnPoolWithCheck(t *testing.T) {
	Convey("ConnPool WithCheck 测试", t, func() {
		Convey("check返回false时创建新连接", func() {
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

			So(createCount.Load(), ShouldEqual, 2)
			So(conn2.id, ShouldNotEqual, conn1.id)
		})
	})
}
