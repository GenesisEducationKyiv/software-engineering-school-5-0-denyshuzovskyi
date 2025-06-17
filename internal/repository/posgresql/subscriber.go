package posgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/sqlutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

type SubscriberRepository struct{}

func NewSubscriberRepository() *SubscriberRepository {
	return &SubscriberRepository{}
}

func (r *SubscriberRepository) Save(ctx context.Context, ex sqlutil.SQLExecutor, subscriber *model.Subscriber) (int32, error) {
	const op = "repository.postgresql.subscriber.Save"
	const query = "INSERT INTO subscriber (email, created_at) VALUES ($1, $2) RETURNING id"
	var id int32
	err := ex.QueryRowContext(
		ctx,
		query,
		subscriber.Email,
		subscriber.CreatedAt.UTC(),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: scan id: %w", op, err)
	}

	return id, nil
}

func (r *SubscriberRepository) FindByEmail(ctx context.Context, ex sqlutil.SQLExecutor, email string) (*model.Subscriber, error) {
	const op = "repository.postgresql.subscriber.FindByEmail"
	const query = `
		SELECT 
			s.id,
			s.email,
			s.created_at
		FROM subscriber s
		WHERE s.email = $1
		LIMIT 1;
	`

	var s model.Subscriber
	err := ex.QueryRowContext(ctx, query, email).Scan(
		&s.Id,
		&s.Email,
		&s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	return &s, nil
}

func (r *SubscriberRepository) FindById(ctx context.Context, ex sqlutil.SQLExecutor, id int32) (*model.Subscriber, error) {
	const op = "repository.postgresql.subscriber.FindById"
	const query = `
		SELECT 
			s.id,
			s.email,
			s.created_at
		FROM subscriber s
		WHERE s.id = $1
		LIMIT 1;
	`

	var s model.Subscriber
	err := ex.QueryRowContext(ctx, query, id).Scan(
		&s.Id,
		&s.Email,
		&s.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	return &s, nil
}
