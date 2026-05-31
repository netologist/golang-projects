package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	r := New()
	defer r.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watch := r.Watch(ctx, "payments")

	dereg1 := r.Register(Instance{ID: "p1", Service: "payments", Addr: "10.0.0.1:80", TTL: time.Second})
	r.Register(Instance{ID: "p2", Service: "payments", Addr: "10.0.0.2:80", TTL: time.Second})

	time.Sleep(10 * time.Millisecond)
	insts, _ := r.Discover("payments")
	fmt.Printf("discovered %d payments instances\n", len(insts))

	dereg1()
	// Drain the latest watch update.
	time.Sleep(10 * time.Millisecond)
	for len(watch) > 0 {
		update := <-watch
		fmt.Printf("watch update: %d instances\n", len(update))
	}
}
