package matchsync

import (
	"context"
	"log/slog"

	"github.com/gabrielevieira/palpitai/backend/internal/domain"
)

type Publisher interface {
	Publish(ctx context.Context, event domain.Event)
}

type LogPublisher struct {
	logger *slog.Logger
}

func (publisher LogPublisher) Publish(_ context.Context, event domain.Event) {
	logger := publisher.logger
	if logger == nil {
		logger = slog.Default()
	}

	logger.Info("realtime event", "name", event.Name, "room", event.Room, "payload", event.Payload)
}
