package main

import (
	"context"
	"testing"
)

type alwaysCommit struct{ committed, aborted bool }

func (a *alwaysCommit) Prepare(_ context.Context, _ string) (Vote, error) { return VoteCommit, nil }
func (a *alwaysCommit) Commit(_ context.Context, _ string) error          { a.committed = true; return nil }

func (a *alwaysCommit) Abort(_ context.Context, _ string) error { a.aborted = true; return nil }

type alwaysAbort struct{}

func (a *alwaysAbort) Prepare(_ context.Context, _ string) (Vote, error) { return VoteAbort, nil }
func (a *alwaysAbort) Commit(_ context.Context, _ string) error          { return nil }
func (a *alwaysAbort) Abort(_ context.Context, _ string) error           { return nil }

func TestCoordinator_allCommit(t *testing.T) {
	p1, p2 := &alwaysCommit{}, &alwaysCommit{}
	c := New(p1, p2)
	if err := c.Execute(context.Background(), "tx-1"); err != nil {
		t.Fatal(err)
	}
	if !p1.committed || !p2.committed {
		t.Error("expected both participants committed")
	}
}

func TestCoordinator_oneAborts(t *testing.T) {
	p1, p2 := &alwaysCommit{}, &alwaysAbort{}
	c := New(p1, p2)
	if err := c.Execute(context.Background(), "tx-2"); err == nil {
		t.Error("expected abort error")
	}
	if p1.committed {
		t.Error("p1 must not commit when another aborts")
	}
	if !p1.aborted {
		t.Error("p1 should have been told to abort")
	}
}
