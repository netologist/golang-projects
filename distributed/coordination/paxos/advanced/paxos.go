package main

import (
	"errors"
	"fmt"
	"sync"
)

// ErrNoQuorum is returned when a phase fails to achieve quorum.
var ErrNoQuorum = errors.New("paxos: no quorum")

// Proposal carries a proposal ID and value.
type Proposal struct {
	ID    int
	Value string
}

// PrepareResponse is the promise returned by an acceptor on Prepare.
type PrepareResponse struct {
	OK       bool
	Accepted Proposal // highest accepted proposal, zero if none
}

// Acceptor implements the Paxos acceptor role.
type Acceptor struct {
	mu         sync.Mutex
	ID         int
	promisedID int
	accepted   Proposal
}

// NewAcceptor creates an acceptor with the given ID.
func NewAcceptor(id int) *Acceptor { return &Acceptor{ID: id} }

// Prepare handles Phase 1. Returns a promise if propID > promisedID.
func (a *Acceptor) Prepare(propID int) PrepareResponse {
	a.mu.Lock()
	defer a.mu.Unlock()
	if propID <= a.promisedID {
		return PrepareResponse{OK: false}
	}
	a.promisedID = propID
	return PrepareResponse{OK: true, Accepted: a.accepted}
}

// Accept handles Phase 2. Accepts if propID >= promisedID.
func (a *Acceptor) Accept(p Proposal) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if p.ID < a.promisedID {
		return false
	}
	a.promisedID = p.ID
	a.accepted = p
	return true
}

// AcceptedValue returns the currently accepted value.
func (a *Acceptor) AcceptedValue() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.accepted.Value
}

// Proposer drives the Paxos protocol.
type Proposer struct {
	ID        int
	acceptors []*Acceptor
	quorum    int
}

// NewProposer creates a proposer with the given acceptors.
func NewProposer(id int, acceptors []*Acceptor) *Proposer {
	return &Proposer{
		ID:        id,
		acceptors: acceptors,
		quorum:    len(acceptors)/2 + 1,
	}
}

// Propose runs the full Paxos protocol for the given value.
// Returns the chosen value (may differ from proposed if another was already accepted).
func (p *Proposer) Propose(propID int, value string) (string, error) {
	// Phase 1: Prepare
	var (
		promises    int
		highestProp Proposal
	)
	for _, a := range p.acceptors {
		resp := a.Prepare(propID)
		if !resp.OK {
			continue
		}
		promises++
		if resp.Accepted.ID > highestProp.ID {
			highestProp = resp.Accepted
		}
	}
	fmt.Printf("  [proposer %d] prepare(id=%d): %d/%d promises\n", p.ID, propID, promises, len(p.acceptors))
	if promises < p.quorum {
		return "", fmt.Errorf("%w in prepare phase", ErrNoQuorum)
	}

	// If a value was already accepted, we must propose it
	if highestProp.Value != "" {
		fmt.Printf("  [proposer %d] adopting prior value %q\n", p.ID, highestProp.Value)
		value = highestProp.Value
	}

	// Phase 2: Accept
	accepts := 0
	prop := Proposal{ID: propID, Value: value}
	for _, a := range p.acceptors {
		if a.Accept(prop) {
			accepts++
		}
	}
	fmt.Printf("  [proposer %d] accept(id=%d, val=%q): %d/%d accepts\n",
		p.ID, propID, value, accepts, len(p.acceptors))
	if accepts < p.quorum {
		return "", fmt.Errorf("%w in accept phase", ErrNoQuorum)
	}
	return value, nil
}
