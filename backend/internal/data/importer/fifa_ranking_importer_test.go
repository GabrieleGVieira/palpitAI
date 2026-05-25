package importer

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/gabrielevieira/palpitai/backend/internal/data/models"
)

var errFakeNotFound = errors.New("not found")

func TestParseFifaRanking(t *testing.T) {
	headers := headerIndex([]string{"rank", "country_full", "country_abrv", "total_points", "previous_points", "rank_change", "confederation", "rank_date"})
	record := []string{"1", "Brazil", "BRA", "1840.77", "1832.69", "0", "CONMEBOL", "2022-12-22"}

	ranking, rawTeam, err := parseFifaRanking(record, headers)
	if err != nil {
		t.Fatal(err)
	}

	if rawTeam != "Brazil" || ranking.Rank != 1 || ranking.TotalPoints == nil || *ranking.TotalPoints != 1840.77 {
		t.Fatalf("unexpected ranking parse: rawTeam=%q ranking=%#v", rawTeam, ranking)
	}
}

func TestFifaRankingImporterDuplicateRankingIsSkipped(t *testing.T) {
	csv := `rank,country_full,total_points,previous_points,rank_change,confederation,rank_date
1,Brazil,1840.77,1832.69,0,CONMEBOL,2022-12-22
1,Brazil,1840.77,1832.69,0,CONMEBOL,2022-12-22
`

	result, err := NewFifaRankingImporter(newFakeTeamFinder(), newFakeRankingsStore()).Import(context.Background(), strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}

	if result.ProcessedCount != 2 || result.InsertedCount != 1 || result.SkippedCount != 1 || result.ErrorCount != 0 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestFifaRankingImporterInvalidLineDoesNotBreakImport(t *testing.T) {
	csv := `rank,country_full,total_points,previous_points,rank_change,confederation,rank_date
1,Brazil,1840.77,1832.69,0,CONMEBOL,2022-12-22
,Brazil,1840.77,1832.69,0,CONMEBOL,2022-12-22
`

	result, err := NewFifaRankingImporter(newFakeTeamFinder(), newFakeRankingsStore()).Import(context.Background(), strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}

	if result.ProcessedCount != 2 || result.InsertedCount != 1 || result.ErrorCount != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

type fakeRankingsStore struct {
	rankings map[string]struct{}
}

func newFakeRankingsStore() *fakeRankingsStore {
	return &fakeRankingsStore{rankings: map[string]struct{}{}}
}

func (s *fakeRankingsStore) Upsert(_ context.Context, ranking models.FifaRanking) (bool, error) {
	key := ranking.TeamID + "|" + ranking.RankingDate.Format("2006-01-02")
	if _, ok := s.rankings[key]; ok {
		return false, nil
	}
	s.rankings[key] = struct{}{}

	return true, nil
}
