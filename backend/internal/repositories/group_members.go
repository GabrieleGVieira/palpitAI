package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/dto"
	"github.com/jackc/pgx/v5"
)

type GroupMembership struct {
	Role   string
	Status string
}

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

func GroupMembershipByUser(ctx context.Context, db Querier, groupID string, userID string) (GroupMembership, error) {
	var membership GroupMembership
	err := db.QueryRow(ctx, `
		select role, status from group_members where group_id = $1 and user_id = $2
	`, groupID, userID).Scan(&membership.Role, &membership.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return GroupMembership{}, ErrNotFound
	}

	return membership, err
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

func ListActiveGroupMembers(ctx context.Context, db Querier, ownerID string, groupID string) ([]dto.GroupMemberResponse, error) {
	rows, err := db.Query(ctx, `
		select
			gm.user_id,
			gm.display_name,
			gm.role,
			gm.joined_at
		from group_members gm
		join groups g on g.id = gm.group_id
		where gm.group_id = $1
			and g.owner_id = $2
			and gm.status = 'active'
		order by
			case when gm.role = 'owner' then 0 else 1 end,
			gm.display_name asc,
			gm.joined_at asc
	`, groupID, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []dto.GroupMemberResponse{}
	for rows.Next() {
		var member dto.GroupMemberResponse
		var joinedAt time.Time
		if err := rows.Scan(&member.UserID, &member.DisplayName, &member.Role, &joinedAt); err != nil {
			return nil, err
		}

		member.JoinedAt = joinedAt
		members = append(members, member)
	}

	return members, rows.Err()
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

func DeleteOwnGroupMembership(ctx context.Context, db Querier, groupID string, userID string) error {
	var deletedGroupID string
	err := db.QueryRow(ctx, `
		delete from group_members
		where group_id = $1
			and user_id = $2
			and role <> 'owner'
			and status = 'active'
		returning group_id
	`, groupID, userID).Scan(&deletedGroupID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return err
}

func DeleteGroupMemberByOwner(ctx context.Context, db Querier, ownerID string, groupID string, memberID string) error {
	var deletedGroupID string
	err := db.QueryRow(ctx, `
		delete from group_members gm
		using groups g
		where g.id = gm.group_id
			and g.owner_id = $1
			and gm.group_id = $2
			and gm.user_id = $3
			and gm.role <> 'owner'
			and gm.status = 'active'
		returning gm.group_id
	`, ownerID, groupID, memberID).Scan(&deletedGroupID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	return err
}

func TransferGroupOwnership(ctx context.Context, db Querier, currentOwnerID string, groupID string, nextOwnerID string) error {
	var transferredGroupID string
	err := db.QueryRow(ctx, `
		with target_member as (
			select gm.group_id, gm.user_id
			from group_members gm
			join groups g on g.id = gm.group_id
			where gm.group_id = $1
				and g.owner_id = $2
				and gm.user_id = $3
				and gm.status = 'active'
				and gm.role = 'member'
			for update
		),
		updated_group as (
			update groups
			set owner_id = $3, updated_at = now()
			where id = (select group_id from target_member)
			returning id
		),
		updated_old_owner as (
			update group_members
			set role = 'member'
			where group_id = (select id from updated_group)
				and user_id = $2
		),
		updated_new_owner as (
			update group_members
			set role = 'owner'
			where group_id = (select id from updated_group)
				and user_id = $3
		)
		select id from updated_group
	`, groupID, currentOwnerID, nextOwnerID).Scan(&transferredGroupID)
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
