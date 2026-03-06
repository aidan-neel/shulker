-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_uuidv7;
ALTER TABLE users ALTER COLUMN id SET DEFAULT uuid_generate_v7();

-- +goose Down
ALTER TABLE users ALTER COLUMN id SET DEFAULT gen_random_uuid();
DROP EXTENSION IF EXISTS pg_uuidv7;
