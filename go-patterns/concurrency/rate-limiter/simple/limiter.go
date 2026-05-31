package main

import "time"

// Limiter is a simple ticker-based limiter: one token per interval.
type Limiter struct {
	ticker *time.Ticker
}

// New allows one operation every interval.
func New(interval time.Duration) *Limiter {
	return &Limiter{ticker: time.NewTicker(interval)}
}

// Wait blocks until the next tick.
func (l *Limiter) Wait() { <-l.ticker.C }

// Stop releases the ticker.
func (l *Limiter) Stop() { l.ticker.Stop() }
