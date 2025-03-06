package ds

import (
	"runtime"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCacheSetAndGet(t *testing.T) {
	Convey("Cache Set And Get", t, func() {
		Convey("No Expire Cache", func() {
			Convey("Set And Get", func() {
				c := NewCache[any](nil)

				c.Set("name", "tyler")
				c.Set("age", 18)

				value, loaded := c.Get("name")
				So(value, ShouldEqual, "tyler")
				So(loaded, ShouldBeTrue)

				value, loaded = c.Get("age")
				So(value, ShouldEqual, 18)
				So(loaded, ShouldBeTrue)
			})
			Convey("Set And ReSet And Get", func() {
				c := NewCache[any](nil)

				c.Set("name", "tyler")
				c.Set("age", 18)

				value, loaded := c.Get("name")
				So(value, ShouldEqual, "tyler")
				So(loaded, ShouldBeTrue)

				value, loaded = c.Get("age")
				So(value, ShouldEqual, 18)
				So(loaded, ShouldBeTrue)

				// reset key value
				c.Set("name", "walter")
				c.Set("age", 10)

				value, loaded = c.Get("name")
				So(value, ShouldEqual, "walter")
				So(loaded, ShouldBeTrue)

				value, loaded = c.Get("age")
				So(value, ShouldEqual, 10)
				So(loaded, ShouldBeTrue)
			})
		})
		Convey("Expire Cache", func() {
			Convey("Set And Get", func() {
				c := NewCache[any](&CacheOpt{Expire: 1 * time.Second})

				c.Set("name", "tyler")
				c.Set("age", 18)

				value, loaded := c.Get("name")
				So(value, ShouldEqual, "tyler")
				So(loaded, ShouldBeTrue)

				value, loaded = c.Get("age")
				So(value, ShouldEqual, 18)
				So(loaded, ShouldBeTrue)

				// 等待key过期
				time.Sleep(2 * time.Second)

				value, loaded = c.Get("name")
				So(value, ShouldBeNil)
				So(loaded, ShouldBeFalse)

				value, loaded = c.Get("age")
				So(value, ShouldBeNil)
				So(loaded, ShouldBeFalse)
			})
			Convey("Set And Reset And Get", func() {
				c := NewCache[any](&CacheOpt{Expire: 1 * time.Second})

				c.Set("name", "tyler")
				value, loaded := c.Get("name")

				So(value, ShouldEqual, "tyler")
				So(loaded, ShouldBeTrue)

				// 等待key过期
				time.Sleep(2 * time.Second)
				value, loaded = c.Get("name")
				So(value, ShouldBeNil)
				So(loaded, ShouldBeFalse)

				// 重置key
				c.Set("name", "walter")
				value, loaded = c.Get("name")

				So(value, ShouldEqual, "walter")
				So(loaded, ShouldBeTrue)

				// 重置key
				c.Set("name", "tyler")
				value, loaded = c.Get("name")
				So(value, ShouldEqual, "tyler")
				So(loaded, ShouldBeTrue)

				// 等待key过期
				time.Sleep(2 * time.Second)
				value, loaded = c.Get("name")
				So(value, ShouldBeNil)
				So(loaded, ShouldBeFalse)
			})
		})
		Convey("Expire And Cleanup Cache", func() {
			c := NewCache[any](&CacheOpt{Expire: 1 * time.Second, Cleanup: 3 * time.Second})

			c.Set("name", "tyler")
			c.Set("age", 18)

			value, loaded := c.Get("name")
			So(value, ShouldEqual, "tyler")
			So(loaded, ShouldBeTrue)

			value, loaded = c.Get("age")
			So(value, ShouldEqual, 18)
			So(loaded, ShouldBeTrue)

			// 等待key过期
			time.Sleep(2 * time.Second)
			value, loaded = c.Get("name")

			So(value, ShouldBeNil)
			So(loaded, ShouldBeFalse)

			value, loaded = c.Get("age")
			So(value, ShouldBeNil)
			So(loaded, ShouldBeFalse)

			So(c.entries["name"], ShouldNotBeNil)
			So(c.entries["age"], ShouldNotBeNil)
			// 等待key清除
			time.Sleep(2 * time.Second)
			So(c.entries["name"], ShouldBeNil)
			So(c.entries["age"], ShouldBeNil)
		})
	})
}

func TestCacheDelete(t *testing.T) {
	Convey("Cache Set", t, func() {
		Convey("No Expire Cache", func() {
			c := NewCache[any](nil)

			c.Set("name", "tyler")
			value, loaded := c.Get("name")

			So(value, ShouldEqual, "tyler")
			So(loaded, ShouldBeTrue)

			c.Delete("name")
			value, loaded = c.Get("name")

			So(value, ShouldBeNil)
			So(loaded, ShouldBeFalse)
		})
		Convey("Expire Cache", func() {
			c := NewCache[any](&CacheOpt{Expire: 1 * time.Second})

			c.Set("name", "tyler")
			value, loaded := c.Get("name")

			So(value, ShouldEqual, "tyler")
			So(loaded, ShouldBeTrue)

			// 清除key
			c.Delete("name")
			value, loaded = c.Get("name")

			So(value, ShouldBeNil)
			So(loaded, ShouldBeFalse)

			// 等待过期
			time.Sleep(2 * time.Second)
			value, loaded = c.Get("name")

			So(value, ShouldBeNil)
			So(loaded, ShouldBeFalse)
		})
		Convey("Expire And Cleanup Cache", func() {
			c := NewCache[any](&CacheOpt{Expire: 1 * time.Second, Cleanup: 3 * time.Second})

			c.Set("name", "tyler")
			value, loaded := c.Get("name")

			So(value, ShouldEqual, "tyler")
			So(loaded, ShouldBeTrue)

			// 清除key
			c.Delete("name")
			value, loaded = c.Get("name")

			So(value, ShouldBeNil)
			So(loaded, ShouldBeFalse)

			// 等待过期
			time.Sleep(2 * time.Second)
			value, loaded = c.Get("name")

			So(value, ShouldBeNil)
			So(loaded, ShouldBeFalse)

			So(c.entries["name"], ShouldBeNil)
			// 等待key清除
			time.Sleep(2 * time.Second)
			So(c.entries["name"], ShouldBeNil)
		})
	})
}

func TestCleanupExit(t *testing.T) {
	Convey("Cache Cleanup Exit", t, func() {
		c := NewCache[any](&CacheOpt{Expire: time.Second, Cleanup: 3 * time.Second})
		c.Set("name", "tyler")
		c = nil
		// 强制 GC 多次，确保回收
		for i := 0; i < 5; i++ {
			runtime.GC()
			time.Sleep(2 * time.Second)
		}
		So(c, ShouldBeNil)
	})
}
