package ds

import (
	"runtime"
	"sync"
	"time"
)

type (
	Cache[T any] struct {
		l       sync.Mutex
		expire  time.Duration
		cleanup time.Duration
		entries map[string]*entry[T]
	}
	entry[T any] struct {
		value      T
		expireTime time.Time
	}
)

func NewCache[T any](expire, cleanup time.Duration) *Cache[T] {
	c := &Cache[T]{
		entries: make(map[string]*entry[T]),
		expire:  max(0, expire),
		cleanup: cleanup,
	}

	if c.cleanup > 0 {
		exit := make(chan struct{})
		go func() {
			t := time.NewTicker(c.cleanup)
			for {
				select {
				case <-t.C:
					c.cleanExpireKey()
				case <-exit:
					return
				}
			}
		}()
		runtime.SetFinalizer(c, func(c *Cache[T]) {
			exit <- struct{}{}
		})
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
		return
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
		e = &entry[T]{value: value}
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
