package repository

import (
	"context"

	"github.com/gabrielevieira/palpitai/backend/internal/data/models"
)

type HistoricalGoalscorersRepository struct {
	db Querier
}

func NewHistoricalGoalscorersRepository(db Querier) *HistoricalGoalscorersRepository {
	return &HistoricalGoalscorersRepository{db: db}
}

func (r *HistoricalGoalscorersRepository) Insert(ctx context.Context, goalscorer models.HistoricalGoalscorer) error {
	_, err := r.db.Exec(ctx, `
		insert into historical_goalscorers (
			match_id, match_date, team_id, scorer, minute, own_goal, penalty, source
		)
		values ($1, $2, $3, nullif($4, ''), $5, $6, $7, $8)
	`, optionalValue(goalscorer.MatchID), goalscorer.MatchDate, optionalValue(goalscorer.TeamID), goalscorer.Scorer,
		optionalValue(goalscorer.Minute), goalscorer.OwnGoal, goalscorer.Penalty, goalscorer.Source)

	return err
}
