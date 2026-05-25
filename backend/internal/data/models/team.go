package models

import "time"

type Team struct {
	ID          string
	Name        string
	CountryCode *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
