-- name: CreateUserBlob :one
INSERT INTO user_blobs (user_id, blob_id, path)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserBlobs :many
SELECT b.* FROM blobs b
JOIN user_blobs ub ON ub.blob_id = b.id
WHERE ub.user_id = $1
ORDER BY ub.created_at DESC;

-- name: GetUserBlobByPath :one
SELECT b.* FROM blobs b
JOIN user_blobs ub ON ub.blob_id = b.id
WHERE ub.user_id = $1 AND ub.path = $2;

-- name: DeleteUserBlob :exec
DELETE FROM user_blobs WHERE user_id = $1 AND blob_id = $2;
