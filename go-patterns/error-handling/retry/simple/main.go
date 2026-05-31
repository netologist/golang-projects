package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {
	attempt := 0
	err := Retry(3, 10*time.Millisecond, func() error {
		attempt++
		if attempt < 3 {
			return errors.New("temporary failure")
		}
		return nil
	})
	fmt.Printf("succeeded after %d attempts, err=%v\n", attempt, err)
}
