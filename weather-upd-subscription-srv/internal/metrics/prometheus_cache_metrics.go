package metrics

import "github.com/prometheus/client_golang/prometheus"

type PrometheusCacheMetrics struct {
	cacheHits   prometheus.Counter
	cacheMisses prometheus.Counter
	cacheErrors prometheus.Counter
}

type CacheMetricOption func(*PrometheusCacheMetrics)

func WithCacheHitsCounter(counter prometheus.Counter) CacheMetricOption {
	return func(p *PrometheusCacheMetrics) {
		p.cacheHits = counter
	}
}

func WithCacheMissesCounter(counter prometheus.Counter) CacheMetricOption {
	return func(p *PrometheusCacheMetrics) {
		p.cacheMisses = counter
	}
}

func WithCacheErrorsCounter(counter prometheus.Counter) CacheMetricOption {
	return func(p *PrometheusCacheMetrics) {
		p.cacheErrors = counter
	}
}

func NewPrometheusCacheMetrics(opts ...CacheMetricOption) *PrometheusCacheMetrics {
	p := &PrometheusCacheMetrics{
		cacheHits: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		}),
		cacheMisses: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		}),
		cacheErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "cache_errors_total",
			Help: "Total number of cache errors",
		}),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *PrometheusCacheMetrics) Register() {
	prometheus.MustRegister(p.cacheHits, p.cacheMisses, p.cacheErrors)
}

func (p *PrometheusCacheMetrics) RecordCacheHit() {
	p.cacheHits.Inc()
}

func (p *PrometheusCacheMetrics) RecordCacheMiss() {
	p.cacheMisses.Inc()
}
func (p *PrometheusCacheMetrics) RecordCacheError() {
	p.cacheErrors.Inc()
}
