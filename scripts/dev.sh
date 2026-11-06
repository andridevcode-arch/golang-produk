#!/bin/sh
# dev.sh — jalankan API server development
# Usage: ./scripts/dev.sh
set -eu
DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$DIR"
go run ./cmd/api "$@"
