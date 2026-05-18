package repositories

import (
	"context"
	"errors"

	"github.com/gabrielevieira/palpitai/backend/internal/dto"
	"github.com/jackc/pgx/v5"
)

func GroupMemberStatus(ctx context.Context, db Querier, groupID string, userID string) (string, error) {
	var status string
	err := db.QueryRow(ctx, `
		select status from group_members where group_id = $1 and user_id = $2
	`, groupID, userID).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}

	return status, err
}

func InsertGroupMember(ctx context.Context, db Querier, groupID string, userID string, status string, displayName string) error {
	_, err := db.Exec(ctx, `
		insert into group_members (group_id, user_id, role, status, display_name)
		values ($1, $2, 'member', $3, $4)
		on conflict (group_id, user_id) do nothing
	`, groupID, userID, status, displayName)

	return err
}

func ListPendingJoinRequests(ctx context.Context, db Querier, ownerID string, groupID string) ([]dto.JoinRequestResponse, error) {
	rows, err := db.Query(ctx, `
		select
			gm.user_id,
			gm.display_name,
			gm.joined_at
		from group_members gm
		join groups g on g.id = gm.group_id
		where gm.group_id = $1
			and g.owner_id = $2
			and gm.status = 'pending'
		order by gm.joined_at asc
	`, groupID, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []dto.JoinRequestResponse{}
	for rows.Next() {
		var request dto.JoinRequestResponse
		if err := rows.Scan(&request.UserID, &request.DisplayName, &request.RequestedAt); err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, rows.Err()
}

func ApprovePendingMember(ctx context.Context, db Querier, groupID string, requesterID string) error {
	var approvedGroupID string
	err := db.QueryRow(ctx, `
		update group_members
		set status = 'active', joined_at = now()
		where group_id = $1 and user_id = $2 and status = 'pending'
		returning group_id
	`, groupID, requesterID).Scan(&approvedGroupID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return err
}

func ActiveGroupMemberExists(ctx context.Context, db Querier, userID string, groupID string) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, `
		select exists (
			select 1
			from group_members
			where group_id = $1
				and user_id = $2
				and status = 'active'
		)
	`, groupID, userID).Scan(&exists)

	return exists, err
}
