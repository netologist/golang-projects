package main

import (
	"context"
	"fmt"
)

// State is the mutable saga state shared across steps.
type State map[string]any

// Step is one saga step with a compensating action.
type Step struct {
	Name       string
	Execute    func(ctx context.Context, state State) error
	Compensate func(ctx context.Context, state State) error
}

// Saga orchestrates steps with reverse-order compensation on failure.
type Saga struct{ steps []Step }

// New creates a saga from ordered steps.
func New(steps ...Step) *Saga { return &Saga{steps: steps} }

// Run executes steps forward; on failure compensates completed steps in reverse.
func (s *Saga) Run(ctx context.Context) error {
	state := State{}
	var executed []int

	for i, step := range s.steps {
		if err := step.Execute(ctx, state); err != nil {
			for j := len(executed) - 1; j >= 0; j-- {
				comp := s.steps[executed[j]]
				if comp.Compensate == nil {
					continue
				}
				if cerr := comp.Compensate(ctx, state); cerr != nil {
					return fmt.Errorf("saga: compensation of %q failed: %w (original: %v)", comp.Name, cerr, err)
				}
			}
			return fmt.Errorf("saga: step %q failed: %w", step.Name, err)
		}
		executed = append(executed, i)
	}
	return nil
}
