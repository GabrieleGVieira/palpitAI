package service

import (
	"context"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/metrics/models"
	"github.com/gabrielevieira/palpitai/backend/internal/metrics/repository"
)

type Service struct {
	teamMetrics  *repository.TeamMetricsRepository
	matchFeature *repository.MatchFeaturesRepository
}

func New(db repository.Querier) *Service {
	return &Service{
		teamMetrics:  repository.NewTeamMetricsRepository(db),
		matchFeature: repository.NewMatchFeaturesRepository(db),
	}
}

func (s *Service) LatestTeamMetricBefore(ctx context.Context, teamID string, matchDate time.Time) (models.TeamMetric, error) {
	return s.teamMetrics.LatestBefore(ctx, teamID, matchDate)
}

func (s *Service) MatchFeatureByMatchID(ctx context.Context, matchID string) (models.MatchFeature, error) {
	return s.matchFeature.FindByMatchID(ctx, matchID)
}

func (s *Service) MatchFeaturesByDateRange(ctx context.Context, fromDate time.Time, toDate time.Time) ([]models.MatchFeature, error) {
	return s.matchFeature.ListByDateRange(ctx, fromDate, toDate)
}
