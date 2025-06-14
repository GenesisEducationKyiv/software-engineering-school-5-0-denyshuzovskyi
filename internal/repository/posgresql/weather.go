package posgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/sqlutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

type WeatherRepository struct {
}

func NewWeatherRepository() *WeatherRepository {
	return &WeatherRepository{}
}

func (r *WeatherRepository) Save(ctx context.Context, ex sqlutil.SQLExecutor, weather *model.Weather) error {
	const op = "repository.postgresql.weather.Save"
	const query = "INSERT INTO weather (location_name, last_updated, fetched_at, temperature, humidity, description) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := ex.ExecContext(
		ctx,
		query,
		weather.LocationName,
		weather.LastUpdated.UTC(),
		weather.FetchedAt.UTC(),
		weather.Temperature,
		weather.Humidity,
		weather.Description,
	)
	if err != nil {
		return fmt.Errorf("%s: scan id: %w", op, err)
	}

	return nil
}

func (r *WeatherRepository) FindLastUpdatedByLocation(
	ctx context.Context,
	ex sqlutil.SQLExecutor,
	location string,
) (*model.Weather, error) {
	const op = "repository.postgresql.weather.FindLastUpdatedByLocation"
	const query = `
		SELECT 
			w.location_name, 
			w.last_updated, 
			w.fetched_at, 
			w.temperature, 
			w.humidity, 
			w.description
		FROM weather w
		WHERE w.location_name = $1
		ORDER BY w.last_updated DESC
		LIMIT 1;
	`

	var w model.Weather
	err := ex.QueryRowContext(ctx, query, location).Scan(
		&w.LocationName,
		&w.LastUpdated,
		&w.FetchedAt,
		&w.Temperature,
		&w.Humidity,
		&w.Description,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	return &w, nil
}
