package main

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

// Config tunes retry behaviour.
type Config struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	Jitter       float64          // fraction of delay added as random jitter (0..1)
	Retryable    func(error) bool // nil means all errors are retryable
}

// DefaultConfig is a sensible production starting point.
var DefaultConfig = Config{
	MaxAttempts:  3,
	InitialDelay: 100 * time.Millisecond,
	MaxDelay:     30 * time.Second,
	Multiplier:   2.0,
	Jitter:       0.1,
	Retryable:    nil,
}

// Do retries fn according to cfg, honouring context cancellation.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	var lastErr error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("retry: cancelled before attempt %d: %w", attempt, err)
		}
		if lastErr = fn(); lastErr == nil {
			return nil
		}
		if cfg.Retryable != nil && !cfg.Retryable(lastErr) {
			return lastErr
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
			return fmt.Errorf("retry: %w", ctx.Err())
		}
	}
	return fmt.Errorf("retry: max attempts (%d) reached: %w", cfg.MaxAttempts, lastErr)
}
