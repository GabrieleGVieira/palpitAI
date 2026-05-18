package repositories

import (
	"context"

	"github.com/gabrielevieira/palpitai/backend/internal/dto"
)

func UserTotalScore(ctx context.Context, db Querier, userID string) (int, error) {
	var totalScore int
	err := db.QueryRow(ctx, `
		select coalesce(sum(p.points), 0)::int
		from predictions p
		join group_members gm on gm.group_id = p.group_id
			and gm.user_id = p.user_id
			and gm.status = 'active'
		where p.user_id = $1
	`, userID).Scan(&totalScore)
	if err != nil {
		return 0, err
	}

	return totalScore, nil
}

func GroupRanking(ctx context.Context, db Querier, groupID string) ([]dto.RankingEntryResponse, error) {
	rows, err := db.Query(ctx, `
		with scores as (
			select
				gm.user_id,
				gm.display_name,
				coalesce(sum(p.points), 0)::int as total_points
			from group_members gm
			left join predictions p on p.group_id = gm.group_id
				and p.user_id = gm.user_id
			where gm.group_id = $1
				and gm.status = 'active'
			group by gm.user_id, gm.display_name
		)
		select
			rank() over (order by total_points desc, display_name asc)::int as position,
			user_id,
			display_name,
			total_points
		from scores
		order by position asc, display_name asc
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ranking := []dto.RankingEntryResponse{}
	for rows.Next() {
		var entry dto.RankingEntryResponse
		if err := rows.Scan(&entry.Position, &entry.UserID, &entry.DisplayName, &entry.TotalPoints); err != nil {
			return nil, err
		}

		ranking = append(ranking, entry)
	}

	return ranking, rows.Err()
}

func UpsertPrediction(ctx context.Context, db Querier, userID string, groupID string, matchID string, request dto.PredictionRequest) (dto.PredictionResponse, error) {
	var prediction dto.PredictionResponse
	err := db.QueryRow(ctx, `
		insert into predictions (group_id, match_id, user_id, home_score, away_score)
		values ($1, $2, $3, $4, $5)
		on conflict (group_id, match_id, user_id)
		do update set
			home_score = excluded.home_score,
			away_score = excluded.away_score,
			updated_at = now()
		returning match_id, home_score, away_score, points, scored_at, updated_at
	`, groupID, matchID, userID, request.HomeScore, request.AwayScore).Scan(
		&prediction.MatchID,
		&prediction.HomeScore,
		&prediction.AwayScore,
		&prediction.Points,
		&prediction.ScoredAt,
		&prediction.UpdatedAt,
	)
	if err != nil {
		return dto.PredictionResponse{}, err
	}

	return prediction, nil
}

func ScoreMatchPredictions(ctx context.Context, db Querier, matchID string, request dto.MatchResultRequest) (int, error) {
	commandTag, err := db.Exec(ctx, `
		update predictions
		set
			points = case
				when home_score = $2 and away_score = $3 then 10
				when sign(home_score - away_score) = sign($2 - $3) then 5
				else 0
			end,
			scored_at = now(),
			updated_at = now()
		where match_id = $1
			and points is distinct from case
				when home_score = $2 and away_score = $3 then 10
				when sign(home_score - away_score) = sign($2 - $3) then 5
				else 0
			end
	`, matchID, request.HomeScore, request.AwayScore)
	if err != nil {
		return 0, err
	}

	return int(commandTag.RowsAffected()), nil
}
