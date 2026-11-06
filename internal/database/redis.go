package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"go-produk/internal/config"
)

// RedisClient adalah struct untuk koneksi ke Redis
type RedisClient struct {
	client *redis.Client
	cfg    config.RedisConfig
}

// NewRedisClient membuat koneksi baru ke Redis
func NewRedisClient(ctx context.Context, cfg config.RedisConfig) (*RedisClient, error) {
	opts, err := redis.ParseURL(cfg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	opts.PoolSize = cfg.PoolSize

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &RedisClient{
		client: client,
		cfg:    cfg,
	}, nil
}

// Close menutup koneksi
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Client mengembalikan Redis client untuk digunakan
func (r *RedisClient) Client() *redis.Client {
	return r.client
}

// Ping test koneksi
func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// AddToBlacklist menyimpan token ke blacklist dengan expiration TTL
func (r *RedisClient) AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	if expiration <= 0 {
		return nil
	}
	key := fmt.Sprintf("blacklist:token:%s", token)
	return r.client.Set(ctx, key, "revoked", expiration).Err()
}

// IsBlacklisted mengecek apakah token ada di daftar blacklist
func (r *RedisClient) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:token:%s", token)
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
