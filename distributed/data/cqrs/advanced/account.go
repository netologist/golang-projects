package main

import (
	"context"
	"sync"
)

// projection is the denormalised read model.
type projection struct {
	mu       sync.RWMutex
	balances map[string]int
}

func newProjection() *projection { return &projection{balances: map[string]int{}} }

// --- Commands ---

// DepositCommand adds funds to an account.
type DepositCommand struct {
	AccountID string
	Amount    int
}

// CommandName identifies the command.
func (DepositCommand) CommandName() string { return "Deposit" }

type depositHandler struct{ proj *projection }

func (h *depositHandler) Handle(_ context.Context, cmd Command) error {
	c := cmd.(DepositCommand)
	h.proj.mu.Lock()
	defer h.proj.mu.Unlock()
	h.proj.balances[c.AccountID] += c.Amount
	return nil
}

// --- Queries ---

// BalanceQuery reads an account balance.
type BalanceQuery struct{ AccountID string }

// QueryName identifies the query.
func (BalanceQuery) QueryName() string { return "Balance" }

type balanceHandler struct{ proj *projection }

func (h *balanceHandler) Handle(_ context.Context, q Query) (any, error) {
	query := q.(BalanceQuery)
	h.proj.mu.RLock()
	defer h.proj.mu.RUnlock()
	return h.proj.balances[query.AccountID], nil
}
