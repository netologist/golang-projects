package main

import (
	"context"
	"errors"
	"fmt"
)

func main() {
	s := New(
		Step{
			Name: "reserve-stock",
			Execute: func(_ context.Context, st State) error {
				st["stock"] = "reserved"
				fmt.Println("reserved stock")
				return nil
			},
			Compensate: func(_ context.Context, _ State) error { fmt.Println("released stock"); return nil },
		},
		Step{
			Name: "charge-payment",
			Execute: func(_ context.Context, _ State) error {
				fmt.Println("charging payment...")
				return errors.New("card declined")
			},
			Compensate: func(_ context.Context, _ State) error { fmt.Println("refunded payment"); return nil },
		},
	)
	fmt.Println("saga result:", s.Run(context.Background()))
}
