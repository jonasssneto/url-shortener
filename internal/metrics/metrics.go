package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HttpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "http",
			Subsystem: "server",
			Name:      "app_request_duration_seconds",
			Help:      "HTTP request latency.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)

func Register() {
	prometheus.MustRegister(HttpDuration)
}
