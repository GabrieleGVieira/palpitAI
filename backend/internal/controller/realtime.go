package controller

import (
	"net/http"

	"github.com/gabrielevieira/palpitai/backend/internal/config"
	"github.com/gabrielevieira/palpitai/backend/internal/usecase"
)

type WebsocketHub interface {
	ServeWS(w http.ResponseWriter, r *http.Request, userID string, rooms []string)
}

type RealtimeService interface {
	RealtimePublisher
	WebsocketHub
}

func RealtimeHandler(cfg config.Config, db usecase.Datastore, hub WebsocketHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if hub == nil {
			writeError(w, http.StatusServiceUnavailable, "Realtime indisponivel.")
			return
		}

		userID, err := userIDFromToken(r, cfg, r.URL.Query().Get("token"))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Informe um token de autenticacao valido.")
			return
		}

		rooms := []string{}
		groupID := r.URL.Query().Get("group_id")
		if groupID != "" {
			if err := usecase.EnsureActiveGroupMember(r.Context(), db, userID, groupID); err != nil {
				writeError(w, http.StatusForbidden, "Você precisa participar deste grupo.")
				return
			}

			rooms = append(rooms, "group:"+groupID)
		}

		hub.ServeWS(w, r, userID, rooms)
	}
}
