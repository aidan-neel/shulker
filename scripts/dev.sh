#!/usr/bin/env bash
set -euo pipefail

# Shulker dev script
# Starts all Rust services via cargo xtask, then the Vite frontend

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Check for tmux and use it if available for a nicer experience
if command -v tmux &>/dev/null; then
  SESSION="shulker"

  tmux new-session -d -s "$SESSION" -x 220 -y 50

  # Pane 0: Rust services
  tmux send-keys -t "$SESSION:0" "cd $ROOT && cargo xtask dev" Enter
  tmux rename-window -t "$SESSION:0" "services"

  # Pane 1: Vite frontend
  tmux new-window -t "$SESSION" -n "web"
  tmux send-keys -t "$SESSION:1" "cd $ROOT/web && npm run dev" Enter

  tmux select-window -t "$SESSION:0"
  tmux attach-session -t "$SESSION"
else
  # No tmux — run services in background, web in foreground
  echo "💡 tip: install tmux for a better dev experience"
  echo ""

  cd "$ROOT"

  echo "🚀 Starting Rust services..."
  cargo xtask dev &
  XTASK_PID=$!

  # Give services a moment before starting frontend
  sleep 5

  echo "🌐 Starting Vite frontend..."
  cd "$ROOT/web" && npm run dev &
  VITE_PID=$!

  trap "kill $XTASK_PID $VITE_PID 2>/dev/null; exit 0" SIGINT SIGTERM

  wait
fi
