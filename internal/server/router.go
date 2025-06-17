package server

import (
	"net/http"
)

type weatherHandler interface {
	GetCurrentWeather(http.ResponseWriter, *http.Request)
}

type subscriptionHandler interface {
	Subscribe(http.ResponseWriter, *http.Request)
	Confirm(http.ResponseWriter, *http.Request)
	Unsubscribe(http.ResponseWriter, *http.Request)
}

func InitRouter(weatherHandler weatherHandler, subscriptionHandler subscriptionHandler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /weather", weatherHandler.GetCurrentWeather)
	router.HandleFunc("POST /subscribe", subscriptionHandler.Subscribe)
	router.HandleFunc("GET /confirm/{token}", subscriptionHandler.Confirm)
	router.HandleFunc("GET /unsubscribe/{token}", subscriptionHandler.Unsubscribe)

	return router
}
