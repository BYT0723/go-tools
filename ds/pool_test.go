package ds

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type testResource struct {
	id     int
	closed atomic.Bool
}

func (t *testResource) Close() error {
	t.closed.Store(true)
	return nil
}

func TestPoolGet(t *testing.T) {
	Convey("Pool Get 测试", t, func() {
		var createCount atomic.Int32

		p := &Pool[int, *testResource]{
			New: func(ctx context.Context, key int) (*testResource, error) {
				createCount.Add(1)
				return &testResource{id: key}, nil
			},
			Identifier: func(key int) string {
				return string(rune(key))
			},
			Destroy: func(ctx context.Context, value *testResource) error {
				return nil
			},
		}

		Convey("Get 不存在的key会创建新资源", func() {
			v, err := p.GetWithCtx(context.Background(), 1)
			So(err, ShouldBeNil)
			So(v, ShouldNotBeNil)
			So(v.id, ShouldEqual, 1)
			So(createCount.Load(), ShouldEqual, 1)
		})

		Convey("Get 已存在的key复用资源", func() {
			v1, _ := p.GetWithCtx(context.Background(), 2)
			v2, _ := p.GetWithCtx(context.Background(), 2)
			So(v1, ShouldEqual, v2)
			So(createCount.Load(), ShouldEqual, 1)
		})
	})
}

func TestPoolPut(t *testing.T) {
	Convey("Pool Put 测试", t, func() {
		var destroyed atomic.Bool

		p := &Pool[int, *testResource]{
			New: func(ctx context.Context, key int) (*testResource, error) {
				return &testResource{id: key}, nil
			},
			Identifier: func(key int) string {
				return string(rune(key))
			},
			Destroy: func(ctx context.Context, value *testResource) error {
				destroyed.Store(true)
				return nil
			},
		}

		Convey("Put 不存在key的borrow count为1时调用Destroy", func() {
			v, _ := p.GetWithCtx(context.Background(), 1)
			So(v, ShouldNotBeNil)

			err := p.PutWithCtx(context.Background(), 1)
			So(err, ShouldBeNil)
		})

		Convey("Put 多次Get后borrow count大于1不触发Destroy", func() {
			destroyed.Store(false)
			v1, _ := p.GetWithCtx(context.Background(), 2)
			v2, _ := p.GetWithCtx(context.Background(), 2)
			So(v1, ShouldEqual, v2)

			_ = p.PutWithCtx(context.Background(), 2)
		})
	})
}

func TestPoolGetDefaultCtx(t *testing.T) {
	Convey("Pool Get 使用默认context", t, func() {
		p := &Pool[int, *testResource]{
			New: func(ctx context.Context, key int) (*testResource, error) {
				return &testResource{id: key}, nil
			},
			Identifier: func(key int) string {
				return string(rune(key))
			},
			Destroy: func(ctx context.Context, value *testResource) error {
				return nil
			},
		}

		v, err := p.Get(context.Background(), 1)
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
	})
}

func TestPoolPutDefaultCtx(t *testing.T) {
	Convey("Pool Put 使用默认context", t, func() {
		p := &Pool[int, *testResource]{
			New: func(ctx context.Context, key int) (*testResource, error) {
				return &testResource{id: key}, nil
			},
			Identifier: func(key int) string {
				return string(rune(key))
			},
			Destroy: func(ctx context.Context, value *testResource) error {
				return nil
			},
		}

		v, _ := p.GetWithCtx(context.Background(), 1)
		So(v, ShouldNotBeNil)
		err := p.Put(1)
		So(err, ShouldBeNil)
	})
}

func TestPoolConcurrent(t *testing.T) {
	Convey("Pool 并发测试", t, func() {
		p := &Pool[int, *testResource]{
			New: func(ctx context.Context, key int) (*testResource, error) {
				return &testResource{id: key}, nil
			},
			Identifier: func(key int) string {
				return string(rune(key))
			},
			Destroy: func(ctx context.Context, value *testResource) error {
				return nil
			},
		}

		var wg sync.WaitGroup
		n := 50
		errCh := make(chan error, n)

		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(k int) {
				defer wg.Done()
				v, err := p.GetWithCtx(context.Background(), k)
				if err != nil {
					errCh <- err
					return
				}
				if v == nil {
					errCh <- nil
				}
				_ = p.Put(k)
			}(i)
		}

		wg.Wait()
		close(errCh)

		for e := range errCh {
			So(e, ShouldBeNil)
		}
	})
}
