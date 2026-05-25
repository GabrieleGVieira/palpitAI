package normalizer

import (
	"context"
	"strings"

	"github.com/gabrielevieira/palpitai/backend/internal/data/repository"
)

type TeamRepository interface {
	FindByName(ctx context.Context, name string) (string, error)
	FindByAlias(ctx context.Context, alias string) (string, error)
	Create(ctx context.Context, name string) (string, error)
	CreateAlias(ctx context.Context, teamID string, alias string) error
}

type TeamNameNormalizer struct {
	teams   TeamRepository
	aliases map[string]string
}

func NewTeamNameNormalizer(teams TeamRepository) *TeamNameNormalizer {
	return &TeamNameNormalizer{
		teams: teams,
		aliases: map[string]string{
			"Brazil":         "Brasil",
			"Germany":        "Alemanha",
			"United States":  "Estados Unidos",
			"USA":            "Estados Unidos",
			"Argentina":      "Argentina",
			"France":         "França",
			"Spain":          "Espanha",
			"Portugal":       "Portugal",
			"England":        "Inglaterra",
			"Netherlands":    "Países Baixos",
			"Holland":        "Países Baixos",
			"Türkiye":        "Turquia",
			"Turkey":         "Turquia",
			"Czechia":        "República Tcheca",
			"Czech Republic": "República Tcheca",
			"South Korea":    "Coreia do Sul",
			"Korea Republic": "Coreia do Sul",
			"North Korea":    "Coreia do Norte",
			"Saudi Arabia":   "Arábia Saudita",
			"Ivory Coast":    "Costa do Marfim",
			"Côte d'Ivoire":  "Costa do Marfim",
			"Mexico":         "México",
			"Japan":          "Japão",
			"Morocco":        "Marrocos",
			"Switzerland":    "Suíça",
			"Belgium":        "Bélgica",
			"Croatia":        "Croácia",
			"Uruguay":        "Uruguai",
		},
	}
}

func (n *TeamNameNormalizer) Normalize(rawName string) string {
	name := cleanName(rawName)
	if canonical, ok := n.aliases[name]; ok {
		return canonical
	}
	for alias, canonical := range n.aliases {
		if strings.EqualFold(alias, name) {
			return canonical
		}
	}

	return name
}

func (n *TeamNameNormalizer) FindOrCreateTeam(ctx context.Context, rawName string) (string, error) {
	rawAlias := cleanName(rawName)
	if rawAlias == "" {
		return "", repository.ErrNotFound
	}

	canonicalName := n.Normalize(rawAlias)
	teamID, err := n.teams.FindByName(ctx, canonicalName)
	if err == nil {
		return teamID, n.ensureAlias(ctx, teamID, rawAlias, canonicalName)
	}
	if err != repository.ErrNotFound {
		return "", err
	}

	teamID, err = n.teams.FindByAlias(ctx, rawAlias)
	if err == nil {
		return teamID, nil
	}
	if err != repository.ErrNotFound {
		return "", err
	}

	teamID, err = n.teams.Create(ctx, canonicalName)
	if err != nil {
		return "", err
	}

	return teamID, n.ensureAlias(ctx, teamID, rawAlias, canonicalName)
}

func (n *TeamNameNormalizer) ensureAlias(ctx context.Context, teamID string, rawAlias string, canonicalName string) error {
	if rawAlias == "" || rawAlias == canonicalName {
		return nil
	}

	return n.teams.CreateAlias(ctx, teamID, rawAlias)
}

func cleanName(rawName string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(rawName)), " ")
}
