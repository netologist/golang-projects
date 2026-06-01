package main

import (
	"fmt"
	"time"
)

type vote int

const (
	Yes vote = iota
	No
)

type participant struct {
	id   int
	vote vote
}

func (p *participant) canCommit() vote  { return p.vote }
func (p *participant) preCommit() bool  { return p.vote == Yes }
func (p *participant) doCommit() string { return fmt.Sprintf("p%d committed", p.id) }
func (p *participant) doAbort() string  { return fmt.Sprintf("p%d aborted", p.id) }

func run3PC(participants []*participant, crashAfterPreCommit bool) {
	fmt.Printf("\n--- 3PC with %d participants (crashAfterPreCommit=%v) ---\n",
		len(participants), crashAfterPreCommit)

	// Phase 1: CanCommit
	allYes := true
	for _, p := range participants {
		v := p.canCommit()
		fmt.Printf("  p%d vote: %v\n", p.id, map[vote]string{Yes: "YES", No: "NO"}[v])
		if v == No {
			allYes = false
		}
	}
	if !allYes {
		for _, p := range participants {
			fmt.Println(" ", p.doAbort())
		}
		return
	}

	// Phase 2: PreCommit
	fmt.Println("  [coordinator] sending PreCommit")
	for _, p := range participants {
		p.preCommit()
	}

	if crashAfterPreCommit {
		fmt.Println("  [coordinator] CRASHED after PreCommit")
		// Participants timeout and abort
		time.Sleep(10 * time.Millisecond)
		for _, p := range participants {
			fmt.Printf("  p%d: coordinator timeout -> abort\n", p.id)
		}
		return
	}

	// Phase 3: DoCommit
	fmt.Println("  [coordinator] sending DoCommit")
	for _, p := range participants {
		fmt.Println(" ", p.doCommit())
	}
}

func main() {
	// Happy path
	run3PC([]*participant{{1, Yes}, {2, Yes}, {3, Yes}}, false)

	// Abort path
	run3PC([]*participant{{1, Yes}, {2, No}, {3, Yes}}, false)

	// Coordinator crash after PreCommit
	run3PC([]*participant{{1, Yes}, {2, Yes}, {3, Yes}}, true)
}
