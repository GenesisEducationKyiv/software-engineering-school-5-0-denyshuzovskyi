//go:build !integration

package weatherprovider

import (
	"log/slog"
	"testing"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	commonerrors "github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/error"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/logger/noophandler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCachingWeatherProvider_CacheMissThenHit(t *testing.T) {
	cacheMock := NewMockCache(t)
	providerMock := NewMockWeatherProvider(t)

	hitc := prometheus.NewCounter(prometheus.CounterOpts{Name: "test_hit", Help: ""})
	missc := prometheus.NewCounter(prometheus.CounterOpts{Name: "test_miss", Help: ""})
	errc := prometheus.NewCounter(prometheus.CounterOpts{Name: "test_err", Help: ""})
	cacheMetrics := metrics.NewPrometheusCacheMetrics(
		metrics.WithCacheHitsCounter(hitc),
		metrics.WithCacheMissesCounter(missc),
		metrics.WithCacheErrorsCounter(errc),
	)

	log := slog.New(noophandler.NewNoOpHandler())

	cp := NewCachingWeatherProvider(cacheMock, providerMock, cacheMetrics, log)

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
	require.InDelta(t, float64(1), testutil.ToFloat64(missc), delta)
	require.InDelta(t, float64(0), testutil.ToFloat64(hitc), delta)

	cacheMock.EXPECT().Get(mock.Anything, key).Return(*expectedWeather, nil).Once()

	// Second call
	result2, err2 := cp.GetCurrentWeather(t.Context(), location)
	require.NoError(t, err2)
	assert.Equal(t, expectedWeather, result2)
	require.InDelta(t, float64(1), testutil.ToFloat64(missc), delta)
	require.InDelta(t, float64(1), testutil.ToFloat64(hitc), delta)
}
