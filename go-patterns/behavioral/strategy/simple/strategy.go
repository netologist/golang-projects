package main

import "sort"

// SortStrategy sorts data in place.
type SortStrategy func([]int)

// Sorter applies a swappable SortStrategy.
type Sorter struct{ strategy SortStrategy }

// SetStrategy chooses the algorithm.
func (s *Sorter) SetStrategy(fn SortStrategy) { s.strategy = fn }

// Sort runs the current strategy.
func (s *Sorter) Sort(data []int) { s.strategy(data) }

// BubbleSort is a simple O(n^2) strategy.
var BubbleSort SortStrategy = func(data []int) {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data)-1-i; j++ {
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
}

// StdSort delegates to the standard library.
var StdSort SortStrategy = func(data []int) { sort.Ints(data) }
