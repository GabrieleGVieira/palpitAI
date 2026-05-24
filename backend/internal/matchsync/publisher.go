package matchsync

import (
	"context"
	"log/slog"

	"github.com/gabrielevieira/palpitai/backend/internal/domain"
)

type Publisher interface {
	// Publish envia um evento de dominio para o destino configurado.
	Publish(ctx context.Context, event domain.Event)
}

type LogPublisher struct {
	logger *slog.Logger
}

func (publisher LogPublisher) Publish(_ context.Context, event domain.Event) {
	// 1. Usa o logger injetado no publisher quando ele existe.
	logger := publisher.logger
	if logger == nil {
		// 2. Se nenhum logger foi configurado, usa o logger padrao para evitar panic.
		logger = slog.Default()
	}

	// 3. Registra o evento com nome, sala e payload para depuracao do fluxo realtime.
	logger.Info("realtime event", "name", event.Name, "room", event.Room, "payload", event.Payload)
}
