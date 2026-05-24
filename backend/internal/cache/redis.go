package cache

import (
	"context"
	"crypto/tls"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(ctx context.Context, redisURL string) (*RedisClient, error) {
	if strings.TrimSpace(redisURL) == "" {
		return nil, errors.New("REDIS_URL is required")
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	opts.TLSConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	opts.DialTimeout = 5 * time.Second
	opts.ReadTimeout = 3 * time.Second
	opts.WriteTimeout = 3 * time.Second
	opts.PoolSize = 5
	opts.MinIdleConns = 1

	client := redis.NewClient(opts)
	return &RedisClient{client: client}, nil
}

func (cache *RedisClient) Publish(ctx context.Context, channel string, payload []byte) error {
	return cache.client.Publish(ctx, channel, payload).Err()
}

func (cache *RedisClient) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return cache.client.Subscribe(ctx, channel)
}

func (cache *RedisClient) Close() error {
	return cache.client.Close()
}
