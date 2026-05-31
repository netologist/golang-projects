package main

import (
	"context"
	"errors"
	"testing"
)

func setup() (*CommandBus, *QueryBus) {
	proj := newProjection()
	cb := NewCommandBus()
	qb := NewQueryBus()
	cb.Register("Deposit", &depositHandler{proj: proj})
	qb.Register("Balance", &balanceHandler{proj: proj})
	return cb, qb
}

func TestCQRS_commandThenQuery(t *testing.T) {
	cb, qb := setup()
	ctx := context.Background()

	if err := cb.Dispatch(ctx, DepositCommand{AccountID: "a", Amount: 100}); err != nil {
		t.Fatal(err)
	}
	if err := cb.Dispatch(ctx, DepositCommand{AccountID: "a", Amount: 50}); err != nil {
		t.Fatal(err)
	}

	res, err := qb.Ask(ctx, BalanceQuery{AccountID: "a"})
	if err != nil {
		t.Fatal(err)
	}
	if res.(int) != 150 {
		t.Errorf("balance: got %v, want 150", res)
	}
}

func TestCQRS_unknownCommand(t *testing.T) {
	cb, _ := setup()
	type unknownCmd struct{}
	err := cb.Dispatch(context.Background(), fakeCmd{})
	if !errors.Is(err, ErrHandlerNotFound) {
		t.Errorf("expected ErrHandlerNotFound, got %v", err)
	}
	_ = unknownCmd{}
}

type fakeCmd struct{}

func (fakeCmd) CommandName() string { return "Nonexistent" }
