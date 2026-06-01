package main

import (
	"fmt"
	"time"
)

func runNode(id string, m *Monitor, stop <-chan struct{}) {
	go func() {
		ticker := time.NewTicker(150 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				m.Beat(id)
			}
		}
	}()
}

func main() {
	m := NewMonitor(
		400*time.Millisecond,
		2,
		func(id string, old, new NodeStatus) {
			fmt.Printf("[monitor] %s: %s → %s\n", id, old, new)
		},
	)

	nodes := []string{"node-1", "node-2", "node-3"}
	stops := make([]chan struct{}, len(nodes))
	for i, n := range nodes {
		m.Register(n)
		stops[i] = make(chan struct{})
		runNode(n, m, stops[i])
	}

	m.Start(100 * time.Millisecond)
	defer m.Stop()

	// Kill node-2 after 1s
	time.AfterFunc(1*time.Second, func() {
		fmt.Println("[demo] stopping node-2")
		close(stops[1])
	})

	time.Sleep(2500 * time.Millisecond)
	for _, n := range nodes {
		fmt.Printf("%-10s alive=%v status=%s\n", n, m.IsAlive(n), m.Status(n))
	}
}
