package main

import (
	"fmt"
	"sync"
	"time"
)

func query(name string, latency time.Duration) int {
	time.Sleep(latency)
	return len(name)
}

func main() {
	backends := map[string]time.Duration{
		"alpha": 10 * time.Millisecond,
		"beta":  20 * time.Millisecond,
		"gamma": 15 * time.Millisecond,
	}

	var wg sync.WaitGroup
	results := make(chan int, len(backends))
	for name, lat := range backends {
		wg.Add(1)
		go func(n string, l time.Duration) {
			defer wg.Done()
			results <- query(n, l)
		}(name, lat)
	}
	wg.Wait()
	close(results)

	total := 0
	for r := range results {
		total += r
	}
	fmt.Println("aggregated:", total)
}
