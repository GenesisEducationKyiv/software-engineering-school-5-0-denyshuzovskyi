package server

import (
	"net/http"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/server/handler"
)

func InitMux(weatherHandler *handler.WeatherHandler, subscriptionHandler *handler.SubscriptionHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather", weatherHandler.GetCurrentWeather)
	mux.HandleFunc("POST /subscribe", subscriptionHandler.Subscribe)
	mux.HandleFunc("GET /confirm/{token}", subscriptionHandler.Confirm)
	mux.HandleFunc("GET /unsubscribe/{token}", subscriptionHandler.Unsubscribe)

	return mux
}
