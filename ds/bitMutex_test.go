package ds

import (
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBitMutexBasic(t *testing.T) {
	Convey("BitMutex 基本功能测试", t, func() {
		Convey("单个锁的锁定和解锁", func() {
			var m BitMutex

			// 锁定位5
			m.Lock(5)
			// 尝试锁定同一个位应该失败
			So(m.TryLock(5), ShouldBeFalse)
			// 解锁
			m.Unlock(5)
			// 解锁后应该可以再次锁定
			So(m.TryLock(5), ShouldBeTrue)
			m.Unlock(5)
		})

		Convey("不同位的独立锁定", func() {
			var m BitMutex

			// 锁定位3
			m.Lock(3)
			// 位10应该可以独立锁定
			So(m.TryLock(10), ShouldBeTrue)
			// 解锁位3
			m.Unlock(3)
			// 解锁位10
			m.Unlock(10)

			// 验证两个位都解锁了
			So(m.TryLock(3), ShouldBeTrue)
			So(m.TryLock(10), ShouldBeTrue)
			m.Unlock(3)
			m.Unlock(10)
		})

		Convey("边界值测试", func() {
			var m BitMutex

			// 测试最小边界
			m.Lock(0)
			So(m.TryLock(0), ShouldBeFalse)
			m.Unlock(0)

			// 测试最大边界
			m.Lock(63)
			So(m.TryLock(63), ShouldBeFalse)
			m.Unlock(63)

			// 解锁后应该可以再次锁定
			So(m.TryLock(0), ShouldBeTrue)
			So(m.TryLock(63), ShouldBeTrue)
			m.Unlock(0)
			m.Unlock(63)
		})
	})
}

func TestBitMutexPanic(t *testing.T) {
	Convey("BitMutex 越界panic测试", t, func() {
		var m BitMutex

		Convey("负索引应该panic", func() {
			So(func() { m.Lock(-1) }, ShouldPanic)
			So(func() { m.Unlock(-1) }, ShouldPanic)
			So(func() { m.TryLock(-1) }, ShouldPanic)
		})

		Convey("大于等于64的索引应该panic", func() {
			So(func() { m.Lock(64) }, ShouldPanic)
			So(func() { m.Unlock(64) }, ShouldPanic)
			So(func() { m.TryLock(64) }, ShouldPanic)

			So(func() { m.Lock(100) }, ShouldPanic)
			So(func() { m.Unlock(100) }, ShouldPanic)
			So(func() { m.TryLock(100) }, ShouldPanic)
		})
	})
}

func TestBitMutexConcurrent(t *testing.T) {
	Convey("BitMutex 并发测试", t, func() {
		Convey("多个goroutine竞争同一个锁", func() {
			var (
				m     BitMutex
				wg    sync.WaitGroup
				total = 100
			)

			// 先锁定位7
			m.Lock(7)

			for i := 0; i < total; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if m.TryLock(7) {
						// 这个应该永远不会成功，因为锁已经被持有
						panic("不应该成功获取锁")
					}
				}()
			}

			wg.Wait()
			// 解锁后，应该可以获取锁
			m.Unlock(7)
			So(m.TryLock(7), ShouldBeTrue)
			m.Unlock(7)
		})

		Convey("多个goroutine锁定不同位", func() {
			var (
				m  BitMutex
				wg sync.WaitGroup
			)

			// 测试10个不同的位
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func(bit int) {
					defer wg.Done()
					m.Lock(bit)
					// 短暂持有锁
					m.Unlock(bit)
				}(i * 6) // 使用间隔的位减少冲突
			}

			wg.Wait()
			// 所有锁都应该已释放
			for i := 0; i < 10; i++ {
				So(m.TryLock(i*6), ShouldBeTrue)
				m.Unlock(i * 6)
			}
		})
	})
}

func TestNFSBitMutex(t *testing.T) {
	Convey("NFSBitMutex 测试", t, func() {
		Convey("基本功能测试", func() {
			var m NFSBitMutex

			// 测试锁定和解锁
			m.Lock(15)
			So(m.TryLock(15), ShouldBeFalse)
			m.Unlock(15)
			So(m.TryLock(15), ShouldBeTrue)
			m.Unlock(15)
		})

		Convey("并发测试", func() {
			var (
				m  NFSBitMutex
				wg sync.WaitGroup
			)

			for i := 0; i < 5; i++ {
				wg.Add(1)
				go func(bit int) {
					defer wg.Done()
					m.Lock(bit)
					m.Unlock(bit)
				}(i * 10)
			}

			wg.Wait()
		})
	})
}

