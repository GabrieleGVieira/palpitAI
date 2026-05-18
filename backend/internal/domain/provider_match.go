package domain

import (
	"errors"
	"fmt"
	"strings"
)

func NormalizeProviderMatch(match ProviderMatch) ProviderMatch {
	match.ExternalID = strings.TrimSpace(match.ExternalID)
	match.HomeTeam = strings.TrimSpace(match.HomeTeam)
	match.AwayTeam = strings.TrimSpace(match.AwayTeam)
	match.Stage = strings.TrimSpace(match.Stage)
	match.Status = strings.ToLower(strings.TrimSpace(match.Status))

	if match.Stage == "" {
		match.Stage = "Copa do Mundo"
	}

	if match.Status == "" {
		match.Status = "scheduled"
	}

	return match
}

func ValidateProviderMatch(match ProviderMatch) error {
	if match.HomeTeam == "" || match.AwayTeam == "" {
		return errors.New("home_team and away_team are required")
	}

	if match.KickoffAt.IsZero() {
		return errors.New("kickoff_at is required")
	}

	switch match.Status {
	case "scheduled", "live", "finished", "postponed", "cancelled":
	default:
		return fmt.Errorf("unsupported status %q", match.Status)
	}

	if (match.HomeScore == nil) != (match.AwayScore == nil) {
		return errors.New("home_score and away_score must be provided together")
	}

	if match.HomeScore != nil && (*match.HomeScore < 0 || *match.HomeScore > 99 || *match.AwayScore < 0 || *match.AwayScore > 99) {
		return errors.New("scores must be between 0 and 99")
	}

	return nil
}
