package main

import "testing"

func TestEpidemic_ConvergenceWithDeadNode(t *testing.T) {
	for trial := 0; trial < 5; trial++ {
		e := New(10, 2)
		e.Infect(0)
		e.Kill(5)

		for round := 0; round < 500; round++ {
			e.Spread()
			if e.Converged() {
				goto ok
			}
		}
		t.Fatal("did not converge")
	ok:
	}
}

func TestEpidemic_DeadNodeStaysDead(t *testing.T) {
	e := New(5, 10)
	e.Infect(0)
	e.Kill(3)

	for round := 0; round < 100; round++ {
		e.Spread()
	}
	if e.StatusOf(3) != Dead {
		t.Fatal("dead node should stay dead")
	}
}

func TestEpidemic_AllLiveInfected(t *testing.T) {
	e := New(10, 3)
	e.Infect(0)
	e.Kill(7)

	for round := 0; round < 500 && !e.Converged(); round++ {
		e.Spread()
	}
	if !e.Converged() {
		t.Fatal("not all live nodes infected")
	}
	if e.StatusOf(7) != Dead {
		t.Fatal("node 7 should be dead")
	}
}
