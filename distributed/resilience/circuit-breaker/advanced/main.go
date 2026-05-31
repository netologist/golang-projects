package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	b := New(Config{FailureThreshold: 3, RecoveryTimeout: 80 * time.Millisecond, HalfOpenProbes: 1},
		prometheus.NewRegistry())
	down := errors.New("service unavailable")

	for i := 1; i <= 5; i++ {
		err := b.Execute(context.Background(), func() error { return down })
		fmt.Printf("call %d: state=%s err=%v\n", i, b.State(), err)
	}

	time.Sleep(100 * time.Millisecond)
	err := b.Execute(context.Background(), func() error { return nil })
	fmt.Printf("probe: state=%s err=%v\n", b.State(), err)
}
