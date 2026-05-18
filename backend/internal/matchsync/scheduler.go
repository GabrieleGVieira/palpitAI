package matchsync

import (
	"context"
	"time"
)

func (syncer *Syncer) Run(ctx context.Context) {
	syncer.runOnce(ctx, syncUpcoming)
	syncer.runOnce(ctx, syncToday)
	syncer.runOnce(ctx, syncLive)

	liveTicker := time.NewTicker(livePollInterval)
	todayTicker := time.NewTicker(todayPollInterval)
	upcomingTicker := time.NewTicker(upcomingPollInterval)
	defer liveTicker.Stop()
	defer todayTicker.Stop()
	defer upcomingTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-liveTicker.C:
			syncer.runOnce(ctx, syncLive)
		case <-todayTicker.C:
			syncer.runOnce(ctx, syncToday)
		case <-upcomingTicker.C:
			syncer.runOnce(ctx, syncUpcoming)
		}
	}
}

func (syncer *Syncer) runOnce(ctx context.Context, kind syncKind) {
	summary, err := syncer.SyncOnce(ctx, kind)
	if err != nil {
		syncer.logger.Warn("match sync failed", "kind", kind, "error", err)
		return
	}

	if summary.ChangedMatches > 0 || summary.CreatedEvents > 0 {
		syncer.logger.Info(
			"matches synced",
			"kind", kind,
			"synced_matches", summary.SyncedMatches,
			"changed_matches", summary.ChangedMatches,
			"created_events", summary.CreatedEvents,
			"updated_live_matches", summary.UpdatedLiveMatches,
			"scored_predictions", summary.ScoredPredictions,
		)
	}
}
