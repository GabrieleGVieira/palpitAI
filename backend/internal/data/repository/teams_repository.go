package repository

import (
	"context"
	"strings"
)

type TeamsRepository struct {
	db Querier
}

func NewTeamsRepository(db Querier) *TeamsRepository {
	return &TeamsRepository{db: db}
}

func (r *TeamsRepository) FindByName(ctx context.Context, name string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `select id from teams where name = $1`, strings.TrimSpace(name)).Scan(&id)
	return id, mapNoRows(err)
}

func (r *TeamsRepository) FindByAlias(ctx context.Context, alias string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `select team_id from team_aliases where alias = $1`, strings.TrimSpace(alias)).Scan(&id)
	return id, mapNoRows(err)
}

func (r *TeamsRepository) Create(ctx context.Context, name string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		insert into teams (name)
		values ($1)
		on conflict (name) do update set updated_at = now()
		returning id
	`, strings.TrimSpace(name)).Scan(&id)

	return id, err
}

func (r *TeamsRepository) CreateAlias(ctx context.Context, teamID string, alias string) error {
	alias = strings.TrimSpace(alias)
	if alias == "" {
		return nil
	}

	_, err := r.db.Exec(ctx, `
		insert into team_aliases (team_id, alias)
		values ($1, $2)
		on conflict (alias) do nothing
	`, teamID, alias)

	return err
}

func (r *TeamsRepository) FindOrCreateByRawName(ctx context.Context, rawName string) (string, error) {
	name := strings.Join(strings.Fields(rawName), " ")
	if name == "" {
		return "", ErrNotFound
	}

	id, err := r.FindByName(ctx, name)
	if err == nil {
		return id, nil
	}
	if err != ErrNotFound {
		return "", err
	}

	id, err = r.FindByAlias(ctx, name)
	if err == nil {
		return id, nil
	}
	if err != ErrNotFound {
		return "", err
	}

	return r.Create(ctx, name)
}
