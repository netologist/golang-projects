package main

import (
	"fmt"
	"sync"
	"time"
)

type store struct {
	mu   sync.RWMutex
	data map[string]string
}

func (s *store) write(k, v string) { s.mu.Lock(); s.data[k] = v; s.mu.Unlock() }
func (s *store) read(k string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[k]
}

func main() {
	primary := &store{data: map[string]string{}}
	replicas := []*store{{data: map[string]string{}}, {data: map[string]string{}}}

	// Write to primary
	primary.write("name", "alice")
	fmt.Println("[primary] wrote name=alice")

	// Async replicate to replicas
	var wg sync.WaitGroup
	for i, r := range replicas {
		wg.Add(1)
		go func(i int, r *store) {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond) // simulate lag
			r.write("name", primary.read("name"))
			fmt.Printf("[replica-%d] replicated name=%s\n", i+1, r.read("name"))
		}(i, r)
	}
	wg.Wait()

	fmt.Printf("[replica-1] read name=%s\n", replicas[0].read("name"))
	fmt.Printf("[replica-2] read name=%s\n", replicas[1].read("name"))
}
