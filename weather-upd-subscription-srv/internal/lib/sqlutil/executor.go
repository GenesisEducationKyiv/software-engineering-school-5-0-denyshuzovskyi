package sqlutil

import (
	"context"
	"database/sql"
)

// todo: mb use github.com/golang-sql/sqlexp@v0.1.0/querier.go
type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
