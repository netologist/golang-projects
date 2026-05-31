package main

import (
	"context"
	"errors"
	"testing"
)

func TestSaga_happyPath(t *testing.T) {
	var order []string
	s := New(
		Step{
			Name:       "reserve-stock",
			Execute:    func(_ context.Context, _ State) error { order = append(order, "exec:reserve"); return nil },
			Compensate: func(_ context.Context, _ State) error { order = append(order, "comp:reserve"); return nil },
		},
		Step{
			Name:       "charge-payment",
			Execute:    func(_ context.Context, _ State) error { order = append(order, "exec:payment"); return nil },
			Compensate: func(_ context.Context, _ State) error { order = append(order, "comp:payment"); return nil },
		},
	)
	if err := s.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(order) != 2 {
		t.Errorf("got %v, want 2 exec entries", order)
	}
}

func TestSaga_compensatesOnFailure(t *testing.T) {
	var compensated []string
	errPay := errors.New("payment declined")

	s := New(
		Step{
			Name:       "reserve-stock",
			Execute:    func(_ context.Context, _ State) error { return nil },
			Compensate: func(_ context.Context, _ State) error { compensated = append(compensated, "reserve"); return nil },
		},
		Step{
			Name:       "charge-payment",
			Execute:    func(_ context.Context, _ State) error { return errPay },
			Compensate: func(_ context.Context, _ State) error { return nil },
		},
	)

	err := s.Run(context.Background())
	if !errors.Is(err, errPay) {
		t.Errorf("expected errPay, got %v", err)
	}
	if len(compensated) != 1 || compensated[0] != "reserve" {
		t.Errorf("expected reserve compensation, got %v", compensated)
	}
}
