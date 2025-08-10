package metrics

import "github.com/prometheus/client_golang/prometheus"

type PrometheusEmailMetrics struct {
	emailSent   *prometheus.CounterVec
	emailFailed *prometheus.CounterVec
}

type EmailMetricOption func(*PrometheusEmailMetrics)

func WithEmailSentCounterVec(counter *prometheus.CounterVec) EmailMetricOption {
	return func(n *PrometheusEmailMetrics) {
		n.emailSent = counter
	}
}

func WithEmailFailedCounterVec(counter *prometheus.CounterVec) EmailMetricOption {
	return func(n *PrometheusEmailMetrics) {
		n.emailFailed = counter
	}
}

func NewPrometheusEmailMetrics(opts ...EmailMetricOption) *PrometheusEmailMetrics {
	n := &PrometheusEmailMetrics{
		emailSent: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "email_sent_total",
			Help: "Total number of emails sent successfully",
		}, []string{"type"}),

		emailFailed: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "email_failed_total",
			Help: "Total number of emails that failed to send",
		}, []string{"type"}),
	}

	for _, opt := range opts {
		opt(n)
	}

	return n
}

func (n *PrometheusEmailMetrics) Register() {
	prometheus.MustRegister(n.emailSent, n.emailFailed)
}

func (n *PrometheusEmailMetrics) Init(types []string) {
	for _, t := range types {
		n.emailSent.WithLabelValues(t).Add(0)
		n.emailFailed.WithLabelValues(t).Add(0)
	}
}

func (n *PrometheusEmailMetrics) RecordEmailSent(emailType string) {
	n.emailSent.WithLabelValues(emailType).Inc()
}

func (n *PrometheusEmailMetrics) RecordEmailFailed(emailType string) {
	n.emailFailed.WithLabelValues(emailType).Inc()
}
