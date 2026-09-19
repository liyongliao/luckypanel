#!/usr/bin/env bash
set -e

echo "==> Building LuckyPanel (server-manager)..."
CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X 'github.com/0xJacky/Nginx-UI/internal/version.Version=1.0.0'" -o server-manager main.go
echo "==> Build completed successfully: ./server-manager"
echo "==> Run LuckyPanel using: ./server-manager run"
