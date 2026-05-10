package ds

import (
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCounterDiff(t *testing.T) {
	Convey("counter Diff 测试", t, func() {
		Convey("首次更新返回0", func() {
			c := NewCounter()
			So(c.Diff(100.0), ShouldEqual, 0)
		})

		Convey("第二次更新返回差值", func() {
			c := NewCounter()
			c.Diff(100.0)
			So(c.Diff(150.0), ShouldEqual, 50.0)
		})

		Convey("值减少返回负差值", func() {
			c := NewCounter()
			c.Diff(150.0)
			So(c.Diff(100.0), ShouldEqual, -50.0)
		})

		Convey("连续多次更新", func() {
			c := NewCounter()
			So(c.Diff(10.0), ShouldEqual, 0)
			So(c.Diff(20.0), ShouldEqual, 10.0)
			So(c.Diff(15.0), ShouldEqual, -5.0)
			So(c.Diff(30.0), ShouldEqual, 15.0)
		})
	})
}

func TestCounterRate(t *testing.T) {
	Convey("counter Rate 测试", t, func() {
		Convey("首次更新返回0", func() {
			c := NewCounter()
			So(c.Rate(100.0), ShouldEqual, 0)
		})
	})
}

func TestCounterRateIn(t *testing.T) {
	Convey("counter RateIn 测试", t, func() {
		Convey("首次更新返回0", func() {
			c := NewCounter()
			So(c.RateIn(100.0, time.Second), ShouldEqual, 0)
		})

		Convey("interval为0或负数返回0", func() {
			c := NewCounter()
			c.Diff(100.0)
			time.Sleep(time.Millisecond)
			So(c.RateIn(200.0, 0), ShouldEqual, 0)
			So(c.RateIn(300.0, -time.Second), ShouldEqual, 0)
		})

		Convey("计算速率", func() {
			c := NewCounter()
			c.Diff(100.0)
			time.Sleep(100 * time.Millisecond)
			rate := c.RateIn(300.0, time.Second)
			// 差值=200, 时间=约0.1秒, 标准化到每秒 ≈ 2000
			So(rate, ShouldBeGreaterThan, 0)
			So(rate, ShouldBeLessThan, 2500.0)
		})
	})
}

func TestMutexCounterDiff(t *testing.T) {
	Convey("mutexCounter Diff 测试", t, func() {
		Convey("首次更新返回0", func() {
			c := NewMutexCounter()
			So(c.Diff(100.0), ShouldEqual, 0)
		})

		Convey("第二次更新返回差值", func() {
			c := NewMutexCounter()
			c.Diff(100.0)
			So(c.Diff(150.0), ShouldEqual, 50.0)
		})
	})
}

func TestMutexCounterRate(t *testing.T) {
	Convey("mutexCounter Rate 测试", t, func() {
		Convey("首次更新返回0", func() {
			c := NewMutexCounter()
			So(c.Rate(100.0), ShouldEqual, 0)
		})
	})
}

func TestMutexCounterRateIn(t *testing.T) {
	Convey("mutexCounter RateIn 测试", t, func() {
		Convey("首次更新返回0", func() {
			c := NewMutexCounter()
			So(c.RateIn(100.0, time.Second), ShouldEqual, 0)
		})

		Convey("interval为0返回0", func() {
			c := NewMutexCounter()
			c.Diff(100.0)
			So(c.RateIn(200.0, 0), ShouldEqual, 0)
		})
	})
}

func TestMutexCounterConcurrent(t *testing.T) {
	Convey("mutexCounter 并发测试", t, func() {
		Convey("多个goroutine并发更新", func() {
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
	Convey("Counter 接口实现验证", t, func() {
		Convey("counter 实现 Counter 接口", func() {
			var c Counter = NewCounter()
			So(c, ShouldNotBeNil)
		})

		Convey("mutexCounter 实现 Counter 接口", func() {
			var c Counter = NewMutexCounter()
			So(c, ShouldNotBeNil)
		})
	})
}
