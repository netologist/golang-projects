package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {
	b := New(3, 100*time.Millisecond)
	flaky := errors.New("downstream down")

	for i := 1; i <= 5; i++ {
		err := b.Execute(func() error { return flaky })
		fmt.Printf("call %d: %v\n", i, err)
	}

	time.Sleep(120 * time.Millisecond)
	err := b.Execute(func() error { return nil })
	fmt.Printf("after cooldown: %v\n", err)
}
