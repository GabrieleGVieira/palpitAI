package importer

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/data/models"
	"github.com/gabrielevieira/palpitai/backend/internal/data/repository"
)

type MatchFinder interface {
	FindByDateAndTeams(ctx context.Context, matchDate time.Time, homeTeamID string, awayTeamID string) (string, error)
}

type HistoricalGoalscorersStore interface {
	Insert(ctx context.Context, goalscorer models.HistoricalGoalscorer) error
}

type HistoricalGoalscorersImporter struct {
	teams       TeamFinder
	matches     MatchFinder
	goalscorers HistoricalGoalscorersStore
}

func NewHistoricalGoalscorersImporter(teams TeamFinder, matches MatchFinder, goalscorers HistoricalGoalscorersStore) *HistoricalGoalscorersImporter {
	return &HistoricalGoalscorersImporter{teams: teams, matches: matches, goalscorers: goalscorers}
}

func (i *HistoricalGoalscorersImporter) ImportFile(ctx context.Context, filePath string) (ImportResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return ImportResult{}, err
	}
	defer file.Close()

	return i.Import(ctx, file)
}

func (i *HistoricalGoalscorersImporter) Import(ctx context.Context, input io.Reader) (ImportResult, error) {
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

		goalscorer, rawHomeTeam, rawAwayTeam, rawTeam, err := parseHistoricalGoalscorer(record, headers)
		if err != nil {
			result.ErrorCount++
			continue
		}

		homeTeamID, err := i.teams.FindOrCreateTeam(ctx, rawHomeTeam)
		if err != nil {
			result.ErrorCount++
			continue
		}
		awayTeamID, err := i.teams.FindOrCreateTeam(ctx, rawAwayTeam)
		if err != nil {
			result.ErrorCount++
			continue
		}
		teamID, err := i.teams.FindOrCreateTeam(ctx, rawTeam)
		if err != nil {
			result.ErrorCount++
			continue
		}

		matchID, err := i.matches.FindByDateAndTeams(ctx, goalscorer.MatchDate, homeTeamID, awayTeamID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				result.SkippedCount++
			} else {
				result.ErrorCount++
			}
			continue
		}

		goalscorer.MatchID = &matchID
		goalscorer.TeamID = &teamID
		if err := i.goalscorers.Insert(ctx, goalscorer); err != nil {
			result.ErrorCount++
			continue
		}
		result.InsertedCount++
	}

	return result, nil
}

func parseHistoricalGoalscorer(record []string, headers map[string]int) (models.HistoricalGoalscorer, string, string, string, error) {
	rawDate := field(record, headers, "date")
	rawHomeTeam := field(record, headers, "home_team")
	rawAwayTeam := field(record, headers, "away_team")
	rawTeam := field(record, headers, "team")

	if err := requireFields(map[string]string{
		"date":      rawDate,
		"home_team": rawHomeTeam,
		"away_team": rawAwayTeam,
		"team":      rawTeam,
	}); err != nil {
		return models.HistoricalGoalscorer{}, "", "", "", err
	}

	matchDate, err := parseDate(rawDate)
	if err != nil {
		return models.HistoricalGoalscorer{}, "", "", "", err
	}
	minute, err := parseOptionalInt(field(record, headers, "minute"))
	if err != nil {
		return models.HistoricalGoalscorer{}, "", "", "", err
	}

	return models.HistoricalGoalscorer{
		MatchDate: matchDate,
		Scorer:    field(record, headers, "scorer"),
		Minute:    minute,
		OwnGoal:   parseBool(field(record, headers, "own_goal")),
		Penalty:   parseBool(field(record, headers, "penalty")),
		Source:    sourceInternationalResults,
	}, rawHomeTeam, rawAwayTeam, rawTeam, nil
}
