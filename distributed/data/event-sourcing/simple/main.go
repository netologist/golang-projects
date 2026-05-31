package main

import "fmt"

type event struct {
	kind   string
	amount int
}

func balance(events []event) int {
	b := 0
	for _, e := range events {
		switch e.kind {
		case "deposit":
			b += e.amount
		case "withdraw":
			b -= e.amount
		}
	}
	return b
}

func main() {
	log := []event{
		{"deposit", 100},
		{"withdraw", 30},
		{"deposit", 50},
	}
	fmt.Println("current balance:", balance(log))
}
