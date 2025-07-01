package weatherprovider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/redis/go-redis/v9"
	"time"
)

type CachingWeatherProvider struct {
	client   *redis.Client
	ttl      time.Duration
	provider WeatherProvider
}

func NewCachingWeatherProvider(client *redis.Client, ttl time.Duration, provider WeatherProvider) *CachingWeatherProvider {
	return &CachingWeatherProvider{
		client:   client,
		ttl:      ttl,
		provider: provider,
	}
}

func (p *CachingWeatherProvider) GetCurrentWeather(ctx context.Context, location string) (*dto.WeatherWithLocationDTO, error) {
	key := "weather:" + location

	data, err := p.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		fmt.Println("Cache miss for key:", key)
	} else if err != nil {
		fmt.Printf("Redis error: %v\n", err)
	} else {
		var weather dto.WeatherWithLocationDTO
		if jsonErr := json.Unmarshal([]byte(data), &weather); jsonErr == nil {
			return &weather, nil
		}
	}

	weather, err := p.provider.GetCurrentWeather(ctx, location)
	if err != nil {
		return nil, err
	}

	serialized, err := json.Marshal(weather)
	if err == nil {
		p.client.Set(ctx, key, serialized, p.ttl)
	}

	return weather, nil
}
