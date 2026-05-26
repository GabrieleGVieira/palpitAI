package models

import "time"

type MatchFeature struct {
	ID                       string
	MatchID                  *string
	MatchDate                time.Time
	HomeTeamID               string
	AwayTeamID               string
	Tournament               *string
	Stage                    *string
	HomeEloScore             *float64
	AwayEloScore             *float64
	EloDiff                  *float64
	HomeAttackScore          *float64
	AwayAttackScore          *float64
	HomeDefenseScore         *float64
	AwayDefenseScore         *float64
	HomeRecentFormScore      *float64
	AwayRecentFormScore      *float64
	HomeFifaRank             *int
	AwayFifaRank             *int
	FifaRankDiff             *int
	HomeAvgGoalsScored       *float64
	AwayAvgGoalsScored       *float64
	HomeAvgGoalsConceded     *float64
	AwayAvgGoalsConceded     *float64
	HomeWorldCupHistoryScore *float64
	AwayWorldCupHistoryScore *float64
	Neutral                  bool
	CreatedAt                time.Time
	UpdatedAt                time.Time
}
