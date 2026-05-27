package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/cache"
	"github.com/gabrielevieira/palpitai/backend/internal/config"
	"github.com/gabrielevieira/palpitai/backend/internal/database"
	"github.com/gabrielevieira/palpitai/backend/internal/matchsync"
	"github.com/gabrielevieira/palpitai/backend/internal/realtime"
)

func main() {
	// 1. Carrega as variaveis de configuracao usadas pelo worker.
	cfg := config.Load()
	// 2. Cria um logger JSON para registrar o ciclo de vida do processo.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 3. Limita o tempo das operacoes de inicializacao, como conectar e migrar o banco.
	startupCtx, startupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startupCancel()

	// 4. Abre o pool de conexoes com o Postgres; se falhar, encerra o worker.
	db, err := database.NewPostgresPool(startupCtx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// 5. Garante que as migracoes necessarias estejam aplicadas antes de sincronizar jogos.
	if err := database.Migrate(startupCtx, db); err != nil {
		logger.Error("database migration failed", "error", err)
		os.Exit(1)
	}

	// 6. Conecta no Redis para publicar eventos que a API repassa aos WebSockets.
	redisClient, err := cache.NewRedisClient(startupCtx, cfg.RedisURL)
	if err != nil {
		logger.Error("redis connection failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Error("redis close failed", "error", err)
		}
	}()

	// 7. Monta o sincronizador; ele fica desabilitado quando o token do provedor nao existe.
	syncer, enabled := matchsync.New(cfg, db, logger)
	if !enabled {
		logger.Error("match sync disabled", "reason", "FOOTBALL_DATA_TOKEN not configured")
		os.Exit(1)
	}
	syncer.SetPublisher(realtime.NewRedisPublisher(redisClient, logger))

	// 8. Cria um contexto cancelavel por Ctrl+C ou SIGTERM para desligamento gracioso.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 9. Inicia o loop de sincronizacao e bloqueia ate o contexto ser cancelado.
	logger.Info("match sync worker started", "provider", "football-data.org", "competition", cfg.FootballDataCompetitionCode)

	go syncer.Run(ctx)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	logger.Info("match sync health server started", "port", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Error("health server failed", "error", err)
		os.Exit(1)
	}

	logger.Info("match sync worker stopped")
}
