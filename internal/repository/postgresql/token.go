package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/lib/sqlutil"
	"github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/model"
)

type TokenRepository struct{}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{}
}

func (r *TokenRepository) Save(ctx context.Context, ex sqlutil.SQLExecutor, token *model.Token) error {
	const op = "repository.postgresql.token.Save"
	const query = "INSERT INTO token (token, subscription_id, type, created_at, expires_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := ex.ExecContext(
		ctx,
		query,
		token.Token,
		token.SubscriptionId,
		token.Type,
		token.CreatedAt.UTC(),
		token.ExpiresAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("%s: scan id: %w", op, err)
	}

	return nil
}

func (r *TokenRepository) FindByToken(ctx context.Context, ex sqlutil.SQLExecutor, token string) (*model.Token, error) {
	const op = "repository.postgresql.token.FindByToken"
	const query = `
		SELECT 
			t.token, 
			t.subscription_id, 
			t.type, 
			t.created_at, 
			t.expires_at
		FROM token t
		WHERE t.token = $1
		LIMIT 1;
	`

	var t model.Token
	err := ex.QueryRowContext(ctx, query, token).Scan(
		&t.Token,
		&t.SubscriptionId,
		&t.Type,
		&t.CreatedAt,
		&t.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	return &t, nil
}

func (r *TokenRepository) FindBySubscriptionIdAndType(ctx context.Context, ex sqlutil.SQLExecutor, subscriptionId int32, tokenType model.TokenType) (*model.Token, error) {
	const op = "repository.postgresql.token.FindBySubscriptionIdAndType"
	const query = `
		SELECT 
			t.token, 
			t.subscription_id, 
			t.type, 
			t.created_at, 
			t.expires_at
		FROM token t
		WHERE t.subscription_id = $1 AND t.type = $2
		LIMIT 1;
	`

	var t model.Token
	err := ex.QueryRowContext(ctx, query, subscriptionId, tokenType).Scan(
		&t.Token,
		&t.SubscriptionId,
		&t.Type,
		&t.CreatedAt,
		&t.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	return &t, nil
}
