package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
		create extension if not exists pgcrypto;

		create table if not exists groups (
			id uuid primary key default gen_random_uuid(),
			owner_id uuid not null,
			name text not null,
			description text not null default '',
			match_scope text not null check (match_scope in ('all', 'selected')),
			selected_teams text[] not null default '{}',
			participant_limit integer check (participant_limit is null or participant_limit > 1),
			is_private boolean not null default true,
			invite_code text not null unique,
			created_at timestamptz not null default now(),
			updated_at timestamptz not null default now()
		);

		create table if not exists group_members (
			group_id uuid not null references groups(id) on delete cascade,
			user_id uuid not null,
			role text not null check (role in ('owner', 'member')),
			joined_at timestamptz not null default now(),
			primary key (group_id, user_id)
		);
	`)

	return err
}
