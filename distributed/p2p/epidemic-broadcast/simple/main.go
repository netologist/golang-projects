package main

import (
	"fmt"
	"math/rand/v2"
)

const total = 10

func main() {
	type Status int
	const (
		Susceptible Status = iota
		Infected
		Dead
	)

	status := make([]Status, total)
	status[0] = Infected // seed
	status[5] = Dead     // dead node

	for round := 1; ; round++ {
		for i := 0; i < total; i++ {
			if status[i] != Infected {
				continue
			}
			for f := 0; f < 2; f++ {
				target := rand.IntN(total)
				if status[target] == Susceptible {
					status[target] = Infected
				}
			}
		}
		// Count live infected
		liveInfected, liveTotal := 0, 0
		for _, s := range status {
			if s != Dead {
				liveTotal++
				if s == Infected {
					liveInfected++
				}
			}
		}
		fmt.Printf("round %d: %d/%d live nodes infected\n", round, liveInfected, liveTotal)
		if liveInfected == liveTotal {
			fmt.Printf("converged in %d rounds (dead node skipped)\n", round)
			break
		}
		if round > 200 {
			fmt.Println("did not converge")
			break
		}
	}
}
