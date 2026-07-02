-- +goose Up
-- Singleton table (id is always 1) holding app-wide config the admin panel
-- can tune at runtime, instead of requiring a redeploy for these values.
CREATE TABLE app_settings (
    id SMALLINT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    max_upload_size_bytes BIGINT NOT NULL DEFAULT 2097152,
    allowed_extensions TEXT[] NOT NULL DEFAULT ARRAY['jpg', 'jpeg', 'png'],
    feature_flags JSONB NOT NULL DEFAULT '{}'::jsonb,
    storage_provider VARCHAR(50) NOT NULL DEFAULT 's3',
    email_provider VARCHAR(50) NOT NULL DEFAULT '',
    email_from_address VARCHAR(255) NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO app_settings (id) VALUES (1);

-- +goose Down
DROP TABLE IF EXISTS app_settings;
