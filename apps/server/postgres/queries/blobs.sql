-- name: CreateBlob :one
INSERT INTO blobs (user_id, hash, filepath, mime_type, size)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetBlob :one
SELECT * FROM blobs WHERE id = $1;

-- name: GetBlobByHash :one
SELECT * FROM blobs WHERE hash = $1;

-- name: GetBlobsByUser :many
SELECT * FROM blobs WHERE user_id = $1 ORDER BY created_at DESC;

-- name: DeleteBlob :exec
DELETE FROM blobs WHERE id = $1;

-- name: UpsertBlob :one
INSERT INTO blobs (user_id, hash, filepath, mime_type, size, ref_count)
VALUES ($1, $2, $3, $4, $5, 1)
ON CONFLICT (hash) DO UPDATE SET ref_count = blobs.ref_count + 1
RETURNING *;

-- name: DecrementBlob :one
UPDATE blobs SET ref_count = ref_count - 1 WHERE hash = $1
RETURNING *;
