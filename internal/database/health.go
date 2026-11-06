package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// StatusHealth status kesehatan koneksi
const (
	StatusUp   = "up"
	StatusDown = "down"
)

// HealthStatus menampung hasil pemeriksaan kesehatan satu database
type HealthStatus struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Latency  int64  `json:"latency_ms"`
	Error    string `json:"error,omitempty"`
	CheckedAt time.Time `json:"checked_at"`
}

// Health memeriksa pool PostgreSQL (Supabase)
func (s *Supabase) Health(ctx context.Context) HealthStatus {
	return checkHealth("postgres", func(ctx context.Context) error {
		return s.pool.Ping(ctx)
	}, ctx)
}

// Health memeriksa klien MongoDB
func (m *Mongo) Health(ctx context.Context) HealthStatus {
	return checkHealth("mongo", func(ctx context.Context) error {
		return m.client.Ping(ctx, readpref.Primary())
	}, ctx)
}

// Health memeriksa klien Redis
func (r *RedisClient) Health(ctx context.Context) HealthStatus {
	return checkHealth("redis", func(ctx context.Context) error {
		return r.client.Ping(ctx).Err()
	}, ctx)
}

// checkHealth mengukur latensi dan menandai status up/down
func checkHealth(name string, fn func(ctx context.Context) error, ctx context.Context) HealthStatus {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	start := time.Now()
	result := HealthStatus{
		Name:      name,
		Status:    StatusDown,
		CheckedAt: time.Now(),
	}

	if err := fn(ctx); err != nil {
		result.Error = err.Error()
	} else {
		result.Status = StatusUp
	}
	result.Latency = time.Since(start).Milliseconds()

	return result
}
