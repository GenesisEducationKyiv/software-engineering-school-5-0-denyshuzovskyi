//go:build !integration

package weatherprovider

import (
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/logger/noophandler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCachingWeatherProvider_CacheMissThenHit(t *testing.T) {
	redisMock := NewMockRedisClient(t)
	providerMock := NewMockWeatherProvider(t)
	metrics := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "weather_cache_requests_total",
			Help: "Total weather cache requests partitioned by result (hit/miss/error).",
		},
		[]string{"result"},
	)
	reg := prometheus.NewRegistry()
	require.NoError(t, reg.Register(metrics))
	log := slog.New(noophandler.NewNoOpHandler())

	cp := NewCachingWeatherProvider(redisMock, 15*time.Minute, providerMock, metrics, log)

	location := "Kyiv"
	key := "weather:" + location

	expectedWeather := &dto.WeatherWithLocationDTO{
		Weather: dto.WeatherDTO{
			Temperature: float32(10),
			Humidity:    float32(10),
			Description: "Clear",
		},
		Location: dto.Location{
			Name: location,
		},
		LastUpdated: time.Now().Unix(),
	}
	weatherJSON, err := json.Marshal(expectedWeather)
	require.NoError(t, err)

	redisMock.EXPECT().Get(mock.Anything, key).Return(redis.NewStringResult("", redis.Nil)).Once()
	providerMock.EXPECT().GetCurrentWeather(mock.Anything, location).Return(expectedWeather, nil).Once()
	redisMock.EXPECT().Set(mock.Anything, key, mock.Anything, 15*time.Minute).Return(redis.NewStatusResult("OK", nil)).Once()

	// First call
	result1, err1 := cp.GetCurrentWeather(t.Context(), location)
	require.NoError(t, err1)
	assert.Equal(t, expectedWeather, result1)
	delta := 0.01
	require.InDelta(t, float64(1), testutil.ToFloat64(metrics.WithLabelValues("miss")), delta)
	require.InDelta(t, float64(0), testutil.ToFloat64(metrics.WithLabelValues("hit")), delta)

	redisMock.EXPECT().Get(mock.Anything, key).Return(redis.NewStringResult(string(weatherJSON), nil)).Once()

	// Second call
	result2, err2 := cp.GetCurrentWeather(t.Context(), location)
	require.NoError(t, err2)
	assert.Equal(t, expectedWeather, result2)
	require.InDelta(t, float64(1), testutil.ToFloat64(metrics.WithLabelValues("miss")), delta)
	require.InDelta(t, float64(1), testutil.ToFloat64(metrics.WithLabelValues("hit")), delta)
}
