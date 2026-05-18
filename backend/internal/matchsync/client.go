package matchsync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/domain"
	"github.com/gabrielevieira/palpitai/backend/internal/dto"
)

func (syncer *Syncer) fetchMatches(ctx context.Context, kind syncKind) ([]domain.ProviderMatch, error) {
	if err := syncer.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	endpoint, err := syncer.matchesURL(kind)
	if err != nil {
		return nil, err
	}

	requestCtx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Auth-Token", syncer.token)

	response, err := syncer.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusTooManyRequests {
		return nil, errors.New("football-data rate limit reached")
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("football-data returned status %d", response.StatusCode)
	}

	var payload dto.FootballDataResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}

	matches := make([]domain.ProviderMatch, 0, len(payload.Matches))
	for _, match := range payload.Matches {
		matches = append(matches, dto.FromFootballDataMatch(match))
	}

	return matches, nil
}

func (syncer *Syncer) matchesURL(kind syncKind) (string, error) {
	parsedURL, err := url.Parse(syncer.baseURL + "/competitions/" + syncer.competitionCode + "/matches")
	if err != nil {
		return "", err
	}

	query := parsedURL.Query()
	if syncer.season != "" {
		query.Set("season", syncer.season)
	}

	now := time.Now().UTC()
	switch kind {
	case syncLive:
		query.Set("status", "LIVE")
	case syncToday:
		today := now.Format(time.DateOnly)
		query.Set("dateFrom", today)
		query.Set("dateTo", today)
	case syncUpcoming:
		query.Set("dateFrom", now.AddDate(0, 0, 1).Format(time.DateOnly))
		query.Set("dateTo", now.Add(upcomingWindow).Format(time.DateOnly))
	default:
		return "", fmt.Errorf("unsupported sync kind %q", kind)
	}

	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
}

func (syncer *Syncer) waitRateLimit(ctx context.Context) error {
	syncer.rateMu.Lock()
	defer syncer.rateMu.Unlock()

	wait := rateLimitGap - time.Since(syncer.lastRequestAt)
	if wait > 0 {
		timer := time.NewTimer(wait)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}

	syncer.lastRequestAt = time.Now()
	return nil
}
