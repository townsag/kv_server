package middleware

import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

/*
metrics goals:
	- log the total number of requests served
		- counter
	- log the number of concurrently processed requests
		- gauge
	- log the latency of requests
		- histogram
*/

type Metrics struct {
	totalRequests prometheus.Counter
	concurrentRequests prometheus.Gauge
	requestLatency prometheus.Histogram
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	metrics := Metrics{
		totalRequests: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "total_requests",
				Help: "total number of requests received across all endpoints",
			},
		),
		concurrentRequests: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "concurrent_requests",
				Help: "number of concurrently served requests across all endpoints",
			},
		),
		requestLatency: prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name: "request_latency",
				Help: "average request latency across all endpoints",
			},
		),
	}

	reg.MustRegister(metrics.totalRequests)
	reg.MustRegister(metrics.concurrentRequests)
	reg.MustRegister(metrics.requestLatency)
	return &metrics
}

func MetricsMiddleware(metrics *Metrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.totalRequests.Add(1.0)
		metrics.concurrentRequests.Inc()
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()
		metrics.concurrentRequests.Dec()
		metrics.requestLatency.Observe(duration)
	})
}

func NewMetricsMiddleware(metrics *Metrics) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return MetricsMiddleware(metrics, h)
	}
}