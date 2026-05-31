package main

import (
	"testing"
	"time"
)

func dep(n int) Deposited {
	return Deposited{BaseEvent: BaseEvent{Type: "Deposited", Time: time.Now()}, Amount: n}
}

func wd(n int) Withdrawn {
	return Withdrawn{BaseEvent: BaseEvent{Type: "Withdrawn", Time: time.Now()}, Amount: n}
}

func TestStore_appendAndRebuild(t *testing.T) {
	s := NewStore()
	if err := s.Append("acc:1", []Event{dep(100), wd(30), dep(50)}, 0); err != nil {
		t.Fatal(err)
	}
	events, err := s.Load("acc:1")
	if err != nil {
		t.Fatal(err)
	}
	acc := Rebuild(events)
	if acc.Balance != 120 {
		t.Errorf("balance: got %d, want 120", acc.Balance)
	}
	if acc.Version != 3 {
		t.Errorf("version: got %d, want 3", acc.Version)
	}
}

func TestStore_optimisticConcurrency(t *testing.T) {
	s := NewStore()
	_ = s.Append("acc:2", []Event{dep(10)}, 0)
	// Stale expected version triggers a conflict.
	err := s.Append("acc:2", []Event{dep(5)}, 0)
	if err == nil {
		t.Error("expected concurrency conflict")
	}
}
