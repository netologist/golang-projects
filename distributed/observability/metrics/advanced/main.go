package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	reg := prometheus.NewRegistry()
	rec := NewRecorder(reg)

	handler := rec.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	srv := httptest.NewServer(handler)
	defer srv.Close()

	for i := 0; i < 5; i++ {
		resp, _ := http.Get(srv.URL)
		resp.Body.Close()
	}

	mfs, _ := reg.Gather()
	for _, mf := range mfs {
		if mf.GetName() == "http_requests_total" {
			fmt.Printf("%s = %v\n", mf.GetName(), mf.GetMetric()[0].GetCounter().GetValue())
		}
	}
}
