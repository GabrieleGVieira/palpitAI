package repository

import (
	"context"

	"github.com/gabrielevieira/palpitai/backend/internal/data/models"
)

type ImportLogsRepository struct {
	db Querier
}

func NewImportLogsRepository(db Querier) *ImportLogsRepository {
	return &ImportLogsRepository{db: db}
}

func (r *ImportLogsRepository) Start(ctx context.Context, importType string, filePath string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		insert into data_import_logs (import_type, file_path, status)
		values ($1, $2, 'running')
		returning id
	`, importType, filePath).Scan(&id)

	return id, err
}

func (r *ImportLogsRepository) Finish(ctx context.Context, id string, status string, summary models.ImportSummary, errorMessage *string) error {
	_, err := r.db.Exec(ctx, `
		update data_import_logs
		set
			status = $2,
			processed_count = $3,
			inserted_count = $4,
			skipped_count = $5,
			error_count = $6,
			error_message = $7,
			finished_at = now()
		where id = $1
	`, id, status, summary.ProcessedCount, summary.InsertedCount, summary.SkippedCount, summary.ErrorCount, errorMessage)

	return err
}
