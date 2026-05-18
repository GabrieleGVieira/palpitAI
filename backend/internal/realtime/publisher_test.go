package realtime

import (
	"context"
	"testing"

	"github.com/gabrielevieira/palpitai/backend/internal/domain"
)

func TestPublishSendsOnlyToSubscribedClients(t *testing.T) {
	hub := NewHub(nil)
	subscribed := newClient("user-1", []string{"group:1"}, nil)
	notSubscribed := newClient("user-2", []string{"group:2"}, nil)
	hub.registerClient(subscribed)
	hub.registerClient(notSubscribed)

	hub.Publish(context.Background(), domain.Event{
		Name: "ranking.updated",
		Room: "group:1",
	})

	select {
	case event := <-subscribed.send:
		if event.Name != "ranking.updated" {
			t.Fatalf("expected ranking.updated, got %q", event.Name)
		}
	default:
		t.Fatal("expected subscribed client to receive event")
	}

	select {
	case event := <-notSubscribed.send:
		t.Fatalf("did not expect event for unsubscribed client: %v", event)
	default:
	}
}

func TestPublishClosesStaleClients(t *testing.T) {
	hub := NewHub(nil)
	stale := newClient("user-1", nil, nil)
	close(stale.send)
	stale.send = make(chan outboundEvent)
	hub.registerClient(stale)

	hub.Publish(context.Background(), domain.Event{Name: "event"})

	hub.mu.RLock()
	_, ok := hub.clients[stale]
	hub.mu.RUnlock()
	if ok {
		t.Fatal("expected stale client to be removed")
	}
}
