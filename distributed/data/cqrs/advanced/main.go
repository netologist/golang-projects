package main

import (
	"context"
	"fmt"
)

func main() {
	proj := newProjection()
	cb := NewCommandBus()
	qb := NewQueryBus()
	cb.Register("Deposit", &depositHandler{proj: proj})
	qb.Register("Balance", &balanceHandler{proj: proj})

	ctx := context.Background()
	_ = cb.Dispatch(ctx, DepositCommand{AccountID: "acc:1", Amount: 200})
	_ = cb.Dispatch(ctx, DepositCommand{AccountID: "acc:1", Amount: 75})

	bal, _ := qb.Ask(ctx, BalanceQuery{AccountID: "acc:1"})
	fmt.Printf("balance of acc:1 = %v\n", bal)
}
