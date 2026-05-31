package main

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

// Config tunes the retry behaviour.
type Config struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	Jitter       float64
	Retryable    func(error) bool
}

// DefaultConfig is a production-friendly default.
var DefaultConfig = Config{
	MaxAttempts:  4,
	InitialDelay: 100 * time.Millisecond,
	MaxDelay:     10 * time.Second,
	Multiplier:   2.0,
	Jitter:       0.2,
}

// Do retries fn and returns the typed result of the first success.
func Do[T any](ctx context.Context, cfg Config, fn func(context.Context) (T, error)) (T, error) {
	var zero T
	var lastErr error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return zero, fmt.Errorf("retry: cancelled before attempt %d: %w", attempt, err)
		}
		v, err := fn(ctx)
		if err == nil {
			return v, nil
		}
		lastErr = err
		if cfg.Retryable != nil && !cfg.Retryable(err) {
			return zero, err
		}
		if attempt == cfg.MaxAttempts-1 {
			break
		}
		delay := time.Duration(math.Min(
			float64(cfg.InitialDelay)*math.Pow(cfg.Multiplier, float64(attempt)),
			float64(cfg.MaxDelay),
		))
		jitter := time.Duration(float64(delay) * cfg.Jitter * rand.Float64())
		select {
		case <-time.After(delay + jitter):
		case <-ctx.Done():
			return zero, fmt.Errorf("retry: %w", ctx.Err())
		}
	}
	return zero, fmt.Errorf("retry: max attempts (%d) reached: %w", cfg.MaxAttempts, lastErr)
}
