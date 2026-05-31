package main

import (
	"context"
	"testing"
	"time"
)

func TestGather_allRespond(t *testing.T) {
	fns := map[string]func(context.Context) (int, error){
		"a": func(_ context.Context) (int, error) { return 1, nil },
		"b": func(_ context.Context) (int, error) { return 2, nil },
		"c": func(_ context.Context) (int, error) { return 3, nil },
	}
	results := Gather(context.Background(), fns)
	if len(results) != 3 {
		t.Errorf("got %d results, want 3", len(results))
	}
}

func TestGather_partialOnDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	fns := map[string]func(context.Context) (int, error){
		"fast": func(_ context.Context) (int, error) { return 1, nil },
		"slow": func(ctx context.Context) (int, error) {
			<-ctx.Done()
			return 0, ctx.Err()
		},
	}
	results := Gather(ctx, fns)
	if len(results) == 0 {
		t.Error("expected at least 1 partial result")
	}
}
