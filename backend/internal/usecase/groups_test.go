package usecase

import (
	"testing"

	"github.com/gabrielevieira/palpitai/backend/internal/apperrors"
	"github.com/gabrielevieira/palpitai/backend/internal/dto"
)

func TestNormalizeCreateGroupRequestRequiresSelectedTeams(t *testing.T) {
	limit := 10
	_, err := NormalizeCreateGroupRequest(dto.CreateGroupRequest{
		Name:             "Copa da firma",
		MatchScope:       "selected",
		ParticipantLimit: &limit,
	})
	if !apperrors.IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestNormalizeCreateGroupRequestClearsTeamsForAllMatches(t *testing.T) {
	request, err := NormalizeCreateGroupRequest(dto.CreateGroupRequest{
		HasUnlimitedParticipants: true,
		MatchScope:               "all",
		Name:                     "Copa da firma",
		SelectedTeams:            []string{"Brasil"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(request.SelectedTeams) != 0 {
		t.Fatalf("expected selected teams to be cleared, got %v", request.SelectedTeams)
	}
	if request.ParticipantLimit != nil {
		t.Fatalf("expected participant limit to be nil")
	}
}

func TestNormalizeInviteCode(t *testing.T) {
	got := NormalizeInviteCode(" abcd-1234 ")
	if got != "ABCD1234" {
		t.Fatalf("expected ABCD1234, got %q", got)
	}
}
