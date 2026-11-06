package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"go-produk/internal/config"
)

// migrationsDir adalah folder yang berisi file migrasi
const migrationsDir = "migrations/postgres"

// fileNameRe mencocokkan format: 000001_nama.up.sql atau nama_folder/000001_nama.up.sql
var fileNameRe = regexp.MustCompile(`^(\d+)_.+\.(up|down)\.sql$`)

type migrationFile struct {
	version int64
	name    string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "migrate:", err)
		os.Exit(1)
	}
}

func run() error {
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("muat konfigurasi: %w", err)
	}

	// Cek apakah DATABASE_URL tersedia
	if strings.TrimSpace(cfg.Postgres.ConnectionString) == "" {
		return fmt.Errorf("DATABASE_URL tidak ditemukan di .env")
	}

	conn, err := pgx.Connect(ctx, cfg.Postgres.ConnectionString)
	if err != nil {
		return fmt.Errorf("hubungkan ke postgres: %w", err)
	}
	defer func() { _ = conn.Close(ctx) }()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	if err := ensureHistoryTable(ctx, conn); err != nil {
		return err
	}

	switch command {
	case "up":
		return runUp(ctx, conn)
	case "down":
		return runDown(ctx, conn)
	case "version":
		printVersion(ctx, conn)
		return nil
	case "force":
		return runForce(ctx, conn, os.Args)
	default:
		return fmt.Errorf(
			"perintah tidak dikenal %q (gunakan: up | down | version | force <versi>)",
			command,
		)
	}
}

func ensureHistoryTable(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    BIGINT PRIMARY KEY,
			dirty      BOOLEAN     NOT NULL DEFAULT false,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	if err != nil {
		return fmt.Errorf("siapkan schema_migrations: %w", err)
	}
	return nil
}

func currentVersion(ctx context.Context, conn *pgx.Conn) int64 {
	var version int64
	err := conn.QueryRow(
		ctx,
		`SELECT COALESCE(MAX(version), 0) FROM schema_migrations WHERE dirty = false`,
	).Scan(&version)
	if err != nil {
		return 0
	}
	return version
}

// findMigrations mencari semua file migrasi di folder migrations (termasuk subfolder)
func findMigrations(dir string) (ups, downs []migrationFile, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("baca folder %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Recursively search in subdirectories
			subUps, subDowns, err := findMigrations(filepath.Join(dir, entry.Name()))
			if err != nil {
				return nil, nil, err
			}
			ups = append(ups, subUps...)
			downs = append(downs, subDowns...)
			continue
		}

		// Skip .bak files
		if strings.HasSuffix(entry.Name(), ".bak") {
			continue
		}

		// Only match files directly in this directory
		m := fileNameRe.FindStringSubmatch(entry.Name())
		if m == nil {
			continue
		}
		version, _ := strconv.ParseInt(m[1], 10, 64)
		file := migrationFile{version: version, name: entry.Name()}
		if m[2] == "up" {
			ups = append(ups, file)
		} else {
			downs = append(downs, file)
		}
	}

	sort.Slice(ups, func(i, j int) bool { return ups[i].version < ups[j].version })
	sort.Slice(downs, func(i, j int) bool { return downs[i].version > downs[j].version })
	return ups, downs, nil
}

func runUp(ctx context.Context, conn *pgx.Conn) error {
	ups, _, err := findMigrations(migrationsDir)
	if err != nil {
		return err
	}
	current := currentVersion(ctx, conn)

	for _, m := range ups {
		if m.version <= current {
			continue
		}
		path := filepath.Join(migrationsDir, m.name)
		fmt.Printf("menerapkan %s ...\n", m.name)
		if err := applyFile(ctx, conn, m.version, path, true); err != nil {
			return err
		}
	}

	fmt.Printf("migrasi up selesai, versi saat ini: %d\n", currentVersion(ctx, conn))
	return nil
}

func runDown(ctx context.Context, conn *pgx.Conn) error {
	_, downs, err := findMigrations(migrationsDir)
	if err != nil {
		return err
	}
	current := currentVersion(ctx, conn)

	var target *migrationFile
	for i := range downs {
		if downs[i].version == current {
			target = &downs[i]
			break
		}
	}
	if target == nil {
		fmt.Printf("tidak ada migrasi untuk di-rollback, versi saat ini: %d\n", current)
		return nil
	}

	fmt.Printf("me-rollback %s ...\n", target.name)
	if err := applyFile(ctx, conn, target.version, filepath.Join(migrationsDir, target.name), false); err != nil {
		return err
	}

	fmt.Printf("migrasi down selesai, versi saat ini: %d\n", currentVersion(ctx, conn))
	return nil
}

func runForce(ctx context.Context, conn *pgx.Conn, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: migrate force <versi>")
	}
	version, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil || version < 0 {
		return fmt.Errorf("versi tidak valid: %q", args[2])
	}

	_, err = conn.Exec(ctx, `DELETE FROM schema_migrations WHERE version > $1`, version)
	if err != nil {
		return fmt.Errorf("force gagal: %w", err)
	}

	fmt.Printf("versi dipaksa kembali ke %d\n", version)
	return nil
}

func printVersion(ctx context.Context, conn *pgx.Conn) {
	fmt.Println(currentVersion(ctx, conn))
}

// applyFile menjalankan satu file SQL di dalam transaksi
// beserta pencatatannya di schema_migrations
func applyFile(ctx context.Context, conn *pgx.Conn, version int64, path string, up bool) error {
	payload, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("baca %s: %w", path, err)
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("mulai transaksi: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, string(payload)); err != nil {
		return fmt.Errorf("jalankan %s: %w", path, err)
	}

	if up {
		_, err = tx.Exec(ctx, `
			INSERT INTO schema_migrations (version, dirty, applied_at)
			VALUES ($1, false, now())
			ON CONFLICT (version) DO UPDATE SET dirty = false, applied_at = now()
		`, version)
	} else {
		_, err = tx.Exec(ctx, `DELETE FROM schema_migrations WHERE version = $1`, version)
	}
	if err != nil {
		return fmt.Errorf("catat schema_migrations: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit %s: %w", path, err)
	}
	return nil
}
