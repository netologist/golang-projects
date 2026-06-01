package main

import "fmt"

type acceptor struct {
	promisedID  int
	acceptedID  int
	acceptedVal string
}

func (a *acceptor) prepare(propID int) (ok bool, accID int, accVal string) {
	if propID > a.promisedID {
		a.promisedID = propID
		return true, a.acceptedID, a.acceptedVal
	}
	return false, 0, ""
}

func (a *acceptor) accept(propID int, val string) bool {
	if propID >= a.promisedID {
		a.promisedID = propID
		a.acceptedID = propID
		a.acceptedVal = val
		return true
	}
	return false
}

func propose(acceptors []*acceptor, propID int, value string) (string, bool) {
	quorum := len(acceptors)/2 + 1

	// Phase 1: Prepare
	promises := 0
	highestID := -1
	highestVal := ""
	for _, a := range acceptors {
		ok, accID, accVal := a.prepare(propID)
		if ok {
			promises++
			if accID > highestID {
				highestID, highestVal = accID, accVal
			}
		}
	}
	fmt.Printf("prepare %d: promises=%d/%d\n", propID, promises, len(acceptors))
	if promises < quorum {
		return "", false
	}

	// Use previously accepted value if any
	if highestVal != "" {
		value = highestVal
		fmt.Printf("adopt previously accepted value %q\n", value)
	}

	// Phase 2: Accept
	accepts := 0
	for _, a := range acceptors {
		if a.accept(propID, value) {
			accepts++
		}
	}
	fmt.Printf("accept  %d: accepts=%d/%d\n", propID, accepts, len(acceptors))
	if accepts < quorum {
		return "", false
	}
	return value, true
}

func main() {
	acceptors := []*acceptor{{}, {}, {}}

	// Proposer 1 with ID 10
	v, ok := propose(acceptors, 10, "alice")
	fmt.Printf("proposer-1 result: ok=%v value=%q\n\n", ok, v)

	// Proposer 2 with higher ID 20
	v2, ok2 := propose(acceptors, 20, "bob")
	fmt.Printf("proposer-2 result: ok=%v value=%q\n", ok2, v2)
}
