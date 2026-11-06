#!/bin/sh
# migrate.sh — wrapper menjalankan migrasi PostgreSQL
# usage:
#   ./scripts/migrate.sh up          terapkan semua migrasi
#   ./scripts/migrate.sh down        rollback satu migrasi terakhir
#   ./scripts/migrate.sh version     lihat versi migrasi saat ini
#   ./scripts/migrate.sh force 3     paksa versi kembali ke 3
set -eu

DIR="$(cd "$(dirname "$0")/.." && pwd)"

cd "$DIR"
go run ./cmd/migrate "$@"
