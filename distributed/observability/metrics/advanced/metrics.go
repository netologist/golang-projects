package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Recorder bundles RED-method instruments.
type Recorder struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	inFlight        *prometheus.GaugeVec
}

// NewRecorder registers RED instruments against reg.
func NewRecorder(reg prometheus.Registerer) *Recorder {
	r := &Recorder{
		requestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests.",
		}, []string{"method", "status"}),
		requestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method"}),
		inFlight: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "In-flight HTTP requests.",
		}, []string{"method"}),
	}
	reg.MustRegister(r.requestsTotal, r.requestDuration, r.inFlight)
	return r
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// Middleware auto-instruments an http.Handler with RED metrics.
func (r *Recorder) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r.inFlight.WithLabelValues(req.Method).Inc()
		defer r.inFlight.WithLabelValues(req.Method).Dec()

		start := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sr, req)

		r.requestDuration.WithLabelValues(req.Method).Observe(time.Since(start).Seconds())
		r.requestsTotal.WithLabelValues(req.Method, strconv.Itoa(sr.status)).Inc()
	})
}
