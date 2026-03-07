-- +goose Up
ALTER TABLE blobs ADD COLUMN ref_count INT NOT NULL DEFAULT 1;

-- +goose Down
ALTER TABLE blobs DROP COLUMN ref_count;
