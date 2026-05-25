package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gabrielevieira/palpitai/backend/internal/data/models"
)

type ExternalApiSnapshotsRepository struct {
	db Querier
}

func NewExternalApiSnapshotsRepository(db Querier) *ExternalApiSnapshotsRepository {
	return &ExternalApiSnapshotsRepository{db: db}
}

func (r *ExternalApiSnapshotsRepository) GetValidSnapshot(ctx context.Context, provider string, endpoint string, now time.Time) (models.ExternalApiSnapshot, error) {
	var snapshot models.ExternalApiSnapshot
	err := r.db.QueryRow(ctx, `
		select id, provider, endpoint, payload_json, fetched_at, expires_at
		from external_api_snapshots
		where provider = $1 and endpoint = $2 and expires_at > $3
		order by fetched_at desc
		limit 1
	`, provider, endpoint, now).Scan(
		&snapshot.ID,
		&snapshot.Provider,
		&snapshot.Endpoint,
		&snapshot.PayloadJSON,
		&snapshot.FetchedAt,
		&snapshot.ExpiresAt,
	)

	return snapshot, mapNoRows(err)
}

func (r *ExternalApiSnapshotsRepository) SaveSnapshot(ctx context.Context, provider string, endpoint string, payload json.RawMessage, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		insert into external_api_snapshots (provider, endpoint, payload_json, expires_at)
		values ($1, $2, $3, $4)
	`, provider, endpoint, payload, expiresAt)

	return err
}
