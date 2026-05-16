# PalpitAI Backend

API inicial em Go para o PalpitAI.

## Requisitos

- Go 1.24+

## Como rodar

```bash
cd backend
cp .env.example .env
make run
```

A API inicia em `http://localhost:3000`.

Configure `DATABASE_URL` no `.env` com a connection string PostgreSQL do Supabase. O backend adiciona `sslmode=require` automaticamente quando a URL nao informa `sslmode`.

## Rotas iniciais

```text
GET /health
GET /api/v1/status
```

As respostas incluem o status da conexao com o banco:

```json
{
  "database": "ok",
  "status": "ok"
}
```

## Comandos

```bash
make run
make test
make fmt
make vet
```
