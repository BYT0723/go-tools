// Package ds provides concurrent data structures and synchronization primitives.
// This file implements bit-based mutexes that allow fine-grained locking using
// individual bits in a uint64 value.
package ds

import (
	"runtime"
	"sync/atomic"
)

const (
	// maxSpin is the maximum number of spins before exponential backoff stops growing
	maxSpin = 64
	// shardCount is the number of shards used in ShardBitMutex to reduce contention
	shardCount = 8
	// bitsPerShard is the number of bits managed by each shard in ShardBitMutex
	bitsPerShard = 64
)

type (
	// BitMutex is a bit-based mutex that uses a single uint64 to manage up to 64 independent locks.
	// Each bit in the uint64 represents a separate lock that can be acquired or released independently.
	// This is useful for fine-grained locking scenarios where you need to lock specific resources
	// without blocking unrelated operations.
	BitMutex struct {
		val uint64 // atomic value where each bit represents a lock state (1 = locked, 0 = unlocked)
	}

	// NFSBitMutex (No False Sharing Bit Mutex) is a version of BitMutex with padding
	// to prevent false sharing in concurrent systems.
	// False sharing occurs when multiple processors frequently access different variables
	// that happen to be on the same cache line, causing unnecessary cache invalidations.
	// The padding ensures this struct occupies its own cache line (typically 64 bytes).
	NFSBitMutex struct {
		_        [64]byte // Padding before to align to cache line boundary
		BitMutex          // Embedded BitMutex for actual functionality
		_        [64]byte // Padding after to ensure cache line isolation
	}

	// ShardBitMutex is a sharded bit mutex that distributes locks across multiple shards
	// to reduce contention in high-concurrency scenarios.
	// Instead of using a single uint64 for all 512 locks, it uses 8 shards (NFSBitMutex)
	// each managing 64 locks, for a total of 512 available locks.
	// Consecutive lock indices are distributed across different shards to minimize
	// contention when multiple goroutines try to acquire adjacent locks.
	ShardBitMutex struct {
		shards [shardCount]NFSBitMutex // Array of sharded mutexes
	}
)

// Lock acquires the lock at the specified bit position.
// The method uses an exponential backoff algorithm with spinning to minimize
// contention while waiting for the lock to become available.
//
// Parameters:
//   - i: The bit index to lock (0 ≤ i < 64). Each bit represents an independent lock.
//
// Panics:
//   - If i is out of range (i < 0 or i ≥ 64).
//
// Algorithm:
//  1. Check if the target bit is already set (locked)
//  2. If locked, perform exponential backoff with runtime.Gosched() calls
//  3. Retry until successful using atomic CompareAndSwap
//  4. The backoff grows exponentially up to maxSpin spins
func (m *BitMutex) Lock(i int) {
	if i >= 64 || i < 0 {
		panic("NFSBitMutex: index out of range")
	}
	var (
		spin = 1
		mask = uint64(1) << i
	)

	for {
		old := atomic.LoadUint64(&m.val)

		// 如果被占用就自旋
		if old&mask != 0 {
			// 指数退避
			for range spin {
				runtime.Gosched()
			}
			if spin < maxSpin {
				spin *= 2
			}
			continue
		}

		if atomic.CompareAndSwapUint64(&m.val, old, old|mask) {
			return
		}
	}
}

// Unlock releases the lock at the specified bit position.
// This method atomically clears the bit at position i, allowing other goroutines
// to acquire the lock.
//
// Parameters:
//   - i: The bit index to unlock (0 ≤ i < 64)
//
// Panics:
//   - If i is out of range (i < 0 or i ≥ 64)
//
// Note: It is the caller's responsibility to ensure that the lock is held
// before calling Unlock. Calling Unlock on an unlocked bit is undefined behavior.
func (m *BitMutex) Unlock(i int) {
	if i >= 64 || i < 0 {
		panic("BitMutex: index out of range")
	}
	atomic.AndUint64(&m.val, ^(uint64(1) << i))
}

// TryLock attempts to acquire the lock at the specified bit position without blocking.
// It returns true if the lock was successfully acquired, false if the lock is already held.
//
// Parameters:
//   - i: The bit index to lock (0 ≤ i < 64)
//
// Returns:
//   - true if the lock was successfully acquired
//   - false if the lock is already held by another goroutine
//
// Panics:
//   - If i is out of range (i < 0 or i ≥ 64)
//
// Note: This method is useful for non-blocking synchronization patterns where
// failure to acquire a lock should not block the calling goroutine.
func (m *BitMutex) TryLock(i int) bool {
	if i >= 64 || i < 0 {
		panic("BitMutex: index out of range")
	}
	mask := uint64(1) << i

	old := atomic.LoadUint64(&m.val)
	if old&mask != 0 {
		return false
	}
	return atomic.CompareAndSwapUint64(&m.val, old, old|mask)
}

