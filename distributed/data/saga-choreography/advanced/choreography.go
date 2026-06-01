package main

import (
	"fmt"
	"sync"
)

// Event carries a typed topic and arbitrary payload.
type Event struct {
	Topic   string
	Payload map[string]any
}

// Handler processes an event and may return follow-up events.
type Handler func(e Event) []Event

// EventBus is a synchronous, in-process topic router.
type EventBus struct {
	mu   sync.RWMutex
	subs map[string][]Handler
}

// NewBus creates an empty EventBus.
func NewBus() *EventBus { return &EventBus{subs: map[string][]Handler{}} }

// Subscribe registers a handler for a topic.
func (b *EventBus) Subscribe(topic string, h Handler) {
	b.mu.Lock()
	b.subs[topic] = append(b.subs[topic], h)
	b.mu.Unlock()
}

// Publish dispatches an event; follow-up events are published recursively.
func (b *EventBus) Publish(e Event) {
	b.mu.RLock()
	handlers := append([]Handler{}, b.subs[e.Topic]...)
	b.mu.RUnlock()
	for _, h := range handlers {
		for _, next := range h(e) {
			b.Publish(next)
		}
	}
}

// StepResult captures outcome of a saga step.
type StepResult struct {
	Step    string
	Success bool
	Err     string
}

// ChoreographySaga wires services on a shared bus with compensation.
type ChoreographySaga struct {
	bus *EventBus
	log []StepResult
	mu  sync.Mutex
}

// NewChoreographySaga creates a saga wired to bus.
func NewChoreographySaga(bus *EventBus) *ChoreographySaga {
	return &ChoreographySaga{bus: bus}
}

// Record appends a step result to the saga log.
func (s *ChoreographySaga) Record(step string, ok bool, errMsg string) {
	s.mu.Lock()
	s.log = append(s.log, StepResult{Step: step, Success: ok, Err: errMsg})
	s.mu.Unlock()
}

// Log returns a copy of recorded results.
func (s *ChoreographySaga) Log() []StepResult {
	s.mu.Lock()
	out := append([]StepResult{}, s.log...)
	s.mu.Unlock()
	return out
}

// PrintLog prints the saga execution log.
func (s *ChoreographySaga) PrintLog() {
	for _, r := range s.Log() {
		status := "OK"
		if !r.Success {
			status = "FAIL: " + r.Err
		}
		fmt.Printf("  [saga-log] %-20s %s\n", r.Step, status)
	}
}
