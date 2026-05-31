package main

import (
	"context"
	"fmt"
	"time"
)

type OrderEvent struct {
	ID    string
	Total int
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	topic := NewTopic[OrderEvent]("orders")
	ch, unsub := topic.Subscribe(ctx, 8)
	defer unsub()

	_ = topic.Publish(ctx, OrderEvent{ID: "A-1", Total: 100})
	_ = topic.Publish(ctx, OrderEvent{ID: "A-2", Total: 250})

	time.Sleep(10 * time.Millisecond)
	for len(ch) > 0 {
		msg := <-ch
		fmt.Printf("received order %s total=%d\n", msg.Payload.ID, msg.Payload.Total)
	}
}
