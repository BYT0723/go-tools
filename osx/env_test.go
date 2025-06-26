package osx

import (
	"os"
	"reflect"
	"testing"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("AGE", "25")
	type args struct {
		key string
		def any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "exist key (string)",
			args: args{
				key: "HOME",
				def: "walter",
			},
			want: "/home/walter",
		},
		{
			name: "not exist key (string)",
			args: args{
				key: "USERNAME",
				def: "walter",
			},
			want: "walter",
		},
		{
			name: "exist key (int)",
			args: args{
				key: "AGE",
				def: 10,
			},
			want: 25,
		},
		{
			name: "not exist key (int)",
			args: args{
				key: "AGE_TEST",
				def: 10,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnv(tt.args.key, tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
