package hashset

import (
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSetLen(t *testing.T) {
	Convey("HashSet Len", t, func() {
		Convey("Empty Set", func() {
			So(NewHashSet[int]().Len(), ShouldEqual, 0)
		})
		Convey("Non-Empty Set", func() {
			So(NewHashSet(1, 2, 3, 4).Len(), ShouldEqual, 4)
		})
	})
}

func TestSetAppend(t *testing.T) {
	Convey("HashSet append", t, func() {
		Convey("Empty Set Append", func() {
			s := NewHashSet[int]()
			s.Append(1, 2, 3, 4)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldEqual, []int{1, 2, 3, 4})
		})
		Convey("Append Self and More", func() {
			s := NewHashSet(1, 2)
			s.Append(1, 2, 3, 4)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldEqual, []int{1, 2, 3, 4})
		})
		Convey("Append Sub", func() {
			s := NewHashSet(1, 2, 3, 4)
			s.Append(1, 2)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldEqual, []int{1, 2, 3, 4})
		})
		Convey("Append Self", func() {
			s := NewHashSet(1, 2, 3, 4)
			s.Append(1, 2, 3, 4)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldEqual, []int{1, 2, 3, 4})
		})
		Convey("Append More", func() {
			s := NewHashSet(1, 2, 3, 4)
			s.Append(5, 6, 7, 8)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldEqual, []int{1, 2, 3, 4, 5, 6, 7, 8})
		})
	})
}

func TestSetValues(t *testing.T) {
	Convey("HashSet Values", t, func() {
		Convey("Empty Set", func() {
			So(NewHashSet[int]().Values(), ShouldBeEmpty)
		})
		Convey("Non-Empty Set", func() {
			vs := NewHashSet(1, 2, 3, 4).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
	})
}

func TestSetRemove(t *testing.T) {
	Convey("HashSet Remove", t, func() {
		Convey("Empty Set", func() {
			s := NewHashSet[int]()
			s.Remove(10)
			So(s.Values(), ShouldBeEmpty)
		})
		Convey("Non-Empty Set remove exist", func() {
			s := NewHashSet(1, 2, 3, 4)
			s.Remove(3)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 4})
		})
		Convey("Non-Empty Set remove not exist", func() {
			s := NewHashSet(1, 2, 3, 4)
			s.Remove(5)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
	})
}

func TestSetContains(t *testing.T) {
	Convey("HashSet Contains", t, func() {
		Convey("Empty Set", func() {
			So(NewHashSet[int]().Contains(10), ShouldBeFalse)
		})
		Convey("Non-Empty Set", func() {
			So(NewHashSet(1, 2, 3, 4).Contains(2), ShouldBeTrue)
			So(NewHashSet(1, 2, 3, 4).Contains(10), ShouldBeFalse)
		})
	})
}

func TestSetUnion(t *testing.T) {
	Convey("HashSet Union", t, func() {
		Convey("Two Empty Set", func() {
			So(NewHashSet[int]().Union(NewHashSet[int]()).Values(), ShouldBeEmpty)
		})
		Convey("Union of non-empty and empty sets", func() {
			s1 := NewHashSet[int]()
			s2 := NewHashSet(1, 2, 3, 4)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Union of non-empty and empty sets 2", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet[int]()
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Union of two sets without intersection", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(6, 7, 8, 9)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4, 6, 7, 8, 9})
		})
		Convey("Union of two sets with intersection", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(3, 4, 5, 6)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4, 5, 6})
		})
		Convey("Union of two identical sets", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(1, 2, 3, 4)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Union of two sets without difference", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(2, 3)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
	})
}

