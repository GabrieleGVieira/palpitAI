package repository

import (
	"context"
	"errors"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/data/models"
	"github.com/jackc/pgx/v5"
)

type HistoricalMatchesRepository struct {
	db Querier
}

func NewHistoricalMatchesRepository(db Querier) *HistoricalMatchesRepository {
	return &HistoricalMatchesRepository{db: db}
}

func (r *HistoricalMatchesRepository) Upsert(ctx context.Context, match models.HistoricalMatch) (string, bool, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		insert into historical_matches (
			match_date, home_team_id, away_team_id, home_score, away_score,
			tournament, city, country, neutral, source
		)
		values ($1, $2, $3, $4, $5, nullif($6, ''), nullif($7, ''), nullif($8, ''), $9, $10)
		on conflict do nothing
		returning id
	`, match.MatchDate, match.HomeTeamID, match.AwayTeamID, match.HomeScore, match.AwayScore,
		match.Tournament, match.City, match.Country, match.Neutral, match.Source).Scan(&id)
	if err == nil {
		return id, true, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return "", false, err
	}

	id, err = r.FindByDateAndTeams(ctx, match.MatchDate, match.HomeTeamID, match.AwayTeamID)
	if err != nil {
		return "", false, err
	}

	return id, false, nil
}

func (r *HistoricalMatchesRepository) FindByDateAndTeams(ctx context.Context, matchDate time.Time, homeTeamID string, awayTeamID string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		select id
		from historical_matches
		where match_date = $1 and home_team_id = $2 and away_team_id = $3
		order by created_at asc
		limit 1
	`, matchDate, homeTeamID, awayTeamID).Scan(&id)

	return id, mapNoRows(err)
}
