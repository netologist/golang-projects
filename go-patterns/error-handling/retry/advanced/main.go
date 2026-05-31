package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	attempt := 0
	cfg := Config{
		MaxAttempts:  5,
		InitialDelay: 20 * time.Millisecond,
		MaxDelay:     time.Second,
		Multiplier:   2.0,
		Jitter:       0.2,
		Retryable:    func(err error) bool { return errors.Is(err, errFlaky) },
	}
	err := Do(context.Background(), cfg, func() error {
		attempt++
		if attempt < 4 {
			fmt.Printf("attempt %d failed\n", attempt)
			return errFlaky
		}
		return nil
	})
	fmt.Printf("final: attempts=%d err=%v\n", attempt, err)
}

var errFlaky = errors.New("flaky upstream")
