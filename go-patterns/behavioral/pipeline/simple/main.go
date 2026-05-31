package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	nums := generate(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	squared := square(ctx, nums)
	even := filter(ctx, squared, func(n int) bool { return n%2 == 0 })
	for n := range even {
		fmt.Println(n)
	}
}
