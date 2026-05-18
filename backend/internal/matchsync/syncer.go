package matchsync

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/config"
	"github.com/gabrielevieira/palpitai/backend/internal/domain"
	"github.com/gabrielevieira/palpitai/backend/internal/repositories"
)

type Syncer struct {
	baseURL         string
	competitionCode string
	db              datastore
	httpClient      *http.Client
	inFlight        sync.Mutex
	lastRequestAt   time.Time
	logger          *slog.Logger
	publisher       Publisher
	rateMu          sync.Mutex
	season          string
	token           string
}

func New(cfg config.Config, db datastore, logger *slog.Logger) (*Syncer, bool) {
	if strings.TrimSpace(cfg.FootballDataToken) == "" {
		return nil, false
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &Syncer{
		baseURL:         strings.TrimRight(cfg.FootballDataAPIBaseURL, "/"),
		competitionCode: cfg.FootballDataCompetitionCode,
		db:              db,
		httpClient:      http.DefaultClient,
		logger:          logger,
		publisher:       LogPublisher{logger: logger},
		season:          cfg.FootballDataSeason,
		token:           cfg.FootballDataToken,
	}, true
}

func (syncer *Syncer) SetPublisher(publisher Publisher) {
	if publisher != nil {
		syncer.publisher = publisher
	}
}

func (syncer *Syncer) SyncOnce(ctx context.Context, kind syncKind) (domain.SyncSummary, error) {
	if !syncer.inFlight.TryLock() {
		return domain.SyncSummary{}, nil
	}
	defer syncer.inFlight.Unlock()

	shouldSync, err := syncer.shouldSync(ctx, kind)
	if err != nil {
		return domain.SyncSummary{}, err
	}
	if !shouldSync {
		return domain.SyncSummary{}, nil
	}

	matches, err := syncer.fetchMatches(ctx, kind)
	if err != nil {
		return domain.SyncSummary{}, err
	}

	summary := domain.SyncSummary{SyncedMatches: len(matches)}
	for _, match := range matches {
		match = domain.NormalizeProviderMatch(match)
		if err := domain.ValidateProviderMatch(match); err != nil {
			syncer.logger.Warn("provider match ignored", "error", err)
			continue
		}

		matchSummary, err := syncer.syncMatch(ctx, match)
		if err != nil {
			return domain.SyncSummary{}, err
		}

		summary.ChangedMatches += matchSummary.ChangedMatches
		summary.CreatedEvents += matchSummary.CreatedEvents
		summary.ScoredPredictions += matchSummary.ScoredPredictions
		summary.UpdatedLiveMatches += matchSummary.UpdatedLiveMatches
	}

	return summary, nil
}

func (syncer *Syncer) shouldSync(ctx context.Context, kind syncKind) (bool, error) {
	if kind != syncLive {
		return true, nil
	}

	return repositories.HasLiveOrSoonMatches(ctx, syncer.db)
}

func (syncer *Syncer) syncMatch(ctx context.Context, match domain.ProviderMatch) (domain.SyncSummary, error) {
	snapshot, err := repositories.MatchSnapshotByProviderMatch(ctx, syncer.db, match)
	if err != nil && !errors.Is(err, repositories.ErrNotFound) {
		return domain.SyncSummary{}, err
	}

	changedRows, matchID, err := repositories.UpsertProviderMatch(ctx, syncer.db, match)
	if err != nil {
		return domain.SyncSummary{}, err
	}

	summary := domain.SyncSummary{ChangedMatches: changedRows}
	if changedRows > 0 {
		syncer.publishMatchChanged(ctx, snapshot, match)
	}

	createdEvents, err := syncer.syncGoals(ctx, matchID, match)
	if err != nil {
		return domain.SyncSummary{}, err
	}
	summary.CreatedEvents = createdEvents

	if match.HomeScore == nil || match.AwayScore == nil || (match.Status != "live" && match.Status != "finished") {
		return summary, nil
	}

	scoredPredictions, err := repositories.ScoreProviderMatchPredictions(ctx, syncer.db, matchID, *match.HomeScore, *match.AwayScore)
	if err != nil {
		return domain.SyncSummary{}, err
	}

	if scoredPredictions > 0 && changedRows > 0 {
		if err := syncer.publishRankingChanged(ctx, matchID, match); err != nil {
			return domain.SyncSummary{}, err
		}
	}

	summary.ScoredPredictions = scoredPredictions
	if match.Status == "live" {
		summary.UpdatedLiveMatches = 1
	}

	return summary, nil
}

func (syncer *Syncer) syncGoals(ctx context.Context, matchID string, match domain.ProviderMatch) (int, error) {
	created := 0
	for _, goal := range match.Goals {
		wasCreated, err := repositories.InsertGoalEvent(ctx, syncer.db, matchID, goal)
		if err != nil {
			return 0, err
		}

		if !wasCreated {
			continue
		}

		created++
		syncer.publisher.Publish(ctx, domain.Event{
			Name: "match.goal",
			Payload: map[string]any{
				"away_score":  goal.AwayScore,
				"away_team":   match.AwayTeam,
				"home_score":  goal.HomeScore,
				"home_team":   match.HomeTeam,
				"match_id":    matchID,
				"minute":      goal.Minute,
				"player_name": goal.PlayerName,
				"team_name":   goal.TeamName,
				"type":        goal.Type,
			},
			Room: "match:" + matchID,
		})
	}

	return created, nil
}
