package ds

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCounterDiff(t *testing.T) {
	t.Run("counter Diff 测试", func(t *testing.T) {
		t.Run("首次更新返回0", func(t *testing.T) {
			c := NewCounter()
			assert.Equal(t, 0.0, c.Diff(100.0))
		})

		t.Run("第二次更新返回差值", func(t *testing.T) {
			c := NewCounter()
			c.Diff(100.0)
			assert.Equal(t, 50.0, c.Diff(150.0))
		})

		t.Run("值减少返回负差值", func(t *testing.T) {
			c := NewCounter()
			c.Diff(150.0)
			assert.Equal(t, -50.0, c.Diff(100.0))
		})

		t.Run("连续多次更新", func(t *testing.T) {
			c := NewCounter()
			assert.Equal(t, 0.0, c.Diff(10.0))
			assert.Equal(t, 10.0, c.Diff(20.0))
			assert.Equal(t, -5.0, c.Diff(15.0))
			assert.Equal(t, 15.0, c.Diff(30.0))
		})
	})
}

func TestCounterRate(t *testing.T) {
	t.Run("counter Rate 测试", func(t *testing.T) {
		t.Run("首次更新返回0", func(t *testing.T) {
			c := NewCounter()
			assert.Equal(t, 0.0, c.Rate(100.0))
		})
	})
}

func TestCounterRateIn(t *testing.T) {
	t.Run("counter RateIn 测试", func(t *testing.T) {
		t.Run("首次更新返回0", func(t *testing.T) {
			c := NewCounter()
			assert.Equal(t, 0.0, c.RateIn(100.0, time.Second))
		})

		t.Run("interval为0或负数返回0", func(t *testing.T) {
			c := NewCounter()
			c.Diff(100.0)
			time.Sleep(time.Millisecond)
			assert.Equal(t, 0.0, c.RateIn(200.0, 0))
			assert.Equal(t, 0.0, c.RateIn(300.0, -time.Second))
		})

		t.Run("计算速率", func(t *testing.T) {
			c := NewCounter()
			c.Diff(100.0)
			time.Sleep(100 * time.Millisecond)
			rate := c.RateIn(300.0, time.Second)
			assert.Greater(t, rate, 0.0)
			assert.Less(t, rate, 2500.0)
		})
	})
}

func TestMutexCounterDiff(t *testing.T) {
	t.Run("mutexCounter Diff 测试", func(t *testing.T) {
		t.Run("首次更新返回0", func(t *testing.T) {
			c := NewMutexCounter()
			assert.Equal(t, 0.0, c.Diff(100.0))
		})

		t.Run("第二次更新返回差值", func(t *testing.T) {
			c := NewMutexCounter()
			c.Diff(100.0)
			assert.Equal(t, 50.0, c.Diff(150.0))
		})
	})
}

func TestMutexCounterRate(t *testing.T) {
	t.Run("mutexCounter Rate 测试", func(t *testing.T) {
		t.Run("首次更新返回0", func(t *testing.T) {
			c := NewMutexCounter()
			assert.Equal(t, 0.0, c.Rate(100.0))
		})
	})
}

func TestMutexCounterRateIn(t *testing.T) {
	t.Run("mutexCounter RateIn 测试", func(t *testing.T) {
		t.Run("首次更新返回0", func(t *testing.T) {
			c := NewMutexCounter()
			assert.Equal(t, 0.0, c.RateIn(100.0, time.Second))
		})

		t.Run("interval为0返回0", func(t *testing.T) {
			c := NewMutexCounter()
			c.Diff(100.0)
			assert.Equal(t, 0.0, c.RateIn(200.0, 0))
		})
	})
}

func TestMutexCounterConcurrent(t *testing.T) {
	t.Run("mutexCounter 并发测试", func(t *testing.T) {
		t.Run("多个goroutine并发更新", func(t *testing.T) {
			c := NewMutexCounter()
			var wg sync.WaitGroup
			n := 100

			for i := 0; i < n; i++ {
				wg.Add(1)
				go func(v float64) {
					defer wg.Done()
					c.Diff(v)
				}(float64(i * 10))
			}

			wg.Wait()
		})
	})
}

func TestCounterInterface(t *testing.T) {
	t.Run("Counter 接口实现验证", func(t *testing.T) {
		t.Run("counter 实现 Counter 接口", func(t *testing.T) {
			var c Counter = NewCounter()
			assert.NotNil(t, c)
		})

		t.Run("mutexCounter 实现 Counter 接口", func(t *testing.T) {
			var c Counter = NewMutexCounter()
			assert.NotNil(t, c)
		})
	})
}
