package server

import (
	"net/http"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/server/handler"
)

func InitRouter(weatherHandler *handler.WeatherHandler, subscriptionHandler *handler.SubscriptionHandler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /weather", weatherHandler.GetCurrentWeather)
	router.HandleFunc("POST /subscribe", subscriptionHandler.Subscribe)
	router.HandleFunc("GET /confirm/{token}", subscriptionHandler.Confirm)
	router.HandleFunc("GET /unsubscribe/{token}", subscriptionHandler.Unsubscribe)

	return router
}
