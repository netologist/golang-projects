package main

import (
	"slices"
	"testing"
)

func TestFilter(t *testing.T) {
	src := FromSlice([]int{1, 2, 3, 4, 5})
	evens := slices.Collect(Filter(src, func(n int) bool { return n%2 == 0 }))
	if len(evens) != 2 || evens[0] != 2 || evens[1] != 4 {
		t.Errorf("got %v, want [2 4]", evens)
	}
}

func TestTake(t *testing.T) {
	src := FromSlice([]int{10, 20, 30, 40, 50})
	got := slices.Collect(Take(src, 3))
	if len(got) != 3 {
		t.Errorf("got %d items, want 3", len(got))
	}
}

func TestMap(t *testing.T) {
	src := FromSlice([]int{1, 2, 3})
	got := slices.Collect(Map(src, func(n int) int { return n * 10 }))
	want := []int{10, 20, 30}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
