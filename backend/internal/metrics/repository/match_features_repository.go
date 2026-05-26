package repository

import (
	"context"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/metrics/models"
)

type MatchFeaturesRepository struct {
	db Querier
}

func NewMatchFeaturesRepository(db Querier) *MatchFeaturesRepository {
	return &MatchFeaturesRepository{db: db}
}

func (r *MatchFeaturesRepository) FindByMatchID(ctx context.Context, matchID string) (models.MatchFeature, error) {
	var feature models.MatchFeature
	err := r.db.QueryRow(ctx, selectMatchFeatureSQL()+`
		where match_id = $1
		order by created_at desc
		limit 1
	`, matchID).Scan(scanMatchFeature(&feature)...)

	return feature, mapNoRows(err)
}

func (r *MatchFeaturesRepository) ListByDateRange(ctx context.Context, fromDate time.Time, toDate time.Time) ([]models.MatchFeature, error) {
	rows, err := r.db.Query(ctx, selectMatchFeatureSQL()+`
		where match_date between $1 and $2
		order by match_date, created_at
	`, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	features := []models.MatchFeature{}
	for rows.Next() {
		var feature models.MatchFeature
		if err := rows.Scan(scanMatchFeature(&feature)...); err != nil {
			return nil, err
		}
		features = append(features, feature)
	}

	return features, rows.Err()
}

func selectMatchFeatureSQL() string {
	return `
		select
			id::text,
			match_id::text,
			match_date,
			home_team_id::text,
			away_team_id::text,
			tournament,
			stage,
			home_elo_score,
			away_elo_score,
			elo_diff,
			home_attack_score,
			away_attack_score,
			home_defense_score,
			away_defense_score,
			home_recent_form_score,
			away_recent_form_score,
			home_fifa_rank,
			away_fifa_rank,
			fifa_rank_diff,
			home_avg_goals_scored,
			away_avg_goals_scored,
			home_avg_goals_conceded,
			away_avg_goals_conceded,
			home_world_cup_history_score,
			away_world_cup_history_score,
			neutral,
			created_at,
			updated_at
		from match_features
	`
}

func scanMatchFeature(feature *models.MatchFeature) []any {
	return []any{
		&feature.ID,
		&feature.MatchID,
		&feature.MatchDate,
		&feature.HomeTeamID,
		&feature.AwayTeamID,
		&feature.Tournament,
		&feature.Stage,
		&feature.HomeEloScore,
		&feature.AwayEloScore,
		&feature.EloDiff,
		&feature.HomeAttackScore,
		&feature.AwayAttackScore,
		&feature.HomeDefenseScore,
		&feature.AwayDefenseScore,
		&feature.HomeRecentFormScore,
		&feature.AwayRecentFormScore,
		&feature.HomeFifaRank,
		&feature.AwayFifaRank,
		&feature.FifaRankDiff,
		&feature.HomeAvgGoalsScored,
		&feature.AwayAvgGoalsScored,
		&feature.HomeAvgGoalsConceded,
		&feature.AwayAvgGoalsConceded,
		&feature.HomeWorldCupHistoryScore,
		&feature.AwayWorldCupHistoryScore,
		&feature.Neutral,
		&feature.CreatedAt,
		&feature.UpdatedAt,
	}
}
