package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func counterValue(t *testing.T, reg *prometheus.Registry, name string) float64 {
	t.Helper()
	mfs, err := reg.Gather()
	if err != nil {
		t.Fatal(err)
	}
	var total float64
	for _, mf := range mfs {
		if mf.GetName() != name {
			continue
		}
		for _, m := range mf.GetMetric() {
			if c := m.GetCounter(); c != nil {
				total += c.GetValue()
			}
		}
	}
	return total
}

func TestRecorder_countsRequests(t *testing.T) {
	reg := prometheus.NewRegistry()
	rec := NewRecorder(reg)

	handler := rec.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		handler.ServeHTTP(httptest.NewRecorder(), req)
	}

	if got := counterValue(t, reg, "http_requests_total"); got != 3 {
		t.Errorf("requests_total: got %v, want 3", got)
	}
}
