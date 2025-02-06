package ds

import (
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSyncSetLen(t *testing.T) {
	Convey("HashSet Len", t, func() {
		Convey("Empty Set", func() {
			So(NewSyncSet[int]().Len(), ShouldEqual, 0)
		})
		Convey("Non-Empty Set", func() {
			So(NewSyncSet(1, 2, 3, 4).Len(), ShouldEqual, 4)
		})
	})
}

func TestSyncSetAppend(t *testing.T) {
	Convey("HashSet append", t, func() {
		Convey("Empty Set Append", func() {
			s := NewSyncSet[int]()
			s.Append(1, 2, 3, 4)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldEqual, []int{1, 2, 3, 4})
		})
		Convey("Append Self and More", func() {
			s := NewSyncSet(1, 2)
			s.Append(1, 2, 3, 4)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldEqual, []int{1, 2, 3, 4})
		})
		Convey("Append Sub", func() {
			s := NewSyncSet(1, 2, 3, 4)
			s.Append(1, 2)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldEqual, []int{1, 2, 3, 4})
		})
		Convey("Append Self", func() {
			s := NewSyncSet(1, 2, 3, 4)
			s.Append(1, 2, 3, 4)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldEqual, []int{1, 2, 3, 4})
		})
		Convey("Append More", func() {
			s := NewSyncSet(1, 2, 3, 4)
			s.Append(5, 6, 7, 8)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldEqual, []int{1, 2, 3, 4, 5, 6, 7, 8})
		})
	})
}

func TestSyncSetValues(t *testing.T) {
	Convey("HashSet Values", t, func() {
		Convey("Empty Set", func() {
			So(NewSyncSet[int]().Values(), ShouldBeEmpty)
		})
		Convey("Non-Empty Set", func() {
			vs := NewSyncSet(1, 2, 3, 4).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
	})
}

func TestSyncSetRemove(t *testing.T) {
	Convey("HashSet Remove", t, func() {
		Convey("Empty Set", func() {
			s := NewSyncSet[int]()
			s.Remove(10)
			So(s.Values(), ShouldBeEmpty)
		})
		Convey("Non-Empty Set remove exist", func() {
			s := NewSyncSet(1, 2, 3, 4)
			s.Remove(3)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 4})
		})
		Convey("Non-Empty Set remove not exist", func() {
			s := NewSyncSet(1, 2, 3, 4)
			s.Remove(5)
			vs := s.Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
	})
}

func TestSyncSetContains(t *testing.T) {
	Convey("HashSet Contains", t, func() {
		Convey("Empty Set", func() {
			So(NewSyncSet[int]().Contains(10), ShouldBeFalse)
		})
		Convey("Non-Empty Set", func() {
			So(NewSyncSet(1, 2, 3, 4).Contains(2), ShouldBeTrue)
			So(NewSyncSet(1, 2, 3, 4).Contains(10), ShouldBeFalse)
		})
	})
}

func TestSyncSetUnion(t *testing.T) {
	Convey("HashSet Union", t, func() {
		Convey("Two Empty Set", func() {
			So(NewSyncSet[int]().Union(NewSyncSet[int]()).Values(), ShouldBeEmpty)
		})
		Convey("Union of non-empty and empty sets", func() {
			s1 := NewSyncSet[int]()
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Union of non-empty and empty sets 2", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet[int]()
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Union of two sets without intersection", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(6, 7, 8, 9)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4, 6, 7, 8, 9})
		})
		Convey("Union of two sets with intersection", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(3, 4, 5, 6)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4, 5, 6})
		})
		Convey("Union of two identical sets", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Union of two sets without difference", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(2, 3)
			vs := s1.Union(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
	})
}

func TestSyncSetIntersection(t *testing.T) {
	Convey("HashSet Intersection", t, func() {
		Convey("Two Empty Set", func() {
			So(NewSyncSet[int]().Intersection(NewSyncSet[int]()).Values(), ShouldBeEmpty)
		})
		Convey("Intersection of non-empty and empty sets", func() {
			s1 := NewSyncSet[int]()
			s2 := NewSyncSet(1, 2, 3, 4)
			So(s1.Intersection(s2).Values(), ShouldBeEmpty)
		})
		Convey("Intersection of non-empty and empty sets 2", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet[int]()
			So(s1.Intersection(s2).Values(), ShouldBeEmpty)
		})
		Convey("Intersection of two sets without intersection", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(6, 7, 8, 9)
			So(s1.Intersection(s2).Values(), ShouldBeEmpty)
		})
		Convey("Intersection of two sets with intersection", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(3, 4, 5, 6)
			vs := s1.Intersection(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{3, 4})
		})
		Convey("Intersection of two identical sets", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.Intersection(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Intersection of two sets without difference", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(2, 3)
			vs := s1.Intersection(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{2, 3})
		})
	})
}

