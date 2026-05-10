package ds

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheSetAndGet(t *testing.T) {
	t.Run("Cache Set And Get", func(t *testing.T) {
		t.Run("No Expire Cache", func(t *testing.T) {
			t.Run("Set And Get", func(t *testing.T) {
				c := NewCache[any](0, 0)

				c.Set("name", "tyler")
				c.Set("age", 18)

				value, loaded := c.Get("name")
				assert.Equal(t, "tyler", value)
				assert.True(t, loaded)

				value, loaded = c.Get("age")
				assert.Equal(t, 18, value)
				assert.True(t, loaded)
			})
			t.Run("Set And ReSet And Get", func(t *testing.T) {
				c := NewCache[any](0, 0)

				c.Set("name", "tyler")
				c.Set("age", 18)

				value, loaded := c.Get("name")
				assert.Equal(t, "tyler", value)
				assert.True(t, loaded)

				value, loaded = c.Get("age")
				assert.Equal(t, 18, value)
				assert.True(t, loaded)

				c.Set("name", "walter")
				c.Set("age", 10)

				value, loaded = c.Get("name")
				assert.Equal(t, "walter", value)
				assert.True(t, loaded)

				value, loaded = c.Get("age")
				assert.Equal(t, 10, value)
				assert.True(t, loaded)
			})
		})
		t.Run("Expire Cache", func(t *testing.T) {
			t.Run("Set And Get", func(t *testing.T) {
				c := NewCache[any](time.Second, 0)

				c.Set("name", "tyler")
				c.Set("age", 18)

				value, loaded := c.Get("name")
				assert.Equal(t, "tyler", value)
				assert.True(t, loaded)

				value, loaded = c.Get("age")
				assert.Equal(t, 18, value)
				assert.True(t, loaded)

				time.Sleep(2 * time.Second)

				value, loaded = c.Get("name")
				assert.Nil(t, value)
				assert.False(t, loaded)

				value, loaded = c.Get("age")
				assert.Nil(t, value)
				assert.False(t, loaded)
			})
			t.Run("Set And Reset And Get", func(t *testing.T) {
				c := NewCache[any](time.Second, 0)

				c.Set("name", "tyler")
				value, loaded := c.Get("name")

				assert.Equal(t, "tyler", value)
				assert.True(t, loaded)

				time.Sleep(2 * time.Second)
				value, loaded = c.Get("name")
				assert.Nil(t, value)
				assert.False(t, loaded)

				c.Set("name", "walter")
				value, loaded = c.Get("name")

				assert.Equal(t, "walter", value)
				assert.True(t, loaded)

				c.Set("name", "tyler")
				value, loaded = c.Get("name")
				assert.Equal(t, "tyler", value)
				assert.True(t, loaded)

				time.Sleep(2 * time.Second)
				value, loaded = c.Get("name")
				assert.Nil(t, value)
				assert.False(t, loaded)
			})
		})
		t.Run("Expire And Cleanup Cache", func(t *testing.T) {
			c := NewCache[any](time.Second, 3*time.Second)

			c.Set("name", "tyler")
			c.Set("age", 18)

			value, loaded := c.Get("name")
			assert.Equal(t, "tyler", value)
			assert.True(t, loaded)

			value, loaded = c.Get("age")
			assert.Equal(t, 18, value)
			assert.True(t, loaded)

			time.Sleep(2 * time.Second)
			value, loaded = c.Get("name")

			assert.Nil(t, value)
			assert.False(t, loaded)

			value, loaded = c.Get("age")
			assert.Nil(t, value)
			assert.False(t, loaded)

			assert.NotNil(t, c.entries["name"])
			assert.NotNil(t, c.entries["age"])
			time.Sleep(2 * time.Second)
			assert.Nil(t, c.entries["name"])
			assert.Nil(t, c.entries["age"])
		})
	})
}

func TestCacheDelete(t *testing.T) {
	t.Run("Cache Set", func(t *testing.T) {
		t.Run("No Expire Cache", func(t *testing.T) {
			c := NewCache[any](0, 0)

			c.Set("name", "tyler")
			value, loaded := c.Get("name")

			assert.Equal(t, "tyler", value)
			assert.True(t, loaded)

			c.Delete("name")
			value, loaded = c.Get("name")

			assert.Nil(t, value)
			assert.False(t, loaded)
		})
		t.Run("Expire Cache", func(t *testing.T) {
			c := NewCache[any](time.Second, 0)

			c.Set("name", "tyler")
			value, loaded := c.Get("name")

			assert.Equal(t, "tyler", value)
			assert.True(t, loaded)

			c.Delete("name")
			value, loaded = c.Get("name")

			assert.Nil(t, value)
			assert.False(t, loaded)

			time.Sleep(2 * time.Second)
			value, loaded = c.Get("name")

			assert.Nil(t, value)
			assert.False(t, loaded)
		})
		t.Run("Expire And Cleanup Cache", func(t *testing.T) {
			c := NewCache[any](time.Second, 3*time.Second)

			c.Set("name", "tyler")
			value, loaded := c.Get("name")

			assert.Equal(t, "tyler", value)
			assert.True(t, loaded)

			c.Delete("name")
			value, loaded = c.Get("name")

			assert.Nil(t, value)
			assert.False(t, loaded)

			time.Sleep(2 * time.Second)
			value, loaded = c.Get("name")

			assert.Nil(t, value)
			assert.False(t, loaded)

			assert.Nil(t, c.entries["name"])
			time.Sleep(2 * time.Second)
			assert.Nil(t, c.entries["name"])
		})
	})
}

func TestCleanupExit(t *testing.T) {
	t.Run("Cache Cleanup Exit", func(t *testing.T) {
		c := NewCache[any](time.Second, 3*time.Second)
		c.Set("name", "tyler")
		c = nil
		for i := 0; i < 5; i++ {
			runtime.GC()
			time.Sleep(2 * time.Second)
		}
		assert.Nil(t, c)
	})
}
