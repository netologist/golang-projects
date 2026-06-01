package main

import "fmt"

func main() {
	e := New(10, 2)
	e.Infect(0) // seed
	e.Kill(5)   // dead node

	fmt.Printf("Epidemic SI model: 10 nodes, 1 dead (node-5), fanout=2, seed=node-0\n")

	statusStr := func(s NodeStatus) string {
		switch s {
		case Infected:
			return "I"
		case Dead:
			return "D"
		default:
			return "S"
		}
	}

	for round := 1; ; round++ {
		e.Spread()
		infected, dead, total := e.Stats()
		row := fmt.Sprintf("  round %2d [%d/%d live infected, %d dead]: ", round, infected, total-dead, dead)
		for i := 0; i < total; i++ {
			row += statusStr(e.StatusOf(i))
		}
		fmt.Println(row)
		if e.Converged() {
			fmt.Printf("Converged in %d rounds\n", round)
			break
		}
		if round > 200 {
			fmt.Println("Did not converge")
			break
		}
	}
}