func TestSetIntersection(t *testing.T) {
	Convey("HashSet Intersection", t, func() {
		Convey("Two Empty Set", func() {
			So(NewHashSet[int]().Intersection(NewHashSet[int]()).Values(), ShouldBeEmpty)
		})
		Convey("Union of non-empty and empty sets", func() {
			s1 := NewHashSet[int]()
			s2 := NewHashSet(1, 2, 3, 4)
			So(s1.Intersection(s2), ShouldBeEmpty)
		})
		Convey("Union of non-empty and empty sets 2", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet[int]()
			So(s1.Intersection(s2).Values(), ShouldBeEmpty)
		})
		Convey("Union of two sets without intersection", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(6, 7, 8, 9)
			So(s1.Intersection(s2).Values(), ShouldBeEmpty)
		})
		Convey("Union of two sets with intersection", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(3, 4, 5, 6)
			vs := s1.Intersection(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{3, 4})
		})
		Convey("Union of two identical sets", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(1, 2, 3, 4)
			vs := s1.Intersection(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Union of two sets without difference", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(2, 3)
			vs := s1.Intersection(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{2, 3})
		})
	})
}

func TestSetDifference(t *testing.T) {
	Convey("HashSet Difference", t, func() {
		Convey("Two Empty Set", func() {
			So(NewHashSet[int]().Difference(NewHashSet[int]()).Values(), ShouldBeEmpty)
		})
		Convey("Union of non-empty and empty sets", func() {
			s1 := NewHashSet[int]()
			s2 := NewHashSet(1, 2, 3, 4)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldBeEmpty)
		})
		Convey("Union of non-empty and empty sets 2", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet[int]()
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Union of two sets without intersection", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(6, 7, 8, 9)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Union of two sets with intersection", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(3, 4, 5, 6)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2})
		})
		Convey("Union of two sets with intersection 2", func() {
			s1 := NewHashSet(3, 4, 5, 6)
			s2 := NewHashSet(1, 2, 3, 4)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{5, 6})
		})
		Convey("Union of two identical sets", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(1, 2, 3, 4)
			So(s1.Difference(s2).Values(), ShouldBeEmpty)
		})
		Convey("Union of two sets without difference", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(2, 3)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 4})
		})
		Convey("Union of two sets without difference 2", func() {
			s1 := NewHashSet(2, 3)
			s2 := NewHashSet(1, 2, 3, 4)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldBeEmpty)
		})
	})
}

func TestSetSymmetricDifference(t *testing.T) {
	Convey("HashSet SymmetricDifference", t, func() {
		Convey("Two Empty Set", func() {
			So(NewHashSet[int]().SymmetricDifference(NewHashSet[int]()).Values(), ShouldBeEmpty)
		})
		Convey("Union of non-empty and empty sets", func() {
			s1 := NewHashSet[int]()
			s2 := NewHashSet(1, 2, 3, 4)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Union of non-empty and empty sets 2", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet[int]()
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Union of two sets without intersection", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(6, 7, 8, 9)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4, 6, 7, 8, 9})
		})
		Convey("Union of two sets with intersection", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(3, 4, 5, 6)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 5, 6})
		})
		Convey("Union of two sets with intersection 2", func() {
			s1 := NewHashSet(3, 4, 5, 6)
			s2 := NewHashSet(1, 2, 3, 4)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 5, 6})
		})
		Convey("Union of two identical sets", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(1, 2, 3, 4)
			So(s1.SymmetricDifference(s2).Values(), ShouldBeEmpty)
		})
		Convey("Union of two sets without difference", func() {
			s1 := NewHashSet(1, 2, 3, 4)
			s2 := NewHashSet(2, 3)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 4})
		})
		Convey("Union of two sets without difference 2", func() {
			s1 := NewHashSet(2, 3)
			s2 := NewHashSet(1, 2, 3, 4)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 4})
		})
	})
}

func TestSetString(t *testing.T) {
	Convey("Set String", t, func() {
		Convey("empty collection", func() {
			So(NewHashSet[int]().String(), ShouldEqual, "[]")
		})
		Convey("non-empty collection", func() {
			So(NewHashSet(1).String(), ShouldEqual, "[1]")
		})
	})
}
