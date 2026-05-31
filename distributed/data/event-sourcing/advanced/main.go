package main

import (
	"fmt"
	"time"
)

func main() {
	s := NewStore()
	now := time.Now()
	_ = s.Append("acc:42", []Event{
		Deposited{BaseEvent{"Deposited", now}, 200},
		Withdrawn{BaseEvent{"Withdrawn", now}, 50},
		Deposited{BaseEvent{"Deposited", now}, 25},
	}, 0)

	events, _ := s.Load("acc:42")
	acc := Rebuild(events)
	fmt.Printf("account balance=%d version=%d\n", acc.Balance, acc.Version)
}
