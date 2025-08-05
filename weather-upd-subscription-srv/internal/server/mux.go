package server

import (
	"net/http"

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

func InitMux(weatherHandler weatherHandler, subscriptionHandler subscriptionHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather", weatherHandler.GetCurrentWeather)
	mux.HandleFunc("POST /subscribe", subscriptionHandler.Subscribe)
	mux.HandleFunc("GET /confirm/{token}", subscriptionHandler.Confirm)
	mux.HandleFunc("GET /unsubscribe/{token}", subscriptionHandler.Unsubscribe)

	mux.Handle("/metrics", promhttp.Handler())

	return mux
}
