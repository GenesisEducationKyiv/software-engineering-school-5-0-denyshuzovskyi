//go:build !integration

package weather

import (
	"database/sql"
	"log/slog"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/logger/noophandler"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/mapper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWeatherService_GetCurrentWeatherForLocation_Twice(t *testing.T) {
	db, sqlmock, err := sqlmock.New()
	require.NoError(t, err)
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	sqlmock.ExpectBegin()
	sqlmock.ExpectCommit()
	sqlmock.ExpectBegin()
	sqlmock.ExpectCommit()

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
	weatherToReturn := mapper.WeatherWithLocationDTOToWeather(weatherWithLocationDTOToReturn)

	weatherProviderMock.EXPECT().GetCurrentWeather(ctx, location).Return(&weatherWithLocationDTOToReturn, nil).Twice()

	weatherRepositoryMock := NewMockWeatherRepository(t)
	weatherRepositoryMock.EXPECT().FindLastUpdatedByLocation(ctx, mock.AnythingOfType("*sql.Tx"), location).Return(nil, nil).Once()
	weatherRepositoryMock.EXPECT().Save(ctx, mock.AnythingOfType("*sql.Tx"), mock.AnythingOfType("*model.Weather")).Return(nil).Once()
	weatherRepositoryMock.EXPECT().FindLastUpdatedByLocation(ctx, mock.AnythingOfType("*sql.Tx"), location).Return(&weatherToReturn, nil).Once()

	log := slog.New(noophandler.NewNoOpHandler())

	weatherService := NewWeatherService(db, weatherProviderMock, weatherRepositoryMock, log)

	for range 2 {
		weatherDTO, err := weatherService.GetCurrentWeatherForLocation(ctx, location)
		require.NoError(t, err)

		delta := 0.01
		require.InDelta(t, weatherWithLocationDTOToReturn.Weather.Temperature, weatherDTO.Temperature, delta)
		require.InDelta(t, weatherWithLocationDTOToReturn.Weather.Humidity, weatherDTO.Humidity, delta)
		require.Equal(t, weatherWithLocationDTOToReturn.Weather.Description, weatherDTO.Description)
	}

	err = sqlmock.ExpectationsWereMet()
	require.NoError(t, err)
}
