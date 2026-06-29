-- +goose Up
-- Defensive backfill: every existing row already has thumbnail_s3_key/s3_key
-- populated alongside thumbnail_url (they're written together in the same
-- INSERT), but extract the key from the URL just in case any row was ever
-- written by older code that only set the URL.
UPDATE analyses
SET thumbnail_s3_key = regexp_replace(thumbnail_url, '^https?://[^/]+/', '')
WHERE thumbnail_s3_key IS NULL OR thumbnail_s3_key = '';

UPDATE thumbnail_versions
SET s3_key = regexp_replace(thumbnail_url, '^https?://[^/]+/', '')
WHERE s3_key IS NULL OR s3_key = '';

ALTER TABLE analyses DROP COLUMN thumbnail_url;
ALTER TABLE thumbnail_versions DROP COLUMN thumbnail_url;

-- +goose Down
ALTER TABLE analyses ADD COLUMN thumbnail_url TEXT;
ALTER TABLE thumbnail_versions ADD COLUMN thumbnail_url TEXT;
