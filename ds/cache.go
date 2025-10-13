package ds

import (
	"context"
	"sync"
	"time"
)

type (
	Cache[T any] struct {
		l       sync.Mutex
		expire  time.Duration
		cleanup time.Duration
		entries map[string]*cacheEntry[T]
		ctx     context.Context
		cf      context.CancelFunc
	}
	cacheEntry[T any] struct {
		value      T
		expireTime time.Time
	}
)

// NewCache 新建缓存
// expire 0 表示不过期
// cleanup 0 表示不定时清理
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

func (c *Cache[T]) Get(key string) (value T, loaded bool) {
	c.l.Lock()
	defer c.l.Unlock()

	e, ok := c.entries[key]
	if !ok || (!e.expireTime.IsZero() && e.expireTime.Before(time.Now())) {
		return value, loaded
	}
	return e.value, true
}

func (c *Cache[T]) Set(key string, value T) {
	c.SetWithExpire(key, value, c.expire)
}

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

func (c *Cache[T]) Delete(key string) {
	c.l.Lock()
	defer c.l.Unlock()
	delete(c.entries, key)
}

// Release releases the goroutine for scheduled cleaning
func (c *Cache[T]) Release() error {
	if c.cf != nil {
		c.cf()
		c.ctx = nil
		c.cf = nil
	}
	return nil
}
