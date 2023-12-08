package util

import (
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestGetDuplicates(t *testing.T) {
	type args struct {
		list []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		//results get sorted by the test, so that a direct comparison is possible
		{name: "no duplicates", args: args{list: []string{"a", "b"}}, want: nil},
		{name: "no duplicates", args: args{list: []string{"a", "a"}}, want: []string{"a"}},
		{name: "no duplicates", args: args{list: []string{"a", "a", "a"}}, want: []string{"a"}},
		{name: "no duplicates", args: args{list: []string{"a", "a", "b", "b"}}, want: []string{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDuplicates(tt.args.list); !reflect.DeepEqual(got, tt.want) {
				sort.Strings(got)
				t.Errorf("GetDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	type args[T any, V any] struct {
		ts []T
		fn func(T) V
	}
	type testCase[T any, V any] struct {
		name string
		args args[T, V]
		want []V
	}
	tests := []testCase[int, string]{
		{
			name: "int to string",
			args: args[int, string]{
				ts: []int{1, 2, 3},
				fn: func(i int) string {
					return strconv.Itoa(i)
				},
			},
			want: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.ts, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
