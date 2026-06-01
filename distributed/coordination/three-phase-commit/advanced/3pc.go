package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Phase3PC enumerates the 3PC protocol phases.
type Phase3PC int

const (
	PhaseInit      Phase3PC = iota
	PhaseCanCommit          // waiting for votes
	PhasePreCommit          // pre-commit sent, waiting for acks
	PhaseCommit             // final commit
	PhaseAbort              // aborted
)

func (ph Phase3PC) String() string {
	return [...]string{"Init", "CanCommit", "PreCommit", "Commit", "Abort"}[ph]
}

// ErrAborted is returned when the transaction is aborted.
var ErrAborted = errors.New("3pc: transaction aborted")

// ErrTimeout is returned on coordinator crash recovery.
var ErrTimeout = errors.New("3pc: coordinator timeout")

// Participant is a single 3PC participant.
type Participant struct {
	mu          sync.Mutex
	ID          int
	phase       Phase3PC
	willVoteYes bool // set by tests to control behavior
}

// NewParticipant creates a participant. willVoteYes controls its CanCommit vote.
func NewParticipant(id int, willVoteYes bool) *Participant {
	return &Participant{ID: id, willVoteYes: willVoteYes}
}

// CanCommit returns true if participant is willing to commit.
func (p *Participant) CanCommit() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.phase = PhaseCanCommit
	return p.willVoteYes
}

// PreCommit transitions participant to pre-committed state.
func (p *Participant) PreCommit() {
	p.mu.Lock()
	p.phase = PhasePreCommit
	p.mu.Unlock()
}

// DoCommit commits the transaction.
func (p *Participant) DoCommit() {
	p.mu.Lock()
	p.phase = PhaseCommit
	p.mu.Unlock()
	fmt.Printf("  [participant %d] COMMITTED\n", p.ID)
}

// Abort aborts the transaction.
func (p *Participant) Abort() {
	p.mu.Lock()
	p.phase = PhaseAbort
	p.mu.Unlock()
	fmt.Printf("  [participant %d] ABORTED\n", p.ID)
}

// Phase returns the participant's current phase.
func (p *Participant) Phase() Phase3PC {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.phase
}

// Coordinator drives the 3PC protocol.
type Coordinator struct {
	participants []*Participant
	Timeout      time.Duration
}

// NewCoordinator creates a coordinator with the given participants and timeout.
func NewCoordinator(participants []*Participant, timeout time.Duration) *Coordinator {
	return &Coordinator{participants: participants, Timeout: timeout}
}

// RunCommit runs the full 3PC protocol. Returns nil on commit, ErrAborted on abort.
func (c *Coordinator) RunCommit() error {
	fmt.Println("[coordinator] Phase 1: CanCommit")
	for _, p := range c.participants {
		if !p.CanCommit() {
			fmt.Printf("  participant %d voted NO\n", p.ID)
			c.abortAll()
			return ErrAborted
		}
		fmt.Printf("  participant %d voted YES\n", p.ID)
	}

	fmt.Println("[coordinator] Phase 2: PreCommit")
	for _, p := range c.participants {
		p.PreCommit()
		fmt.Printf("  participant %d pre-committed\n", p.ID)
	}

	fmt.Println("[coordinator] Phase 3: DoCommit")
	for _, p := range c.participants {
		p.DoCommit()
	}
	return nil
}

// RunWithCrash simulates coordinator crash after PreCommit.
// Participants detect timeout and abort.
func (c *Coordinator) RunWithCrash() error {
	fmt.Println("[coordinator] Phase 1: CanCommit")
	for _, p := range c.participants {
		if !p.CanCommit() {
			c.abortAll()
			return ErrAborted
		}
	}

	fmt.Println("[coordinator] Phase 2: PreCommit (coordinator will crash after this)")
	for _, p := range c.participants {
		p.PreCommit()
	}

	fmt.Printf("[coordinator] CRASHED after PreCommit (timeout=%s)\n", c.Timeout)
	// Simulate participants waiting for DoCommit, then timing out
	time.Sleep(c.Timeout)
	for _, p := range c.participants {
		if p.Phase() == PhasePreCommit {
			p.Abort()
			fmt.Printf("  participant %d: timed out waiting for DoCommit -> ABORT\n", p.ID)
		}
	}
	return ErrTimeout
}

func (c *Coordinator) abortAll() {
	for _, p := range c.participants {
		p.Abort()
	}
}
