package main

import "fmt"

func main() {
	bs := []*Backend{NewBackend("10.0.0.1"), NewBackend("10.0.0.2"), NewBackend("10.0.0.3")}

	fmt.Println("round-robin:")
	rr := &RoundRobin{}
	for i := 0; i < 5; i++ {
		b, _ := rr.Next(bs)
		fmt.Println("  ->", b.Addr)
	}

	fmt.Println("least-connections:")
	bs[0].ActiveConns.Store(5)
	bs[1].ActiveConns.Store(1)
	bs[2].ActiveConns.Store(3)
	lc := &LeastConns{}
	b, _ := lc.Next(bs)
	fmt.Println("  ->", b.Addr, "(fewest active conns)")
}
