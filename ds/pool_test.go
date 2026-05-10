package ds

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
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
	t.Run("Pool Get 测试", func(t *testing.T) {
		t.Run("Get 不存在的key会创建新资源", func(t *testing.T) {
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

			v, err := p.GetWithCtx(context.Background(), 1)
			assert.Nil(t, err)
			assert.NotNil(t, v)
			assert.Equal(t, 1, v.id)
			assert.Equal(t, int32(1), createCount.Load())
		})

		t.Run("Get 已存在的key复用资源", func(t *testing.T) {
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

			v1, _ := p.GetWithCtx(context.Background(), 2)
			v2, _ := p.GetWithCtx(context.Background(), 2)
			assert.Equal(t, v1, v2)
			assert.Equal(t, int32(1), createCount.Load())
		})
	})
}

func TestPoolPut(t *testing.T) {
	t.Run("Pool Put 测试", func(t *testing.T) {
		t.Run("Put 不存在key的borrow count为1时调用Destroy", func(t *testing.T) {
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

			v, _ := p.GetWithCtx(context.Background(), 1)
			assert.NotNil(t, v)

			err := p.PutWithCtx(context.Background(), 1)
			assert.Nil(t, err)
		})

		t.Run("Put 多次Get后borrow count大于1不触发Destroy", func(t *testing.T) {
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

			v1, _ := p.GetWithCtx(context.Background(), 2)
			v2, _ := p.GetWithCtx(context.Background(), 2)
			assert.Equal(t, v1, v2)

			_ = p.PutWithCtx(context.Background(), 2)
		})
	})
}

func TestPoolGetDefaultCtx(t *testing.T) {
	t.Run("Pool Get 使用默认context", func(t *testing.T) {
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
		assert.Nil(t, err)
		assert.NotNil(t, v)
	})
}

func TestPoolPutDefaultCtx(t *testing.T) {
	t.Run("Pool Put 使用默认context", func(t *testing.T) {
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
		assert.NotNil(t, v)
		err := p.Put(1)
		assert.Nil(t, err)
	})
}

func TestPoolConcurrent(t *testing.T) {
	t.Run("Pool 并发测试", func(t *testing.T) {
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
			assert.Nil(t, e)
		}
	})
}
