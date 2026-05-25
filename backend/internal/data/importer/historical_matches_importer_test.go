package importer

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/data/models"
)

func TestParseHistoricalMatch(t *testing.T) {
	headers := headerIndex([]string{"date", "home_team", "away_team", "home_score", "away_score", "tournament", "city", "country", "neutral"})
	record := []string{"2022-12-18", "Argentina", "France", "3", "3", "FIFA World Cup", "Lusail", "Qatar", "TRUE"}

	match, err := parseHistoricalMatch(record, headers)
	if err != nil {
		t.Fatal(err)
	}

	if match.HomeScore != 3 || match.AwayScore != 3 || !match.Neutral || match.Tournament != "FIFA World Cup" {
		t.Fatalf("unexpected match parse: %#v", match)
	}
}

func TestHistoricalMatchesImporterInvalidLineDoesNotBreakImport(t *testing.T) {
	csv := `date,home_team,away_team,home_score,away_score,tournament,city,country,neutral
2022-12-18,Argentina,France,3,3,FIFA World Cup,Lusail,Qatar,TRUE
,Argentina,France,3,3,FIFA World Cup,Lusail,Qatar,TRUE
`

	result, err := NewHistoricalMatchesImporter(newFakeTeamFinder(), newFakeMatchesStore()).Import(context.Background(), strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}

	if result.ProcessedCount != 2 || result.InsertedCount != 1 || result.ErrorCount != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestHistoricalMatchesImporterDuplicateMatchIsSkipped(t *testing.T) {
	csv := `date,home_team,away_team,home_score,away_score,tournament,city,country,neutral
2022-12-18,Argentina,France,3,3,FIFA World Cup,Lusail,Qatar,TRUE
2022-12-18,Argentina,France,3,3,FIFA World Cup,Lusail,Qatar,TRUE
`

	result, err := NewHistoricalMatchesImporter(newFakeTeamFinder(), newFakeMatchesStore()).Import(context.Background(), strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}

	if result.ProcessedCount != 2 || result.InsertedCount != 1 || result.SkippedCount != 1 || result.ErrorCount != 0 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

type fakeTeamFinder struct {
	ids map[string]string
}

func newFakeTeamFinder() *fakeTeamFinder {
	return &fakeTeamFinder{ids: map[string]string{}}
}

func (f *fakeTeamFinder) FindOrCreateTeam(_ context.Context, rawName string) (string, error) {
	if id, ok := f.ids[rawName]; ok {
		return id, nil
	}
	f.ids[rawName] = rawName + "-id"
	return f.ids[rawName], nil
}

type fakeMatchesStore struct {
	matches map[string]string
}

func newFakeMatchesStore() *fakeMatchesStore {
	return &fakeMatchesStore{matches: map[string]string{}}
}

func (s *fakeMatchesStore) Upsert(_ context.Context, match models.HistoricalMatch) (string, bool, error) {
	key := match.MatchDate.Format(time.DateOnly) + "|" + match.HomeTeamID + "|" + match.AwayTeamID + "|" + match.Tournament
	if id, ok := s.matches[key]; ok {
		return id, false, nil
	}
	s.matches[key] = key

	return key, true, nil
}

func (s *fakeMatchesStore) FindByDateAndTeams(_ context.Context, matchDate time.Time, homeTeamID string, awayTeamID string) (string, error) {
	for key, id := range s.matches {
		if strings.HasPrefix(key, matchDate.Format(time.DateOnly)+"|"+homeTeamID+"|"+awayTeamID+"|") {
			return id, nil
		}
	}

	return "", errFakeNotFound
}
