package weatherprovider

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Get(context.Context, string) *redis.StringCmd
	Set(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
}

type CachingWeatherProvider struct {
	client   RedisClient
	ttl      time.Duration
	provider WeatherProvider
	metrics  *prometheus.CounterVec
	log      *slog.Logger
}

func NewCachingWeatherProvider(client RedisClient, ttl time.Duration, provider WeatherProvider, metrics *prometheus.CounterVec, logger *slog.Logger) *CachingWeatherProvider {
	return &CachingWeatherProvider{
		client:   client,
		ttl:      ttl,
		provider: provider,
		metrics:  metrics,
		log:      logger,
	}
}

func (p *CachingWeatherProvider) GetCurrentWeather(ctx context.Context, location string) (*dto.WeatherWithLocationDTO, error) {
	key := "weather:" + location

	data, err := p.client.Get(ctx, key).Result()
	if err == nil {
		var weather dto.WeatherWithLocationDTO
		if err := json.Unmarshal([]byte(data), &weather); err == nil {
			p.metrics.WithLabelValues("hit").Inc()
			p.log.Debug("cache hit", "key", key)
			return &weather, nil
		}
		p.log.Warn("failed to unmarshal cached weather", "key", key, "error", err)
	} else if errors.Is(err, redis.Nil) {
		p.metrics.WithLabelValues("miss").Inc()
		p.log.Debug("cache miss", "key", key)
	} else {
		p.metrics.WithLabelValues("error").Inc()
		p.log.Error("redis error", "key", key, "error", err)
	}

	weather, err := p.provider.GetCurrentWeather(ctx, location)
	if err != nil {
		p.log.Error("provider error", "location", location, "error", err)
		return nil, err
	}

	serialized, err := json.Marshal(weather)
	if err != nil {
		p.log.Warn("failed to marshal weather for cache", "key", key, "error", err)
	} else {
		if cacheErr := p.client.Set(ctx, key, serialized, p.ttl).Err(); cacheErr != nil {
			p.log.Warn("failed to set cache", "key", key, "error", cacheErr)
		}
	}

	return weather, nil
}
