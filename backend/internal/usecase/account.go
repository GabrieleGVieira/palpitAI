package usecase

import (
	"context"

	"github.com/gabrielevieira/palpitai/backend/internal/apperrors"
	"github.com/gabrielevieira/palpitai/backend/internal/repositories"
)

var ErrAccountOwnsGroups = apperrors.NewConflict("account owns groups")

type AccountUsecase struct {
	db Datastore
}

func NewAccountUsecase(db Datastore) AccountUsecase {
	return AccountUsecase{db: db}
}

func (uc AccountUsecase) DeleteAccount(ctx context.Context, userID string) error {
	return DeleteAccount(ctx, uc.db, userID)
}

func DeleteAccount(ctx context.Context, db Datastore, userID string) error {
	return withTx(ctx, db, func(tx repositories.Querier) error {
		ownedGroups, err := repositories.UserOwnedGroupCount(ctx, tx, userID)
		if err != nil {
			return err
		}
		if ownedGroups > 0 {
			return ErrAccountOwnsGroups
		}

		return repositories.AnonymizeAccountData(ctx, tx, userID)
	})
}
