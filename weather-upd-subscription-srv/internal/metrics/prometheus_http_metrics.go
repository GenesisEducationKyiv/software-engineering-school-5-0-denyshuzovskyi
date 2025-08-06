package metrics

import "github.com/prometheus/client_golang/prometheus"

type PrometheusHTTPMetrics struct {
	requestCount    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	inFlight        prometheus.Gauge
}

type HTTPOption func(*PrometheusHTTPMetrics)

func WithRequestCount(counter *prometheus.CounterVec) HTTPOption {
	return func(m *PrometheusHTTPMetrics) {
		m.requestCount = counter
	}
}

func WithRequestDuration(histogram *prometheus.HistogramVec) HTTPOption {
	return func(m *PrometheusHTTPMetrics) {
		m.requestDuration = histogram
	}
}

func WithInFlightGauge(gauge prometheus.Gauge) HTTPOption {
	return func(m *PrometheusHTTPMetrics) {
		m.inFlight = gauge
	}
}

func NewPrometheusHTTPMetrics(opts ...HTTPOption) *PrometheusHTTPMetrics {
	m := &PrometheusHTTPMetrics{
		requestCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
		inFlight: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_in_flight_requests",
				Help: "Current number of in-flight HTTP requests",
			},
		),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *PrometheusHTTPMetrics) Register() {
	prometheus.MustRegister(m.requestCount, m.requestDuration, m.inFlight)
}

func (m *PrometheusHTTPMetrics) IncInFlight() {
	m.inFlight.Inc()
}

func (m *PrometheusHTTPMetrics) DecInFlight() {
	m.inFlight.Dec()
}

func (m *PrometheusHTTPMetrics) ObserveRequest(method, path, status string, durationSeconds float64) {
	m.requestCount.WithLabelValues(method, path, status).Inc()
	m.requestDuration.WithLabelValues(method, path, status).Observe(durationSeconds)
}
