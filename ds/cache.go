package ds

import (
	"context"
	"sync"
	"time"
)

type (
	// Cache is a thread-safe in-memory cache with optional expiration and cleanup.
	// It stores key-value pairs with configurable expiration times and automatic
	// cleanup of expired entries.
	//
	// Type parameters:
	//   - T: The type of values stored in the cache
	Cache[T any] struct {
		l       sync.Mutex
		expire  time.Duration
		cleanup time.Duration
		entries map[string]*cacheEntry[T]
		ctx     context.Context
		cf      context.CancelFunc
	}

	// cacheEntry represents a single entry in the cache with its value and expiration time.
	cacheEntry[T any] struct {
		value      T
		expireTime time.Time
	}
)

// NewCache creates a new cache with the specified expiration and cleanup intervals.
//
// Parameters:
//   - expire: Duration after which entries expire. Use 0 for no expiration.
//   - cleanup: Interval for automatic cleanup of expired entries. Use 0 for no automatic cleanup.
//
// Returns:
//   - *Cache[T]: A new cache instance
//
// Example:
//   // Create a cache with 1-minute expiration and 5-minute cleanup interval
//   cache := NewCache[string](time.Minute, 5*time.Minute)
func NewCache[T any](expire, cleanup time.Duration) *Cache[T] {
	c := &Cache[T]{
		entries: make(map[string]*cacheEntry[T]),
		expire:  max(0, expire),
		cleanup: max(0, cleanup),
	}
	c.ctx, c.cf = context.WithCancel(context.Background())

	if c.cleanup > 0 {
		go func() {
			t := time.NewTicker(c.cleanup)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					c.cleanExpireKey()
				case <-c.ctx.Done():
					return
				}
			}
		}()
	}
	return c
}

// cleanExpireKey removes expired entries from the cache.
// This method is called automatically during cleanup cycles.
func (c *Cache[T]) cleanExpireKey() {
	c.l.Lock()
	defer c.l.Unlock()

	now := time.Now()
	for k, e := range c.entries {
		if !e.expireTime.IsZero() && e.expireTime.Before(now) {
			delete(c.entries, k)
		}
	}
}

// Get retrieves a value from the cache by key.
//
// Parameters:
//   - key: The key to look up
//
// Returns:
//   - value: The retrieved value (zero value if not found)
//   - loaded: True if the key was found and not expired
func (c *Cache[T]) Get(key string) (value T, loaded bool) {
	c.l.Lock()
	defer c.l.Unlock()

	e, ok := c.entries[key]
	if !ok || (!e.expireTime.IsZero() && e.expireTime.Before(time.Now())) {
		return value, loaded
	}
	return e.value, true
}

// Set stores a value in the cache with the default expiration time.
//
// Parameters:
//   - key: The key to store the value under
//   - value: The value to store
func (c *Cache[T]) Set(key string, value T) {
	c.SetWithExpire(key, value, c.expire)
}

// SetWithExpire stores a value in the cache with a custom expiration time.
//
// Parameters:
//   - key: The key to store the value under
//   - value: The value to store
//   - expire: Custom expiration duration for this entry
func (c *Cache[T]) SetWithExpire(key string, value T, expire time.Duration) {
	c.l.Lock()
	defer c.l.Unlock()

	e, ok := c.entries[key]
	if ok {
		e.value = value
		if expire > 0 {
			e.expireTime = time.Now().Add(expire)
		}
	} else {
		e = &cacheEntry[T]{value: value}
		if expire > 0 {
			e.expireTime = time.Now().Add(expire)
		}
		c.entries[key] = e
	}
}

// Delete removes an entry from the cache.
//
// Parameters:
//   - key: The key to delete
func (c *Cache[T]) Delete(key string) {
	c.l.Lock()
	defer c.l.Unlock()
	delete(c.entries, key)
}

// Release stops the cleanup goroutine and releases resources.
// Call this method when the cache is no longer needed to prevent goroutine leaks.
//
// Returns:
//   - error: Always returns nil
func (c *Cache[T]) Release() error {
	if c.cf != nil {
		c.cf()
		c.ctx = nil
		c.cf = nil
	}
	return nil
}
