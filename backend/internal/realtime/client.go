package realtime

import "github.com/gorilla/websocket"

func newClient(userID string, rooms []string, conn *websocket.Conn) *client {
	return &client{
		conn:   conn,
		rooms:  clientRooms(userID, rooms),
		send:   make(chan outboundEvent, clientSendBuffer),
		userID: userID,
	}
}

func clientRooms(userID string, rooms []string) map[string]struct{} {
	clientRooms := map[string]struct{}{
		matchesRoom:             {},
		userRoomPrefix + userID: {},
	}

	if len(rooms) == 0 {
		clientRooms[rankingsRoom] = struct{}{}
	}

	for _, room := range rooms {
		if room != "" {
			clientRooms[room] = struct{}{}
		}
	}

	return clientRooms
}

func (client *client) subscribesTo(room string) bool {
	if room == "" {
		return true
	}

	_, ok := client.rooms[room]
	return ok
}
