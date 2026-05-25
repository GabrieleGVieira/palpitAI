package normalizer

import (
	"context"
	"testing"

	"github.com/gabrielevieira/palpitai/backend/internal/data/repository"
)

func TestTeamNameNormalizerNormalize(t *testing.T) {
	normalizer := NewTeamNameNormalizer(newFakeTeamsRepo())

	tests := map[string]string{
		"  Brazil  ":        "Brasil",
		"United   States":   "Estados Unidos",
		"brazil":            "Brasil",
		"Holland":           "Países Baixos",
		"Côte d'Ivoire":     "Costa do Marfim",
		"Unknown   Country": "Unknown Country",
	}

	for raw, want := range tests {
		if got := normalizer.Normalize(raw); got != want {
			t.Fatalf("Normalize(%q) = %q, want %q", raw, got, want)
		}
	}
}

func TestTeamNameNormalizerFindOrCreateTeamCreatesAlias(t *testing.T) {
	repo := newFakeTeamsRepo()
	normalizer := NewTeamNameNormalizer(repo)

	id, err := normalizer.FindOrCreateTeam(context.Background(), "Brazil")
	if err != nil {
		t.Fatal(err)
	}
	if repo.names["Brasil"] != id {
		t.Fatalf("canonical team was not created")
	}
	if repo.aliases["Brazil"] != id {
		t.Fatalf("alias was not created")
	}
}

type fakeTeamsRepo struct {
	nextID  int
	names   map[string]string
	aliases map[string]string
}

func newFakeTeamsRepo() *fakeTeamsRepo {
	return &fakeTeamsRepo{
		names:   map[string]string{},
		aliases: map[string]string{},
	}
}

func (r *fakeTeamsRepo) FindByName(_ context.Context, name string) (string, error) {
	id, ok := r.names[name]
	if !ok {
		return "", repository.ErrNotFound
	}

	return id, nil
}

func (r *fakeTeamsRepo) FindByAlias(_ context.Context, alias string) (string, error) {
	id, ok := r.aliases[alias]
	if !ok {
		return "", repository.ErrNotFound
	}

	return id, nil
}

func (r *fakeTeamsRepo) Create(_ context.Context, name string) (string, error) {
	if id, ok := r.names[name]; ok {
		return id, nil
	}
	r.nextID++
	id := name
	r.names[name] = id

	return id, nil
}

func (r *fakeTeamsRepo) CreateAlias(_ context.Context, teamID string, alias string) error {
	r.aliases[alias] = teamID
	return nil
}
