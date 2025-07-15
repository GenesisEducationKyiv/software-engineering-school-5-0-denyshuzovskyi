package weatherprovider

import (
	"context"
	"log/slog"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/weather-upd-subscription-srv/internal/dto"
)

type LoggingWeatherProvider struct {
	name       string
	provider   WeatherProvider
	weatherLog *slog.Logger
	log        *slog.Logger
}

func NewLoggingWeatherProvider(name string, provider WeatherProvider, weatherLog *slog.Logger, log *slog.Logger) *LoggingWeatherProvider {
	return &LoggingWeatherProvider{
		name:       name,
		provider:   provider,
		weatherLog: weatherLog,
		log:        log,
	}
}

func (c *LoggingWeatherProvider) GetCurrentWeather(ctx context.Context, location string) (*dto.WeatherWithLocationDTO, error) {
	weatherWithLocationDTO, err := c.provider.GetCurrentWeather(ctx, location)
	if err != nil {
		return nil, err
	}
	c.weatherLog.Info("current weather", "provider", c.name, "weather", weatherWithLocationDTO)

	return weatherWithLocationDTO, err
}
