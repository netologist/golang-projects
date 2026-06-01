package main

import (
	"context"
	"errors"
	"fmt"
)

func main() {
	ctx := context.Background()

	// 1. Primary succeeds
	f := Fallback[string]{
		Primary:  func(_ context.Context) (string, error) { return "live data", nil },
		Fallback: func(_ context.Context) (string, error) { return "cached data", nil },
	}
	v, _ := f.Execute(ctx)
	fmt.Printf("primary ok   → %s\n", v)

	// 2. Primary fails → fallback cache
	f2 := Fallback[string]{
		Primary:  func(_ context.Context) (string, error) { return "", errors.New("upstream down") },
		Fallback: func(_ context.Context) (string, error) { return "stale:alice", nil },
	}
	v2, _ := f2.Execute(ctx)
	fmt.Printf("primary fail → %s\n", v2)

	// 3. Chain: primary→replica→default
	fail := func(_ context.Context) (string, error) { return "", errors.New("fail") }
	defaultFn := func(_ context.Context) (string, error) { return "default", nil }
	v3, _ := Chain(ctx, fail, fail, defaultFn)
	fmt.Printf("chain        → %s\n", v3)
}
