package ds

import (
	"sync"
	"testing"
)

const (
	benchmarkIterations = 10000
	benchmarkGoroutines = 8
)

func BenchmarkBitMutexSingleLockUnlock(b *testing.B) {
	var m BitMutex
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < benchmarkIterations; j++ {
			m.Lock(j % 64)
			m.Unlock(j % 64)
		}
	}
}

func BenchmarkBitMutexSingleTryLock(b *testing.B) {
	var m BitMutex
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < benchmarkIterations; j++ {
			if m.TryLock(j % 64) {
				m.Unlock(j % 64)
			}
		}
	}
}

func BenchmarkBitMutexConcurrentSameBit(b *testing.B) {
	var m BitMutex
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for g := 0; g < benchmarkGoroutines; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < benchmarkIterations/benchmarkGoroutines; j++ {
					m.Lock(5) // 所有goroutine竞争同一个位
					m.Unlock(5)
				}
			}()
		}
		wg.Wait()
	}
}

func BenchmarkBitMutexConcurrentDifferentBits(b *testing.B) {
	var m BitMutex
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for g := 0; g < benchmarkGoroutines; g++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				// 每个goroutine使用不同的位（间隔8个位）
				bit := goroutineID * 8
				for j := 0; j < benchmarkIterations/benchmarkGoroutines; j++ {
					m.Lock(bit % 64)
					m.Unlock(bit % 64)
				}
			}(g)
		}
		wg.Wait()
	}
}

func BenchmarkNFSSingleLockUnlock(b *testing.B) {
	var m NFSBitMutex
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < benchmarkIterations; j++ {
			m.Lock(j % 64)
			m.Unlock(j % 64)
		}
	}
}

func BenchmarkNFSConcurrentSameBit(b *testing.B) {
	var m NFSBitMutex
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for g := 0; g < benchmarkGoroutines; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < benchmarkIterations/benchmarkGoroutines; j++ {
					m.Lock(10)
					m.Unlock(10)
				}
			}()
		}
		wg.Wait()
	}
}

func BenchmarkShardBitMutexSingleLockUnlock(b *testing.B) {
	var m ShardBitMutex
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < benchmarkIterations; j++ {
			m.Lock(j % 512)
			m.Unlock(j % 512)
		}
	}
}

func BenchmarkShardBitMutexConsequentialIndices(b *testing.B) {
	var m ShardBitMutex
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for g := 0; g < benchmarkGoroutines; g++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				// 连续索引（分布在不同的分片上）
				start := goroutineID * 10
				for j := 0; j < benchmarkIterations/benchmarkGoroutines; j++ {
					idx := (start + j) % 512
					m.Lock(idx)
					m.Unlock(idx)
				}
			}(g)
		}
		wg.Wait()
	}
}

func BenchmarkShardBitMutexSameShard(b *testing.B) {
	var m ShardBitMutex
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for g := 0; g < benchmarkGoroutines; g++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				// 所有goroutine使用同一个分片的不同位
				// 索引：0, 8, 16, 24, 32, 40, 48, 56 都在分片0
				bit := goroutineID * 8
				for j := 0; j < benchmarkIterations/benchmarkGoroutines; j++ {
					m.Lock(bit % 512)
					m.Unlock(bit % 512)
				}
			}(g)
		}
		wg.Wait()
	}
}

func BenchmarkShardBitMutexDifferentShards(b *testing.B) {
	var m ShardBitMutex
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for g := 0; g < benchmarkGoroutines; g++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				// 每个goroutine使用不同分片的相同位位置
				// goroutine 0: 位0 (分片0, 位0)
				// goroutine 1: 位1 (分片1, 位0)
				// goroutine 2: 位2 (分片2, 位0)
				// 等等
				bit := goroutineID
				for j := 0; j < benchmarkIterations/benchmarkGoroutines; j++ {
					m.Lock(bit % 512)
					m.Unlock(bit % 512)
				}
			}(g)
		}
		wg.Wait()
	}
}

func BenchmarkCompareBitMutexVsShard(b *testing.B) {
	b.Run("BitMutex-Concurrent-SameBit", func(b *testing.B) {
		var m BitMutex
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for g := 0; g < benchmarkGoroutines; g++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < benchmarkIterations/benchmarkGoroutines; j++ {
						m.Lock(7)
						m.Unlock(7)
					}
				}()
			}
			wg.Wait()
		}
	})

	b.Run("ShardBitMutex-Concurrent-SameBit", func(b *testing.B) {
		var m ShardBitMutex
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for g := 0; g < benchmarkGoroutines; g++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for j := 0; j < benchmarkIterations/benchmarkGoroutines; j++ {
						m.Lock(7)
						m.Unlock(7)
					}
				}()
			}
			wg.Wait()
		}
	})
}

