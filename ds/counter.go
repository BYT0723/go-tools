package ds

import (
	"sync"
	"time"
)

var (
	_ Counter = (*mutexCounter)(nil)
	_ Counter = (*counter)(nil)
)

type Counter interface {
	// Diff calculates the difference between the new value and the previous value.
	// It updates the counter with the new value and returns the difference.
	//
	// Parameters:
	//   - val: The new value to update the counter with
	//
	// Returns:
	//   - The difference between the new value and the previous value
	//   - Returns 0 if this is the first update
	Diff(float64) float64

	// Rate calculates the rate of change per second.
	// Equivalent to RateIn(val, time.Second).
	//
	// Parameters:
	//   - val: The new value to update the counter with
	//
	// Returns:
	//   - The rate of change per second
	//   - Returns 0 if this is the first update or time elapsed is 0
	Rate(float64) float64

	// RateIn calculates the rate of change per specified time interval.
	// Formula: (new_value - old_value) / (time_elapsed / interval)
	//
	// Parameters:
	//   - val: The new value to update the counter with
	//   - interval: The time interval to normalize the rate to (e.g., time.Second for per-second rate)
	//
	// Returns:
	//   - The rate of change per specified interval
	//   - Returns 0 if this is the first update, interval <= 0, or time elapsed is 0
	RateIn(float64, time.Duration) float64
}

// counter is a non-thread-safe implementation of the Counter interface.
// It stores the current value and the timestamp of the last update.
//
// Use this implementation only when:
//   - You have exclusive access to the counter (single goroutine)
//   - Performance is critical and you want to avoid mutex overhead
//
// Warning: Concurrent access to counter will cause data races.
type counter struct {
	val  float64   // Current value stored in the counter
	last time.Time // Timestamp of the last update
}

// Diff calculates the difference between the new value and the previous value.
// It updates the counter with the new value and returns the difference.
//
// Parameters:
//   - val: The new value to update the counter with
//
// Returns:
//   - The difference between the new value and the previous value
//   - Returns 0 if this is the first update
//
// Example:
//
//	counter := NewCounter()
//	diff1 := counter.Diff(100.0) // Returns 0 (first update)
//	diff2 := counter.Diff(150.0) // Returns 50.0 (150 - 100)
//	diff3 := counter.Diff(120.0) // Returns -30.0 (120 - 150)
func (c *counter) Diff(val float64) (value float64) {
	if !c.last.IsZero() {
		value = val - c.val
	}

	c.val = val
	c.last = time.Now()

	return
}

// Rate calculates the rate of change per second.
// Equivalent to RateIn(val, time.Second).
//
// Parameters:
//   - val: The new value to update the counter with
//
// Returns:
//   - The rate of change per second
//   - Returns 0 if this is the first update or time elapsed is 0
//
// Example:
//
//	counter := NewCounter{}
//	counter.Diff(100.0) // First update
//	time.Sleep(time.Second)
//	rate := counter.Rate(200.0) // Returns approximately 100.0 (100 bytes per second)
func (c *counter) Rate(val float64) float64 {
	return c.RateIn(val, time.Second)
}

// RateIn calculates the rate of change per specified time interval.
// Formula: (new_value - old_value) / (time_elapsed / interval)
//
// Parameters:
//   - val: The new value to update the counter with
//   - interval: The time interval to normalize the rate to (e.g., time.Second for per-second rate)
//
// Returns:
//   - The rate of change per specified interval
//   - Returns 0 if this is the first update, interval <= 0, or time elapsed is 0
//
// Example:
//
//	counter := NewCounter{}
//	counter.Diff(100.0) // First update
//	time.Sleep(100 * time.Millisecond)
//	ratePerSec := counter.RateIn(200.0, time.Second) // Returns approximately 1000.0 (1000 bytes per second)
//	ratePerMin := counter.RateIn(200.0, time.Minute) // Returns approximately 60000.0 (60000 bytes per minute)
func (c *counter) RateIn(val float64, interval time.Duration) (value float64) {
	now := time.Now()

	if !c.last.IsZero() && interval > 0 {
		if dur := float64(now.Sub(c.last)) / float64(interval); dur > 0 {
			// Calculate rate per interval using floating point division
			value = (val - c.val) / dur
		}
	}

	c.val = val
	c.last = now
	return
}

// mutexCounter is a thread-safe implementation of the Counter interface.
// It wraps counter with a mutex to provide safe concurrent access.
//
// Use this implementation when:
//   - Multiple goroutines need to access the counter concurrently
//   - Thread safety is required
//
// Performance considerations:
//   - Mutex overhead may impact performance in high-throughput scenarios
//   - For single-threaded use, prefer counter for better performance
type mutexCounter struct {
	counter
	mu sync.Mutex
}

// Diff calculates the difference between the new value and the previous value.
// This is a thread-safe version of counter.Diff().
//
// Parameters:
//   - val: The new value to update the counter with
//
// Returns:
//   - The difference between the new value and the previous value
//   - Returns 0 if this is the first update
//
// Thread safety: This method is thread-safe.
func (c *mutexCounter) Diff(val float64) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counter.Diff(val)
}

// Rate calculates the rate of change per second.
// This is a thread-safe version of counter.Rate().
//
// Parameters:
//   - val: The new value to update the counter with
//
// Returns:
//   - The rate of change per second
//   - Returns 0 if this is the first update or time elapsed is 0
//
// Thread safety: This method is thread-safe.
func (c *mutexCounter) Rate(val float64) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counter.Rate(val)
}

// RateIn calculates the rate of change per specified time interval.
// This is a thread-safe version of counter.RateIn().
//
// Parameters:
//   - val: The new value to update the counter with
//   - interval: The time interval to normalize the rate to (e.g., time.Second for per-second rate)
//
// Returns:
//   - The rate of change per specified interval
//   - Returns 0 if this is the first update, interval <= 0, or time elapsed is 0
//
// Thread safety: This method is thread-safe.
func (c *mutexCounter) RateIn(val float64, interval time.Duration) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counter.RateIn(val, interval)
}

func NewCounter() *counter           { return &counter{} }
func NewMutexCounter() *mutexCounter { return &mutexCounter{} }
