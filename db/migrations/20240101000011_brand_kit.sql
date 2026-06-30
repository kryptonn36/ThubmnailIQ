-- +goose Up
ALTER TABLE workspaces
  ADD COLUMN brand_primary_color VARCHAR(7) DEFAULT '#6366F1',
  ADD COLUMN brand_secondary_color VARCHAR(7) DEFAULT '#8B5CF6',
  ADD COLUMN brand_font VARCHAR(100) DEFAULT 'Inter';

-- +goose Down
ALTER TABLE workspaces
  DROP COLUMN brand_font,
  DROP COLUMN brand_secondary_color,
  DROP COLUMN brand_primary_color;
