package repository

import (
	"context"
	"errors"

	"github.com/gabrielevieira/palpitai/backend/internal/data/models"
	"github.com/jackc/pgx/v5"
)

type FifaRankingsRepository struct {
	db Querier
}

func NewFifaRankingsRepository(db Querier) *FifaRankingsRepository {
	return &FifaRankingsRepository{db: db}
}

func (r *FifaRankingsRepository) Upsert(ctx context.Context, ranking models.FifaRanking) (bool, error) {
	err := r.db.QueryRow(ctx, `
		insert into fifa_rankings (
			team_id, ranking_date, rank, total_points, previous_points,
			rank_change, confederation, source
		)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		on conflict (team_id, ranking_date) do nothing
		returning id
	`, ranking.TeamID, ranking.RankingDate, ranking.Rank, optionalValue(ranking.TotalPoints), optionalValue(ranking.PreviousPoints),
		optionalValue(ranking.RankChange), optionalValue(ranking.Confederation), ranking.Source).Scan(new(string))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	return false, err
}
