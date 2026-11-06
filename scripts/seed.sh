#!/bin/sh
# seed.sh — wrapper menjalankan seeder admin
set -eu

DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$DIR"
go run ./cmd/seed "$@"
