-- name: UpsertBlob :one
INSERT INTO blobs (hash, mime_type, size)
VALUES ($1, $2, $3)
ON CONFLICT (hash) DO UPDATE SET ref_count = blobs.ref_count + 1
RETURNING *;

-- name: GetBlob :one
SELECT * FROM blobs WHERE id = $1;

-- name: GetBlobByHash :one
SELECT * FROM blobs WHERE hash = $1;

-- name: DeleteBlob :exec
DELETE FROM blobs WHERE id = $1;

-- name: DecrementBlob :one
UPDATE blobs SET ref_count = ref_count - 1 WHERE hash = $1
RETURNING *;
