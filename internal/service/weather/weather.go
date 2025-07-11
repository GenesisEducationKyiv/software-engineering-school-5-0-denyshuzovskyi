package weather

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/dto"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/sqlutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/mapper"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

type WeatherProvider interface {
	GetCurrentWeather(context.Context, string) (*dto.WeatherWithLocationDTO, error)
}

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
	lastWeather, err := s.weatherRepository.FindLastUpdatedByLocation(ctx, s.db, location)
	if err != nil {
		return nil, fmt.Errorf("fetch weather from cache: %w", err)
	}

	if lastWeather == nil || lastWeather.LastUpdated.Add(15*time.Minute).Before(time.Now().UTC()) {
		weatherWithLocationDTO, err := s.weatherProvider.GetCurrentWeather(ctx, location)
		if err != nil {
			return nil, fmt.Errorf("update weather: %w", err)
		}
		weather := mapper.WeatherWithLocationDTOToWeather(*weatherWithLocationDTO)
		if err := s.weatherRepository.Save(ctx, s.db, &weather); err != nil {
			return nil, fmt.Errorf("save weather: %w", err)
		}
		lastWeather = &weather
	}

	weatherDto := mapper.WeatherToWeatherDTO(*lastWeather)

	return &weatherDto, nil
}
