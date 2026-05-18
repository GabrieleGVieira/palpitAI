package repositories

import (
	"context"
	"errors"

	"github.com/gabrielevieira/palpitai/backend/internal/apperrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrNotFound = apperrors.ErrNotFound

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Datastore interface {
	Querier
	Ping(ctx context.Context) error
}

type TxRunner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
