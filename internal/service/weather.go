package service

import (
	"context"
	"database/sql"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/sqlutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/mapper"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
	"log/slog"
	"time"
)

type WeatherRepository interface {
	Save(context.Context, sqlutil.SQLExecutor, *model.Weather) error
	FindLastUpdatedByLocation(context.Context, sqlutil.SQLExecutor, string) (*model.Weather, error)
}

type WeatherService struct {
	db                 *sql.DB
	weatherProvider    WeatherProvider
	locationRepository LocationRepository
	weatherRepository  WeatherRepository
	log                *slog.Logger
}

func NewWeatherService(db *sql.DB, weatherProvider WeatherProvider, locationRepository LocationRepository, weatherRepository WeatherRepository, log *slog.Logger) *WeatherService {
	return &WeatherService{
		db:                 db,
		weatherProvider:    weatherProvider,
		locationRepository: locationRepository,
		weatherRepository:  weatherRepository,
		log:                log,
	}
}

func (s *WeatherService) GetCurrentWeatherForLocation(ctx context.Context, location string) (*dto.WeatherDTO, error) {
	weather, err := s.weatherProvider.GetCurrentWeather(location)
	if err != nil {
		return nil, err
	}
	weatherDto := mapper.WeatherToWeatherDTO(weather.Weather)

	err = sqlutil.WithTx(ctx, s.db, nil, func(tx *sql.Tx) error {
		loc, errIn := s.locationRepository.FindByName(ctx, tx, weather.Location.Name)
		if errIn != nil {
			return errIn
		}

		var locId int32
		if loc != nil {
			locId = loc.Id
		} else {
			locId, errIn = s.locationRepository.Save(ctx, tx, &weather.Location)
			if errIn != nil {
				return errIn
			}
		}

		weather.Weather.LocationId = locId
		weather.Weather.FetchedAt = time.Now().UTC()

		lastWeather, errIn := s.weatherRepository.FindLastUpdatedByLocation(ctx, tx, weather.Location.Name)
		if errIn != nil {
			return errIn
		}
		if lastWeather != nil && lastWeather.LastUpdated.Equal(weather.LastUpdated) {
			s.log.Info("last weather update is already saved")
			return nil
		}

		errIn = s.weatherRepository.Save(ctx, tx, &weather.Weather)
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
