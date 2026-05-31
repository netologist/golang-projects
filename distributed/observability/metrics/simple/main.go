package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	reg := prometheus.NewRegistry()
	requests := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total requests.",
	})
	reg.MustRegister(requests)

	for i := 0; i < 3; i++ {
		requests.Inc()
	}

	mfs, _ := reg.Gather()
	for _, mf := range mfs {
		fmt.Printf("%s = %v\n", mf.GetName(), mf.GetMetric()[0].GetCounter().GetValue())
	}
}
