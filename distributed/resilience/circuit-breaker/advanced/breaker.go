package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ErrOpen is returned when the breaker is open.
var ErrOpen = errors.New("circuit open")

// State is the breaker state.
type State int32

// Breaker states.
const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	}
	return "unknown"
}

// Config tunes the breaker.
type Config struct {
	FailureThreshold int
	RecoveryTimeout  time.Duration
	HalfOpenProbes   int
}

// DefaultConfig is a reasonable starting point.
var DefaultConfig = Config{
	FailureThreshold: 5,
	RecoveryTimeout:  10 * time.Second,
	HalfOpenProbes:   1,
}

// Breaker is a production-grade circuit breaker with Prometheus metrics.
type Breaker struct {
	cfg        Config
	mu         sync.Mutex
	state      State
	failures   int
	probes     int
	openedAt   time.Time
	tripsTotal prometheus.Counter
	stateGauge *prometheus.GaugeVec
}

// New creates a breaker registered against reg.
func New(cfg Config, reg prometheus.Registerer) *Breaker {
	tripsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "circuit_breaker_trips_total",
		Help: "Total number of times the breaker tripped to open.",
	})
	stateGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "circuit_breaker_state",
		Help: "Current breaker state (1 = active).",
	}, []string{"state"})
	reg.MustRegister(tripsTotal, stateGauge)
	stateGauge.WithLabelValues("closed").Set(1)

	return &Breaker{cfg: cfg, tripsTotal: tripsTotal, stateGauge: stateGauge}
}

// Execute runs fn unless the breaker is open, recording the outcome.
func (b *Breaker) Execute(_ context.Context, fn func() error) error {
	if err := b.allow(); err != nil {
		return err
	}
	err := fn()
	b.record(err)
	return err
}

func (b *Breaker) allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateClosed:
		return nil
	case StateOpen:
		if time.Since(b.openedAt) < b.cfg.RecoveryTimeout {
			return fmt.Errorf("execute: %w", ErrOpen)
		}
		b.transition(StateHalfOpen)
		return nil
	case StateHalfOpen:
		return nil
	}
	return nil
}

func (b *Breaker) record(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err != nil {
		b.failures++
		if b.state == StateHalfOpen || b.failures >= b.cfg.FailureThreshold {
			b.transition(StateOpen)
			b.openedAt = time.Now()
			b.tripsTotal.Inc()
		}
		return
	}
	if b.state == StateHalfOpen {
		b.probes++
		if b.probes >= b.cfg.HalfOpenProbes {
			b.transition(StateClosed)
		}
		return
	}
	b.failures = 0
}

func (b *Breaker) transition(s State) {
	b.stateGauge.WithLabelValues(b.state.String()).Set(0)
	b.state = s
	b.probes = 0
	b.failures = 0
	b.stateGauge.WithLabelValues(s.String()).Set(1)
}

// State returns the current breaker state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
