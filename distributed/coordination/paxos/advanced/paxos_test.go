package main

import (
	"errors"
	"testing"
)

func makeAcceptors(n int) []*Acceptor {
	acc := make([]*Acceptor, n)
	for i := range acc {
		acc[i] = NewAcceptor(i)
	}
	return acc
}

func TestPaxos_SingleProposerReachesConsensus(t *testing.T) {
	acceptors := makeAcceptors(3)
	p := NewProposer(1, acceptors)
	v, err := p.Propose(10, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if v != "alice" {
		t.Fatalf("want alice, got %q", v)
	}
}

func TestPaxos_HigherIDPreempts(t *testing.T) {
	acceptors := makeAcceptors(3)
	p1 := NewProposer(1, acceptors)
	p2 := NewProposer(2, acceptors)

	// p1 proposes
	v1, err := p1.Propose(10, "alice")
	if err != nil {
		t.Fatal(err)
	}

	// p2 uses higher ID; Paxos safety requires it adopts v1
	v2, err := p2.Propose(20, "bob")
	if err != nil {
		t.Fatal(err)
	}
	if v1 != v2 {
		t.Fatalf("safety violated: p1=%q p2=%q", v1, v2)
	}
}

func TestPaxos_QuorumFailure(t *testing.T) {
	// Only 1 acceptor, quorum = 1, should succeed
	// But test: proposer with no acceptors fails
	p := NewProposer(1, []*Acceptor{})
	// quorum for 0 acceptors is 1, so it should fail
	// Actually NewProposer sets quorum = 0/2+1 = 1, but no acceptors respond.
	// Let's test with proposer using 3 acceptors but low propID (blocked by prior promise)
	acceptors := makeAcceptors(3)
	p2 := NewProposer(2, acceptors)
	// First, have acceptors promise to a higher ID
	for _, a := range acceptors {
		a.Prepare(100)
	}
	// Now p2 with lower ID should fail prepare
	_, err := p2.Propose(5, "fail")
	if !errors.Is(err, ErrNoQuorum) {
		t.Fatalf("want ErrNoQuorum, got %v", err)
	}
	_ = p
}
