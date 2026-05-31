package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	l := NewMemoryLocker()
	ctx := context.Background()

	tok, err := l.Lock(ctx, "job:nightly", time.Second)
	fmt.Printf("acquired token=%s err=%v\n", tok, err)

	if err := l.Extend(ctx, "job:nightly", tok, 2*time.Second); err == nil {
		fmt.Println("lease extended")
	}

	if err := l.Unlock(ctx, "job:nightly", tok); err == nil {
		fmt.Println("released")
	}
}
