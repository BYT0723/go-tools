package ds

import "fmt"

// Set is a generic interface for collections of unique elements.
// It provides standard set operations like union, intersection, difference, etc.
//
// Type parameters:
//   - T: The element type, must be comparable
type Set[T comparable] interface {
	fmt.Stringer
	// Len returns the number of elements in the set.
	Len() int
	// Append adds one or more elements to the set.
	// Duplicate elements are ignored.
	Append(...T)
	// Remove removes one or more elements from the set.
	// Non-existent elements are ignored.
	Remove(...T) bool
	// Contains checks if an element exists in the set.
	Contains(T) bool
	// Values returns a slice containing all elements in the set.
	Values() []T
	// Union returns a new set containing all elements from both sets.
	Union(Set[T]) Set[T]
	// Intersection returns a new set containing elements present in both sets.
	Intersection(Set[T]) Set[T]
	// Difference returns a new set containing elements in this set but not in the other.
	Difference(Set[T]) Set[T]
	// SymmetricDifference returns a new set containing elements present in either set but not both.
	SymmetricDifference(Set[T]) Set[T]
}
