package models

import (
	"encoding/json"
	"time"
)

type ExternalApiSnapshot struct {
	ID          string
	Provider    string
	Endpoint    string
	PayloadJSON json.RawMessage
	FetchedAt   time.Time
	ExpiresAt   time.Time
}
