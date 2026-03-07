-- +goose Up
CREATE TABLE blobs (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    hash        TEXT NOT NULL UNIQUE,
    filepath    TEXT NOT NULL,
    mime_type   TEXT NOT NULL,
    size        BIGINT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE blobs;
