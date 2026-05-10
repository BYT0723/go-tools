package ds

import (
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
			defer c.Release()

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

			// 条目过期但尚未被 cleanup goroutine 清除
			_, loaded = c.Get("name")
			assert.False(t, loaded)
			_, loaded = c.Get("age")
			assert.False(t, loaded)

			// 等待 cleanup 周期
			time.Sleep(2 * time.Second)
			_, loaded = c.Get("name")
			assert.False(t, loaded)
			_, loaded = c.Get("age")
			assert.False(t, loaded)
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

			_, loaded = c.Get("name")
			assert.False(t, loaded)
			time.Sleep(2 * time.Second)
			_, loaded = c.Get("name")
			assert.False(t, loaded)
		})
	})
}

func TestCleanupExit(t *testing.T) {
	t.Run("Cache Release 清理goroutine", func(t *testing.T) {
		c := NewCache[any](time.Second, 3*time.Second)
		c.Set("name", "tyler")
		err := c.Release()
		assert.Nil(t, err)
	})
}
