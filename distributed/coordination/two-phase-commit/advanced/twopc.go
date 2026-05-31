package main

import (
	"context"
	"errors"
	"fmt"
)

// Vote is a participant's phase-1 decision.
type Vote int

// Votes.
const (
	VoteCommit Vote = iota
	VoteAbort
)

// Participant takes part in a 2PC transaction.
type Participant interface {
	Prepare(ctx context.Context, txID string) (Vote, error)
	Commit(ctx context.Context, txID string) error
	Abort(ctx context.Context, txID string) error
}

// Coordinator drives the 2PC protocol across participants.
type Coordinator struct{ participants []Participant }

// New creates a coordinator over the given participants.
func New(participants ...Participant) *Coordinator {
	return &Coordinator{participants: participants}
}

// Execute runs the two-phase commit and returns nil only if all committed.
func (c *Coordinator) Execute(ctx context.Context, txID string) error {
	// Phase 1: Prepare / collect votes.
	for _, p := range c.participants {
		vote, err := p.Prepare(ctx, txID)
		if err != nil || vote == VoteAbort {
			c.abortAll(ctx, txID)
			if err != nil {
				return fmt.Errorf("2pc prepare: %w", err)
			}
			return errors.New("2pc: a participant voted abort")
		}
	}

	// Phase 2: Commit.
	var commitErr error
	for _, p := range c.participants {
		if err := p.Commit(ctx, txID); err != nil {
			commitErr = fmt.Errorf("2pc commit: %w", err)
		}
	}
	return commitErr
}

func (c *Coordinator) abortAll(ctx context.Context, txID string) {
	for _, p := range c.participants {
		_ = p.Abort(ctx, txID)
	}
}
