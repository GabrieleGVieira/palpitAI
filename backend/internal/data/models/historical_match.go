package models

import "time"

type HistoricalMatch struct {
	ID         string
	MatchDate  time.Time
	HomeTeamID string
	AwayTeamID string
	HomeScore  int
	AwayScore  int
	Tournament string
	City       string
	Country    string
	Neutral    bool
	Source     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
