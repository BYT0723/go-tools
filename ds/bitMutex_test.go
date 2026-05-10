package ds

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitMutexBasic(t *testing.T) {
	t.Run("BitMutex 基本功能测试", func(t *testing.T) {
		t.Run("单个锁的锁定和解锁", func(t *testing.T) {
			var m BitMutex

			m.Lock(5)
			assert.False(t, m.TryLock(5))
			m.Unlock(5)
			assert.True(t, m.TryLock(5))
			m.Unlock(5)
		})

		t.Run("不同位的独立锁定", func(t *testing.T) {
			var m BitMutex

			m.Lock(3)
			assert.True(t, m.TryLock(10))
			m.Unlock(3)
			m.Unlock(10)

			assert.True(t, m.TryLock(3))
			assert.True(t, m.TryLock(10))
			m.Unlock(3)
			m.Unlock(10)
		})

		t.Run("边界值测试", func(t *testing.T) {
			var m BitMutex

			m.Lock(0)
			assert.False(t, m.TryLock(0))
			m.Unlock(0)

			m.Lock(63)
			assert.False(t, m.TryLock(63))
			m.Unlock(63)

			assert.True(t, m.TryLock(0))
			assert.True(t, m.TryLock(63))
			m.Unlock(0)
			m.Unlock(63)
		})
	})
}

func TestBitMutexPanic(t *testing.T) {
	t.Run("BitMutex 越界panic测试", func(t *testing.T) {
		var m BitMutex

		t.Run("负索引应该panic", func(t *testing.T) {
			assert.Panics(t, func() { m.Lock(-1) })
			assert.Panics(t, func() { m.Unlock(-1) })
			assert.Panics(t, func() { m.TryLock(-1) })
		})

		t.Run("大于等于64的索引应该panic", func(t *testing.T) {
			assert.Panics(t, func() { m.Lock(64) })
			assert.Panics(t, func() { m.Unlock(64) })
			assert.Panics(t, func() { m.TryLock(64) })

			assert.Panics(t, func() { m.Lock(100) })
			assert.Panics(t, func() { m.Unlock(100) })
			assert.Panics(t, func() { m.TryLock(100) })
		})
	})
}

func TestBitMutexConcurrent(t *testing.T) {
	t.Run("BitMutex 并发测试", func(t *testing.T) {
		t.Run("多个goroutine竞争同一个锁", func(t *testing.T) {
			var (
				m     BitMutex
				wg    sync.WaitGroup
				total = 100
			)

			m.Lock(7)

			for i := 0; i < total; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if m.TryLock(7) {
						panic("不应该成功获取锁")
					}
				}()
			}

			wg.Wait()
			m.Unlock(7)
			assert.True(t, m.TryLock(7))
			m.Unlock(7)
		})

		t.Run("多个goroutine锁定不同位", func(t *testing.T) {
			var (
				m  BitMutex
				wg sync.WaitGroup
			)

			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func(bit int) {
					defer wg.Done()
					m.Lock(bit)
					m.Unlock(bit)
				}(i * 6)
			}

			wg.Wait()
			for i := 0; i < 10; i++ {
				assert.True(t, m.TryLock(i*6))
				m.Unlock(i * 6)
			}
		})
	})
}

