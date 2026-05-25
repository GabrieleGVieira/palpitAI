# Data Pipeline - Etapa 1

## Objetivo

A Etapa 1 importa, normaliza e persiste dados históricos e externos para preparar o banco da Etapa 2. Esta etapa não calcula métricas, não gera previsões e não chama IA.

## Fontes de dados

- International football results from 1872 - 2026
  - `results.csv`
  - `goalscorers.csv`
  - `shootouts.csv`, se necessário em uma expansão futura
  - `former_names.csv`, se necessário para ampliar aliases
- FIFA Ranking Historical
  - `ranking_fifa_historical.csv`

## Tabelas

- `teams`: cadastro normalizado de seleções.
- `team_aliases`: nomes alternativos recebidos nas fontes históricas.
- `historical_matches`: jogos históricos com data, placar, torneio e local.
- `historical_goalscorers`: gols vinculados ao jogo histórico quando ele existe.
- `fifa_rankings`: ranking FIFA histórico por seleção e data.
- `data_import_logs`: execução, status e contadores de cada importação.
- `external_api_snapshots`: cache persistente de respostas externas futuras, sem Redis.

## Migrations

O SQL versionado está em:

```bash
migrations/202605240001_create_data_pipeline_tables.sql
```

No fluxo atual do backend, `database.Migrate` também cria essas tabelas. Para rodar pelo app:

```bash
go run ./cmd/api
```

Ou execute o arquivo SQL diretamente no Supabase/PostgreSQL usando o painel SQL ou `psql`.

## Importação

Configure `DATABASE_URL` no ambiente ou em `.env`. Os comandos também executam `database.Migrate` antes de importar.

Importar partidas:

```bash
go run ./cmd/import-history --type=matches --file=./data/results.csv
```

Importar gols:

```bash
go run ./cmd/import-history --type=goalscorers --file=./data/goalscorers.csv
```

Importar ranking FIFA:

```bash
go run ./cmd/import-history --type=fifa-ranking --file=./data/ranking_fifa_historical.csv
```

Cada comando imprime `processed`, `inserted`, `skipped` e `errors`, e grava o resultado em `data_import_logs`.

## Normalização de seleções

O `TeamNameNormalizer` aplica trim, compacta múltiplos espaços e resolve aliases conhecidos para nomes canônicos em português, por exemplo:

- `Brazil` -> `Brasil`
- `USA` e `United States` -> `Estados Unidos`
- `Netherlands` e `Holland` -> `Países Baixos`
- `Türkiye` e `Turkey` -> `Turquia`
- `Côte d'Ivoire` e `Ivory Coast` -> `Costa do Marfim`

Quando um nome bruto é diferente do nome canônico, o importador salva o nome bruto em `team_aliases` apontando para o registro de `teams`. A lista de aliases pode crescer usando `former_names.csv` ou regras novas sem alterar as tabelas.

## Tratamento de erros e duplicidades

- Linhas inválidas não interrompem o arquivo inteiro.
- Partidas duplicadas são ignoradas pelo índice de `match_date`, `home_team_id`, `away_team_id` e `tournament`.
- Rankings duplicados são ignorados pelo índice de `team_id` e `ranking_date`.
- Gols sem partida correspondente são ignorados e entram em `skipped`.

## Próximos passos da Etapa 2

- Calcular métricas históricas por seleção.
- Agregar estatísticas por janela temporal, mando/neutro, torneio e força do adversário.
- Usar `fifa_rankings` como feature histórica.
- Criar camada de leitura otimizada para previsão.
- Só depois disso integrar modelos ou IA.

