package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// ConnPool adalah interface minimal yang dibutuhkan WithinTx.
// Supabase.Pool() dan pgx.Tx sama-sama memenuhi interface ini.
type ConnPool interface {
	BeginTx(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error)
}

// WithinTx menjalankan fn di dalam satu transaksi PostgreSQL.
// Commit jika fn sukses, rollback jika fn error atau panic.
// Dipakai untuk operasi yang harus konsisten seperti pembaruan
// kuota, pesanan, pembayaran, dan entitlement.
func WithinTx(ctx context.Context, pool ConnPool, fn func(ctx context.Context, tx pgx.Tx) error) (err error) {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			err = fmt.Errorf("panik di dalam transaksi: %v", r)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

// WithinTxWithRetry menjalankan WithinTx dengan retry bila
// terjadi serialization failure atau deadlock.
func WithinTxWithRetry(ctx context.Context, pool ConnPool, maxRetries int, fn func(ctx context.Context, tx pgx.Tx) error) error {
	if maxRetries <= 0 {
		maxRetries = 3
	}

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if err := WithinTx(ctx, pool, fn); err != nil {
			if IsSerializationFailure(err) {
				lastErr = err
				continue
			}
			return err
		}
		return nil
	}
	return fmt.Errorf("gagal setelah %d kali retry: %w", maxRetries, lastErr)
}
