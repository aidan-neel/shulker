# Architecture

## Overview

Shulker is a microservice-style monorepo. The frontend only ever talks to the API gateway, which in turn calls internal gRPC services for auth and storage.

```
[Browser / Web]
      |
      | HTTP (port 3000)
      v
  [ api ] ── gRPC (50052) ──> [ auth ]
      |
      └── gRPC (50051) ──> [ storage ]
              |
              └── SQLite + local filesystem
```

All three Rust services share the `common` crate for database access, JWT handling, password hashing, and models.

---

## Services

### `api` — HTTP Gateway (port 3000)

Built with Axum. Handles all HTTP traffic from the frontend and translates it into gRPC calls. Responsible for:

- Routing and CORS (allows `localhost:5173`)
- JWT middleware on protected routes
- Forwarding auth requests to `auth` via `TokenServiceClient`
- Forwarding upload requests to `storage` via `UploadServiceClient`
- Setting HTTP-only cookies on registration

State includes a DB pool, both gRPC clients, and a shared `JWTTokenService`.

### `auth` — Token Service (gRPC, port 50052)

Handles credential verification and token issuance. Implements the `TokenService` proto:

- **`GetToken`** — verifies email/password with Argon2, issues access + refresh JWTs, persists refresh token to DB
- **`RefreshToken`** — validates refresh token against DB, checks expiry, rotates and reissues both tokens

### `storage` — Upload Service (gRPC, port 50051)

Handles file uploads via a client-streaming gRPC call (`UploadFile`). Saves files to `storage/files/<user_id>/` and records metadata in SQLite.

### `common` — Shared Library

Not a service — a Rust library crate used by all three services. Contains:

- `db/connection.rs` — SQLite connection pool via `r2d2`
- `db/schema.rs` — DB initialization (creates tables on startup)
- `db/queries.rs` — async DB query functions (run via `spawn_blocking`)
- `jwt.rs` — `JWTTokenService` with HS256 access/refresh keys
- `hash.rs` — Argon2 password hashing
- `models/` — `User`, `File`, `Refresh` structs
- `utils.rs` — gRPC status to HTTP status code mapping

### `web` — React Frontend (port 5173)

Vite + React + TypeScript frontend. Talks only to the API gateway at `localhost:3000`.

---

## Database

SQLite at `data/db.sqlite3`, initialized automatically on first run by `common::db::schema::init_db`.

| Table     | Purpose                                            |
| --------- | -------------------------------------------------- |
| `user`    | User accounts (email, display_name, password_hash) |
| `refresh` | Active refresh tokens, one per user (upsert)       |
| `file`    | Uploaded file metadata (path, name, size, user_id) |

---

## Protobuf

Schemas are in `proto/` and compiled at build time by each service's `build.rs` using `tonic-build`.

| File           | Package  | Service                                      |
| -------------- | -------- | -------------------------------------------- |
| `token.proto`  | `auth`   | `TokenService` — login/refresh               |
| `upload.proto` | `upload` | `UploadService` — file streaming             |
| `user.proto`   | `auth`   | `UserService` — (defined, unused in gateway) |

---

## Auth Flow

```
Register:
  POST /auth/register
    → check user doesn't exist (SQLite)
    → hash password (Argon2)
    → insert user (SQLite)
    → call auth::GetToken (gRPC)
      → verify password, issue JWT pair
      → store refresh token (SQLite)
    → set HTTP-only cookies
    → return token pair

Login:
  POST /auth/login
    → call auth::GetToken (gRPC)
      → fetch user by email
      → verify password (Argon2)
      → issue JWT pair
      → store refresh token (SQLite)
    → return token pair

Refresh:
  POST /auth/refresh
    → call auth::RefreshToken (gRPC)
      → look up refresh token in DB
      → check expiry
      → rotate: issue new pair, replace DB record
    → return new token pair
```
