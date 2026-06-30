-- +goose Up
ALTER TABLE analyses
  ADD COLUMN actual_ctr FLOAT,
  ADD COLUMN published_at TIMESTAMPTZ;

-- +goose Down
ALTER TABLE analyses
  DROP COLUMN published_at,
  DROP COLUMN actual_ctr;
