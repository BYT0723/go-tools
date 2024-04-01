package set

import (
	"fmt"
	"testing"
)

func TestUnion(t *testing.T) {
	type T any
	type args struct {
		s1 *Set[T]
		s2 *Set[T]
	}
	tests := []struct {
		name string
		args args
		want *Set[T]
	}{
		{
			name: "Union of two empty sets",
			args: args{
				s1: NewSet[T](),
				s2: NewSet[T](),
			},
			want: NewSet[T](),
		},
		{
			name: "Union of an empty set and a non-empty set with common elements",
			args: args{},
			want: NewSet[T](1, 2, 3, 4, 5),
		},
		{
			name: "Union of two non-empty sets without common elements",
			args: args{
				s1: NewSet[T](1, 2),
				s2: NewSet[T](3, 4),
			},
			want: NewSet[T](1, 2, 3, 4),
		},
		// Add more test cases here
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Union(tt.args.s1, tt.args.s2); !got.Equal(tt.want) {
				t.Errorf("Union() = %v, want %v", got.Values(), tt.want.Values())
			}
		})
	}
}

func TestUnionFunc(t *testing.T) {
	type T any
	type args struct {
		s1         *Set[T]
		s2         *Set[T]
		identifier func(T) string
	}
	tests := []struct {
		name string
		args args
		want *Set[T]
	}{
		{
			name: "Union of two empty sets",
			args: args{
				s1:         NewSet[T](),
				s2:         NewSet[T](),
				identifier: func(t T) string { return fmt.Sprint(t) },
			},
			want: NewSet[T](),
		},
		{
			name: "Union of an empty set and a non-empty set with common elements",
			args: args{
				s1:         NewSet[T](),
				s2:         NewSet[T](3, 4, 5),
				identifier: func(t T) string { return fmt.Sprint(t) },
			},
			want: NewSet[T](3, 4, 5),
		},
		{
			name: "Union of two non-empty sets without common elements",
			args: args{
				s1:         NewSet[T](1, 2),
				s2:         NewSet[T](3, 4),
				identifier: func(t T) string { return fmt.Sprint(t) },
			},
			want: NewSet[T](1, 2, 3, 4),
		},
		// Add more test cases here
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnionFunc(tt.args.s1, tt.args.s2, tt.args.identifier); !got.Equal(tt.want) {
				t.Errorf("UnionFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersection(t *testing.T) {
	type T any
	type args struct {
		s1 *Set[T]
		s2 *Set[T]
	}
	tests := []struct {
		name string
		args args
		want *Set[T]
	}{
		// TODO: Add test cases.
		{
			name: "Intersection of two empty sets",
			args: args{
				s1: NewSet[T](),
				s2: NewSet[T](),
			},
			want: NewSet[T](),
		},
		{
			name: "Intersection of one and one empty sets with common elements",
			args: args{
				s1: NewSet[T](),
				s2: NewSet[T](3, 4, 5),
			},
			want: NewSet[T](),
		},
		{
			name: "Intersection of two non-empty sets without common elements",
			args: args{
				s1: NewSet[T](1, 2, 3),
				s2: NewSet[T](3, 4, 5),
			},
			want: NewSet[T](3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Intersection(tt.args.s1, tt.args.s2); !got.Equal(tt.want) {
				t.Errorf("Intersection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntersectionFunc(t *testing.T) {
	type T any
	type args struct {
		s1         *Set[T]
		s2         *Set[T]
		identifier func(T) string
	}
	tests := []struct {
		name string
		args args
		want *Set[T]
	}{
		// TODO: Add test cases.
		{
			name: "Intersection of two empty sets",
			args: args{
				s1:         NewSet[T](),
				s2:         NewSet[T](),
				identifier: func(t T) string { return fmt.Sprint(t) },
			},
			want: NewSet[T](),
		},
		{
			name: "Intersection of one and one empty sets with common elements",
			args: args{
				s1:         NewSet[T](),
				s2:         NewSet[T](3, 4, 5),
				identifier: func(t T) string { return fmt.Sprint(t) },
			},
			want: NewSet[T](),
		},
		{
			name: "Intersection of two non-empty sets without common elements",
			args: args{
				s1:         NewSet[T](1, 2, 3),
				s2:         NewSet[T](2, 3, 4),
				identifier: func(t T) string { return fmt.Sprint(t) },
			},
			want: NewSet[T](2, 3, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntersectionFunc(tt.args.s1, tt.args.s2, tt.args.identifier); !got.Equal(tt.want) {
				t.Errorf("IntersectionFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	type T any
	type args struct {
		s1 *Set[T]
		s2 *Set[T]
	}
	tests := []struct {
		name string
		args args
		want *Set[T]
	}{
		// TODO: Add test cases.
		{
			name: "Intersection of two empty sets",
			args: args{
				s1: NewSet[T](),
				s2: NewSet[T](),
			},
			want: NewSet[T](),
		},
		{
			name: "Intersection of one and one empty sets with common elements",
			args: args{
				s1: NewSet[T](),
				s2: NewSet[T](3, 4, 5),
			},
			want: NewSet[T](3, 4, 5),
		},
		{
			name: "Intersection of two non-empty sets without common elements",
			args: args{
				s1: NewSet[T](1, 2, 3),
				s2: NewSet[T](2, 3, 4),
			},
			want: NewSet[T](1, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Difference(tt.args.s1, tt.args.s2); !got.Equal(tt.want) {
				t.Errorf("Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDifferenceFunc(t *testing.T) {
	type T any
	type args struct {
		s1         *Set[T]
		s2         *Set[T]
		identifier func(T) string
	}
	tests := []struct {
		name string
		args args
		want *Set[T]
	}{
		// TODO: Add test cases.
		{
			name: "Intersection of two empty sets",
			args: args{
				s1:         NewSet[T](),
				s2:         NewSet[T](),
				identifier: func(t T) string { return fmt.Sprint(t) },
			},
			want: NewSet[T](),
		},
		{
			name: "Intersection of one and one empty sets with common elements",
			args: args{
				s1:         NewSet[T](),
				s2:         NewSet[T](3, 4, 5),
				identifier: func(t T) string { return fmt.Sprint(t) },
			},
			want: NewSet[T](3, 4, 5),
		},
		{
			name: "Intersection of two non-empty sets without common elements",
			args: args{
				s1:         NewSet[T](1, 2, 3),
				s2:         NewSet[T](2, 3, 4),
				identifier: func(t T) string { return fmt.Sprint(t) },
			},
			want: NewSet[T](1, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DifferenceFunc(tt.args.s1, tt.args.s2, tt.args.identifier); !got.Equal(tt.want) {
				t.Errorf("DifferenceFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
