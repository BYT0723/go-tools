package set

type Set[T comparable] interface {
	// 集合长度
	Len() int
	// 添加元素
	Append(...T)
	// 移除元素
	Remove(T) bool
	// 判断元素是否存在
	Contains(T) bool
	// 集合元素的切片
	Values() []T
	// 并集
	Union(Set[T]) Set[T]
	// 交集
	Intersection(Set[T]) Set[T]
	// 差集
	Difference(Set[T]) Set[T]
	// 对称差集
	SymmetricDifference(Set[T]) Set[T]
}
