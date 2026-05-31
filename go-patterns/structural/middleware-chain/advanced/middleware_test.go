package main

import (
	"context"
	"testing"
)

func TestChain_order(t *testing.T) {
	var order []string
	mw := func(label string) Middleware {
		return func(next HandlerFunc) HandlerFunc {
			return func(ctx context.Context, req Request) (Response, error) {
				order = append(order, label+":before")
				resp, err := next(ctx, req)
				order = append(order, label+":after")
				return resp, err
			}
		}
	}
	h := Chain(
		func(_ context.Context, r Request) (Response, error) { return Response{Body: r.Body}, nil },
		mw("A"), mw("B"),
	)
	h(context.Background(), Request{Body: "test"})
	want := []string{"A:before", "B:before", "B:after", "A:after"}
	if len(order) != len(want) {
		t.Fatalf("got %v, want %v", order, want)
	}
	for i, w := range want {
		if order[i] != w {
			t.Errorf("order[%d]: got %s, want %s", i, order[i], w)
		}
	}
}

func TestRecovery_catchesPanic(t *testing.T) {
	h := Chain(
		func(_ context.Context, _ Request) (Response, error) { panic("boom") },
		Recovery,
	)
	resp, err := h(context.Background(), Request{})
	if err == nil {
		t.Error("expected error")
	}
	if resp.StatusCode != 500 {
		t.Errorf("got %d, want 500", resp.StatusCode)
	}
}
