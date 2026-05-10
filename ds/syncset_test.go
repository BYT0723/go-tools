package ds

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncSetLen(t *testing.T) {
	t.Run("HashSet Len", func(t *testing.T) {
		t.Run("Empty Set", func(t *testing.T) {
			assert.Equal(t, 0, NewSyncSet[int]().Len())
		})
		t.Run("Non-Empty Set", func(t *testing.T) {
			assert.Equal(t, 4, NewSyncSet(1, 2, 3, 4).Len())
		})
	})
}

func TestSyncSetAppend(t *testing.T) {
	t.Run("HashSet append", func(t *testing.T) {
		t.Run("Empty Set Append", func(t *testing.T) {
			s := NewSyncSet[int]()
			s.Append(1, 2, 3, 4)
			vs := s.Values()
			sort.Ints(vs)
			assert.Equal(t, []int{1, 2, 3, 4}, vs)
		})
		t.Run("Append Self and More", func(t *testing.T) {
			s := NewSyncSet(1, 2)
			s.Append(1, 2, 3, 4)
			vs := s.Values()
			sort.Ints(vs)
			assert.Equal(t, []int{1, 2, 3, 4}, vs)
		})
		t.Run("Append Sub", func(t *testing.T) {
			s := NewSyncSet(1, 2, 3, 4)
			s.Append(1, 2)
			vs := s.Values()
			sort.Ints(vs)
			assert.Equal(t, []int{1, 2, 3, 4}, vs)
		})
		t.Run("Append Self", func(t *testing.T) {
			s := NewSyncSet(1, 2, 3, 4)
			s.Append(1, 2, 3, 4)
			vs := s.Values()
			sort.Ints(vs)
			assert.Equal(t, []int{1, 2, 3, 4}, vs)
		})
		t.Run("Append More", func(t *testing.T) {
			s := NewSyncSet(1, 2, 3, 4)
			s.Append(5, 6, 7, 8)
			vs := s.Values()
			sort.Ints(vs)
			assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, vs)
		})
	})
}

func TestSyncSetValues(t *testing.T) {
	t.Run("HashSet Values", func(t *testing.T) {
		t.Run("Empty Set", func(t *testing.T) {
			assert.Empty(t, NewSyncSet[int]().Values())
		})
		t.Run("Non-Empty Set", func(t *testing.T) {
			vs := NewSyncSet(1, 2, 3, 4).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4}, vs)
		})
	})
}

func TestSyncSetRemove(t *testing.T) {
	t.Run("HashSet Remove", func(t *testing.T) {
		t.Run("Empty Set", func(t *testing.T) {
			s := NewSyncSet[int]()
			s.Remove(10)
			assert.Empty(t, s.Values())
		})
		t.Run("Non-Empty Set remove exist", func(t *testing.T) {
			s := NewSyncSet(1, 2, 3, 4)
			s.Remove(3)
			vs := s.Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 4}, vs)
		})
		t.Run("Non-Empty Set remove not exist", func(t *testing.T) {
			s := NewSyncSet(1, 2, 3, 4)
			s.Remove(5)
			vs := s.Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4}, vs)
		})
	})
}

func TestSyncSetContains(t *testing.T) {
	t.Run("HashSet Contains", func(t *testing.T) {
		t.Run("Empty Set", func(t *testing.T) {
			assert.False(t, NewSyncSet[int]().Contains(10))
		})
		t.Run("Non-Empty Set", func(t *testing.T) {
			assert.True(t, NewSyncSet(1, 2, 3, 4).Contains(2))
			assert.False(t, NewSyncSet(1, 2, 3, 4).Contains(10))
		})
	})
}

func TestSyncSetUnion(t *testing.T) {
	t.Run("HashSet Union", func(t *testing.T) {
		t.Run("Two Empty Set", func(t *testing.T) {
			assert.Empty(t, NewSyncSet[int]().Union(NewSyncSet[int]()).Values())
		})
		t.Run("Union of non-empty and empty sets", func(t *testing.T) {
			s1 := NewSyncSet[int]()
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4}, vs)
		})
		t.Run("Union of non-empty and empty sets 2", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet[int]()
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4}, vs)
		})
		t.Run("Union of two sets without intersection", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(6, 7, 8, 9)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4, 6, 7, 8, 9}, vs)
		})
		t.Run("Union of two sets with intersection", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(3, 4, 5, 6)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4, 5, 6}, vs)
		})
		t.Run("Union of two identical sets", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4}, vs)
		})
		t.Run("Union of two sets without difference", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(2, 3)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4}, vs)
		})
	})
}

