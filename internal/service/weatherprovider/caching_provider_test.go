//go:build !integration

package weatherprovider

import (
	"log/slog"
	"testing"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/error"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/logger/noophandler"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCachingWeatherProvider_CacheMissThenHit(t *testing.T) {
	cacheMock := NewMockCache(t)
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

	cp := NewCachingWeatherProvider(cacheMock, providerMock, metrics, log)

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

	cacheMock.EXPECT().Get(mock.Anything, key).Return(dto.WeatherWithLocationDTO{}, commonerrors.ErrCacheMiss).Once()
	providerMock.EXPECT().GetCurrentWeather(mock.Anything, location).Return(expectedWeather, nil).Once()
	cacheMock.EXPECT().Set(mock.Anything, key, mock.Anything).Return(nil).Once()

	// First call
	result1, err1 := cp.GetCurrentWeather(t.Context(), location)
	require.NoError(t, err1)
	assert.Equal(t, expectedWeather, result1)
	delta := 0.01
	require.InDelta(t, float64(1), testutil.ToFloat64(metrics.WithLabelValues("miss")), delta)
	require.InDelta(t, float64(0), testutil.ToFloat64(metrics.WithLabelValues("hit")), delta)

	cacheMock.EXPECT().Get(mock.Anything, key).Return(*expectedWeather, nil).Once()

	// Second call
	result2, err2 := cp.GetCurrentWeather(t.Context(), location)
	require.NoError(t, err2)
	assert.Equal(t, expectedWeather, result2)
	require.InDelta(t, float64(1), testutil.ToFloat64(metrics.WithLabelValues("miss")), delta)
	require.InDelta(t, float64(1), testutil.ToFloat64(metrics.WithLabelValues("hit")), delta)
}
