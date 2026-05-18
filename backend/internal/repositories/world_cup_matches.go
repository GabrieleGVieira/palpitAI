package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/domain"
	"github.com/gabrielevieira/palpitai/backend/internal/dto"
	"github.com/jackc/pgx/v5"
)

func ListGroupMatches(ctx context.Context, db Querier, groupID string, userID string) ([]dto.MatchResponse, error) {
	rows, err := db.Query(ctx, `
		select
			m.id,
			m.home_team,
			m.away_team,
			m.stage,
			m.status,
			m.kickoff_at,
			m.home_score,
			m.away_score,
			m.finished_at,
			p.home_score,
			p.away_score,
			p.points,
			p.scored_at,
			p.updated_at
		from world_cup_matches m
		join groups g on g.id = $1
		left join predictions p on p.group_id = g.id and p.match_id = m.id and p.user_id = $2
		where g.match_scope = 'all'
			or m.home_team = any(g.selected_teams)
			or m.away_team = any(g.selected_teams)
		order by m.kickoff_at asc
	`, groupID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := []dto.MatchResponse{}
	for rows.Next() {
		var match dto.MatchResponse
		var homeScore *int
		var awayScore *int
		var finalHomeScore *int
		var finalAwayScore *int
		var finishedAt *time.Time
		var points *int
		var scoredAt *time.Time
		var updatedAt *time.Time

		if err := rows.Scan(
			&match.ID,
			&match.HomeTeam,
			&match.AwayTeam,
			&match.Stage,
			&match.Status,
			&match.KickoffAt,
			&finalHomeScore,
			&finalAwayScore,
			&finishedAt,
			&homeScore,
			&awayScore,
			&points,
			&scoredAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		match.FinalHomeScore = finalHomeScore
		match.FinalAwayScore = finalAwayScore
		match.FinishedAt = finishedAt

		if homeScore != nil && awayScore != nil && updatedAt != nil {
			match.MyPrediction = &dto.PredictionResponse{
				AwayScore: *awayScore,
				HomeScore: *homeScore,
				MatchID:   match.ID,
				Points:    points,
				ScoredAt:  scoredAt,
				UpdatedAt: *updatedAt,
			}
		}

		matches = append(matches, match)
	}

	return matches, rows.Err()
}

func MatchKickoffForGroup(ctx context.Context, db Querier, groupID string, matchID string) (time.Time, error) {
	var kickoffAt time.Time
	err := db.QueryRow(ctx, `
		select m.kickoff_at
		from world_cup_matches m
		join groups g on g.id = $1
		where m.id = $2
			and (
				g.match_scope = 'all'
				or m.home_team = any(g.selected_teams)
				or m.away_team = any(g.selected_teams)
			)
	`, groupID, matchID).Scan(&kickoffAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, ErrNotFound
	}

	return kickoffAt, err
}

func UpdateMatchResult(ctx context.Context, db Querier, matchID string, request dto.MatchResultRequest) error {
	_, err := db.Exec(ctx, `
		update world_cup_matches
		set
			home_score = $2,
			away_score = $3,
			status = 'finished',
			finished_at = now(),
			last_synced_at = now()
		where id = $1
	`, matchID, request.HomeScore, request.AwayScore)

	return err
}

func MatchDetailsByID(ctx context.Context, db Querier, matchID string) (domain.MatchDetails, error) {
	var details domain.MatchDetails
	err := db.QueryRow(ctx, `
		select home_team, away_team
		from world_cup_matches
		where id = $1
	`, matchID).Scan(&details.HomeTeam, &details.AwayTeam)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.MatchDetails{}, ErrNotFound
	}

	return details, err
}
