package util

import (
	"github.com/stretchr/testify/assert"
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
		// results get sorted by the test, so that a direct comparison is possible
		{name: "no duplicates", args: args{list: []string{"a", "b"}}, want: nil},
		{name: "no duplicates", args: args{list: []string{"a", "a"}}, want: []string{"a"}},
		{name: "no duplicates", args: args{list: []string{"a", "a", "a"}}, want: []string{"a"}},
		{name: "no duplicates", args: args{list: []string{"a", "a", "b", "b"}}, want: []string{"a", "b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDuplicates(tt.args.list)
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
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

func TestGroupBy(t *testing.T) {
	type args[V any, K comparable] struct {
		elements    []V
		keySelector func(V) K
	}
	type testCase[V any, K comparable] struct {
		name string
		args args[V, K]
		want map[K][]V
	}
	tests := []testCase[testStruct, int]{
		{
			name: "should group by number",
			args: args[testStruct, int]{
				elements: []testStruct{
					{status: 1, name: "a"},
					{status: 1, name: "aa"},
					{status: 2, name: "b"},
					{status: 3, name: "c"},
				},
				keySelector: func(t testStruct) int {
					return t.status
				},
			},
			want: map[int][]testStruct{
				1: {{status: 1, name: "a"}, {status: 1, name: "aa"}},
				2: {{status: 2, name: "b"}},
				3: {{status: 3, name: "c"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GroupBy(tt.args.elements, tt.args.keySelector), "GroupBy(%v, %v)", tt.args.elements, tt.args.keySelector)
		})
	}
}

type testStruct struct {
	name   string
	status int
}
