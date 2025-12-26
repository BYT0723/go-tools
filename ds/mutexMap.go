package ds

import (
	"maps"
	"slices"
	"sync"
)

var _ Map[int, int] = (*MutexMap[int, int])(nil)

type MutexMap[K comparable, V any] struct {
	l       sync.Mutex
	entries map[K]V
}

func NewMutexMap[K comparable, V any]() *MutexMap[K, V] {
	return &MutexMap[K, V]{
		entries: make(map[K]V),
	}
}

func (m *MutexMap[K, V]) Store(key K, value V) {
	m.l.Lock()
	defer m.l.Unlock()
	m.entries[key] = value
}

func (m *MutexMap[K, V]) Load(key K) (value V, ok bool) {
	m.l.Lock()
	defer m.l.Unlock()
	value, ok = m.entries[key]
	return
}

func (m *MutexMap[K, V]) Delete(key K) bool {
	m.l.Lock()
	defer m.l.Unlock()
	delete(m.entries, key)
	return true
}

func (m *MutexMap[K, V]) Swap(key K, value V) (pre V, loaded bool) {
	m.l.Lock()
	defer m.l.Unlock()
	pre, loaded = m.entries[key]
	m.entries[key] = value
	return
}

func (m *MutexMap[K, V]) Range(iterator func(key K, value V) bool) {
	m.l.Lock()
	defer m.l.Unlock()

	for k, v := range m.entries {
		if !iterator(k, v) {
			break
		}
	}
}

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

func (m *MutexMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	m.l.Lock()
	defer m.l.Unlock()

	value, loaded = m.entries[key]
	if loaded {
		delete(m.entries, key)
	}
	return
}

func (m *MutexMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	return m.CompareFnAndSwap(key, func(a, b V) bool { return any(a) == any(b) }, old, new)
}

func (m *MutexMap[K, V]) CompareAndDelete(key K, old V) bool {
	return m.CompareFnAndDelete(key, func(a, b V) bool { return any(a) == any(b) }, old)
}

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

func (m *MutexMap[K, V]) Keys() []K {
	m.l.Lock()
	defer m.l.Unlock()
	return slices.Collect(maps.Keys(m.entries))
}

func (m *MutexMap[K, V]) Values() []V {
	m.l.Lock()
	defer m.l.Unlock()
	return slices.Collect(maps.Values(m.entries))
}

func (m *MutexMap[K, V]) Len() int {
	m.l.Lock()
	defer m.l.Unlock()
	return len(m.entries)
}

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
