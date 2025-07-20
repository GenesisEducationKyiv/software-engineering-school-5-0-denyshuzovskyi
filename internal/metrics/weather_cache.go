package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	WeatherCacheRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "weather_cache_requests_total",
			Help: "Total weather cache requests partitioned by result (hit/miss/error)",
		},
		[]string{"result"},
	)
)

func Init() {
	prometheus.MustRegister(WeatherCacheRequests)
}
