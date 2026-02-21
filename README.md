# Shulker

A personal cloud storage project built to learn gRPC, Rust, and React.

## Architecture

Shulker is a monorepo with three Rust services and a React frontend:

| Service   | Port    | Description                                                |
| --------- | ------- | ---------------------------------------------------------- |
| `api`     | `3000`  | HTTP gateway (Axum) — the only thing the frontend talks to |
| `auth`    | `50052` | gRPC service handling tokens and authentication            |
| `storage` | `50051` | gRPC service handling file uploads                         |
| `web`     | `5173`  | React + Vite frontend                                      |

The API gateway communicates with `auth` and `storage` internally over gRPC. Protobuf schemas live in `proto/` and are shared across services via `buf`.

## Prerequisites

The easiest way to get started is with the provided devcontainer, which handles everything automatically. You'll need either:

- [VS Code](https://code.visualstudio.com/) + [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers), or
- [Docker](https://www.docker.com/) + the [devcontainer CLI](https://github.com/devcontainers/cli)

If running locally without a devcontainer, you'll need:

- `Rust (stable)`
- `Node.js 22+`
- `protoc` (protobuf compiler)
- [`buf`](https://buf.build/docs/installation)

## Getting Started

### With devcontainer (recommended)

Open the project in VS Code and click **Reopen in Container** when prompted. The postCreate script will install all dependencies automatically.

### Without devcontainer

```bash
# Install JS dependencies
cd web && npm install && cd ..

# Fetch Rust dependencies
cargo fetch
```

## Running

Shulker uses `xtask` as its task runner. You'll need two terminals:

**Terminal 1 — Rust services:**

```bash
cargo xtask dev
```

This builds and starts `storage`, `auth`, and `api` in the correct order, with gRPC services coming up before the gateway.

**Terminal 2 — Frontend:**

```bash
cd web && npm run dev
```

Then open [http://localhost:5173](http://localhost:5173).

## Project Structure

```
shulker/
├── api/          # Axum HTTP gateway
├── auth/         # gRPC auth/token service
├── storage/      # gRPC file upload service
├── common/       # Shared Rust library (DB, models, JWT, hashing)
├── proto/        # Protobuf definitions
├── web/          # React + Vite frontend
├── xtask/        # Cargo task runner
└── data/         # SQLite database
```

## Protobuf

Schemas are defined in `proto/` and compiled at build time via `build.rs` in each service. To regenerate after changing a `.proto` file:

```bash
buf generate
```

## Notes

- The SQLite database lives at `data/db.sqlite3` and is initialized automatically on first run.
- Uploaded files are stored locally in `storage/files/`.
- This is a learning project — not intended for production use.
