package main

import (
	"context"
	"errors"
	"fmt"
)

func main() {
	n := 0
	v, err := Do(context.Background(), DefaultConfig, func(_ context.Context) (int, error) {
		n++
		if n < 3 {
			return 0, errors.New("flaky")
		}
		return 99, nil
	})
	fmt.Printf("result=%d attempts=%d err=%v\n", v, n, err)
}
