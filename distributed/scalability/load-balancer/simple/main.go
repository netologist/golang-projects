package main

import "fmt"

type roundRobin struct {
	backends []string
	idx      int
}

func (rr *roundRobin) next() string {
	b := rr.backends[rr.idx%len(rr.backends)]
	rr.idx++
	return b
}

func main() {
	lb := &roundRobin{backends: []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}}
	for i := 0; i < 6; i++ {
		fmt.Println("route to", lb.next())
	}
}
