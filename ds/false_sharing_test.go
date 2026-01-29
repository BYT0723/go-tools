// false_sharing_test.go
package ds

import (
	"sync"
	"sync/atomic"
	"testing"
)

const (
	cacheLine = 64
	loops     = 1_000_000
)

/*
❌ 有伪共享
两个 counter 极大概率落在同一个 cache line
*/
type NoPadding struct {
	a uint64
	b uint64
}

/*
✅ 消除伪共享
a 和 b 被 padding 隔离
*/
type WithPadding struct {
	a uint64
	_ [cacheLine]byte
	b uint64
}

/*
基线：单变量
*/
type Single struct {
	a uint64
}

/*
测试在非伪共享的情况下，cas的性能差距
*/

func BenchmarkFalseSharing(b *testing.B) {
	var s NoPadding
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			for j := 0; j < loops; j++ {
				atomic.AddUint64(&s.a, 1)
			}
		}()

		go func() {
			defer wg.Done()
			for j := 0; j < loops; j++ {
				atomic.AddUint64(&s.b, 1)
			}
		}()

		wg.Wait()
	}
}

func BenchmarkNoFalseSharing(b *testing.B) {
	var s WithPadding
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			for j := 0; j < loops; j++ {
				atomic.AddUint64(&s.a, 1)
			}
		}()

		go func() {
			defer wg.Done()
			for j := 0; j < loops; j++ {
				atomic.AddUint64(&s.b, 1)
			}
		}()

		wg.Wait()
	}
}

func BenchmarkSingle(b *testing.B) {
	var s Single
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < loops*2; j++ {
			atomic.AddUint64(&s.a, 1)
		}
	}
}
