#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BROWSER_DIR="$ROOT/browser"
RELEASE_DIR="$ROOT/release"

REMOTE="gusmontana@dailies"
REMOTE_TMP="~/dailies.new"

BINARY_NAME="dailies-linux-amd64"
BINARY_PATH="$RELEASE_DIR/$BINARY_NAME"

echo "==> Building frontend"
cd "$BROWSER_DIR"

npm ci
npm run build

echo "==> Building Linux AMD64 binary"
cd "$ROOT"

rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

GOOS=linux \
GOARCH=amd64 \
CGO_ENABLED=0 \
go build \
  -trimpath \
  -ldflags="-s -w" \
  -o "$BINARY_PATH" \
  .

echo "==> Built:"
ls -lh "$BINARY_PATH"

echo "==> Uploading binary to $REMOTE"
scp "$BINARY_PATH" "$REMOTE:$REMOTE_TMP"

echo "==> Installing binary and restarting service"

ssh "$REMOTE" bash <<'EOF'
set -euo pipefail

REMOTE_TMP="$HOME/dailies.new"
REMOTE_APP="/opt/dailies/dailies"

echo "==> Installing new binary"

sudo chmod 755 "$REMOTE_TMP"
sudo chown dailies:dailies "$REMOTE_TMP"

sudo mv "$REMOTE_TMP" "$REMOTE_APP"

echo "==> Restarting dailies"

sudo systemctl restart dailies

sleep 1

echo "==> Checking service"

if ! sudo systemctl is-active --quiet dailies; then
    echo "ERROR: dailies failed to start"
    echo
    sudo systemctl status dailies --no-pager
    exit 1
fi

echo "==> Dailies is running"
echo
sudo systemctl status dailies --no-pager
EOF

echo
echo "Deployment successful"
