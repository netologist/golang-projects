package main

import "time"

// Retry calls fn up to maxAttempts times with a fixed delay between tries.
func Retry(maxAttempts int, delay time.Duration, fn func() error) error {
	var err error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if err = fn(); err == nil {
			return nil
		}
		if attempt < maxAttempts-1 {
			time.Sleep(delay)
		}
	}
	return err
}
