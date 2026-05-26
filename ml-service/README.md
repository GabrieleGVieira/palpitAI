# PalpitAI Metrics Engine

Etapa 2 calcula métricas agregadas de seleções e features de partidas para serem consumidas pelo backend Go na Etapa 3.

## Fonte histórica local

Os CSVs em `ml-service/data` são usados como fonte histórica estática e não são persistidos integralmente no PostgreSQL. Os scripts carregam esses arquivos uma vez por execução, calculam as métricas em memória e salvam somente os agregados em `team_metrics`, `team_metric_snapshots` e `match_features`.

O arquivo `data/team_aliases.json` centraliza os aliases usados para casar nomes dos CSVs, nomes traduzidos do backend e times existentes no banco. Ele foi montado a partir do mapper em `backend/internal/utils/mapper.go` e de `data/former_names.csv`.

## Configuração

```bash
cd ml-service
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

Variáveis:

```bash
export DATABASE_URL='postgresql://...'
```

## Calcular métricas de seleções

```bash
python ml-service/app/scripts/calculate_team_metrics.py --metric-date=2026-06-01
```

O script carrega `data/results.csv`, `data/goalscorers.csv` e `data/shootouts.csv`, normaliza os nomes com `data/team_aliases.json`, filtra jogos até `metric-date`, calcula métricas apenas para os times-alvo do mapper e faz upsert em `team_metrics`. Partidas contra seleções fora do mapper continuam entrando no histórico com IDs locais em memória, sem serem salvas como times no banco. Também grava snapshots agregados em `team_metric_snapshots` com componentes como profundidade de artilharia, dependência de pênaltis e aproveitamento em disputas por pênaltis.

## Calcular features de partidas

```bash
python ml-service/app/scripts/calculate_match_features.py --from-date=2026-06-01 --to-date=2026-07-31
```

O script busca partidas alvo em `world_cup_matches`, resolve os nomes com `data/team_aliases.json`, usa as métricas mais recentes antes da data do jogo, busca o ranking FIFA anterior à data da partida em `data/ranking_fifa_historical.csv` e salva em `match_features`.

## Job diário

```bash
python ml-service/app/scripts/run_daily_metrics_job.py --feature-to-date=2026-07-31
```

Por padrão, o job usa a data atual como `metric-date`, inclui partidas finalizadas de `world_cup_matches` no histórico em memória, recalcula `team_metrics` e recalcula `match_features` a partir do dia seguinte até `feature-to-date`. Para reprocessar uma data específica:

```bash
python ml-service/app/scripts/run_daily_metrics_job.py --metric-date=2026-06-20 --feature-to-date=2026-07-31
```

## Tabelas preenchidas

- `team_metrics`
- `team_metric_snapshots`
- `match_features`

## Próximos passos

Na Etapa 3, o backend pode consumir `match_features` para alimentar modelos estatísticos ou ML real. Esta etapa não chama LLM, não gera explicações e não produz previsões finais.
