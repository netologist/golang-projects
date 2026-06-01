package main

import (
	"fmt"
	"time"
)

const heartbeatInterval = 200 * time.Millisecond
const missThreshold = 500 * time.Millisecond

func startNode(id int, beat chan<- int, stop <-chan struct{}) {
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-time.After(heartbeatInterval):
				beat <- id
			}
		}
	}()
}

func main() {
	beat := make(chan int, 10)
	stops := make([]chan struct{}, 3)
	for i := range stops {
		stops[i] = make(chan struct{})
		startNode(i+1, beat, stops[i])
	}

	last := make(map[int]time.Time)
	stop := time.After(2 * time.Second)

	// stop node 2 after 800ms
	time.AfterFunc(800*time.Millisecond, func() { close(stops[1]) })

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

loop:
	for {
		select {
		case id := <-beat:
			last[id] = time.Now()
		case <-ticker.C:
			for id := 1; id <= 3; id++ {
				if t, ok := last[id]; ok && time.Since(t) > missThreshold {
					fmt.Printf("[monitor] node %d DOWN (missed heartbeat)\n", id)
					delete(last, id)
				}
			}
		case <-stop:
			break loop
		}
	}
	fmt.Println("monitor exiting")
}
