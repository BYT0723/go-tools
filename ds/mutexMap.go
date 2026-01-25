package ds

import (
	"maps"
	"slices"
	"sync"
)

var _ Map[int, int] = (*MutexMap[int, int])(nil)

// MutexMap is a thread-safe map implementation using a mutex for synchronization.
// It provides a simpler alternative to SyncMap with mutex-based concurrency control.
type MutexMap[K comparable, V any] struct {
	l       sync.Mutex
	entries map[K]V
}

// NewMutexMap creates a new MutexMap instance.
func NewMutexMap[K comparable, V any]() *MutexMap[K, V] {
	return &MutexMap[K, V]{
		entries: make(map[K]V),
	}
}

// Store stores a value for a key.
func (m *MutexMap[K, V]) Store(key K, value V) {
	m.l.Lock()
	defer m.l.Unlock()
	m.entries[key] = value
}

// Load returns the value stored for a key, or false if no value is present.
func (m *MutexMap[K, V]) Load(key K) (value V, ok bool) {
	m.l.Lock()
	defer m.l.Unlock()
	value, ok = m.entries[key]
	return
}

// Delete deletes the value for a key.
func (m *MutexMap[K, V]) Delete(key K) bool {
	m.l.Lock()
	defer m.l.Unlock()
	delete(m.entries, key)
	return true
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (m *MutexMap[K, V]) Swap(key K, value V) (pre V, loaded bool) {
	m.l.Lock()
	defer m.l.Unlock()
	pre, loaded = m.entries[key]
	m.entries[key] = value
	return
}

// Range calls iterator sequentially for each key and value present in the map.
// If iterator returns false, range stops the iteration.
func (m *MutexMap[K, V]) Range(iterator func(key K, value V) bool) {
	m.l.Lock()
	defer m.l.Unlock()

	for k, v := range m.entries {
		if !iterator(k, v) {
			break
		}
	}
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *MutexMap[K, V]) LoadOrStore(key K, new V) (actual V, loaded bool) {
	m.l.Lock()
	defer m.l.Unlock()

	actual, loaded = m.entries[key]
	if !loaded {
		m.entries[key] = new
		actual = new
	}
	return
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *MutexMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	m.l.Lock()
	defer m.l.Unlock()

	value, loaded = m.entries[key]
	if loaded {
		delete(m.entries, key)
	}
	return
}

// CompareAndSwap swaps the value for a key if the current value equals old.
// Returns true if the swap was performed.
func (m *MutexMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	return m.CompareFnAndSwap(key, func(a, b V) bool { return any(a) == any(b) }, old, new)
}

// CompareAndDelete deletes the entry for a key if its value equals old.
// Returns true if the entry was deleted.
func (m *MutexMap[K, V]) CompareAndDelete(key K, old V) bool {
	return m.CompareFnAndDelete(key, func(a, b V) bool { return any(a) == any(b) }, old)
}

// CompareFnAndSwap swaps the value for a key using a custom comparison function.
// The function fn is called with the current value and old value.
// Returns true if fn returns true and the swap was performed.
func (m *MutexMap[K, V]) CompareFnAndSwap(key K, fn func(V, V) bool, old, new V) bool {
	m.l.Lock()
	defer m.l.Unlock()

	v, ok := m.entries[key]
	if !ok {
		return false
	}
	if !fn(v, old) {
		return false
	}
	m.entries[key] = new
	return true
}

// CompareFnAndDelete deletes the entry for a key using a custom comparison function.
// The function fn is called with the current value and old value.
// Returns true if fn returns true and the entry was deleted.
func (m *MutexMap[K, V]) CompareFnAndDelete(key K, fn func(V, V) bool, old V) bool {
	m.l.Lock()
	defer m.l.Unlock()

	v, ok := m.entries[key]
	if !ok {
		return false
	}
	if !fn(v, old) {
		return false
	}
	delete(m.entries, key)
	return true
}

// Keys returns a slice containing all keys in the map.
func (m *MutexMap[K, V]) Keys() []K {
	m.l.Lock()
	defer m.l.Unlock()
	return slices.Collect(maps.Keys(m.entries))
}

// Values returns a slice containing all values in the map.
func (m *MutexMap[K, V]) Values() []V {
	m.l.Lock()
	defer m.l.Unlock()
	return slices.Collect(maps.Values(m.entries))
}

// Len returns the number of entries in the map.
func (m *MutexMap[K, V]) Len() int {
	m.l.Lock()
	defer m.l.Unlock()
	return len(m.entries)
}

// Filter returns a new Map containing only entries that satisfy the filter function.
func (m *MutexMap[K, V]) Filter(filter func(K, V) bool) Map[K, V] {
	m.l.Lock()
	defer m.l.Unlock()

	result := NewMutexMap[K, V]()
	for k, v := range m.entries {
		if filter(k, v) {
			result.Store(k, v)
		}
	}
	return result
}
