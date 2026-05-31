package main

import (
	"context"
	"errors"
	"testing"
)

func TestRunStage_happyPath(t *testing.T) {
	in := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		in <- i
	}
	close(in)

	out, errs := RunStage(context.Background(), in, 2, func(_ context.Context, n int) (int, error) {
		return n * 2, nil
	})

	var results []int
	for v := range out {
		results = append(results, v)
	}
	if len(results) != 5 {
		t.Errorf("got %d results, want 5", len(results))
	}
	for e := range errs {
		t.Errorf("unexpected error: %v", e.Err)
	}
}

func TestRunStage_errorIsolation(t *testing.T) {
	in := make(chan int, 3)
	in <- 1
	in <- 2
	in <- 3
	close(in)

	errBoom := errors.New("boom")
	out, errs := RunStage(context.Background(), in, 1, func(_ context.Context, n int) (int, error) {
		if n == 2 {
			return 0, errBoom
		}
		return n, nil
	})

	for range out {
	}
	var errCount int
	for e := range errs {
		if !errors.Is(e.Err, errBoom) {
			t.Errorf("unexpected: %v", e.Err)
		}
		errCount++
	}
	if errCount != 1 {
		t.Errorf("got %d errors, want 1", errCount)
	}
}
