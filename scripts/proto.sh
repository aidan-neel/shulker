#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if ! command -v buf &>/dev/null; then
  echo "❌ buf not found. Install it or reopen in devcontainer."
  exit 1
fi

echo "🔄 Regenerating protobuf files..."
buf generate

echo "✅ Done. You may need to rebuild: cargo build"
