#!/usr/bin/env bash
set -euo pipefail

sudo apt-get update

# protobuf compiler
sudo apt-get install -y protobuf-compiler

# common build deps for Rust crates w/ native deps
sudo apt-get install -y pkg-config build-essential clang lldb lld

# buf CLI
BUF_VERSION="1.50.0"
curl -sSL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-Linux-x86_64" -o /tmp/buf
sudo mv /tmp/buf /usr/local/bin/buf
sudo chmod +x /usr/local/bin/buf

# JS deps for web frontend
if [ -f web/package.json ]; then
  cd web && npm ci && cd ..
fi

# Rust deps
if [ -f Cargo.toml ]; then
  cargo fetch
fi
