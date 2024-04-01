package set

// Clone 克隆目标Set
func Clone[T any](src *Set[T]) *Set[T] {
	return NewSet(src.Values()...)
}

// Union 返回两个集合的并集
func Union[T any](s1, s2 *Set[T]) *Set[T] {
	unionSet := NewSetFunc(s1.identifier)

	s1.rwMux.RLock()
	for _, v := range s1.item {
		unionSet.Append(v)
	}
	s1.rwMux.RUnlock()

	s2.rwMux.RLock()
	for _, v := range s2.item {
		unionSet.Append(v)
	}
	s2.rwMux.RUnlock()

	return unionSet
}

// Union 返回两个集合的并集
func UnionFunc[T any](s1, s2 *Set[T], identifier func(T) string) *Set[T] {
	unionSet := NewSetFunc(identifier)

	s1.rwMux.RLock()
	for _, v := range s1.item {
		unionSet.Append(v)
	}
	s1.rwMux.RUnlock()

	s2.rwMux.RLock()
	for _, v := range s2.item {
		unionSet.Append(v)
	}
	s2.rwMux.RUnlock()

	return unionSet
}

// Intersection 返回两个集合的交集
func Intersection[T any](s1, s2 *Set[T]) *Set[T] {
	intersectionSet := NewSetFunc(s1.identifier)
	s1.rwMux.RLock()
	s2.rwMux.RLock()
	defer s1.rwMux.RUnlock()
	defer s2.rwMux.RUnlock()

	for _, v := range s1.item {
		if s2.Contains(v) {
			intersectionSet.Append(v)
		}
	}
	return intersectionSet
}

// Intersection 返回两个集合的交集
func IntersectionFunc[T any](s1, s2 *Set[T], identifier func(T) string) *Set[T] {
	intersectionSet := NewSetFunc(identifier)
	s1.rwMux.RLock()
	s2.rwMux.RLock()
	defer s1.rwMux.RUnlock()
	defer s2.rwMux.RUnlock()

	for _, v := range s1.item {
		if s2.Contains(v) {
			intersectionSet.Append(v)
		}
	}
	return intersectionSet
}

// Difference 返回两个集合的差集
func Difference[T any](s1, s2 *Set[T]) *Set[T] {
	differenceSet := NewSetFunc(s1.identifier)

	s1.rwMux.RLock()
	for _, v := range s1.item {
		if !s2.Contains(v) {
			differenceSet.Append(v)
		}
	}
	s1.rwMux.RUnlock()

	s2.rwMux.RLock()
	for _, v := range s2.item {
		if !s1.Contains(v) {
			differenceSet.Append(v)
		}
	}
	s2.rwMux.RUnlock()

	return differenceSet
}

// Difference 返回两个集合的差集
func DifferenceFunc[T any](s1, s2 *Set[T], identifier func(T) string) *Set[T] {
	differenceSet := NewSetFunc(identifier)

	s1.rwMux.RLock()
	for _, v := range s1.item {
		if !s2.Contains(v) {
			differenceSet.Append(v)
		}
	}
	s1.rwMux.RUnlock()

	s2.rwMux.RLock()
	for _, v := range s2.item {
		if !s1.Contains(v) {
			differenceSet.Append(v)
		}
	}
	s2.rwMux.RUnlock()

	return differenceSet
}
