package main

import (
	"context"
	"fmt"
)

func main() {
	s := NewStore()
	ctx := context.Background()

	_, _ = s.SaveWithOutbox(ctx, "order:1", []byte("..."), "order", []byte(`{"type":"OrderCreated","id":1}`))
	_, _ = s.SaveWithOutbox(ctx, "order:2", []byte("..."), "order", []byte(`{"type":"OrderCreated","id":2}`))

	n, err := s.Relay(ctx, 10, func(e *OutboxEvent) error {
		fmt.Printf("publishing %s: %s\n", e.ID, e.Payload)
		return nil
	})
	fmt.Printf("relayed %d events, err=%v\n", n, err)
}
