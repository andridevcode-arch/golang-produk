#!/bin/sh
# test.sh — jalankan semua tests
# Usage: ./scripts/test.sh
set -eu
DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$DIR"
go test ./... "$@"
