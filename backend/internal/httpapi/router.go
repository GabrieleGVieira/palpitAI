package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/config"
)

type statusResponse struct {
	App       string `json:"app"`
	Database  string `json:"database"`
	Env       string `json:"env"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type databasePinger interface {
	Ping(ctx context.Context) error
}

func NewRouter(cfg config.Config, db databasePinger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler(db))
	mux.HandleFunc("GET /api/v1/status", statusHandler(cfg, db))

	return mux
}

func healthHandler(db databasePinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{
				"database": "not_configured",
				"status":   "degraded",
			})
			return
		}

		if err := db.Ping(r.Context()); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{
				"database": "unavailable",
				"status":   "degraded",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"database": "ok",
			"status":   "ok",
		})
	}
}

func statusHandler(cfg config.Config, db databasePinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		databaseStatus := "ok"
		responseStatus := "ok"

		if db == nil {
			databaseStatus = "not_configured"
			responseStatus = "degraded"
		} else if err := db.Ping(r.Context()); err != nil {
			databaseStatus = "unavailable"
			responseStatus = "degraded"
		}

		writeJSON(w, http.StatusOK, statusResponse{
			App:       "palpitai-api",
			Database:  databaseStatus,
			Env:       cfg.Env,
			Status:    responseStatus,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
