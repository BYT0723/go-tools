package ds

import "sync"

var _ Map[int, int] = (*SyncMap[int, int])(nil)

// SyncMap Generic wrapper for sync.Map
type SyncMap[K comparable, V any] struct {
	entries sync.Map
}

// NewSyncMap creates a new SyncMap instance.
func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{}
}

// Store stores a value for a key.
func (m *SyncMap[K, V]) Store(key K, value V) {
	m.entries.Store(key, value)
}

// Load returns the value stored for a key, or false if no value is present.
func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.entries.Load(key)
	if !ok {
		return
	}
	value, ok = v.(V)
	return
}

// Delete deletes the value for a key.
func (m *SyncMap[K, V]) Delete(key K) bool {
	m.entries.Delete(key)
	return true
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (m *SyncMap[K, V]) Swap(key K, value V) (pre V, loaded bool) {
	previous, loaded := m.entries.Swap(key, value)
	pre, _ = previous.(V)
	return
}

// Range calls iterator sequentially for each key and value present in the map.
// If iterator returns false, range stops the iteration.
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

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *SyncMap[K, V]) LoadOrStore(key K, new V) (actual V, loaded bool) {
	v, loaded := m.entries.LoadOrStore(key, new)
	actual, _ = v.(V)
	return
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.entries.LoadAndDelete(key)
	value, _ = v.(V)
	return
}

// CompareAndSwap swaps the value for a key if the current value equals old.
// Returns true if the swap was performed.
func (m *SyncMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	return m.entries.CompareAndSwap(key, old, new)
}

// CompareAndDelete deletes the entry for a key if its value equals old.
// Returns true if the entry was deleted.
func (m *SyncMap[K, V]) CompareAndDelete(key K, old V) bool {
	return m.entries.CompareAndDelete(key, old)
}

// CompareFnAndSwap swaps the value for a key using a custom comparison function.
// The function fn is called with the current value and old value.
// Returns true if fn returns true and the swap was performed.
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

// CompareFnAndDelete deletes the entry for a key using a custom comparison function.
// The function fn is called with the current value and old value.
// Returns true if fn returns true and the entry was deleted.
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

// Keys returns a slice containing all keys in the map.
func (m *SyncMap[K, V]) Keys() []K {
	var keys []K
	m.entries.Range(func(key, _ any) bool {
		if k, ok := key.(K); ok {
			keys = append(keys, k)
		}
		return true
	})
	return keys
}

// Values returns a slice containing all values in the map.
func (m *SyncMap[K, V]) Values() []V {
	var values []V
	m.entries.Range(func(_, value any) bool {
		if v, ok := value.(V); ok {
			values = append(values, v)
		}
		return true
	})
	return values
}

// Filter returns a new Map containing only entries that satisfy the filter function.
func (m *SyncMap[K, V]) Filter(filter func(K, V) bool) Map[K, V] {
	result := NewSyncMap[K, V]()

	m.entries.Range(func(key, value any) bool {
		k, _ := key.(K)
		v, _ := value.(V)
		if filter(k, v) {
			result.Store(k, v)
		}
		return true
	})
	return result
}
