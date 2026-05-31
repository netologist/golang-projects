package main

import (
	"fmt"
	"sync"
)

type registry struct {
	mu        sync.RWMutex
	instances map[string][]string
}

func newRegistry() *registry { return &registry{instances: map[string][]string{}} }

func (r *registry) register(service, addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.instances[service] = append(r.instances[service], addr)
}

func (r *registry) discover(service string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.instances[service]
}

func main() {
	reg := newRegistry()
	reg.register("payments", "10.0.0.1:8080")
	reg.register("payments", "10.0.0.2:8080")

	fmt.Println("payments instances:", reg.discover("payments"))
}
