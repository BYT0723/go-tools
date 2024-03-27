package set

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Collection, T is interact{} type
type Set[T any] struct {
	n          int32
	item       map[string]T   // Element storage map, key is the unique identifier of the element, value is the element
	identifier func(T) string // Generate unique identification function for elements，
	rwMux      sync.RWMutex
}

// Create a new Set with element type T
func NewSet[T any](items ...T) *Set[T] {
	res := &Set[T]{
		item:       map[string]T{},
		identifier: func(t T) string { return fmt.Sprint(t) },
	}
	res.Append(items...)
	return res
}

// Create a new Set with element type T
// identifier, Generate unique identification function for elements, used for custom deduplication
func NewSetFunc[T any](identifier func(T) string, items ...T) *Set[T] {
	return &Set[T]{
		item:       map[string]T{},
		identifier: identifier,
	}
}

func (s *Set[T]) Length() int {
	return int(s.n)
}

func (s *Set[T]) Append(items ...T) {
	s.rwMux.Lock()
	defer s.rwMux.Unlock()
	for _, v := range items {
		if _, ok := s.item[s.identifier(v)]; !ok {
			atomic.AddInt32(&(s.n), 1)
		}
		s.item[s.identifier(v)] = v
	}
}

func (s *Set[T]) Values() (values []T) {
	s.rwMux.RLock()
	defer s.rwMux.RUnlock()
	values = []T{}
	for _, v := range s.item {
		values = append(values, v)
	}
	return
}

func (s *Set[T]) Remove(value T) {
	s.rwMux.Lock()
	defer s.rwMux.Unlock()
	if _, ok := s.item[s.identifier(value)]; ok {
		delete(s.item, s.identifier(value))
		atomic.AddInt32(&(s.n), -1)
	}
}

func (s *Set[T]) Contains(value T) (exist bool) {
	s.rwMux.RLock()
	defer s.rwMux.RUnlock()

	v2, exist := s.item[s.identifier(value)]
	fmt.Printf("v1: %v\tv2: %v\n", value, v2)
	return
}

// Union 返回两个集合的并集
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	unionSet := NewSetFunc(s.identifier)

	s.rwMux.RLock()
	other.rwMux.RLock()
	defer s.rwMux.RUnlock()
	defer other.rwMux.RUnlock()

	for _, v := range s.item {
		unionSet.Append(v)
	}

	for _, v := range other.item {
		unionSet.Append(v)
	}

	return unionSet
}

// Intersection 返回两个集合的交集
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	intersectionSet := NewSetFunc(s.identifier)
	s.rwMux.RLock()
	other.rwMux.RLock()
	defer s.rwMux.RUnlock()
	defer other.rwMux.RUnlock()

	for _, v := range s.item {
		if other.Contains(v) {
			intersectionSet.Append(v)
		}
	}

	return intersectionSet
}

// Difference 返回两个集合的差集
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	differenceSet := NewSetFunc(s.identifier)
	s.rwMux.RLock()
	other.rwMux.RLock()
	defer s.rwMux.RUnlock()
	defer other.rwMux.RUnlock()

	for _, v := range s.item {
		if !other.Contains(v) {
			differenceSet.Append(v)
		}
	}
	return differenceSet
}

func (s *Set[T]) Equal(other *Set[T]) bool {
	if s.Length() != other.Length() {
		return false
	}
	s.rwMux.RLock()
	other.rwMux.RLock()
	defer s.rwMux.RUnlock()
	defer other.rwMux.RUnlock()

	for _, v := range s.item {
		if !other.Contains(v) {
			return false
		}
	}
	return true
}
