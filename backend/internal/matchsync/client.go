package matchsync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/domain"
	"github.com/gabrielevieira/palpitai/backend/internal/dto"
)

func (syncer *Syncer) fetchMatches(ctx context.Context, kind syncKind) ([]domain.ProviderMatch, error) {
	// 1. Respeita o intervalo minimo entre chamadas para nao estourar o limite do provedor.
	if err := syncer.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	// 2. Monta a URL de consulta de acordo com o tipo de sincronizacao.
	endpoint, err := syncer.matchesURL(kind)
	if err != nil {
		return nil, err
	}
	syncer.logger.Info("fetching provider matches", "kind", kind, "url", endpoint)

	// 3. Cria um contexto com timeout apenas para a chamada HTTP externa.
	requestCtx, cancel := context.WithTimeout(ctx, defaultRequestTimeout)
	defer cancel()

	// 4. Cria a requisicao autenticada para a API football-data.org.
	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Auth-Token", syncer.token)

	// 5. Executa a chamada e garante fechamento do corpo da resposta.
	response, err := syncer.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 6. Trata limites e status HTTP inesperados como erros do provedor.
	if response.StatusCode == http.StatusTooManyRequests {
		return nil, errors.New("football-data rate limit reached")
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("football-data returned status %d", response.StatusCode)
	}

	// 7. Decodifica o JSON bruto no DTO que representa a resposta do football-data.
	var payload dto.FootballDataResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}

	// 8. Converte cada jogo do DTO externo para o modelo interno usado pelo dominio.
	matches := make([]domain.ProviderMatch, 0, len(payload.Matches))
	for _, match := range payload.Matches {
		matches = append(matches, dto.FromFootballDataMatch(match))
	}

	syncer.logger.Info("provider matches fetched", "kind", kind, "matches", len(matches))
	return matches, nil
}

func (syncer *Syncer) matchesURL(kind syncKind) (string, error) {
	// 1. Parte da URL base configurada e aponta para as partidas da competicao.
	parsedURL, err := url.Parse(syncer.baseURL + "/competitions/" + syncer.competitionCode + "/matches")
	if err != nil {
		return "", err
	}

	// 2. Prepara a query string e inclui a temporada quando ela foi configurada.
	query := parsedURL.Query()
	if syncer.season != "" {
		query.Set("season", syncer.season)
	}

	// 3. Calcula filtros de data/status em UTC para manter consistencia com a API externa.
	now := time.Now().UTC()
	switch kind {
	case syncLive:
		// 4. Para o polling live, busca uma janela recente por data em vez de status=LIVE.
		// Isso captura partidas que estavam ao vivo no banco, mas ja viraram finished no provedor.
		query.Set("dateFrom", now.Add(-liveRecentWindow).Format(time.DateOnly))
		query.Set("dateTo", now.Format(time.DateOnly))
	case syncToday:
		// 5. Para jogos de hoje, usa a mesma data no inicio e no fim do intervalo.
		today := now.Format(time.DateOnly)
		query.Set("dateFrom", today)
		query.Set("dateTo", today)
	case syncUpcoming:
		// 6. Para proximos jogos, consulta de amanha ate a janela futura configurada.
		query.Set("dateFrom", now.AddDate(0, 0, 1).Format(time.DateOnly))
		query.Set("dateTo", now.Add(upcomingWindow).Format(time.DateOnly))
	default:
		// 7. Rejeita tipos desconhecidos para evitar uma consulta ampla sem querer.
		return "", fmt.Errorf("unsupported sync kind %q", kind)
	}

	// 8. Codifica os parametros e retorna a URL final.
	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
}

func (syncer *Syncer) waitRateLimit(ctx context.Context) error {
	// 1. Serializa o controle de rate limit para que chamadas concorrentes compartilhem o mesmo relogio.
	syncer.rateMu.Lock()
	defer syncer.rateMu.Unlock()

	// 2. Calcula quanto ainda falta para atingir o intervalo minimo entre requisicoes.
	wait := rateLimitGap - time.Since(syncer.lastRequestAt)
	if wait > 0 {
		// 3. Espera o tempo restante, mas respeita cancelamento do contexto.
		timer := time.NewTimer(wait)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}

	// 4. Registra o horario da requisicao que esta prestes a ser feita.
	syncer.lastRequestAt = time.Now()
	return nil
}
