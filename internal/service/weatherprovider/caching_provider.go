package weatherprovider

import (
	"context"
	"errors"
	"log/slog"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/error"
)

type Cache interface {
	Set(context.Context, string, dto.WeatherWithLocationDTO) error
	Get(context.Context, string) (dto.WeatherWithLocationDTO, error)
}

type CacheMetrics interface {
	RecordCacheHit()
	RecordCacheMiss()
	RecordCacheError()
}

type CachingWeatherProvider struct {
	cache    Cache
	provider WeatherProvider
	metrics  CacheMetrics
	log      *slog.Logger
}

func NewCachingWeatherProvider(cache Cache, provider WeatherProvider, metrics CacheMetrics, log *slog.Logger) *CachingWeatherProvider {
	return &CachingWeatherProvider{
		cache:    cache,
		provider: provider,
		metrics:  metrics,
		log:      log,
	}
}

func (p *CachingWeatherProvider) GetCurrentWeather(ctx context.Context, location string) (*dto.WeatherWithLocationDTO, error) {
	key := "weather:" + location

	weatherWithLocation, err := p.cache.Get(ctx, key)
	if err == nil {
		p.metrics.RecordCacheHit()
		p.log.Debug("cache hit", "key", key)
		return &weatherWithLocation, nil
	} else if errors.Is(err, commonerrors.ErrCacheMiss) {
		p.metrics.RecordCacheMiss()
		p.log.Debug("cache miss", "key", key)
	} else {
		p.metrics.RecordCacheError()
		p.log.Error("cache error", "key", key, "error", err)
	}

	weatherWithLocationNew, err := p.provider.GetCurrentWeather(ctx, location)
	if err != nil {
		p.log.Error("provider error", "location", location, "error", err)
		return nil, err
	}

	if err = p.cache.Set(ctx, key, *weatherWithLocationNew); err != nil {
		p.log.Warn("failed to set cache", "key", key, "error", err)
	}

	return weatherWithLocationNew, nil
}
