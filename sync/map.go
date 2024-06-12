package sync

import "sync"

// Generic wrapper for sync.Map
type SyncMap[T1, T2 any] struct {
	entry sync.Map
}

func NewSyncMap[T1 comparable, T2 any](kvs map[T1]T2) *SyncMap[T1, T2] {
	res := new(SyncMap[T1, T2])

	for k, v := range kvs {
		res.Store(k, v)
	}

	return res
}

func (m *SyncMap[T1, T2]) Store(key T1, value T2) {
	m.entry.Store(key, value)
}

func (m *SyncMap[T1, T2]) Load(key T1) (value T2, ok bool) {
	v, ok := m.entry.Load(key)
	if !ok {
		return
	}
	value, ok = v.(T2)
	return
}

func (m *SyncMap[T1, T2]) Delete(key T1) {
	m.entry.Delete(key)
}

func (m *SyncMap[T1, T2]) Swap(key T1, value T2) (pre T2, loaded bool) {
	previous, loaded := m.entry.Swap(key, value)
	pre, _ = previous.(T2)
	return
}

func (m *SyncMap[T1, T2]) Range(iterator func(key T1, value T2) bool) {
	m.entry.Range(func(key, value any) bool {
		k, ok := key.(T1)
		if !ok {
			return false
		}
		v, ok := value.(T2)
		if !ok {
			return false
		}
		return iterator(k, v)
	})
}

func (m *SyncMap[T1, T2]) LoadOrStore(key T1, new T2) (actual T2, loaded bool) {
	v, loaded := m.entry.LoadOrStore(key, new)
	actual, _ = v.(T2)
	return
}

func (m *SyncMap[T1, T2]) LoadAndDelete(key T1) (value T2, loaded bool) {
	v, loaded := m.entry.LoadAndDelete(key)
	value, _ = v.(T2)
	return
}

func (m *SyncMap[T1, T2]) CompareAndSwap(key T1, old, new T2) bool {
	return m.entry.CompareAndSwap(key, old, new)
}

func (m *SyncMap[T1, T2]) CompareAndDelete(key T1, old T2) bool {
	return m.entry.CompareAndDelete(key, old)
}