func BenchmarkTryLockSuccessRate(b *testing.B) {
	b.Run("BitMutex-HighContention", func(b *testing.B) {
		var m BitMutex
		// 先锁定一些位
		for i := 0; i < 32; i++ {
			m.Lock(i * 2)
		}

		b.ResetTimer()
		success := 0
		total := 0

		for i := 0; i < b.N; i++ {
			for j := 0; j < 100; j++ {
				total++
				if m.TryLock((j*3 + 1) % 64) {
					success++
					m.Unlock((j*3 + 1) % 64)
				}
			}
		}

		b.ReportMetric(float64(success)/float64(total)*100, "%success")
	})

	b.Run("ShardBitMutex-HighContention", func(b *testing.B) {
		var m ShardBitMutex
		// 先锁定一些位
		for i := 0; i < 128; i++ {
			m.Lock(i * 4)
		}

		b.ResetTimer()
		success := 0
		total := 0

		for i := 0; i < b.N; i++ {
			for j := 0; j < 100; j++ {
				total++
				if m.TryLock((j*5 + 1) % 512) {
					success++
					m.Unlock((j*5 + 1) % 512)
				}
			}
		}

		b.ReportMetric(float64(success)/float64(total)*100, "%success")
	})
}

func BenchmarkMemoryAccessPatterns(b *testing.B) {
	// 测试不同内存访问模式下的性能
	b.Run("BitMutex-LocalizedAccess", func(b *testing.B) {
		var m BitMutex
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for g := 0; g < benchmarkGoroutines; g++ {
				wg.Add(1)
				go func(goroutineID int) {
					defer wg.Done()
					// 每个goroutine访问局部的一组位
					base := goroutineID * 8
					for j := 0; j < benchmarkIterations/benchmarkGoroutines; j++ {
						bit := (base + (j % 8)) % 64
						m.Lock(bit)
						m.Unlock(bit)
					}
				}(g)
			}
			wg.Wait()
		}
	})

	b.Run("ShardBitMutex-LocalizedAccess", func(b *testing.B) {
		var m ShardBitMutex
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var wg sync.WaitGroup
			for g := 0; g < benchmarkGoroutines; g++ {
				wg.Add(1)
				go func(goroutineID int) {
					defer wg.Done()
					// 每个goroutine访问局部的一组位
					base := goroutineID * 64
					for j := 0; j < benchmarkIterations/benchmarkGoroutines; j++ {
						bit := (base + (j % 64)) % 512
						m.Lock(bit)
						m.Unlock(bit)
					}
				}(g)
			}
			wg.Wait()
		}
	})
}

// goos: linux
// goarch: amd64
// pkg: github.com/BYT0723/go-tools/ds
// cpu: 12th Gen Intel(R) Core(TM) i7-12700H
// BenchmarkBitMutexSingleLockUnlock-20                                 8672     121076 ns/op
// BenchmarkBitMutexSingleTryLock-20                                    9898     121829 ns/op
// BenchmarkBitMutexConcurrentSameBit-20                                4926     252847 ns/op
// BenchmarkBitMutexConcurrentDifferentBits-20                          2444     456714 ns/op
// BenchmarkNFSSingleLockUnlock-20                                      9470     118771 ns/op
// BenchmarkNFSConcurrentSameBit-20                                     5863     204827 ns/op
// BenchmarkShardBitMutexSingleLockUnlock-20                            9642     111949 ns/op
// BenchmarkShardBitMutexConsequentialIndices-20                        5121     264077 ns/op
// BenchmarkShardBitMutexSameShard-20                                   2606     409007 ns/op
// BenchmarkShardBitMutexDifferentShards-20                             24830    48599  ns/op
// BenchmarkCompareBitMutexVsShard/BitMutex-Concurrent-SameBit-20       4771     254379 ns/op
// BenchmarkCompareBitMutexVsShard/ShardBitMutex-Concurrent-SameBit-20  6105     204021 ns/op
// BenchmarkTryLockSuccessRate/BitMutex-HighContention-20               1900521  637.2  ns/op 50.00 %success
// BenchmarkTryLockSuccessRate/ShardBitMutex-HighContention-20          1616216  715.2  ns/op 75.00 %success
// BenchmarkMemoryAccessPatterns/BitMutex-LocalizedAccess-20            2371     515204 ns/op
// BenchmarkMemoryAccessPatterns/ShardBitMutex-LocalizedAccess-20       4158     246724 ns/op
