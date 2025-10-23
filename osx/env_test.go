package osx

import (
	"os"
	"reflect"
	"testing"
)

func TestGetEnv(t *testing.T) {
	type (
		args[T any] struct {
			key string
			def T
		}
		item[T any] struct {
			name string
			args args[T]
			want T
		}
		Integer int
	)

	stringTests := []item[string]{
		{
			name: "exist key (string)",
			args: args[string]{
				key: "HOME",
				def: "HOME",
			},
			want: os.Getenv("HOME"),
		},
		{
			name: "not exist key (string)",
			args: args[string]{
				key: "USERNAME",
				def: "default",
			},
			want: "default",
		},
	}

	_ = os.Setenv("AGE", "25")
	intTests := []item[int]{
		{
			name: "exist key (int)",
			args: args[int]{
				key: "AGE",
				def: 10,
			},
			want: 25,
		},
		{
			name: "not exist key (int)",
			args: args[int]{
				key: "AGE_TEST",
				def: 10,
			},
			want: 10,
		},
	}

	IntegerTests := []item[Integer]{
		{
			name: "exist key (custom base type)",
			args: args[Integer]{
				key: "AGE",
				def: 10,
			},
			want: 25,
		},
		{
			name: "not exist key (custom base type)",
			args: args[Integer]{
				key: "AGE_TEST",
				def: 10,
			},
			want: 10,
		},
	}

	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnv(tt.args.key, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}

	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnv(tt.args.key, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}

	for _, tt := range IntegerTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnv(tt.args.key, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
