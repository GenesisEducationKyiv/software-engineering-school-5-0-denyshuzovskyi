package service

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/sqlutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/mapper"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

type WeatherRepository interface {
	Save(context.Context, sqlutil.SQLExecutor, *model.Weather) error
	FindLastUpdatedByLocation(context.Context, sqlutil.SQLExecutor, string) (*model.Weather, error)
}

type WeatherService struct {
	db                *sql.DB
	weatherProvider   WeatherProvider
	weatherRepository WeatherRepository
	log               *slog.Logger
}

func NewWeatherService(db *sql.DB, weatherProvider WeatherProvider, weatherRepository WeatherRepository, log *slog.Logger) *WeatherService {
	return &WeatherService{
		db:                db,
		weatherProvider:   weatherProvider,
		weatherRepository: weatherRepository,
		log:               log,
	}
}

func (s *WeatherService) GetCurrentWeatherForLocation(ctx context.Context, location string) (*dto.WeatherDTO, error) {
	weather, err := s.weatherProvider.GetCurrentWeather(ctx, location)
	if err != nil {
		return nil, err
	}
	weatherDto := mapper.WeatherToWeatherDTO(*weather)

	err = sqlutil.WithTx(ctx, s.db, nil, func(tx *sql.Tx) error {
		lastWeather, errIn := s.weatherRepository.FindLastUpdatedByLocation(ctx, tx, weather.LocationName)
		if errIn != nil {
			return errIn
		}
		if lastWeather != nil && lastWeather.LastUpdated.Equal(weather.LastUpdated) {
			s.log.Info("last weather update is already saved")
			return nil
		}

		weather.FetchedAt = time.Now().UTC()
		errIn = s.weatherRepository.Save(ctx, tx, weather)
		if errIn != nil {
			return errIn
		}

		return nil
	})
	if err != nil {
		s.log.Error("rolled back transaction because of", "error", err)
	} else {
		s.log.Info("transaction commited successfully")
	}

	return &weatherDto, nil
}
