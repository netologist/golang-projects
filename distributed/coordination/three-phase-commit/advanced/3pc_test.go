package main

import (
	"errors"
	"testing"
	"time"
)

func TestThreePC_SuccessfulCommit(t *testing.T) {
	participants := []*Participant{
		NewParticipant(1, true),
		NewParticipant(2, true),
		NewParticipant(3, true),
	}
	c := NewCoordinator(participants, 50*time.Millisecond)
	if err := c.RunCommit(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range participants {
		if p.Phase() != PhaseCommit {
			t.Errorf("participant %d: want Commit, got %s", p.ID, p.Phase())
		}
	}
}

func TestThreePC_AbortOnVoteNo(t *testing.T) {
	participants := []*Participant{
		NewParticipant(1, true),
		NewParticipant(2, false), // votes NO
		NewParticipant(3, true),
	}
	c := NewCoordinator(participants, 50*time.Millisecond)
	err := c.RunCommit()
	if !errors.Is(err, ErrAborted) {
		t.Fatalf("want ErrAborted, got %v", err)
	}
}

func TestThreePC_CrashRecovery(t *testing.T) {
	participants := []*Participant{
		NewParticipant(1, true),
		NewParticipant(2, true),
		NewParticipant(3, true),
	}
	c := NewCoordinator(participants, 20*time.Millisecond)
	err := c.RunWithCrash()
	if !errors.Is(err, ErrTimeout) {
		t.Fatalf("want ErrTimeout, got %v", err)
	}
	for _, p := range participants {
		if p.Phase() != PhaseAbort {
			t.Errorf("participant %d: want Abort, got %s", p.ID, p.Phase())
		}
	}
}
