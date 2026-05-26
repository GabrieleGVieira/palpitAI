package models

import "time"

type TeamMetric struct {
	ID                   string
	TeamID               string
	MetricDate           time.Time
	EloScore             *float64
	AttackScore          *float64
	DefenseScore         *float64
	RecentFormScore      *float64
	WorldCupHistoryScore *float64
	KnockoutScore        *float64
	GroupStageScore      *float64
	AvgGoalsScored       *float64
	AvgGoalsConceded     *float64
	WinRate              *float64
	DrawRate             *float64
	LossRate             *float64
	MatchesPlayed        int
	Source               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
