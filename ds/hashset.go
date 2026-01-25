package ds

import (
	"fmt"
)

// HashSet is a simple set implementation using Go's built-in map.
// It provides O(1) average time complexity for basic operations.
//
// Type parameters:
//   - T: The element type, must be comparable
type HashSet[T comparable] map[T]struct{}

var _ Set[int] = HashSet[int](nil)

// NewHashSet creates a new HashSet with optional initial elements.
//
// Parameters:
//   - items: Optional initial elements to add to the set
//
// Returns:
//   - HashSet[T]: A new HashSet instance
func NewHashSet[T comparable](items ...T) HashSet[T] {
	res := make(HashSet[T])
	res.Append(items...)
	return res
}

// Len returns the number of elements in the set.
func (s HashSet[T]) Len() int {
	return len(s)
}

// Append adds one or more elements to the set.
// Duplicate elements are ignored.
func (s HashSet[T]) Append(items ...T) {
	for _, v := range items {
		s[v] = struct{}{}
	}
}

// Values returns a slice containing all elements in the set.
func (s HashSet[T]) Values() (values []T) {
	values = make([]T, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return
}

// Remove removes one or more elements from the set.
// Non-existent elements are ignored.
func (s HashSet[T]) Remove(values ...T) bool {
	for _, v := range values {
		delete(s, v)
	}
	return true
}

// Contains checks if an element exists in the set.
func (s HashSet[T]) Contains(value T) (exist bool) {
	_, exist = s[value]
	return
}

// String returns a string representation of the set.
func (s HashSet[T]) String() string {
	return fmt.Sprint(s.Values())
}

// Union returns a new set containing all elements from both sets.
func (s HashSet[T]) Union(s1 Set[T]) Set[T] {
	result := make(HashSet[T])

	for v := range s {
		result[v] = struct{}{}
	}
	for _, v := range s1.Values() {
		result[v] = struct{}{}
	}
	return result
}

// Intersection returns a new set containing elements present in both sets.
func (s HashSet[T]) Intersection(s1 Set[T]) Set[T] {
	result := make(HashSet[T])

	for _, v := range s1.Values() {
		if s.Contains(v) {
			result[v] = struct{}{}
		}
	}
	return result
}

// Difference returns a new set containing elements in this set but not in the other.
func (s HashSet[T]) Difference(s1 Set[T]) Set[T] {
	result := make(HashSet[T])

	for v := range s {
		if !s1.Contains(v) {
			result[v] = struct{}{}
		}
	}
	return result
}

// SymmetricDifference returns a new set containing elements present in either set but not both.
func (s HashSet[T]) SymmetricDifference(s1 Set[T]) Set[T] {
	result := make(HashSet[T])

	for v := range s {
		if !s1.Contains(v) {
			result[v] = struct{}{}
		}
	}

	for _, v := range s1.Values() {
		if !s.Contains(v) {
			result[v] = struct{}{}
		}
	}
	return result
}
