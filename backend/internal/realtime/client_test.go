package realtime

import "testing"

func TestClientRoomsDefaultSubscriptions(t *testing.T) {
	rooms := clientRooms("user-id", nil)

	expectedRooms := []string{
		matchesRoom,
		rankingsRoom,
		userRoomPrefix + "user-id",
	}

	for _, room := range expectedRooms {
		if _, ok := rooms[room]; !ok {
			t.Fatalf("expected room %q", room)
		}
	}
}

func TestClientRoomsWithGroupDoesNotSubscribeToGlobalRankings(t *testing.T) {
	rooms := clientRooms("user-id", []string{"group:123"})

	if _, ok := rooms[rankingsRoom]; ok {
		t.Fatalf("expected group-scoped client to skip global rankings room")
	}
	if _, ok := rooms["group:123"]; !ok {
		t.Fatalf("expected group room")
	}
}

func TestClientSubscribesToBroadcast(t *testing.T) {
	client := newClient("user-id", nil, nil)

	if !client.subscribesTo("") {
		t.Fatal("expected client to subscribe to broadcast events")
	}
}
