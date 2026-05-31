package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func slowCall(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := slowCall(ctx, 200*time.Millisecond)
	fmt.Println("timed out:", errors.Is(err, context.DeadlineExceeded), err)
}
