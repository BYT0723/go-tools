package ds

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncMapStoreLoad(t *testing.T) {
	t.Run("SyncMap Store/Load 测试", func(t *testing.T) {
		t.Run("Store 后 Load 成功", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			v, ok := m.Load("key1")
			assert.True(t, ok)
			assert.Equal(t, 100, v)
		})

		t.Run("Load 不存在的key返回false", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			v, ok := m.Load("nonexistent")
			assert.False(t, ok)
			assert.Equal(t, 0, v)
		})
	})
}

func TestSyncMapDelete(t *testing.T) {
	t.Run("SyncMap Delete 测试", func(t *testing.T) {
		t.Run("Delete 已存在的key返回true", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			assert.True(t, m.Delete("key1"))
			_, ok := m.Load("key1")
			assert.False(t, ok)
		})

		t.Run("Delete 不存在的key仍返回true", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			assert.True(t, m.Delete("key1"))
		})
	})
}

func TestSyncMapSwap(t *testing.T) {
	t.Run("SyncMap Swap 测试", func(t *testing.T) {
		t.Run("Swap 已存在的key返回旧值", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			old, loaded := m.Swap("key1", 200)
			assert.True(t, loaded)
			assert.Equal(t, 100, old)
			v, _ := m.Load("key1")
			assert.Equal(t, 200, v)
		})

		t.Run("Swap 不存在的key loaded返回false", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			old, loaded := m.Swap("key1", 200)
			assert.False(t, loaded)
			assert.Equal(t, 0, old)
			v, _ := m.Load("key1")
			assert.Equal(t, 200, v)
		})
	})
}

func TestSyncMapRange(t *testing.T) {
	t.Run("SyncMap Range 测试", func(t *testing.T) {
		t.Run("Range 遍历所有元素", func(t *testing.T) {
			m := NewSyncMap[int, string]()
			m.Store(1, "a")
			m.Store(2, "b")
			m.Store(3, "c")

			result := make(map[int]string)
			m.Range(func(k int, v string) bool {
				result[k] = v
				return true
			})
			assert.Equal(t, 3, len(result))
			assert.Equal(t, "a", result[1])
			assert.Equal(t, "b", result[2])
			assert.Equal(t, "c", result[3])
		})

		t.Run("Range 提前终止", func(t *testing.T) {
			m := NewSyncMap[int, int]()
			for i := 0; i < 10; i++ {
				m.Store(i, i*10)
			}

			count := 0
			m.Range(func(k, v int) bool {
				count++
				return count < 5
			})
			assert.Equal(t, 5, count)
		})
	})
}

func TestSyncMapLoadOrStore(t *testing.T) {
	t.Run("SyncMap LoadOrStore 测试", func(t *testing.T) {
		t.Run("LoadOrStore 不存在的key", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			v, loaded := m.LoadOrStore("key1", 100)
			assert.False(t, loaded)
			assert.Equal(t, 100, v)
		})

		t.Run("LoadOrStore 已存在的key", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			v, loaded := m.LoadOrStore("key1", 200)
			assert.True(t, loaded)
			assert.Equal(t, 100, v)
		})
	})
}

func TestSyncMapLoadAndDelete(t *testing.T) {
	t.Run("SyncMap LoadAndDelete 测试", func(t *testing.T) {
		t.Run("LoadAndDelete 已存在的key", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			v, loaded := m.LoadAndDelete("key1")
			assert.True(t, loaded)
			assert.Equal(t, 100, v)
			_, ok := m.Load("key1")
			assert.False(t, ok)
		})

		t.Run("LoadAndDelete 不存在的key", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			v, loaded := m.LoadAndDelete("key1")
			assert.False(t, loaded)
			assert.Equal(t, 0, v)
		})
	})
}

func TestSyncMapCompareAndSwap(t *testing.T) {
	t.Run("SyncMap CompareAndSwap 测试", func(t *testing.T) {
		t.Run("CompareAndSwap 值匹配", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			assert.True(t, m.CompareAndSwap("key1", 100, 200))
			v, _ := m.Load("key1")
			assert.Equal(t, 200, v)
		})

		t.Run("CompareAndSwap 值不匹配", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			assert.False(t, m.CompareAndSwap("key1", 999, 200))
			v, _ := m.Load("key1")
			assert.Equal(t, 100, v)
		})
	})
}

