package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	fast, err := Do(context.Background(), 100*time.Millisecond, func(_ context.Context) (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "done", nil
	})
	fmt.Printf("fast: val=%q err=%v\n", fast, err)

	_, err = Do(context.Background(), 20*time.Millisecond, func(ctx context.Context) (string, error) {
		<-ctx.Done()
		return "", ctx.Err()
	})
	fmt.Printf("slow: err=%v\n", err)
}
