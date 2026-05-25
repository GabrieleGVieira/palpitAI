package importer

import (
	"context"
	"encoding/csv"
	"io"
	"os"

	"github.com/gabrielevieira/palpitai/backend/internal/data/models"
)

type TeamFinder interface {
	FindOrCreateTeam(ctx context.Context, rawName string) (string, error)
}

type HistoricalMatchesStore interface {
	Upsert(ctx context.Context, match models.HistoricalMatch) (string, bool, error)
}

type HistoricalMatchesImporter struct {
	teams   TeamFinder
	matches HistoricalMatchesStore
}

func NewHistoricalMatchesImporter(teams TeamFinder, matches HistoricalMatchesStore) *HistoricalMatchesImporter {
	return &HistoricalMatchesImporter{teams: teams, matches: matches}
}

func (i *HistoricalMatchesImporter) ImportFile(ctx context.Context, filePath string) (ImportResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return ImportResult{}, err
	}
	defer file.Close()

	return i.Import(ctx, file)
}

func (i *HistoricalMatchesImporter) Import(ctx context.Context, input io.Reader) (ImportResult, error) {
	reader := csv.NewReader(input)
	reader.FieldsPerRecord = -1

	headers, err := readHeader(reader)
	if err != nil {
		return ImportResult{}, err
	}

	var result ImportResult
	for {
		record, err := nextRecord(reader)
		if isEOF(err) {
			break
		}
		result.ProcessedCount++
		if err != nil {
			result.ErrorCount++
			continue
		}

		match, err := parseHistoricalMatch(record, headers)
		if err != nil {
			result.ErrorCount++
			continue
		}

		homeTeamID, err := i.teams.FindOrCreateTeam(ctx, field(record, headers, "home_team"))
		if err != nil {
			result.ErrorCount++
			continue
		}
		awayTeamID, err := i.teams.FindOrCreateTeam(ctx, field(record, headers, "away_team"))
		if err != nil {
			result.ErrorCount++
			continue
		}

		match.HomeTeamID = homeTeamID
		match.AwayTeamID = awayTeamID
		_, inserted, err := i.matches.Upsert(ctx, match)
		if err != nil {
			result.ErrorCount++
			continue
		}
		if inserted {
			result.InsertedCount++
		} else {
			result.SkippedCount++
		}
	}

	return result, nil
}

func parseHistoricalMatch(record []string, headers map[string]int) (models.HistoricalMatch, error) {
	rawDate := field(record, headers, "date")
	rawHomeTeam := field(record, headers, "home_team")
	rawAwayTeam := field(record, headers, "away_team")
	rawHomeScore := field(record, headers, "home_score")
	rawAwayScore := field(record, headers, "away_score")

	if err := requireFields(map[string]string{
		"date":       rawDate,
		"home_team":  rawHomeTeam,
		"away_team":  rawAwayTeam,
		"home_score": rawHomeScore,
		"away_score": rawAwayScore,
	}); err != nil {
		return models.HistoricalMatch{}, err
	}

	matchDate, err := parseDate(rawDate)
	if err != nil {
		return models.HistoricalMatch{}, err
	}
	homeScore, err := parseInt(rawHomeScore)
	if err != nil {
		return models.HistoricalMatch{}, err
	}
	awayScore, err := parseInt(rawAwayScore)
	if err != nil {
		return models.HistoricalMatch{}, err
	}

	return models.HistoricalMatch{
		MatchDate:  matchDate,
		HomeScore:  homeScore,
		AwayScore:  awayScore,
		Tournament: field(record, headers, "tournament"),
		City:       field(record, headers, "city"),
		Country:    field(record, headers, "country"),
		Neutral:    parseBool(field(record, headers, "neutral")),
		Source:     sourceInternationalResults,
	}, nil
}
