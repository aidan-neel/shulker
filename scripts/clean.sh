#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "🧹 Cleaning Shulker..."

read -r -p "Also delete data/db.sqlite3? This wipes your local database. [y/N] " confirm
if [[ "$confirm" =~ ^[Yy]$ ]]; then
  rm -f "$ROOT/data/db.sqlite3"
  echo "  🗑  data/db.sqlite3"
fi

rm -rf "$ROOT/target"
echo "  🗑  target/"

rm -rf "$ROOT/web/node_modules"
echo "  🗑  web/node_modules/"

rm -rf "$ROOT/storage/files"/*
echo "  🗑  storage/files/*"

echo "✅ Clean."
