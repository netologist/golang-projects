package main

import (
	"errors"
	"testing"
)

func TestResult_mapChain(t *testing.T) {
	r := Ok(5).
		Map(func(n int) int { return n * 2 }).
		Map(func(n int) int { return n + 1 })
	val, err := r.Unwrap()
	if err != nil {
		t.Fatal(err)
	}
	if val != 11 {
		t.Errorf("got %d, want 11", val)
	}
}

func TestResult_errShortCircuits(t *testing.T) {
	boom := errors.New("boom")
	r := Err[int](boom).Map(func(n int) int { return n * 100 })
	if _, err := r.Unwrap(); !errors.Is(err, boom) {
		t.Errorf("expected boom, got %v", err)
	}
}

func TestResult_flatMap(t *testing.T) {
	safeDiv := func(n int) Result[int] {
		if n == 0 {
			return Err[int](errors.New("div by zero"))
		}
		return Ok(100 / n)
	}
	r := Ok(5).FlatMap(safeDiv)
	if v, _ := r.Unwrap(); v != 20 {
		t.Errorf("got %d, want 20", v)
	}
}

func TestResult_or(t *testing.T) {
	r := Err[string](errors.New("nope"))
	if got := r.Or("default"); got != "default" {
		t.Errorf("got %s, want default", got)
	}
}
