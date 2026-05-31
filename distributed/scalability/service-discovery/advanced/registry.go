package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Instance is a registered service instance.
type Instance struct {
	ID      string
	Service string
	Addr    string
	TTL     time.Duration
	expiry  time.Time
}

// Registry is an in-process service registry with TTL and watches.
type Registry struct {
	mu        sync.RWMutex
	services  map[string]map[string]*Instance
	watchers  map[string][]chan []Instance
	reapEvery time.Duration
	stop      chan struct{}
}

// New creates a registry and starts its reaper.
func New() *Registry {
	r := &Registry{
		services:  map[string]map[string]*Instance{},
		watchers:  map[string][]chan []Instance{},
		reapEvery: 50 * time.Millisecond,
		stop:      make(chan struct{}),
	}
	go r.reaper()
	return r
}

// Close stops the reaper.
func (r *Registry) Close() { close(r.stop) }

// Register adds an instance and returns a deregister func.
func (r *Registry) Register(inst Instance) func() {
	inst.expiry = time.Now().Add(inst.TTL)
	r.mu.Lock()
	if r.services[inst.Service] == nil {
		r.services[inst.Service] = map[string]*Instance{}
	}
	r.services[inst.Service][inst.ID] = &inst
	r.mu.Unlock()
	r.notify(inst.Service)

	return func() {
		r.mu.Lock()
		delete(r.services[inst.Service], inst.ID)
		r.mu.Unlock()
		r.notify(inst.Service)
	}
}

// Heartbeat renews an instance's TTL.
func (r *Registry) Heartbeat(service, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	svc, ok := r.services[service]
	if !ok {
		return fmt.Errorf("service %q not found", service)
	}
	inst, ok := svc[id]
	if !ok {
		return fmt.Errorf("instance %q not found", id)
	}
	inst.expiry = time.Now().Add(inst.TTL)
	return nil
}

// Discover returns the current instances of a service.
func (r *Registry) Discover(service string) ([]Instance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	svc, ok := r.services[service]
	if !ok {
		return nil, fmt.Errorf("service %q not found", service)
	}
	out := make([]Instance, 0, len(svc))
	for _, inst := range svc {
		out = append(out, *inst)
	}
	return out, nil
}

// Watch returns a channel that receives the instance list on each change.
func (r *Registry) Watch(ctx context.Context, service string) <-chan []Instance {
	ch := make(chan []Instance, 1)
	r.mu.Lock()
	r.watchers[service] = append(r.watchers[service], ch)
	r.mu.Unlock()
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		ws := r.watchers[service]
		for i, w := range ws {
			if w == ch {
				r.watchers[service] = append(ws[:i], ws[i+1:]...)
				break
			}
		}
		r.mu.Unlock()
	}()
	return ch
}

func (r *Registry) notify(service string) {
	instances, _ := r.Discover(service)
	r.mu.RLock()
	watchers := append([]chan []Instance(nil), r.watchers[service]...)
	r.mu.RUnlock()
	for _, ch := range watchers {
		select {
		case ch <- instances:
		default:
		}
	}
}

func (r *Registry) reaper() {
	ticker := time.NewTicker(r.reapEvery)
	defer ticker.Stop()
	for {
		select {
		case <-r.stop:
			return
		case <-ticker.C:
			r.mu.Lock()
			changed := map[string]bool{}
			for svc, instances := range r.services {
				for id, inst := range instances {
					if time.Now().After(inst.expiry) {
						delete(instances, id)
						changed[svc] = true
					}
				}
			}
			r.mu.Unlock()
			for svc := range changed {
				r.notify(svc)
			}
		}
	}
}
