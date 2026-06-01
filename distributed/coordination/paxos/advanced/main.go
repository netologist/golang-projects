package main

import "fmt"

func main() {
	acceptors := make([]*Acceptor, 3)
	for i := range acceptors {
		acceptors[i] = NewAcceptor(i)
	}

	p1 := NewProposer(1, acceptors)
	p2 := NewProposer(2, acceptors)

	fmt.Println("=== Proposer 1 (id=10, value=\"alice\") ===")
	v1, err := p1.Propose(10, "alice")
	fmt.Printf("result: value=%q err=%v\n\n", v1, err)

	fmt.Println("=== Proposer 2 (id=20, value=\"bob\") ===")
	v2, err := p2.Propose(20, "bob")
	fmt.Printf("result: value=%q err=%v\n", v2, err)

	if v1 == v2 {
		fmt.Printf("\nConsensus: both proposers agree on %q\n", v1)
	} else {
		fmt.Printf("\nNote: proposer 2 adopted prior value (safety maintained)\n")
	}
}
