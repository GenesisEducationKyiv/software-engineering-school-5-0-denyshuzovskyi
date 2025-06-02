package posgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/sqlutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

type LocationRepository struct{}

func NewLocationRepository() *LocationRepository {
	return &LocationRepository{}
}

func (r *LocationRepository) Save(ctx context.Context, ex sqlutil.SQLExecutor, location *model.Location) (int32, error) {
	const op = "repository.postgresql.location.Save"
	const query = "INSERT INTO location (name) VALUES ($1) RETURNING id"
	var id int32
	err := ex.QueryRowContext(
		ctx,
		query,
		location.Name,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: scan id: %w", op, err)
	}

	return id, nil
}

func (r *LocationRepository) FindByName(ctx context.Context, ex sqlutil.SQLExecutor, name string) (*model.Location, error) {
	const op = "repository.postgresql.location.FindByName"
	const query = `
		SELECT 
			l.id,
			l.name
		FROM location l
		WHERE l.name = $1
		LIMIT 1;
	`

	var l model.Location
	err := ex.QueryRowContext(ctx, query, name).Scan(
		&l.Id,
		&l.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	return &l, nil
}

func (r *LocationRepository) FindById(ctx context.Context, ex sqlutil.SQLExecutor, id int32) (*model.Location, error) {
	const op = "repository.postgresql.location.FindById"
	const query = `
		SELECT 
			l.id,
			l.name
		FROM location l
		WHERE l.id = $1
		LIMIT 1;
	`

	var l model.Location
	err := ex.QueryRowContext(ctx, query, id).Scan(
		&l.Id,
		&l.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	return &l, nil
}
