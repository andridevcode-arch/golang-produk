#!/bin/sh
# build.sh — compile semua package Go
# Usage: ./scripts/build.sh
set -eu
DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$DIR"
echo "Building..."
go build ./...
echo "Build successful."
