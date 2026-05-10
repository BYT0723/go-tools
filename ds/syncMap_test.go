package ds

import (
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSyncMapStoreLoad(t *testing.T) {
	Convey("SyncMap Store/Load 测试", t, func() {
		Convey("Store 后 Load 成功", func() {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			v, ok := m.Load("key1")
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, 100)
		})

		Convey("Load 不存在的key返回false", func() {
			m := NewSyncMap[string, int]()
			v, ok := m.Load("nonexistent")
			So(ok, ShouldBeFalse)
			So(v, ShouldEqual, 0)
		})
	})
}

func TestSyncMapDelete(t *testing.T) {
	Convey("SyncMap Delete 测试", t, func() {
		Convey("Delete 已存在的key返回true", func() {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			So(m.Delete("key1"), ShouldBeTrue)
			_, ok := m.Load("key1")
			So(ok, ShouldBeFalse)
		})

		Convey("Delete 不存在的key仍返回true", func() {
			m := NewSyncMap[string, int]()
			So(m.Delete("key1"), ShouldBeTrue)
		})
	})
}

func TestSyncMapSwap(t *testing.T) {
	Convey("SyncMap Swap 测试", t, func() {
		Convey("Swap 已存在的key返回旧值", func() {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			old, loaded := m.Swap("key1", 200)
			So(loaded, ShouldBeTrue)
			So(old, ShouldEqual, 100)
			v, _ := m.Load("key1")
			So(v, ShouldEqual, 200)
		})

		Convey("Swap 不存在的key loaded返回false", func() {
			m := NewSyncMap[string, int]()
			old, loaded := m.Swap("key1", 200)
			So(loaded, ShouldBeFalse)
			So(old, ShouldEqual, 0)
			v, _ := m.Load("key1")
			So(v, ShouldEqual, 200)
		})
	})
}

func TestSyncMapRange(t *testing.T) {
	Convey("SyncMap Range 测试", t, func() {
		Convey("Range 遍历所有元素", func() {
			m := NewSyncMap[int, string]()
			m.Store(1, "a")
			m.Store(2, "b")
			m.Store(3, "c")

			result := make(map[int]string)
			m.Range(func(k int, v string) bool {
				result[k] = v
				return true
			})
			So(len(result), ShouldEqual, 3)
			So(result[1], ShouldEqual, "a")
			So(result[2], ShouldEqual, "b")
			So(result[3], ShouldEqual, "c")
		})

		Convey("Range 提前终止", func() {
			m := NewSyncMap[int, int]()
			for i := 0; i < 10; i++ {
				m.Store(i, i*10)
			}

			count := 0
			m.Range(func(k, v int) bool {
				count++
				return count < 5
			})
			So(count, ShouldEqual, 5)
		})
	})
}

func TestSyncMapLoadOrStore(t *testing.T) {
	Convey("SyncMap LoadOrStore 测试", t, func() {
		Convey("LoadOrStore 不存在的key", func() {
			m := NewSyncMap[string, int]()
			v, loaded := m.LoadOrStore("key1", 100)
			So(loaded, ShouldBeFalse)
			So(v, ShouldEqual, 100)
		})

		Convey("LoadOrStore 已存在的key", func() {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			v, loaded := m.LoadOrStore("key1", 200)
			So(loaded, ShouldBeTrue)
			So(v, ShouldEqual, 100)
		})
	})
}

func TestSyncMapLoadAndDelete(t *testing.T) {
	Convey("SyncMap LoadAndDelete 测试", t, func() {
		Convey("LoadAndDelete 已存在的key", func() {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			v, loaded := m.LoadAndDelete("key1")
			So(loaded, ShouldBeTrue)
			So(v, ShouldEqual, 100)
			_, ok := m.Load("key1")
			So(ok, ShouldBeFalse)
		})

		Convey("LoadAndDelete 不存在的key", func() {
			m := NewSyncMap[string, int]()
			v, loaded := m.LoadAndDelete("key1")
			So(loaded, ShouldBeFalse)
			So(v, ShouldEqual, 0)
		})
	})
}

func TestSyncMapCompareAndSwap(t *testing.T) {
	Convey("SyncMap CompareAndSwap 测试", t, func() {
		Convey("CompareAndSwap 值匹配", func() {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			So(m.CompareAndSwap("key1", 100, 200), ShouldBeTrue)
			v, _ := m.Load("key1")
			So(v, ShouldEqual, 200)
		})

		Convey("CompareAndSwap 值不匹配", func() {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			So(m.CompareAndSwap("key1", 999, 200), ShouldBeFalse)
			v, _ := m.Load("key1")
			So(v, ShouldEqual, 100)
		})
	})
}

