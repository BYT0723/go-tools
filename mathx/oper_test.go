package mathx

import (
	"reflect"
	"testing"
)

func TestDivCeil(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "zero1", args: args{a: 0, b: 1}, want: 0},
		{name: "zero2", args: args{a: 0, b: -1}, want: 0},
		{name: "p / p", args: args{a: 3, b: 2}, want: 2},
		{name: "p / n", args: args{a: 3, b: -2}, want: -1},
		{name: "n / n", args: args{a: -3, b: -2}, want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DivCeil[int](tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DivCeil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDivFloor(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "zero1", args: args{a: 0, b: 1}, want: 0},
		{name: "zero2", args: args{a: 0, b: -1}, want: 0},
		{name: "p / p", args: args{a: 3, b: 2}, want: 1},
		{name: "p / n", args: args{a: 3, b: -2}, want: -2},
		{name: "n / n", args: args{a: -3, b: -2}, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DivFloor[int](tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DivFloor() = %v, want %v", got, tt.want)
			}
		})
	}
}
