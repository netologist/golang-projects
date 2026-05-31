package main

import "fmt"

func main() {
	r := New(150)
	for _, n := range []string{"cache-1", "cache-2", "cache-3"} {
		r.Add(n)
	}

	keys := []string{"session:abc", "session:def", "session:ghi", "session:jkl"}
	fmt.Println("initial placement:")
	for _, k := range keys {
		n, _ := r.Get(k)
		fmt.Printf("  %s -> %s\n", k, n)
	}

	fmt.Println("after adding cache-4:")
	r.Add("cache-4")
	for _, k := range keys {
		n, _ := r.Get(k)
		fmt.Printf("  %s -> %s\n", k, n)
	}
}
