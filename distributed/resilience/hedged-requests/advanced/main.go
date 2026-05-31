package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	var attempt atomic.Int32
	val, err := Do(context.Background(), 20*time.Millisecond, 3, func(_ context.Context) (string, error) {
		n := attempt.Add(1)
		// First attempt is slow; hedged attempts are fast.
		if n == 1 {
			time.Sleep(100 * time.Millisecond)
		} else {
			time.Sleep(5 * time.Millisecond)
		}
		return fmt.Sprintf("served-by-attempt-%d", n), nil
	})
	fmt.Printf("result=%q err=%v\n", val, err)
}
