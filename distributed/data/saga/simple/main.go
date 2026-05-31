package main

import (
	"errors"
	"fmt"
)

type step struct {
	name       string
	execute    func() error
	compensate func()
}

func runSaga(steps []step) error {
	var completed []step
	for _, s := range steps {
		if err := s.execute(); err != nil {
			for i := len(completed) - 1; i >= 0; i-- {
				completed[i].compensate()
			}
			return fmt.Errorf("saga failed at %q: %w", s.name, err)
		}
		completed = append(completed, s)
	}
	return nil
}

func main() {
	steps := []step{
		{"reserve", func() error { fmt.Println("reserve stock"); return nil }, func() { fmt.Println("release stock") }},
		{"charge", func() error { fmt.Println("charge card"); return errors.New("declined") }, func() { fmt.Println("refund card") }},
		{"ship", func() error { fmt.Println("ship order"); return nil }, func() { fmt.Println("cancel shipment") }},
	}
	fmt.Println("result:", runSaga(steps))
}
