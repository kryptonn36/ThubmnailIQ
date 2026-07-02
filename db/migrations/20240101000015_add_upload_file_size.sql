-- +goose Up
ALTER TABLE analyses ADD COLUMN file_size_bytes BIGINT;
ALTER TABLE thumbnail_versions ADD COLUMN file_size_bytes BIGINT;

-- +goose Down
ALTER TABLE thumbnail_versions DROP COLUMN IF EXISTS file_size_bytes;
ALTER TABLE analyses DROP COLUMN IF EXISTS file_size_bytes;
