//go:build !integration

package weather

import (
	"log/slog"
	"testing"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/logger/noophandler"
	"github.com/stretchr/testify/require"
)

func TestWeatherService_GetCurrentWeatherForLocation_Twice(t *testing.T) {
	ctx := t.Context()
	location := "Kyiv"

	weatherProviderMock := NewMockWeatherProvider(t)
	weatherWithLocationDTOToReturn := dto.WeatherWithLocationDTO{
		Weather: dto.WeatherDTO{
			Temperature: float32(23),
			Humidity:    float32(43),
			Description: "Enjoy",
		},
		Location: dto.Location{
			Name: location,
		},
		LastUpdated: time.Now().UTC().Unix(),
	}
	weatherProviderMock.EXPECT().GetCurrentWeather(ctx, location).Return(&weatherWithLocationDTOToReturn, nil).Once()

	log := slog.New(noophandler.NewNoOpHandler())
	weatherService := NewWeatherService(weatherProviderMock, log)

	weatherDTO, err := weatherService.GetCurrentWeatherForLocation(ctx, location)
	require.NoError(t, err)

	delta := 0.01
	require.InDelta(t, weatherWithLocationDTOToReturn.Weather.Temperature, weatherDTO.Temperature, delta)
	require.InDelta(t, weatherWithLocationDTOToReturn.Weather.Humidity, weatherDTO.Humidity, delta)
	require.Equal(t, weatherWithLocationDTOToReturn.Weather.Description, weatherDTO.Description)
}
