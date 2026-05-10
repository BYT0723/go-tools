package ds

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheSetAndGet(t *testing.T) {
	t.Run("No Expire", func(t *testing.T) {
		c := NewCache[any](0, 0)

		c.Set("name", "tyler")
		v, ok := c.Get("name")
		assert.Equal(t, "tyler", v)
		assert.True(t, ok)

		// reset
		c.Set("name", "walter")
		v, ok = c.Get("name")
		assert.Equal(t, "walter", v)
		assert.True(t, ok)

		// missing key
		v, ok = c.Get("nope")
		assert.Nil(t, v)
		assert.False(t, ok)
	})

	t.Run("Expire", func(t *testing.T) {
		c := NewCache[any](50*time.Millisecond, 0)

		c.Set("name", "tyler")
		v, ok := c.Get("name")
		assert.Equal(t, "tyler", v)
		assert.True(t, ok)

		time.Sleep(100 * time.Millisecond)
		v, ok = c.Get("name")
		assert.Nil(t, v)
		assert.False(t, ok)

		// reset after expire
		c.Set("name", "walter")
		v, ok = c.Get("name")
		assert.Equal(t, "walter", v)
		assert.True(t, ok)
	})

	t.Run("Cleanup", func(t *testing.T) {
		c := NewCache[any](50*time.Millisecond, 100*time.Millisecond)
		defer c.Release()

		c.Set("name", "tyler")
		c.Set("age", 18)

		// 等待过期
		time.Sleep(100 * time.Millisecond)

		// Get 返回 expired
		_, ok := c.Get("name")
		assert.False(t, ok)
		_, ok = c.Get("age")
		assert.False(t, ok)
	})
}

func TestCacheDelete(t *testing.T) {
	t.Run("No Expire", func(t *testing.T) {
		c := NewCache[any](0, 0)

		c.Set("name", "tyler")
		v, ok := c.Get("name")
		assert.Equal(t, "tyler", v)
		assert.True(t, ok)

		c.Delete("name")
		v, ok = c.Get("name")
		assert.Nil(t, v)
		assert.False(t, ok)
	})

	t.Run("Expire", func(t *testing.T) {
		c := NewCache[any](50*time.Millisecond, 0)

		c.Set("name", "tyler")
		c.Delete("name")
		v, ok := c.Get("name")
		assert.Nil(t, v)
		assert.False(t, ok)

		// 即使等待过期后仍无数据
		time.Sleep(100 * time.Millisecond)
		v, ok = c.Get("name")
		assert.Nil(t, v)
		assert.False(t, ok)
	})
}

func TestCacheRelease(t *testing.T) {
	c := NewCache[any](time.Second, 100*time.Millisecond)
	c.Set("name", "tyler")
	err := c.Release()
	assert.Nil(t, err)
}
