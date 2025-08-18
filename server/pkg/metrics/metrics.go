package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_requests_total",
			Help: "Total number of requests",
		},
		[]string{"service", "handler", "method", "code"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_request_duration_seconds",
			Help:    "Request duration seconds",
			Buckets: prometheus.DefBuckets, // you can customize
		},
		[]string{"service", "handler", "method"},
	)

	ErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_errors_total",
			Help: "Total number of errors",
		},
		[]string{"service", "handler", "code"},
	)

	InflightRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_inflight_requests",
			Help: "Current inflight requests",
		},
		[]string{"service", "handler"},
	)
)

func Register() {
	prometheus.MustRegister(RequestsTotal, RequestDuration, ErrorsTotal, InflightRequests)
}
