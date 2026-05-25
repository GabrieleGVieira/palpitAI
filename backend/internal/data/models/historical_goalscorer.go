package models

import "time"

type HistoricalGoalscorer struct {
	ID        string
	MatchID   *string
	MatchDate time.Time
	TeamID    *string
	Scorer    string
	Minute    *int
	OwnGoal   bool
	Penalty   bool
	Source    string
	CreatedAt time.Time
}
