package server

import (
	"net/http"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/server/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type weatherHandler interface {
	GetCurrentWeather(http.ResponseWriter, *http.Request)
}

type subscriptionHandler interface {
	Subscribe(http.ResponseWriter, *http.Request)
	Confirm(http.ResponseWriter, *http.Request)
	Unsubscribe(http.ResponseWriter, *http.Request)
}

type httpMetrics interface {
	IncInFlight()
	DecInFlight()
	ObserveRequest(string, string, string, float64)
}

func InitMux(weatherHandler weatherHandler, subscriptionHandler subscriptionHandler, httpMetrics httpMetrics) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /weather",
		middleware.HTTPMetricsMiddleware(httpMetrics, http.HandlerFunc(weatherHandler.GetCurrentWeather)),
	)
	mux.Handle("POST /subscribe",
		middleware.HTTPMetricsMiddleware(httpMetrics, http.HandlerFunc(subscriptionHandler.Subscribe)),
	)
	mux.Handle("GET /confirm/{token}",
		middleware.HTTPMetricsMiddleware(httpMetrics, http.HandlerFunc(subscriptionHandler.Confirm)),
	)
	mux.Handle("GET /unsubscribe/{token}",
		middleware.HTTPMetricsMiddleware(httpMetrics, http.HandlerFunc(subscriptionHandler.Unsubscribe)),
	)

	mux.Handle("/metrics", promhttp.Handler())

	return mux
}
