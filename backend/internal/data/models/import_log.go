package models

import "time"

type ImportLog struct {
	ID             string
	ImportType     string
	FilePath       string
	Status         string
	ProcessedCount int
	InsertedCount  int
	SkippedCount   int
	ErrorCount     int
	ErrorMessage   *string
	StartedAt      time.Time
	FinishedAt     *time.Time
}

type ImportSummary struct {
	ProcessedCount int
	InsertedCount  int
	SkippedCount   int
	ErrorCount     int
}
