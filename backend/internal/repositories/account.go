package repositories

import "context"

const DeletedUserDisplayName = "Usuário excluído"

func UserOwnedGroupCount(ctx context.Context, db Querier, userID string) (int, error) {
	var count int
	err := db.QueryRow(ctx, `
		select count(*)::int
		from groups
		where owner_id = $1
	`, userID).Scan(&count)

	return count, err
}

func AnonymizeAccountData(ctx context.Context, db Querier, userID string) error {
	if _, err := db.Exec(ctx, `
		update group_members
		set
			display_name = $2,
			status = 'deleted'
		where user_id = $1
	`, userID, DeletedUserDisplayName); err != nil {
		return err
	}

	_, err := db.Exec(ctx, `
		update groups
		set updated_at = now()
		where owner_id = $1
	`, userID)

	return err
}
