package weather

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
)

type WeatherProvider interface {
	GetCurrentWeather(context.Context, string) (*dto.WeatherWithLocationDTO, error)
}

type WeatherService struct {
	weatherProvider WeatherProvider
	log             *slog.Logger
}

func NewWeatherService(weatherProvider WeatherProvider, log *slog.Logger) *WeatherService {
	return &WeatherService{
		weatherProvider: weatherProvider,
		log:             log,
	}
}

func (s *WeatherService) GetCurrentWeatherForLocation(ctx context.Context, location string) (*dto.WeatherDTO, error) {
	weatherWithLocationDTO, err := s.weatherProvider.GetCurrentWeather(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("get current weather: %w", err)
	}

	return &weatherWithLocationDTO.Weather, nil
}
