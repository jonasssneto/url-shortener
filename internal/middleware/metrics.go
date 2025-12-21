package middleware

import (
	"main/internal/metrics"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

var ()

func Metrics() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			metrics.HttpDuration.WithLabelValues(
				r.Method,
				r.URL.Path,
				http.StatusText(ww.Status()),
			).Observe(time.Since(start).Seconds())

			metrics.HttpRequestsTotal.WithLabelValues(
				r.Method,
				r.URL.Path,
				http.StatusText(ww.Status()),
			).Inc()
		})
	}
}