// Lock acquires the lock at the specified bit position in the sharded mutex.
// This method works similarly to BitMutex.Lock but distributes locks across shards
// to reduce contention. Consecutive lock indices are mapped to different shards.
//
// Parameters:
//   - i: The bit index to lock (0 ≤ i < 512). Each bit represents an independent lock.
//
// Panics:
//   - If i is out of range (i < 0 or i ≥ 512)
//
// Algorithm:
//  1. Map the global index i to a shard index and bit index using i2Index
//  2. Check if the target bit is already set (locked) in the corresponding shard
//  3. If locked, perform exponential backoff with runtime.Gosched() calls
//  4. Retry until successful using atomic CompareAndSwap
//  5. The backoff grows exponentially up to maxSpin spins
//
// Note: The sharding strategy maps consecutive indices to different shards to
// minimize contention when multiple goroutines try to acquire adjacent locks.
func (m *ShardBitMutex) Lock(i int) {
	if i < 0 || i >= shardCount*bitsPerShard {
		panic("ShardBitMutex: index out of range")
	}

	var (
		spin = 1
		// 连续的锁被分配到不同的shard中，降低竞争
		shardIndex, bitIndex = i2Index(i)
		mask                 = uint64(1) << bitIndex
	)

	for {
		old := atomic.LoadUint64(&m.shards[shardIndex].val)

		// 如果被占用就自旋
		if old&mask != 0 {
			// 指数退避
			for range spin {
				runtime.Gosched()
			}
			if spin < maxSpin {
				spin *= 2
			}
			continue
		}

		if atomic.CompareAndSwapUint64(&(m.shards[shardIndex].val), old, old|mask) {
			return
		}
	}
}

// Unlock releases the lock at the specified bit position in the sharded mutex.
// This method atomically clears the bit at position i, allowing other goroutines
// to acquire the lock.
//
// Parameters:
//   - i: The bit index to unlock (0 ≤ i < 512)
//
// Panics:
//   - If i is out of range (i < 0 or i ≥ 512)
//
// Note: It is the caller's responsibility to ensure that the lock is held
// before calling Unlock. Calling Unlock on an unlocked bit is undefined behavior.
func (m *ShardBitMutex) Unlock(i int) {
	if i < 0 || i >= shardCount*bitsPerShard {
		panic("ShardBitMutex: index out of range")
	}
	shardIndex, bitIndex := i2Index(i)
	atomic.AndUint64(&(m.shards[shardIndex].val), ^(uint64(1) << uint64(bitIndex)))
}

// TryLock attempts to acquire the lock at the specified bit position without blocking.
// It returns true if the lock was successfully acquired, false if the lock is already held.
//
// Parameters:
//   - i: The bit index to lock (0 ≤ i < 512)
//
// Returns:
//   - true if the lock was successfully acquired
//   - false if the lock is already held by another goroutine
//
// Panics:
//   - If i is out of range (i < 0 or i ≥ 512)
//
// Note: This method is useful for non-blocking synchronization patterns where
// failure to acquire a lock should not block the calling goroutine.
func (m *ShardBitMutex) TryLock(i int) bool {
	if i < 0 || i >= shardCount*bitsPerShard {
		panic("ShardBitMutex: index out of range")
	}

	shardIndex, bitIndex := i2Index(i)

	mask := uint64(1) << uint64(bitIndex)

	old := atomic.LoadUint64(&(m.shards[shardIndex].val))
	if old&mask != 0 {
		return false
	}
	return atomic.CompareAndSwapUint64(&(m.shards[shardIndex].val), old, old|mask)
}

// i2Index maps a global lock index to a shard index and bit index within that shard.
// This function implements the sharding strategy used by ShardBitMutex to distribute
// consecutive lock indices across different shards, reducing contention.
//
// Parameters:
//   - i: The global lock index (0 ≤ i < 512)
//
// Returns:
//   - shardIndex: The index of the shard containing the lock (0 ≤ shardIndex < shardCount)
//   - bitIndex: The bit position within the shard (0 ≤ bitIndex < bitsPerShard)
//
// Algorithm:
//   - shardIndex = i % shardCount (distributes consecutive indices across shards)
//   - bitIndex = i / shardCount (groups indices with the same remainder together)
//
// Example:
//   - i2Index(0) returns (0, 0)  // First lock in first shard
//   - i2Index(1) returns (1, 0)  // First lock in second shard
//   - i2Index(8) returns (0, 1)  // Second lock in first shard
func i2Index(i int) (shardIndex, bitIndex int) {
	shardIndex = i % shardCount
	bitIndex = i / shardCount
	return
}

// Example usage demonstrates how to use the bit-based mutexes in practice.
//
// Example 1: Basic BitMutex usage
//
//	func main() {
//		var mutex ds.BitMutex
//
//		// Lock bit 5
//		mutex.Lock(5)
//		// Perform critical section operations
//		// ...
//		// Unlock bit 5
//		mutex.Unlock(5)
//
//		// Try to lock bit 10 without blocking
//		if mutex.TryLock(10) {
//			// Lock acquired successfully
//			defer mutex.Unlock(10)
//			// Perform operations
//		}
//	}
//
// Example 2: ShardBitMutex for high-concurrency scenarios
//
//	func main() {
//		var shardedMutex ds.ShardBitMutex
//
//		// Multiple goroutines can lock adjacent indices with minimal contention
//		go func() {
//			shardedMutex.Lock(0)
//			defer shardedMutex.Unlock(0)
//			// Process resource 0
//		}()
//
//		go func() {
//			shardedMutex.Lock(1)
//			defer shardedMutex.Unlock(1)
//			// Process resource 1 (different shard, less contention)
//		}()
//
//		// Wait for goroutines to complete
//	}
//
// Example 3: NFSBitMutex for cache-line optimization
//
//	func main() {
//		var nfsMutex ds.NFSBitMutex
//
//		// Use in performance-critical sections where false sharing is a concern
//		nfsMutex.Lock(3)
//		defer nfsMutex.Unlock(3)
//
//		// Critical section with reduced cache contention
//	}
//
// Note: Always ensure proper unlocking, preferably using defer statements
// to prevent deadlocks in case of panics or early returns.
