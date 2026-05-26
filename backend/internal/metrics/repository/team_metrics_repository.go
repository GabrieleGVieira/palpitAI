package repository

import (
	"context"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/metrics/models"
)

type TeamMetricsRepository struct {
	db Querier
}

func NewTeamMetricsRepository(db Querier) *TeamMetricsRepository {
	return &TeamMetricsRepository{db: db}
}

func (r *TeamMetricsRepository) LatestBefore(ctx context.Context, teamID string, matchDate time.Time) (models.TeamMetric, error) {
	var metric models.TeamMetric
	err := r.db.QueryRow(ctx, `
		select
			id::text,
			team_id::text,
			metric_date,
			elo_score,
			attack_score,
			defense_score,
			recent_form_score,
			world_cup_history_score,
			knockout_score,
			group_stage_score,
			avg_goals_scored,
			avg_goals_conceded,
			win_rate,
			draw_rate,
			loss_rate,
			matches_played,
			source,
			created_at,
			updated_at
		from team_metrics
		where team_id = $1 and metric_date < $2
		order by metric_date desc
		limit 1
	`, teamID, matchDate).Scan(
		&metric.ID,
		&metric.TeamID,
		&metric.MetricDate,
		&metric.EloScore,
		&metric.AttackScore,
		&metric.DefenseScore,
		&metric.RecentFormScore,
		&metric.WorldCupHistoryScore,
		&metric.KnockoutScore,
		&metric.GroupStageScore,
		&metric.AvgGoalsScored,
		&metric.AvgGoalsConceded,
		&metric.WinRate,
		&metric.DrawRate,
		&metric.LossRate,
		&metric.MatchesPlayed,
		&metric.Source,
		&metric.CreatedAt,
		&metric.UpdatedAt,
	)

	return metric, mapNoRows(err)
}
