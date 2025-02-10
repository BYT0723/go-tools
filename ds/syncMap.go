package ds

import "sync"

var _ Map[int, int] = (*SyncMap[int, int])(nil)

// Generic wrapper for sync.Map
type SyncMap[K comparable, V any] struct {
	entries sync.Map
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{}
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.entries.Store(key, value)
}

func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.entries.Load(key)
	if !ok {
		return
	}
	value, ok = v.(V)
	return
}

func (m *SyncMap[K, V]) Delete(key K) bool {
	m.entries.Delete(key)
	return true
}

func (m *SyncMap[K, V]) Swap(key K, value V) (pre V, loaded bool) {
	previous, loaded := m.entries.Swap(key, value)
	pre, _ = previous.(V)
	return
}

func (m *SyncMap[K, V]) Range(iterator func(key K, value V) bool) {
	m.entries.Range(func(key, value any) bool {
		k, ok := key.(K)
		if !ok {
			return false
		}
		v, ok := value.(V)
		if !ok {
			return false
		}
		return iterator(k, v)
	})
}

func (m *SyncMap[K, V]) LoadOrStore(key K, new V) (actual V, loaded bool) {
	v, loaded := m.entries.LoadOrStore(key, new)
	actual, _ = v.(V)
	return
}

func (m *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.entries.LoadAndDelete(key)
	value, _ = v.(V)
	return
}

func (m *SyncMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	return m.entries.CompareAndSwap(key, old, new)
}

func (m *SyncMap[K, V]) CompareAndDelete(key K, old V) bool {
	return m.entries.CompareAndDelete(key, old)
}

func (m *SyncMap[K, V]) CompareFnAndSwap(key K, fn func(V, V) bool, old, new V) bool {
	value, ok := m.entries.Load(key)
	if !ok {
		return false
	}
	if !fn(value.(V), old) {
		return false
	}
	m.entries.Store(key, new)
	return true
}

func (m *SyncMap[K, V]) CompareFnAndDelete(key K, fn func(V, V) bool, old V) bool {
	value, ok := m.entries.Load(key)
	if !ok {
		return false
	}
	if !fn(value.(V), old) {
		return false
	}
	m.entries.Delete(key)
	return true
}
