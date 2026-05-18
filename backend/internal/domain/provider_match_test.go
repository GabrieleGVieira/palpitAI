package domain

import (
	"testing"
	"time"
)

func TestNormalizeProviderMatch(t *testing.T) {
	match := NormalizeProviderMatch(ProviderMatch{
		AwayTeam:  " Croácia ",
		HomeTeam:  " Brasil ",
		KickoffAt: time.Now(),
		Status:    "LIVE",
	})

	if match.HomeTeam != "Brasil" {
		t.Fatalf("expected trimmed home team, got %q", match.HomeTeam)
	}
	if match.Status != "live" {
		t.Fatalf("expected normalized status live, got %q", match.Status)
	}
	if match.Stage != "Copa do Mundo" {
		t.Fatalf("expected default stage, got %q", match.Stage)
	}
}

func TestValidateProviderMatchRejectsPartialScore(t *testing.T) {
	homeScore := 1
	match := ProviderMatch{
		AwayTeam:  "Croácia",
		HomeScore: &homeScore,
		HomeTeam:  "Brasil",
		KickoffAt: time.Now(),
		Status:    "live",
	}

	if err := ValidateProviderMatch(match); err == nil {
		t.Fatal("expected partial score validation error")
	}
}
