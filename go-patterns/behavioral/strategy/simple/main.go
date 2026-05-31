package main

import "fmt"

func main() {
	s := &Sorter{}

	a := []int{5, 2, 4, 1, 3}
	s.SetStrategy(BubbleSort)
	s.Sort(a)
	fmt.Println("bubble:", a)

	b := []int{9, 7, 8, 6}
	s.SetStrategy(StdSort)
	s.Sort(b)
	fmt.Println("std:", b)
}
