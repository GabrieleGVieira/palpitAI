package realtime

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

func (hub *Hub) readPump(client *client) {
	defer hub.closeClient(client)

	client.conn.SetReadLimit(readLimit)
	_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		return client.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		if _, _, err := client.conn.NextReader(); err != nil {
			return
		}
	}
}

func (hub *Hub) writePump(client *client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		hub.closeClient(client)
	}()

	for {
		select {
		case event, ok := <-client.send:
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			payload, err := json.Marshal(event)
			if err != nil {
				hub.logger.Warn("websocket event marshal failed", "error", err)
				continue
			}

			if err := client.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		case <-ticker.C:
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
