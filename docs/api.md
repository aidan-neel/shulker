# API Reference

The API service is an Axum HTTP gateway running on port `3000`. It's the only service the frontend talks to — it proxies requests to the internal gRPC services.

## Base URL

```
http://localhost:3000
```

## Authentication

Auth routes are public. All `/storage` routes require a valid JWT access token passed as a Bearer token in the `Authorization` header.

```
Authorization: Bearer <access_token>
```

---

## Auth Routes

### POST /auth/register

Creates a new user account and returns a token pair. Also sets `access_token` and `refresh_token` as HTTP-only cookies.

**Request body:**

```json
{
  "email": "user@example.com",
  "display_name": "Aidan",
  "password": "hunter2"
}
```

**Response:**

```json
{
  "message": "Successfully registered",
  "token": {
    "access_token": "<jwt>",
    "refresh_token": "<jwt>"
  }
}
```

**Errors:**

- `409 Conflict` — email already in use

---

### POST /auth/login

Authenticates an existing user and returns a token pair.

**Request body:**

```json
{
  "email": "user@example.com",
  "password": "hunter2"
}
```

**Response:**

```json
{
  "message": "Successfully signed in",
  "token": {
    "access_token": "<jwt>",
    "refresh_token": "<jwt>"
  }
}
```

**Errors:**

- `401 Unauthorized` — invalid credentials

---

### POST /auth/refresh

Issues a new token pair from a valid refresh token. Rotates the refresh token (old one is replaced in the DB).

**Request body:**

```json
{
  "refresh_token": "<jwt>"
}
```

**Response:**

```json
{
  "message": "Successfully refreshed",
  "token": {
    "access_token": "<jwt>",
    "refresh_token": "<jwt>"
  }
}
```

**Errors:**

- `401 Unauthorized` — token invalid, not found, or expired

---

## Storage Routes

> All storage routes require JWT authentication.

### POST /storage/upload

Uploads a file. Streams the data to the storage gRPC service and records metadata in SQLite.

_Details depend on `api/src/routes/storage.rs` implementation._

---

## Token Lifetimes

| Token         | Lifetime   |
| ------------- | ---------- |
| Access token  | 15 minutes |
| Refresh token | 14 days    |

Refresh tokens are stored hashed in the database with one active token per user (upsert on conflict).