func TestShardBitMutex(t *testing.T) {
	Convey("ShardBitMutex 测试", t, func() {
		Convey("基本功能测试", func() {
			var m ShardBitMutex

			// 测试边界值
			m.Lock(0)   // 第一个分片的第一个位
			So(m.TryLock(0), ShouldBeFalse)
			m.Unlock(0)

			m.Lock(511) // 最后一个分片的最后一个位
			So(m.TryLock(511), ShouldBeFalse)
			m.Unlock(511)

			// 测试中间值
			m.Lock(256)
			So(m.TryLock(256), ShouldBeFalse)
			m.Unlock(256)
		})

		Convey("分片策略测试", func() {
			// 验证i2Index函数的分片策略
			// 连续索引应该分布到不同分片
			testCases := []struct {
				i          int
				shardIndex int
				bitIndex   int
			}{
				{0, 0, 0},   // 第一个分片第一个位
				{1, 1, 0},   // 第二个分片第一个位
				{2, 2, 0},   // 第三个分片第一个位
				{7, 7, 0},   // 第八个分片第一个位
				{8, 0, 1},   // 第一个分片第二个位
				{9, 1, 1},   // 第二个分片第二个位
				{15, 7, 1},  // 第八个分片第二个位
				{500, 4, 62}, // 第五个分片第63位
				{511, 7, 63}, // 第八个分片第64位
			}

			for _, tc := range testCases {
				shardIndex, bitIndex := i2Index(tc.i)
				So(shardIndex, ShouldEqual, tc.shardIndex)
				So(bitIndex, ShouldEqual, tc.bitIndex)
			}
		})

		Convey("越界panic测试", func() {
			var m ShardBitMutex

			Convey("负索引应该panic", func() {
				So(func() { m.Lock(-1) }, ShouldPanic)
				So(func() { m.Unlock(-1) }, ShouldPanic)
				So(func() { m.TryLock(-1) }, ShouldPanic)
			})

			Convey("大于等于512的索引应该panic", func() {
				So(func() { m.Lock(512) }, ShouldPanic)
				So(func() { m.Unlock(512) }, ShouldPanic)
				So(func() { m.TryLock(512) }, ShouldPanic)

				So(func() { m.Lock(1000) }, ShouldPanic)
				So(func() { m.Unlock(1000) }, ShouldPanic)
				So(func() { m.TryLock(1000) }, ShouldPanic)
			})
		})

		Convey("并发测试 - 连续索引", func() {
			var (
				m  ShardBitMutex
				wg sync.WaitGroup
			)

			// 测试连续索引（应该分布到不同分片）
			for i := 0; i < 16; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					m.Lock(idx)
					// 短暂持有
					m.Unlock(idx)
				}(i)
			}

			wg.Wait()

			// 验证所有锁都已释放
			for i := 0; i < 16; i++ {
				So(m.TryLock(i), ShouldBeTrue)
				m.Unlock(i)
			}
		})

		Convey("并发测试 - 同一分片不同位", func() {
			var (
				m  ShardBitMutex
				wg sync.WaitGroup
			)

			// 这些索引都在同一个分片（分片0）
			// 0, 8, 16, 24, 32, 40, 48, 56 都在分片0
			indices := []int{0, 8, 16, 24, 32, 40, 48, 56}

			for _, idx := range indices {
				wg.Add(1)
				go func(bit int) {
					defer wg.Done()
					m.Lock(bit)
					m.Unlock(bit)
				}(idx)
			}

			wg.Wait()
		})
	})
}

func TestTryLockBehavior(t *testing.T) {
	Convey("TryLock 行为测试", t, func() {
		Convey("BitMutex TryLock", func() {
			var m BitMutex

			// 初始状态应该可以获取锁
			So(m.TryLock(20), ShouldBeTrue)
			// 再次尝试应该失败
			So(m.TryLock(20), ShouldBeFalse)
			// 解锁后应该可以再次获取
			m.Unlock(20)
			So(m.TryLock(20), ShouldBeTrue)
			m.Unlock(20)
		})

		Convey("ShardBitMutex TryLock", func() {
			var m ShardBitMutex

			// 测试多个不同位置的TryLock
			for i := 0; i < 50; i += 5 {
				So(m.TryLock(i), ShouldBeTrue)
				So(m.TryLock(i), ShouldBeFalse)
				m.Unlock(i)
				So(m.TryLock(i), ShouldBeTrue)
				m.Unlock(i)
			}
		})
	})
}

func TestMixedOperations(t *testing.T) {
	Convey("混合操作测试", t, func() {
		Convey("BitMutex 混合Lock/TryLock", func() {
			var m BitMutex

			// 使用Lock锁定一些位
			m.Lock(1)
			m.Lock(2)
			m.Lock(3)

			// 验证这些位被锁定
			So(m.TryLock(1), ShouldBeFalse)
			So(m.TryLock(2), ShouldBeFalse)
			So(m.TryLock(3), ShouldBeFalse)

			// 其他位应该可以锁定
			So(m.TryLock(4), ShouldBeTrue)
			So(m.TryLock(5), ShouldBeTrue)

			// 解锁所有位
			m.Unlock(1)
			m.Unlock(2)
			m.Unlock(3)
			m.Unlock(4)
			m.Unlock(5)

			// 验证所有位都解锁了
			for i := 1; i <= 5; i++ {
				So(m.TryLock(i), ShouldBeTrue)
				m.Unlock(i)
			}
		})
	})
}

