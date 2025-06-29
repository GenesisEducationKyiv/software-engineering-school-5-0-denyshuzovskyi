package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/error"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/httputil"
)

type WeatherService interface {
	GetCurrentWeatherForLocation(context.Context, string) (*dto.WeatherDTO, error)
}

type LocationValidator interface {
	ValidateLocation(string) error
}

type WeatherHandler struct {
	weatherService    WeatherService
	locationValidator LocationValidator
	log               *slog.Logger
}

func NewWeatherHandler(weatherService WeatherService, locationValidator LocationValidator, log *slog.Logger) *WeatherHandler {
	return &WeatherHandler{
		weatherService:    weatherService,
		locationValidator: locationValidator,
		log:               log,
	}
}

func (h *WeatherHandler) GetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	var location = query.Get("city")
	if location == "" {
		http.Error(w, "Missing city parameter", http.StatusBadRequest)
		h.log.Error("no query parameter 'city' found")
		return
	}
	if err := h.locationValidator.ValidateLocation(location); err != nil {
		http.Error(w, "City validation failed", http.StatusBadRequest)
		h.log.Error("invalid city", "error", err)
		return
	}

	weatherDto, err := h.weatherService.GetCurrentWeatherForLocation(r.Context(), location)
	if err != nil {
		if errors.Is(err, commonerrors.ErrLocationNotFound) {
			http.Error(w, "City not found", http.StatusNotFound)
			h.log.Error("couldn't get weatherDto for provided location", "location", location, "error", err)
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