func TestSyncMapCompareAndDelete(t *testing.T) {
	t.Run("SyncMap CompareAndDelete 测试", func(t *testing.T) {
		t.Run("CompareAndDelete 值匹配", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			assert.True(t, m.CompareAndDelete("key1", 100))
			_, ok := m.Load("key1")
			assert.False(t, ok)
		})
	})
}

func TestSyncMapCompareFnAndSwap(t *testing.T) {
	t.Run("SyncMap CompareFnAndSwap 测试", func(t *testing.T) {
		t.Run("自定义比较函数匹配", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			assert.True(t, m.CompareFnAndSwap("key1", func(a, b int) bool { return a == b }, 100, 200))
			v, _ := m.Load("key1")
			assert.Equal(t, 200, v)
		})

		t.Run("自定义比较函数不匹配", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			assert.False(t, m.CompareFnAndSwap("key1", func(a, b int) bool { return a > b }, 100, 200))
			v, _ := m.Load("key1")
			assert.Equal(t, 100, v)
		})

		t.Run("key 不存在返回false", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			assert.False(t, m.CompareFnAndSwap("key1", func(a, b int) bool { return true }, 100, 200))
		})
	})
}

func TestSyncMapCompareFnAndDelete(t *testing.T) {
	t.Run("SyncMap CompareFnAndDelete 测试", func(t *testing.T) {
		t.Run("自定义比较函数匹配", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			m.Store("key1", 100)
			assert.True(t, m.CompareFnAndDelete("key1", func(a, b int) bool { return a == b }, 100))
			_, ok := m.Load("key1")
			assert.False(t, ok)
		})

		t.Run("key 不存在返回false", func(t *testing.T) {
			m := NewSyncMap[string, int]()
			assert.False(t, m.CompareFnAndDelete("key1", func(a, b int) bool { return true }, 100))
		})
	})
}

func TestSyncMapKeys(t *testing.T) {
	t.Run("SyncMap Keys 测试", func(t *testing.T) {
		t.Run("返回所有key", func(t *testing.T) {
			m := NewSyncMap[int, string]()
			m.Store(1, "a")
			m.Store(2, "b")

			keys := m.Keys()
			assert.Equal(t, 2, len(keys))
			assert.Contains(t, keys, 1)
			assert.Contains(t, keys, 2)
		})

		t.Run("空map返回空slice", func(t *testing.T) {
			m := NewSyncMap[int, string]()
			assert.Empty(t, m.Keys())
		})
	})
}

func TestSyncMapValues(t *testing.T) {
	t.Run("SyncMap Values 测试", func(t *testing.T) {
		t.Run("返回所有value", func(t *testing.T) {
			m := NewSyncMap[int, string]()
			m.Store(1, "a")
			m.Store(2, "b")

			values := m.Values()
			assert.Equal(t, 2, len(values))
			assert.Contains(t, values, "a")
			assert.Contains(t, values, "b")
		})

		t.Run("空map返回空slice", func(t *testing.T) {
			m := NewSyncMap[int, string]()
			assert.Empty(t, m.Values())
		})
	})
}

func TestSyncMapFilter(t *testing.T) {
	t.Run("SyncMap Filter 测试", func(t *testing.T) {
		t.Run("Filter 过滤结果", func(t *testing.T) {
			m := NewSyncMap[int, int]()
			m.Store(1, 10)
			m.Store(2, 20)
			m.Store(3, 30)

			filtered := m.Filter(func(k, v int) bool {
				return v > 15
			})

			_, ok := filtered.Load(1)
			assert.False(t, ok)
			v2, ok := filtered.Load(2)
			assert.True(t, ok)
			assert.Equal(t, 20, v2)
			v3, ok := filtered.Load(3)
			assert.True(t, ok)
			assert.Equal(t, 30, v3)
		})
	})
}

func TestSyncMapConcurrent(t *testing.T) {
	t.Run("SyncMap 并发测试", func(t *testing.T) {
		t.Run("并发 Store 和 Load", func(t *testing.T) {
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
				assert.True(t, ok)
				assert.Equal(t, i*10, v)
			}
		})
	})
}

func TestSyncMapInterface(t *testing.T) {
	t.Run("SyncMap 实现 Map 接口", func(t *testing.T) {
		var m Map[int, int] = NewSyncMap[int, int]()
		m.Store(1, 10)
		v, ok := m.Load(1)
		assert.True(t, ok)
		assert.Equal(t, 10, v)
	})
}
