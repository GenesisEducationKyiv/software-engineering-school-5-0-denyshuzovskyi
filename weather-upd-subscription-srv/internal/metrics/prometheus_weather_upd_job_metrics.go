package metrics

import "github.com/prometheus/client_golang/prometheus"

type PrometheusWeatherUpdJobMetrics struct {
	jobRuns              *prometheus.CounterVec
	notificationsSent    *prometheus.CounterVec
	notificationsFailed  *prometheus.CounterVec
	subscriptionsHandled *prometheus.CounterVec
	jobDuration          *prometheus.HistogramVec
}

type WeatherUpdJobMetricOption func(*PrometheusWeatherUpdJobMetrics)

func WithJobRunsCounterVec(counter *prometheus.CounterVec) WeatherUpdJobMetricOption {
	return func(m *PrometheusWeatherUpdJobMetrics) {
		m.jobRuns = counter
	}
}

func WithNotificationsSentCounterVec(counter *prometheus.CounterVec) WeatherUpdJobMetricOption {
	return func(m *PrometheusWeatherUpdJobMetrics) {
		m.notificationsSent = counter
	}
}

func WithNotificationsFailedCounterVec(counter *prometheus.CounterVec) WeatherUpdJobMetricOption {
	return func(m *PrometheusWeatherUpdJobMetrics) {
		m.notificationsFailed = counter
	}
}

func WithSubscriptionsHandledCounterVec(counter *prometheus.CounterVec) WeatherUpdJobMetricOption {
	return func(m *PrometheusWeatherUpdJobMetrics) {
		m.subscriptionsHandled = counter
	}
}

func WithJobDurationHistogramVec(histogram *prometheus.HistogramVec) WeatherUpdJobMetricOption {
	return func(m *PrometheusWeatherUpdJobMetrics) {
		m.jobDuration = histogram
	}
}

func NewPrometheusWeatherUpdJobMetrics(opts ...WeatherUpdJobMetricOption) *PrometheusWeatherUpdJobMetrics {
	m := &PrometheusWeatherUpdJobMetrics{
		jobRuns: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "weather_update_job_runs_total",
			Help: "Total number of weather update jobs executed",
		}, []string{"frequency"}),

		notificationsSent: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "weather_notifications_sent_total",
			Help: "Total number of weather notifications successfully sent",
		}, []string{"frequency"}),

		notificationsFailed: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "weather_notifications_failed_total",
			Help: "Total number of weather notifications that failed to send",
		}, []string{"frequency"}),

		subscriptionsHandled: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "weather_subscriptions_handled_total",
			Help: "Total number of subscriptions processed per weather update job",
		}, []string{"frequency"}),

		jobDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "weather_update_job_duration_seconds",
			Help:    "Duration of weather update jobs in seconds",
			Buckets: prometheus.DefBuckets,
		}, []string{"frequency"}),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *PrometheusWeatherUpdJobMetrics) Register() {
	prometheus.MustRegister(
		m.jobRuns,
		m.notificationsSent,
		m.notificationsFailed,
		m.subscriptionsHandled,
		m.jobDuration,
	)
}

func (m *PrometheusWeatherUpdJobMetrics) Init(freqs []string) {
	for _, f := range freqs {
		m.jobRuns.WithLabelValues(f).Add(0)
		m.notificationsSent.WithLabelValues(f).Add(0)
		m.notificationsFailed.WithLabelValues(f).Add(0)
		m.subscriptionsHandled.WithLabelValues(f).Add(0)
		m.jobDuration.WithLabelValues(f).Observe(0)
	}
}

func (m *PrometheusWeatherUpdJobMetrics) RecordJobRun(freq string) {
	m.jobRuns.WithLabelValues(freq).Inc()
}

func (m *PrometheusWeatherUpdJobMetrics) RecordNotificationSent(freq string) {
	m.notificationsSent.WithLabelValues(freq).Inc()
}

func (m *PrometheusWeatherUpdJobMetrics) RecordNotificationFailed(freq string) {
	m.notificationsFailed.WithLabelValues(freq).Inc()
}

func (m *PrometheusWeatherUpdJobMetrics) RecordSubscriptionsHandled(freq string, count int) {
	m.subscriptionsHandled.WithLabelValues(freq).Add(float64(count))
}

func (m *PrometheusWeatherUpdJobMetrics) ObserveJobDuration(freq string, durationSeconds float64) {
	m.jobDuration.WithLabelValues(freq).Observe(durationSeconds)
}
