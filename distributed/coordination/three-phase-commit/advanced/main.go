package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("========== 3-Phase Commit Demo ==========")

	fmt.Println("\n--- Happy Path ---")
	p1 := []*Participant{NewParticipant(1, true), NewParticipant(2, true), NewParticipant(3, true)}
	c1 := NewCoordinator(p1, 50*time.Millisecond)
	err := c1.RunCommit()
	fmt.Printf("result: %v\n", err)

	fmt.Println("\n--- Abort Path (participant-2 votes NO) ---")
	p2 := []*Participant{NewParticipant(1, true), NewParticipant(2, false), NewParticipant(3, true)}
	c2 := NewCoordinator(p2, 50*time.Millisecond)
	err = c2.RunCommit()
	fmt.Printf("result: %v\n", err)

	fmt.Println("\n--- Crash Recovery (coordinator crashes after PreCommit) ---")
	p3 := []*Participant{NewParticipant(1, true), NewParticipant(2, true), NewParticipant(3, true)}
	c3 := NewCoordinator(p3, 30*time.Millisecond)
	err = c3.RunWithCrash()
	fmt.Printf("result: %v\n", err)
}