func TestSyncSetDifference(t *testing.T) {
	Convey("HashSet Difference", t, func() {
		Convey("Two Empty Set", func() {
			So(NewSyncSet[int]().Difference(NewSyncSet[int]()).Values(), ShouldBeEmpty)
		})
		Convey("Difference of non-empty and empty sets", func() {
			s1 := NewSyncSet[int]()
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldBeEmpty)
		})
		Convey("Difference of non-empty and empty sets 2", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet[int]()
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Difference of two sets without intersection", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(6, 7, 8, 9)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("Difference of two sets with intersection", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(3, 4, 5, 6)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2})
		})
		Convey("Difference of two sets with intersection 2", func() {
			s1 := NewSyncSet(3, 4, 5, 6)
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{5, 6})
		})
		Convey("Difference of two identical sets", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(1, 2, 3, 4)
			So(s1.Difference(s2).Values(), ShouldBeEmpty)
		})
		Convey("Difference of two sets without difference", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(2, 3)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 4})
		})
		Convey("Difference of two sets without difference 2", func() {
			s1 := NewSyncSet(2, 3)
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.Difference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldBeEmpty)
		})
	})
}

func TestSyncSetSymmetricDifference(t *testing.T) {
	Convey("HashSet SymmetricDifference", t, func() {
		Convey("Two Empty Set", func() {
			So(NewSyncSet[int]().SymmetricDifference(NewSyncSet[int]()).Values(), ShouldBeEmpty)
		})
		Convey("SymmetricDifference of non-empty and empty sets", func() {
			s1 := NewSyncSet[int]()
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("SymmetricDifference of non-empty and empty sets 2", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet[int]()
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4})
		})
		Convey("SymmetricDifference of two sets without intersection", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(6, 7, 8, 9)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 3, 4, 6, 7, 8, 9})
		})
		Convey("SymmetricDifference of two sets with intersection", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(3, 4, 5, 6)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 5, 6})
		})
		Convey("SymmetricDifference of two sets with intersection 2", func() {
			s1 := NewSyncSet(3, 4, 5, 6)
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 2, 5, 6})
		})
		Convey("SymmetricDifference of two identical sets", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(1, 2, 3, 4)
			So(s1.SymmetricDifference(s2).Values(), ShouldBeEmpty)
		})
		Convey("SymmetricDifference of two sets without difference", func() {
			s1 := NewSyncSet(1, 2, 3, 4)
			s2 := NewSyncSet(2, 3)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 4})
		})
		Convey("SymmetricDifference of two sets without difference 2", func() {
			s1 := NewSyncSet(2, 3)
			s2 := NewSyncSet(1, 2, 3, 4)
			vs := s1.SymmetricDifference(s2).Values()
			sort.Ints(vs)
			So(vs, ShouldResemble, []int{1, 4})
		})
	})
}

func TestSyncSetString(t *testing.T) {
	Convey("Set String", t, func() {
		Convey("empty collection", func() {
			So(NewSyncSet[int]().String(), ShouldEqual, "[]")
		})
		Convey("non-empty collection", func() {
			So(NewSyncSet(1).String(), ShouldEqual, "[1]")
		})
	})
}
