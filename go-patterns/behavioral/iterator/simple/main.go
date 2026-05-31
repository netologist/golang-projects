package main

import "fmt"

func main() {
	for v := range SliceIter([]string{"a", "b", "c"}) {
		fmt.Println(v)
	}
}
