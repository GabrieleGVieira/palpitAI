package realtime

import (
	"context"

	"github.com/gabrielevieira/palpitai/backend/internal/domain"
)

type Publisher interface {
	Publish(ctx context.Context, event domain.Event)
}

func (hub *Hub) Publish(_ context.Context, event domain.Event) {
	outbound := outboundEvent{
		Name:    event.Name,
		Payload: event.Payload,
		Room:    event.Room,
	}

	staleClients := []*client{}

	hub.mu.RLock()
	for client := range hub.clients {
		if !client.subscribesTo(event.Room) {
			continue
		}

		select {
		case client.send <- outbound:
		default:
			staleClients = append(staleClients, client)
		}
	}
	hub.mu.RUnlock()

	for _, client := range staleClients {
		hub.closeClient(client)
	}
}
