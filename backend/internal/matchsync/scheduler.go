package matchsync

import (
	"context"
	"time"
)

func (syncer *Syncer) Run(ctx context.Context) {
	// 1. Faz uma primeira carga imediata para nao esperar o primeiro tick dos agendadores.
	syncer.runOnce(ctx, syncUpcoming)
	syncer.runOnce(ctx, syncToday)
	syncer.runOnce(ctx, syncLive)

	// 2. Cria tickers independentes porque jogos ao vivo precisam ser atualizados com mais frequencia.
	liveTicker := time.NewTicker(livePollInterval)
	todayTicker := time.NewTicker(todayPollInterval)
	upcomingTicker := time.NewTicker(upcomingPollInterval)
	defer liveTicker.Stop()
	defer todayTicker.Stop()
	defer upcomingTicker.Stop()

	// 3. Mantem o worker rodando ate receber cancelamento ou algum ticker disparar.
	for {
		select {
		case <-ctx.Done():
			// 4. Ao receber cancelamento, encerra o loop e deixa os defers limparem os tickers.
			return
		case <-liveTicker.C:
			// 5. Sincroniza partidas ao vivo no intervalo mais curto.
			syncer.runOnce(ctx, syncLive)
		case <-todayTicker.C:
			// 6. Sincroniza as partidas do dia em um intervalo intermediario.
			syncer.runOnce(ctx, syncToday)
		case <-upcomingTicker.C:
			// 7. Sincroniza partidas futuras em um intervalo maior.
			syncer.runOnce(ctx, syncUpcoming)
		}
	}
}

func (syncer *Syncer) runOnce(ctx context.Context, kind syncKind) {
	// 1. Executa uma sincronizacao pontual para o tipo recebido.
	syncer.logger.Info("match sync started", "kind", kind)
	summary, err := syncer.SyncOnce(ctx, kind)
	if err != nil {
		// 2. Registra falha como warning e deixa o proximo tick tentar novamente.
		syncer.logger.Warn("match sync failed", "kind", kind, "error", err)
		return
	}

	// 3. Registra o fim de toda execucao, mesmo quando o provedor nao retornou mudancas.
	syncer.logger.Info(
		"matches synced",
		"kind", kind,
		"synced_matches", summary.SyncedMatches,
		"changed_matches", summary.ChangedMatches,
		"created_events", summary.CreatedEvents,
		"updated_live_matches", summary.UpdatedLiveMatches,
		"scored_predictions", summary.ScoredPredictions,
	)
}
