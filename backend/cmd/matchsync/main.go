package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/config"
	"github.com/gabrielevieira/palpitai/backend/internal/database"
	"github.com/gabrielevieira/palpitai/backend/internal/matchsync"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	startupCtx, startupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startupCancel()

	db, err := database.NewPostgresPool(startupCtx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.Migrate(startupCtx, db); err != nil {
		logger.Error("database migration failed", "error", err)
		os.Exit(1)
	}

	syncer, enabled := matchsync.New(cfg, db, logger)
	if !enabled {
		logger.Error("match sync disabled", "reason", "FOOTBALL_DATA_TOKEN not configured")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("match sync worker started", "provider", "football-data.org", "competition", cfg.FootballDataCompetitionCode)
	syncer.Run(ctx)
	logger.Info("match sync worker stopped")
}
