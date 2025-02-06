package ds

import (
	"fmt"
	"sync"
)

type (
	SyncSet[T comparable] struct {
		entries sync.Map
	}
)

var _ Set[int] = (*SyncSet[int])(nil)

// Create a new Set with element type T
func NewSyncSet[T comparable](items ...T) *SyncSet[T] {
	var result SyncSet[T]
	result.Append(items...)
	return &result
}

// 集合长度
func (s *SyncSet[T]) Len() int {
	var n int
	s.entries.Range(func(key, value any) bool {
		n++
		return true
	})
	return n
}

func (s *SyncSet[T]) String() string {
	return fmt.Sprint(s.Values())
}

// 添加元素
func (s *SyncSet[T]) Append(values ...T) {
	for _, v := range values {
		s.entries.Store(v, struct{}{})
	}
}

// 移除元素
func (s *SyncSet[T]) Remove(values ...T) bool {
	for _, v := range values {
		s.entries.Delete(v)
	}
	return true
}

// 判断元素是否存在
func (s *SyncSet[T]) Contains(v T) bool {
	_, ok := s.entries.Load(v)
	return ok
}

// 集合元素的切片
func (s *SyncSet[T]) Values() []T {
	var result []T
	s.entries.Range(func(key, _ any) bool {
		if v, ok := key.(T); ok {
			result = append(result, v)
		}
		return true
	})
	return result
}

// 并集
func (s *SyncSet[T]) Union(s1 Set[T]) Set[T] {
	var result SyncSet[T]
	for _, v := range s.Values() {
		result.entries.Store(v, struct{}{})
	}
	for _, v := range s1.Values() {
		result.entries.Store(v, struct{}{})
	}
	return &result
}

// 交集
func (s *SyncSet[T]) Intersection(s1 Set[T]) Set[T] {
	var result SyncSet[T]
	for _, v := range s.Values() {
		if s1.Contains(v) {
			result.entries.Store(v, struct{}{})
		}
	}
	return &result
}

// 差集
func (s *SyncSet[T]) Difference(s1 Set[T]) Set[T] {
	var result SyncSet[T]
	for _, v := range s.Values() {
		if !s1.Contains(v) {
			result.entries.Store(v, struct{}{})
		}
	}
	return &result
}

// 对称差集
func (s *SyncSet[T]) SymmetricDifference(s1 Set[T]) Set[T] {
	var result SyncSet[T]
	for _, v := range s.Values() {
		if !s1.Contains(v) {
			result.entries.Store(v, struct{}{})
		}
	}

	for _, v := range s1.Values() {
		if !s.Contains(v) {
			result.entries.Store(v, struct{}{})
		}
	}
	return &result
}
