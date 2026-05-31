package main

import "fmt"

// write model
type accounts struct{ balances map[string]int }

func (a *accounts) deposit(id string, amount int) { a.balances[id] += amount }

// read model (projection kept in sync by the command side)
type balanceView struct{ data map[string]int }

func (v *balanceView) get(id string) int { return v.data[id] }

func main() {
	write := &accounts{balances: map[string]int{}}
	read := &balanceView{data: map[string]int{}}

	// command
	write.deposit("acc:1", 100)
	read.data["acc:1"] = write.balances["acc:1"] // project

	// query
	fmt.Println("balance:", read.get("acc:1"))
}
