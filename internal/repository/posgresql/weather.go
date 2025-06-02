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
	const query = "INSERT INTO weather (location_id, last_updated, fetched_at, temperature, humidity, description) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := ex.ExecContext(
		ctx,
		query,
		weather.LocationId,
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
			w.location_id, 
			w.last_updated, 
			w.fetched_at, 
			w.temperature, 
			w.humidity, 
			w.description
		FROM weather w
		JOIN location l ON w.location_id = l.id
		WHERE l.name = $1
		ORDER BY w.last_updated DESC
		LIMIT 1;
	`

	var w model.Weather
	err := ex.QueryRowContext(ctx, query, location).Scan(
		&w.LocationId,
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
