package usecase

import (
	"context"

	"github.com/gabrielevieira/palpitai/backend/internal/repositories"
)

func withTx(ctx context.Context, db Datastore, fn func(repositories.Querier) error) error {
	runner, ok := db.(repositories.TxRunner)
	if !ok {
		return fn(db)
	}

	tx, err := runner.Begin(ctx)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
