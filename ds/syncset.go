package ds

import (
	"fmt"
	"maps"
	"slices"
	"sync"
)

type (
	// SyncSet is a thread-safe set implementation for comparable types.
	SyncSet[T comparable] struct {
		mutex   sync.RWMutex
		entries map[T]struct{}
	}
)

// Ensure SyncSet implements Set[int] interface (assumed to exist).
var _ Set[int] = (*SyncSet[int])(nil)

// NewSyncSet creates a new SyncSet and optionally adds initial items.
func NewSyncSet[T comparable](items ...T) *SyncSet[T] {
	result := &SyncSet[T]{
		entries: make(map[T]struct{}),
	}
	result.Append(items...)
	return result
}

// Len returns the number of elements in the set.
func (s *SyncSet[T]) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.entries)
}

// String returns the string representation of the set's elements.
func (s *SyncSet[T]) String() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return fmt.Sprint(s.Values())
}

// Append adds one or more elements to the set. Duplicates are ignored.
func (s *SyncSet[T]) Append(values ...T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.entries == nil {
		s.entries = make(map[T]struct{})
	}
	for _, v := range values {
		s.entries[v] = struct{}{}
	}
}

// Remove deletes one or more elements from the set. Non-existent elements are ignored.
func (s *SyncSet[T]) Remove(values ...T) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, v := range values {
		delete(s.entries, v)
	}
	return true
}

// Contains reports whether the set contains the specified element.
func (s *SyncSet[T]) Contains(v T) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, ok := s.entries[v]
	return ok
}

// Values returns a slice containing all elements in the set (unordered).
func (s *SyncSet[T]) Values() []T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return slices.Collect(maps.Keys(s.entries))
}

// Union returns a new set containing all elements from the current set and another set.
func (s *SyncSet[T]) Union(s1 Set[T]) Set[T] {
	result := NewSyncSet[T](s.Values()...)
	result.Append(s1.Values()...)
	return result
}

// Intersection returns a new set containing elements present in both sets.
func (s *SyncSet[T]) Intersection(s1 Set[T]) Set[T] {
	result := NewSyncSet[T]()
	for _, v := range s.Values() {
		if s1.Contains(v) {
			result.Append(v)
		}
	}
	return result
}

// Difference returns a new set containing elements in the current set but not in the other.
func (s *SyncSet[T]) Difference(s1 Set[T]) Set[T] {
	result := NewSyncSet[T]()
	for _, v := range s.Values() {
		if !s1.Contains(v) {
			result.Append(v)
		}
	}
	return result
}

// SymmetricDifference returns a new set containing elements present in either of the sets but not both.
func (s *SyncSet[T]) SymmetricDifference(s1 Set[T]) Set[T] {
	result := NewSyncSet[T]()
	for _, v := range s.Values() {
		if !s1.Contains(v) {
			result.Append(v)
		}
	}
	for _, v := range s1.Values() {
		if !s.Contains(v) {
			result.Append(v)
		}
	}
	return result
}
