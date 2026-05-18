package realtime

import (
	"log/slog"
	"sync"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/dto"
	"github.com/gorilla/websocket"
)

const (
	clientSendBuffer = 16
	pingPeriod       = 45 * time.Second
	pongWait         = 60 * time.Second
	readLimit        = 512
	writeWait        = 10 * time.Second
)

const (
	matchesRoom    = "matches"
	rankingsRoom   = "rankings"
	userRoomPrefix = "user:"
)

type Hub struct {
	clients  map[*client]struct{}
	logger   *slog.Logger
	mu       sync.RWMutex
	upgrader websocket.Upgrader
}

type client struct {
	conn   *websocket.Conn
	rooms  map[string]struct{}
	send   chan outboundEvent
	userID string
}

type outboundEvent = dto.RealtimeEvent
