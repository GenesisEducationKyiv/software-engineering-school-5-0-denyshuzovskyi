package handler

import (
	"context"
	"errors"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/error"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/httputil"
	"log/slog"
	"net/http"
)

type WeatherService interface {
	GetCurrentWeatherForLocation(context.Context, string) (*dto.WeatherDTO, error)
}

type WeatherHandler struct {
	weatherService WeatherService
	log            *slog.Logger
}

func NewWeatherHandler(weatherService WeatherService, log *slog.Logger) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
		log:            log,
	}
}

func (h *WeatherHandler) GetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var location = query.Get("city")
	if location == "" {
		http.Error(w, "Missing city parameter", http.StatusBadRequest)
		h.log.Info("no query parameter 'city' found")
		return
	}

	weatherDto, err := h.weatherService.GetCurrentWeatherForLocation(r.Context(), location)
	if err != nil {
		if errors.Is(err, commonerrors.ErrLocationNotFound) {
			http.Error(w, "City not found", http.StatusNotFound)
			h.log.Info("couldn't get weatherDto for provided location", "location", location)
			return
		}
		http.Error(w, "", http.StatusInternalServerError)
		h.log.Error("error getting weatherDto", "error", err)
		return
	}

	err = httputil.WriteJSON(w, weatherDto)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		h.log.Error("error writing response", "error", err)
		return
	}
}