func TestSyncMapCompareAndDelete(t *testing.T) {
	Convey("SyncMap CompareAndDelete 测试", t, func() {
		Convey("CompareAndDelete 值匹配", func() {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			So(m.CompareAndDelete("key1", 100), ShouldBeTrue)
			_, ok := m.Load("key1")
			So(ok, ShouldBeFalse)
		})
	})
}

func TestSyncMapCompareFnAndSwap(t *testing.T) {
	Convey("SyncMap CompareFnAndSwap 测试", t, func() {
		Convey("自定义比较函数匹配", func() {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			So(m.CompareFnAndSwap("key1", func(a, b int) bool { return a == b }, 100, 200), ShouldBeTrue)
			v, _ := m.Load("key1")
			So(v, ShouldEqual, 200)
		})

		Convey("自定义比较函数不匹配", func() {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			So(m.CompareFnAndSwap("key1", func(a, b int) bool { return a > b }, 100, 200), ShouldBeFalse)
			v, _ := m.Load("key1")
			So(v, ShouldEqual, 100)
		})

		Convey("key 不存在返回false", func() {
			m := NewSyncMap[string, int]()
			So(m.CompareFnAndSwap("key1", func(a, b int) bool { return true }, 100, 200), ShouldBeFalse)
		})
	})
}

func TestSyncMapCompareFnAndDelete(t *testing.T) {
	Convey("SyncMap CompareFnAndDelete 测试", t, func() {
		Convey("自定义比较函数匹配", func() {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			So(m.CompareFnAndDelete("key1", func(a, b int) bool { return a == b }, 100), ShouldBeTrue)
			_, ok := m.Load("key1")
			So(ok, ShouldBeFalse)
		})

		Convey("key 不存在返回false", func() {
			m := NewSyncMap[string, int]()
			So(m.CompareFnAndDelete("key1", func(a, b int) bool { return true }, 100), ShouldBeFalse)
		})
	})
}

func TestSyncMapKeys(t *testing.T) {
	Convey("SyncMap Keys 测试", t, func() {
		Convey("返回所有key", func() {
			m := NewSyncMap[int, string]()
			m.Store(1, "a")
			m.Store(2, "b")

			keys := m.Keys()
			So(len(keys), ShouldEqual, 2)
			So(keys, ShouldContain, 1)
			So(keys, ShouldContain, 2)
		})

		Convey("空map返回空slice", func() {
			m := NewSyncMap[int, string]()
			So(m.Keys(), ShouldBeEmpty)
		})
	})
}

func TestSyncMapValues(t *testing.T) {
	Convey("SyncMap Values 测试", t, func() {
		Convey("返回所有value", func() {
			m := NewSyncMap[int, string]()
			m.Store(1, "a")
			m.Store(2, "b")

			values := m.Values()
			So(len(values), ShouldEqual, 2)
			So(values, ShouldContain, "a")
			So(values, ShouldContain, "b")
		})

		Convey("空map返回空slice", func() {
			m := NewSyncMap[int, string]()
			So(m.Values(), ShouldBeEmpty)
		})
	})
}

func TestSyncMapFilter(t *testing.T) {
	Convey("SyncMap Filter 测试", t, func() {
		Convey("Filter 过滤结果", func() {
			m := NewSyncMap[int, int]()
			m.Store(1, 10)
			m.Store(2, 20)
			m.Store(3, 30)

			filtered := m.Filter(func(k, v int) bool {
				return v > 15
			})

			_, ok := filtered.Load(1)
			So(ok, ShouldBeFalse)
			v2, ok := filtered.Load(2)
			So(ok, ShouldBeTrue)
			So(v2, ShouldEqual, 20)
			v3, ok := filtered.Load(3)
			So(ok, ShouldBeTrue)
			So(v3, ShouldEqual, 30)
		})
	})
}

func TestSyncMapConcurrent(t *testing.T) {
	Convey("SyncMap 并发测试", t, func() {
		Convey("并发 Store 和 Load", func() {
			m := NewSyncMap[int, int]()
			var wg sync.WaitGroup
			n := 100

			for i := 0; i < n; i++ {
				wg.Add(1)
				go func(k int) {
					defer wg.Done()
					m.Store(k, k*10)
				}(i)
			}

			wg.Wait()

			for i := 0; i < n; i++ {
				v, ok := m.Load(i)
				So(ok, ShouldBeTrue)
				So(v, ShouldEqual, i*10)
			}
		})
	})
}

func TestSyncMapInterface(t *testing.T) {
	Convey("SyncMap 实现 Map 接口", t, func() {
		var m Map[int, int] = NewSyncMap[int, int]()
		m.Store(1, 10)
		v, ok := m.Load(1)
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, 10)
	})
}