func TestSyncSetIntersection(t *testing.T) {
	t.Run("HashSet Intersection", func(t *testing.T) {
		t.Run("Two Empty Set", func(t *testing.T) {
			assert.Empty(t, NewSyncSet[int]().Intersection(NewSyncSet[int]()).Values())
		})
		t.Run("Intersection of non-empty and empty sets", func(t *testing.T) {
			s1 := NewSyncSet[int]()
			s2 := NewSyncSet(1, 2, 3, 4)
			assert.Empty(t, s1.Intersection(s2).Values())
		})
		t.Run("Intersection of non-empty and empty sets 2", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet[int]()
			assert.Empty(t, s1.Intersection(s2).Values())
		})
		t.Run("Intersection of two sets without intersection", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(6, 7, 8, 9)
			assert.Empty(t, s1.Intersection(s2).Values())
		})
		t.Run("Intersection of two sets with intersection", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(3, 4, 5, 6)
			vs := s1.Intersection(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{3, 4}, vs)
		})
		t.Run("Intersection of two identical sets", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.Intersection(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4}, vs)
		})
		t.Run("Intersection of two sets without difference", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(2, 3)
			vs := s1.Intersection(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{2, 3}, vs)
		})
	})
}

func TestSyncSetDifference(t *testing.T) {
	t.Run("HashSet Difference", func(t *testing.T) {
		t.Run("Two Empty Set", func(t *testing.T) {
			assert.Empty(t, NewSyncSet[int]().Difference(NewSyncSet[int]()).Values())
		})
		t.Run("Difference of non-empty and empty sets", func(t *testing.T) {
			s1 := NewSyncSet[int]()
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			assert.Empty(t, vs)
		})
		t.Run("Difference of non-empty and empty sets 2", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet[int]()
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4}, vs)
		})
		t.Run("Difference of two sets without intersection", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(6, 7, 8, 9)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4}, vs)
		})
		t.Run("Difference of two sets with intersection", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(3, 4, 5, 6)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2}, vs)
		})
		t.Run("Difference of two sets with intersection 2", func(t *testing.T) {
			s1 := NewSyncSet(3, 4, 5, 6)
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{5, 6}, vs)
		})
		t.Run("Difference of two identical sets", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(1, 2, 3, 4)
			assert.Empty(t, s1.Difference(s2).Values())
		})
		t.Run("Difference of two sets without difference", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(2, 3)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 4}, vs)
		})
		t.Run("Difference of two sets without difference 2", func(t *testing.T) {
			s1 := NewSyncSet(2, 3)
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			assert.Empty(t, vs)
		})
	})
}

func TestSyncSetSymmetricDifference(t *testing.T) {
	t.Run("HashSet SymmetricDifference", func(t *testing.T) {
		t.Run("Two Empty Set", func(t *testing.T) {
			assert.Empty(t, NewSyncSet[int]().SymmetricDifference(NewSyncSet[int]()).Values())
		})
		t.Run("SymmetricDifference of non-empty and empty sets", func(t *testing.T) {
			s1 := NewSyncSet[int]()
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4}, vs)
		})
		t.Run("SymmetricDifference of non-empty and empty sets 2", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet[int]()
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4}, vs)
		})
		t.Run("SymmetricDifference of two sets without intersection", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(6, 7, 8, 9)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 3, 4, 6, 7, 8, 9}, vs)
		})
		t.Run("SymmetricDifference of two sets with intersection", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(3, 4, 5, 6)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 5, 6}, vs)
		})
		t.Run("SymmetricDifference of two sets with intersection 2", func(t *testing.T) {
			s1 := NewSyncSet(3, 4, 5, 6)
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 2, 5, 6}, vs)
		})
		t.Run("SymmetricDifference of two identical sets", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(1, 2, 3, 4)
			assert.Empty(t, s1.SymmetricDifference(s2).Values())
		})
		t.Run("SymmetricDifference of two sets without difference", func(t *testing.T) {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(2, 3)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 4}, vs)
		})
		t.Run("SymmetricDifference of two sets without difference 2", func(t *testing.T) {
			s1 := NewSyncSet(2, 3)
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			assert.EqualValues(t, []int{1, 4}, vs)
		})
	})
}

func TestSyncSetString(t *testing.T) {
	t.Run("Set String", func(t *testing.T) {
		t.Run("empty collection", func(t *testing.T) {
			assert.Equal(t, "[]", NewSyncSet[int]().String())
		})
		t.Run("non-empty collection", func(t *testing.T) {
			assert.Equal(t, "[1]", NewSyncSet(1).String())
		})
	})
}
