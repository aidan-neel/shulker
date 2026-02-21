#!/usr/bin/env bash
set -euo pipefail

# Shulker manual setup script (for non-devcontainer environments)
# Tested on Ubuntu/Debian. For other distros, adjust apt commands accordingly.

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "🔧 Setting up Shulker..."

# --- System deps ---
if command -v apt-get &>/dev/null; then
  sudo apt-get update
  sudo apt-get install -y \
    protobuf-compiler \
    pkg-config \
    build-essential \
    clang \
    lldb \
    lld \
    curl
else
  echo "⚠️  Non-Debian system detected. Make sure you have installed:"
  echo "   protoc, pkg-config, clang, lld, curl"
fi

# --- Rust ---
if ! command -v cargo &>/dev/null; then
  echo "📦 Installing Rust..."
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
  source "$HOME/.cargo/env"
else
  echo "✅ Rust already installed ($(rustc --version))"
fi

# --- Node ---
if ! command -v node &>/dev/null; then
  echo "⚠️  Node.js not found. Install Node 22+ from https://nodejs.org"
  exit 1
else
  echo "✅ Node already installed ($(node --version))"
fi

# --- buf ---
if ! command -v buf &>/dev/null; then
  echo "📦 Installing buf..."
  BUF_VERSION="1.50.0"
  curl -sSL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-Linux-x86_64" -o /tmp/buf
  sudo mv /tmp/buf /usr/local/bin/buf
  sudo chmod +x /usr/local/bin/buf
  echo "✅ buf installed"
else
  echo "✅ buf already installed ($(buf --version))"
fi

# --- JS deps ---
echo "📦 Installing web dependencies..."
cd "$ROOT/web" && npm install && cd "$ROOT"

# --- Rust deps ---
echo "📦 Fetching Rust dependencies..."
cd "$ROOT" && cargo fetch

echo ""
echo "✅ Setup complete. Run the project with:"
echo "   cargo xtask dev        # Rust services"
echo "   cd web && npm run dev  # Frontend"
echo ""
echo "   Or use: bash scripts/dev.sh"
