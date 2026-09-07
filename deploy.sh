#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BROWSER_DIR="$ROOT/browser"
DIST_DIR="$ROOT/pkg/server/dist"
RELEASE_DIR="$ROOT/release"

BINARY_NAME="dailies"
BINARY_PATH="$RELEASE_DIR/${BINARY_NAME}-linux-amd64"

echo "==> Cleaning previous release"
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

echo "==> Building frontend"
cd "$BROWSER_DIR"
npm ci
npm run build

echo "==> Building Linux AMD64 binary"
cd "$ROOT"

GOOS=linux \
GOARCH=amd64 \
CGO_ENABLED=0 \
go build \
  -trimpath \
  -ldflags="-s -w" \
  -o "$BINARY_PATH" \
  .

echo "==> Verifying binary"
file "$BINARY_PATH"

echo
echo "Release built:"
echo "  $BINARY_PATH"
echo
echo "Frontend embedded into binary:"
echo "  $DIST_DIR"
echo
echo "Size:"
du -h "$BINARY_PATH"

echo
echo "Done."
