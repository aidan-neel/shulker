-- +goose Up
ALTER TABLE blobs DROP COLUMN user_id;
ALTER TABLE blobs DROP COLUMN filepath;

CREATE TABLE user_blobs (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    blob_id    UUID NOT NULL REFERENCES blobs(id) ON DELETE CASCADE,
    path       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, path)
);

-- +goose Down
DROP TABLE user_blobs;
ALTER TABLE blobs ADD COLUMN user_id UUID REFERENCES users(id);
ALTER TABLE blobs ADD COLUMN filepath TEXT;
