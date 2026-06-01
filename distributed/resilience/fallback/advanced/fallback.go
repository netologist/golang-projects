package main

import (
	"context"
	"errors"
	"fmt"
)

// ErrAllFailed is returned when every provider in a chain fails.
var ErrAllFailed = errors.New("fallback: all providers failed")

// Provider is a function that attempts to return a value.
type Provider[T any] func(ctx context.Context) (T, error)

// Fallback tries primary; on error invokes the fallback provider.
type Fallback[T any] struct {
	Primary  Provider[T]
	Fallback Provider[T]
}

// Execute calls Primary; if it fails, calls Fallback.
func (f Fallback[T]) Execute(ctx context.Context) (T, error) {
	if v, err := f.Primary(ctx); err == nil {
		return v, nil
	}
	return f.Fallback(ctx)
}

// Chain executes providers in order, returning the first success.
func Chain[T any](ctx context.Context, providers ...Provider[T]) (T, error) {
	var zero T
	var errs []error
	for _, p := range providers {
		v, err := p(ctx)
		if err == nil {
			return v, nil
		}
		errs = append(errs, err)
	}
	return zero, fmt.Errorf("%w: %v", ErrAllFailed, errs)
}
