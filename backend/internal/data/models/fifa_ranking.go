package models

import "time"

type FifaRanking struct {
	ID             string
	TeamID         string
	RankingDate    time.Time
	Rank           int
	TotalPoints    *float64
	PreviousPoints *float64
	RankChange     *int
	Confederation  *string
	Source         string
	CreatedAt      time.Time
}
