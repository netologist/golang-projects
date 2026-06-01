package main

import (
	"fmt"
	"math/rand/v2"
)

const (
	nodeCount = 10
	fanout    = 2
)

func main() {
	infected := make([]bool, nodeCount)
	infected[0] = true // seed

	for round := 1; ; round++ {
		newInfections := 0
		for i := 0; i < nodeCount; i++ {
			if !infected[i] {
				continue
			}
			for f := 0; f < fanout; f++ {
				target := rand.IntN(nodeCount)
				if !infected[target] {
					infected[target] = true
					newInfections++
				}
			}
		}

		count := 0
		for _, v := range infected {
			if v {
				count++
			}
		}
		fmt.Printf("round %d: %d/%d infected\n", round, count, nodeCount)
		if count == nodeCount {
			fmt.Printf("converged in %d rounds\n", round)
			break
		}
		if round > 100 {
			fmt.Println("did not converge")
			break
		}
	}
}
