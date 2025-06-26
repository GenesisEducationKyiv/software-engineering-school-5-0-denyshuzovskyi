package weather

import (
	"context"
	"log/slog"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
)

type ChainWeatherProvider struct {
	providers []WeatherProvider
	log       *slog.Logger
}

func NewChainWeatherProvider(log *slog.Logger, providers ...WeatherProvider) *ChainWeatherProvider {
	return &ChainWeatherProvider{
		log:       log,
		providers: providers,
	}
}

func (c *ChainWeatherProvider) GetCurrentWeather(ctx context.Context, location string) (*dto.WeatherWithLocationDTO, error) {
	var lastErr error
	for i, provider := range c.providers {
		weatherWithLocationDTO, err := provider.GetCurrentWeather(ctx, location)
		if err == nil {
			return weatherWithLocationDTO, nil
		}
		c.log.Error("provider failed", "idx", i, "error", err)
		lastErr = err
	}
	return nil, lastErr
}
