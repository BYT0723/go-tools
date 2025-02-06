package ds

import (
	"fmt"
)

type HashSet[T comparable] map[T]struct{}

var _ Set[int] = HashSet[int](nil)

// Create a new Set with element type T
func NewHashSet[T comparable](items ...T) HashSet[T] {
	res := make(HashSet[T])
	res.Append(items...)
	return res
}

func (s HashSet[T]) Len() int {
	return len(s)
}

func (s HashSet[T]) Append(items ...T) {
	for _, v := range items {
		s[v] = struct{}{}
	}
}

func (s HashSet[T]) Values() (values []T) {
	values = make([]T, 0, len(s))
	for v := range s {
		values = append(values, v)
	}
	return
}

func (s HashSet[T]) Remove(values ...T) bool {
	for _, v := range values {
		delete(s, v)
	}
	return true
}

func (s HashSet[T]) Contains(value T) (exist bool) {
	_, exist = s[value]
	return
}

func (s HashSet[T]) String() string {
	return fmt.Sprint(s.Values())
}

// Union 返回两个集合的并集
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

// Intersection 返回两个集合的交集
func (s HashSet[T]) Intersection(s1 Set[T]) Set[T] {
	result := make(HashSet[T])

	for _, v := range s1.Values() {
		if s.Contains(v) {
			result[v] = struct{}{}
		}
	}
	return result
}

// Difference 返回两个集合的对称差集
func (s HashSet[T]) Difference(s1 Set[T]) Set[T] {
	result := make(HashSet[T])

	for v := range s {
		if !s1.Contains(v) {
			result[v] = struct{}{}
		}
	}
	return result
}

// Difference 返回两个集合的差集
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
