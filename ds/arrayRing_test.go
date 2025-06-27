package ds

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArrayRingBasic(t *testing.T) {
	r := NewArrayRingWithSize[int](5)

	require.Equal(t, 0, r.Len())
	require.Equal(t, []int{}, r.Values())

	// 插入 3 个
	r.Push(1)
	r.Push(2)
	r.Push(3)

	require.Equal(t, 3, r.Len())
	require.Equal(t, []int{1, 2, 3}, r.Values())

	// 插入 2 个，刚好填满
	r.Push(4)
	r.Push(5)

	require.Equal(t, 5, r.Len())
	require.Equal(t, []int{1, 2, 3, 4, 5}, r.Values())

	// 插入 2 个，触发覆盖（环回）
	r.Push(6)
	r.Push(7)

	require.Equal(t, 5, r.Len())
	require.Equal(t, []int{3, 4, 5, 6, 7}, r.Values())
}

func TestArrayRingWrapAround(t *testing.T) {
	r := NewArrayRingWithSize[string](3)

	r.Push("a")
	r.Push("b")
	r.Push("c")

	require.Equal(t, []string{"a", "b", "c"}, r.Values())

	// 覆盖最旧的 "a"
	r.Push("d")

	require.Equal(t, []string{"b", "c", "d"}, r.Values())
}
