package importer

import (
	"context"
	"encoding/csv"
	"io"
	"os"

	"github.com/gabrielevieira/palpitai/backend/internal/data/models"
)

type FifaRankingsStore interface {
	Upsert(ctx context.Context, ranking models.FifaRanking) (bool, error)
}

type FifaRankingImporter struct {
	teams    TeamFinder
	rankings FifaRankingsStore
}

func NewFifaRankingImporter(teams TeamFinder, rankings FifaRankingsStore) *FifaRankingImporter {
	return &FifaRankingImporter{teams: teams, rankings: rankings}
}

func (i *FifaRankingImporter) ImportFile(ctx context.Context, filePath string) (ImportResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return ImportResult{}, err
	}
	defer file.Close()

	return i.Import(ctx, file)
}

func (i *FifaRankingImporter) Import(ctx context.Context, input io.Reader) (ImportResult, error) {
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

		ranking, rawTeam, err := parseFifaRanking(record, headers)
		if err != nil {
			result.ErrorCount++
			continue
		}

		teamID, err := i.teams.FindOrCreateTeam(ctx, rawTeam)
		if err != nil {
			result.ErrorCount++
			continue
		}
		ranking.TeamID = teamID

		inserted, err := i.rankings.Upsert(ctx, ranking)
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

func parseFifaRanking(record []string, headers map[string]int) (models.FifaRanking, string, error) {
	rawRank := field(record, headers, "rank")
	rawTeam := field(record, headers, "country_full", "country", "team", "name")
	rawDate := field(record, headers, "rank_date", "ranking_date", "date")

	if err := requireFields(map[string]string{
		"rank":         rawRank,
		"country_full": rawTeam,
		"rank_date":    rawDate,
	}); err != nil {
		return models.FifaRanking{}, "", err
	}

	rank, err := parseInt(rawRank)
	if err != nil {
		return models.FifaRanking{}, "", err
	}
	rankingDate, err := parseDate(rawDate)
	if err != nil {
		return models.FifaRanking{}, "", err
	}
	totalPoints, err := parseOptionalFloat(field(record, headers, "total_points", "total_points_2"))
	if err != nil {
		return models.FifaRanking{}, "", err
	}
	previousPoints, err := parseOptionalFloat(field(record, headers, "previous_points"))
	if err != nil {
		return models.FifaRanking{}, "", err
	}
	rankChange, err := parseOptionalInt(field(record, headers, "rank_change"))
	if err != nil {
		return models.FifaRanking{}, "", err
	}

	var confederation *string
	if value := field(record, headers, "confederation", "confederation_code"); value != "" {
		confederation = &value
	}

	return models.FifaRanking{
		RankingDate:    rankingDate,
		Rank:           rank,
		TotalPoints:    totalPoints,
		PreviousPoints: previousPoints,
		RankChange:     rankChange,
		Confederation:  confederation,
		Source:         sourceFifaRanking,
	}, rawTeam, nil
}
