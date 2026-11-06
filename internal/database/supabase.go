package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"go-produk/internal/config"
)

// Supabase adalah struct untuk koneksi ke Supabase PostgreSQL
type Supabase struct {
	pool *pgxpool.Pool
	cfg  config.PostgresConfig
}

// NewSupabase membuat koneksi baru ke Supabase
func NewSupabase(ctx context.Context, cfg config.PostgresConfig) (*Supabase, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}
	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create supabase pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("supabase ping: %w", err)
	}

	return &Supabase{
		pool: pool,
		cfg:  cfg,
	}, nil
}

// Close menutup koneksi dengan graceful shutdown
func (s *Supabase) Close() {
	s.pool.Close()
}

// Pool mengembalikan connection pool untuk digunakan
func (s *Supabase) Pool() *pgxpool.Pool {
	return s.pool
}

// ConnectionString mengembalikan connection string untuk postgres
func (s *Supabase) ConnectionString() string {
	return s.cfg.ConnectionString
}

// Ping test koneksi
func (s *Supabase) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}