func TestNFSBitMutex(t *testing.T) {
	t.Run("NFSBitMutex 测试", func(t *testing.T) {
		t.Run("基本功能测试", func(t *testing.T) {
			var m NFSBitMutex

			m.Lock(15)
			assert.False(t, m.TryLock(15))
			m.Unlock(15)
			assert.True(t, m.TryLock(15))
			m.Unlock(15)
		})

		t.Run("并发测试", func(t *testing.T) {
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
	t.Run("ShardBitMutex 测试", func(t *testing.T) {
		t.Run("基本功能测试", func(t *testing.T) {
			var m ShardBitMutex

			m.Lock(0)
			assert.False(t, m.TryLock(0))
			m.Unlock(0)

			m.Lock(511)
			assert.False(t, m.TryLock(511))
			m.Unlock(511)

			m.Lock(256)
			assert.False(t, m.TryLock(256))
			m.Unlock(256)
		})

		t.Run("分片策略测试", func(t *testing.T) {
			testCases := []struct {
				i          int
				shardIndex int
				bitIndex   int
			}{
				{0, 0, 0},
				{1, 1, 0},
				{2, 2, 0},
				{7, 7, 0},
				{8, 0, 1},
				{9, 1, 1},
				{15, 7, 1},
				{500, 4, 62},
				{511, 7, 63},
			}

			for _, tc := range testCases {
				shardIndex, bitIndex := i2Index(tc.i)
				assert.Equal(t, tc.shardIndex, shardIndex)
				assert.Equal(t, tc.bitIndex, bitIndex)
			}
		})

		t.Run("越界panic测试", func(t *testing.T) {
			var m ShardBitMutex

			t.Run("负索引应该panic", func(t *testing.T) {
				assert.Panics(t, func() { m.Lock(-1) })
				assert.Panics(t, func() { m.Unlock(-1) })
				assert.Panics(t, func() { m.TryLock(-1) })
			})

			t.Run("大于等于512的索引应该panic", func(t *testing.T) {
				assert.Panics(t, func() { m.Lock(512) })
				assert.Panics(t, func() { m.Unlock(512) })
				assert.Panics(t, func() { m.TryLock(512) })

				assert.Panics(t, func() { m.Lock(1000) })
				assert.Panics(t, func() { m.Unlock(1000) })
				assert.Panics(t, func() { m.TryLock(1000) })
			})
		})

		t.Run("并发测试 - 连续索引", func(t *testing.T) {
			var (
				m  ShardBitMutex
				wg sync.WaitGroup
			)

			for i := 0; i < 16; i++ {
				wg.Add(1)
				go func(idx int) {
					defer wg.Done()
					m.Lock(idx)
					m.Unlock(idx)
				}(i)
			}

			wg.Wait()

			for i := 0; i < 16; i++ {
				assert.True(t, m.TryLock(i))
				m.Unlock(i)
			}
		})

		t.Run("并发测试 - 同一分片不同位", func(t *testing.T) {
			var (
				m  ShardBitMutex
				wg sync.WaitGroup
			)

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
	t.Run("TryLock 行为测试", func(t *testing.T) {
		t.Run("BitMutex TryLock", func(t *testing.T) {
			var m BitMutex

			assert.True(t, m.TryLock(20))
			assert.False(t, m.TryLock(20))
			m.Unlock(20)
			assert.True(t, m.TryLock(20))
			m.Unlock(20)
		})

		t.Run("ShardBitMutex TryLock", func(t *testing.T) {
			var m ShardBitMutex

			for i := 0; i < 50; i += 5 {
				assert.True(t, m.TryLock(i))
				assert.False(t, m.TryLock(i))
				m.Unlock(i)
				assert.True(t, m.TryLock(i))
				m.Unlock(i)
			}
		})
	})
}

func TestMixedOperations(t *testing.T) {
	t.Run("混合操作测试", func(t *testing.T) {
		t.Run("BitMutex 混合Lock/TryLock", func(t *testing.T) {
			var m BitMutex

			m.Lock(1)
			m.Lock(2)
			m.Lock(3)

			assert.False(t, m.TryLock(1))
			assert.False(t, m.TryLock(2))
			assert.False(t, m.TryLock(3))

			assert.True(t, m.TryLock(4))
			assert.True(t, m.TryLock(5))

			m.Unlock(1)
			m.Unlock(2)
			m.Unlock(3)
			m.Unlock(4)
			m.Unlock(5)

			for i := 1; i <= 5; i++ {
				assert.True(t, m.TryLock(i))
				m.Unlock(i)
			}
		})
	})
}
