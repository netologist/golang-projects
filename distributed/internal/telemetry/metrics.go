// Package telemetry provides shared observability bootstrap helpers for the
// distributed pattern examples.
package telemetry

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewRegistry creates a fresh Prometheus registry (not the default global).
// Use a local registry in examples so patterns don't pollute each other.
func NewRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}

// Handler returns an HTTP handler for /metrics using the given registry.
func Handler(reg *prometheus.Registry) http.Handler {
	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
}
