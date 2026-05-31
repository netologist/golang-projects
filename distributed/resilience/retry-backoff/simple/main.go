package main

import (
	"errors"
	"fmt"
	"time"
)

func retry(maxAttempts int, base time.Duration, fn func() error) error {
	var err error
	for i := 0; i < maxAttempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if i < maxAttempts-1 {
			delay := base * (1 << i) // exponential
			fmt.Printf("attempt %d failed, backing off %s\n", i+1, delay)
			time.Sleep(delay)
		}
	}
	return err
}

func main() {
	n := 0
	err := retry(4, 10*time.Millisecond, func() error {
		n++
		if n < 3 {
			return errors.New("temporary")
		}
		return nil
	})
	fmt.Printf("done after %d attempts: %v\n", n, err)
}
