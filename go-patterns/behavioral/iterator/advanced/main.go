package main

import (
	"fmt"
	"slices"
)

func main() {
	src := FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	// Compose: keep evens, square them, take 3.
	pipe := Take(Map(Filter(src, func(n int) bool { return n%2 == 0 }), func(n int) int { return n * n }), 3)
	fmt.Println(slices.Collect(pipe)) // [4 16 36]
}
