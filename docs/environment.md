# Environment Variables

Shulker uses a `.env` file loaded via `dotenvy`. Create a `.env` in the project root before running.

## Required Variables

| Variable             | Used By       | Description                                   |
| -------------------- | ------------- | --------------------------------------------- |
| `JWT_ACCESS_SECRET`  | `api`, `auth` | Secret key for signing access tokens (HS256)  |
| `JWT_REFRESH_SECRET` | `api`, `auth` | Secret key for signing refresh tokens (HS256) |

## Example `.env`

```env
JWT_ACCESS_SECRET=some-long-random-access-secret
JWT_REFRESH_SECRET=some-long-random-refresh-secret
```

## Notes

- Use different values for each secret — they sign different token types and should not be interchangeable
- In production, generate these with something like `openssl rand -hex 32`
- The `secure` flag on cookies is currently set to `false` — flip it to `true` if serving over HTTPS
