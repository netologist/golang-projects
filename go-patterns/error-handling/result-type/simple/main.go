package main

import "fmt"

func main() {
	r := Ok(3).Map(func(n int) int { return n + 4 }).Map(func(n int) int { return n * 2 })
	val, err := r.Unwrap()
	fmt.Printf("val=%d err=%v\n", val, err)
}
