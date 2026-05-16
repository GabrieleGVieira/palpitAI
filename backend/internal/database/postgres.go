package database

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	cfg, err := pgxpool.ParseConfig(withSSLMode(databaseURL))
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 5
	cfg.MinConns = 1
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func withSSLMode(databaseURL string) string {
	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		return databaseURL
	}

	if parsedURL.Query().Has("sslmode") {
		return databaseURL
	}

	query := parsedURL.Query()
	query.Set("sslmode", "require")
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String()
}
