package matchsync

import (
	"context"
	"fmt"

	"github.com/gabrielevieira/palpitai/backend/internal/domain"
	"github.com/gabrielevieira/palpitai/backend/internal/repositories"
)

func (syncer *Syncer) publishMatchChanged(ctx context.Context, previous domain.MatchSnapshot, match domain.ProviderMatch) {
	payload := map[string]any{
		"away_score":      match.AwayScore,
		"away_team":       match.AwayTeam,
		"external_id":     match.ExternalID,
		"home_score":      match.HomeScore,
		"home_team":       match.HomeTeam,
		"kickoff_at":      match.KickoffAt,
		"message":         resultMessage(match.HomeTeam, match.AwayTeam, match.HomeScore, match.AwayScore),
		"previous_score":  scorePair(previous.HomeScore, previous.AwayScore),
		"previous_status": previous.Status,
		"status":          match.Status,
	}

	syncer.publisher.Publish(ctx, domain.Event{
		Name:    "match.updated",
		Payload: payload,
		Room:    "matches",
	})

	if match.Status == "finished" {
		syncer.publisher.Publish(ctx, domain.Event{
			Name:    "match.finished",
			Payload: payload,
			Room:    "matches",
		})
	}
}

func (syncer *Syncer) publishRankingChanged(ctx context.Context, matchID string, match domain.ProviderMatch) error {
	groups, err := repositories.AffectedGroupsByMatch(ctx, syncer.db, matchID)
	if err != nil {
		return err
	}

	for _, group := range groups {
		payload := map[string]any{
			"away_score": match.AwayScore,
			"away_team":  match.AwayTeam,
			"group_id":   group.ID,
			"group_name": group.Name,
			"home_score": match.HomeScore,
			"home_team":  match.HomeTeam,
			"match_id":   matchID,
			"message":    "Ranking do grupo " + group.Name + " atualizado",
		}

		syncer.publisher.Publish(ctx, domain.Event{
			Name:    "ranking.updated",
			Payload: payload,
			Room:    "rankings",
		})
		syncer.publisher.Publish(ctx, domain.Event{
			Name:    "ranking.updated",
			Payload: payload,
			Room:    "group:" + group.ID,
		})
	}

	return nil
}

func scorePair(homeScore *int, awayScore *int) map[string]*int {
	return map[string]*int{
		"away": awayScore,
		"home": homeScore,
	}
}

func resultMessage(homeTeam string, awayTeam string, homeScore *int, awayScore *int) string {
	if homeScore == nil || awayScore == nil {
		return homeTeam + " x " + awayTeam + " - resultado final lancado"
	}

	return fmt.Sprintf("%s %dx%d %s - resultado final lancado", homeTeam, *homeScore, *awayScore, awayTeam)
}
