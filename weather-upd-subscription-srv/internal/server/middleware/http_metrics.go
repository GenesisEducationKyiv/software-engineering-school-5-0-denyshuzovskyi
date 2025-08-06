package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type HTTPMetrics interface {
	IncInFlight()
	DecInFlight()
	ObserveRequest(string, string, string, float64)
}

func HTTPMetricsMiddleware(metrics HTTPMetrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.IncInFlight()
		defer metrics.DecInFlight()

		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		duration := time.Since(start).Seconds()

		var url string
		if r.Pattern != "" {
			url = r.Pattern
		} else {
			url = r.URL.Path
		}

		metrics.ObserveRequest(r.Method, url, fmt.Sprintf("%d", rec.status), duration)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
