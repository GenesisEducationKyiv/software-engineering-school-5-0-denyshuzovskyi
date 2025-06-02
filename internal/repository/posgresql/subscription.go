package posgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/sqlutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

type SubscriptionRepository struct{}

func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{}
}

func (r *SubscriptionRepository) Save(ctx context.Context, ex sqlutil.SQLExecutor, subscription *model.Subscription) (int32, error) {
	const op = "repository.postgresql.subscription.Save"
	const query = "INSERT INTO subscription (subscriber_id, location_id, frequency, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	var id int32
	err := ex.QueryRowContext(
		ctx,
		query,
		subscription.SubscriberId,
		subscription.LocationId,
		subscription.Frequency,
		subscription.Status,
		subscription.CreatedAt.UTC(),
		subscription.UpdatedAt.UTC(),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: scan id: %w", op, err)
	}

	return id, nil
}

func (r *SubscriptionRepository) FindBySubscriberIdAndLocationId(ctx context.Context, ex sqlutil.SQLExecutor, subscriberId int32, location_id int32) (*model.Subscription, error) {
	const op = "repository.postgresql.subscriber.FindBySubscriberIdAndLocationId"
	const query = `
		SELECT 
			s.id,
			s.subscriber_id,
			s.location_id,
			s.frequency,
			s.status,
			s.created_at,
			s.updated_at
		FROM subscription s
		WHERE s.subscriber_id = $1 AND s.location_id = $2
		LIMIT 1;
	`

	var s model.Subscription
	err := ex.QueryRowContext(ctx, query, subscriberId, location_id).Scan(
		&s.Id,
		&s.SubscriberId,
		&s.LocationId,
		&s.Frequency,
		&s.Status,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	return &s, nil
}

func (r *SubscriptionRepository) FindById(ctx context.Context, ex sqlutil.SQLExecutor, id int32) (*model.Subscription, error) {
	const op = "repository.postgresql.subscription.FindById"
	const query = `
		SELECT 
			s.id,
			s.subscriber_id,
			s.location_id,
			s.frequency,
			s.status,
			s.created_at,
			s.updated_at
		FROM subscription s
		WHERE s.id = $1
		LIMIT 1;
	`

	var s model.Subscription
	err := ex.QueryRowContext(ctx, query, id).Scan(
		&s.Id,
		&s.SubscriberId,
		&s.LocationId,
		&s.Frequency,
		&s.Status,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	return &s, nil
}

func (r *SubscriptionRepository) DeleteById(ctx context.Context, ex sqlutil.SQLExecutor, id int32) error {
	const op = "repository.postgresql.subscription.DeleteById"
	const query = `
		DELETE FROM subscription
		WHERE id = $1;
	`

	_, err := ex.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: delete failed: %w", op, err)
	}
	return nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, ex sqlutil.SQLExecutor, subscription *model.Subscription) (*model.Subscription, error) {
	const op = "repository.postgresql.subscription.Update"
	const query = `
		UPDATE subscription
		SET subscriber_id = $1,
		    location_id = $2,
		    frequency = $3,
		    status = $4,
		    updated_at = $5
		WHERE id = $6
		RETURNING 
		    id,
		    subscriber_id,
		    location_id,
		    frequency,
		    status,
		    created_at,
		    updated_at;
	`

	var updated model.Subscription
	err := ex.QueryRowContext(ctx, query,
		subscription.SubscriberId,
		subscription.LocationId,
		subscription.Frequency,
		subscription.Status,
		subscription.UpdatedAt,
		subscription.Id,
	).Scan(
		&updated.Id,
		&updated.SubscriberId,
		&updated.LocationId,
		&updated.Frequency,
		&updated.Status,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: update failed: %w", op, err)
	}
	return &updated, nil
}

func (r *SubscriptionRepository) FindAllByFrequencyAndConfirmedStatus(ctx context.Context, ex sqlutil.SQLExecutor, frequency model.Frequency) (subscriptions []*model.Subscription, err error) {
	const op = "repository.postgresql.subscription.FindAllByFrequencyAndConfirmedStatus"
	const query = `
		SELECT 
			s.id,
			s.subscriber_id,
			s.location_id,
			s.frequency,
			s.status,
			s.created_at,
			s.updated_at
		FROM subscription s
		WHERE s.frequency = $1 AND s.status = 'confirmed';
	`

	rows, err := ex.QueryContext(ctx, query, frequency)
	if err != nil {
		err = fmt.Errorf("%s: query failed: %w", op, err)

		return
	}
	defer func(rows *sql.Rows) {
		cerr := rows.Close()
		err = errors.Join(err, cerr)
	}(rows)

	for rows.Next() {
		var s model.Subscription
		err = rows.Scan(
			&s.Id,
			&s.SubscriberId,
			&s.LocationId,
			&s.Frequency,
			&s.Status,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			err = fmt.Errorf("%s: scan failed: %w", op, err)

			return
		}
		subscriptions = append(subscriptions, &s)
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("%s: rows iteration error: %w", op, err)

		return
	}

	return
}
