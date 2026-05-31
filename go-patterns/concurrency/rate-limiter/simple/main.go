package main

import (
	"fmt"
	"time"
)

func main() {
	l := New(50 * time.Millisecond)
	defer l.Stop()

	start := time.Now()
	for i := 1; i <= 4; i++ {
		l.Wait()
		fmt.Printf("op %d at %v\n", i, time.Since(start).Round(10*time.Millisecond))
	}
}
