package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	b := NewBroker()

	// A simple in-process responder routes replies by correlation ID.
	reply, err := b.Request(context.Background(), time.Second, func(id string) {
		go func() {
			time.Sleep(10 * time.Millisecond)
			b.HandleReply(id, []byte("result for "+id))
		}()
	})
	fmt.Printf("reply=%q err=%v\n", reply, err)
}
