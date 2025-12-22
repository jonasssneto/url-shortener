package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "http",
			Subsystem: "server",
			Name:      "request_duration_seconds",
			Help:      "HTTP request latency.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "http",
			Subsystem: "server",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	UrlsRedirected = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "http",
			Subsystem: "server",
			Name:      "urls_redirected_total",
			Help:      "Total number of URLs redirected.",
		},
		[]string{"status"},
	)

	UrlsCreated = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "http",
			Subsystem: "server",
			Name:      "urls_created_total",
			Help:      "Total number of URLs created.",
		},
		[]string{"status"},
	)
)

func Register() {
	prometheus.MustRegister(HttpDuration, HttpRequestsTotal, UrlsRedirected, UrlsCreated)
}
