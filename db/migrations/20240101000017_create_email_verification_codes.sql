-- +goose Up
-- email_verification_codes holds short-lived OTP codes emailed to users to
-- confirm ownership of their address. Only a SHA-256 hash of the code is
-- stored (never the plaintext), mirroring how refresh_tokens/api_keys are kept.
CREATE TABLE email_verification_codes (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash   VARCHAR(64) NOT NULL,
    attempts    INT NOT NULL DEFAULT 0,
    expires_at  TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Fast lookup of a user's currently-active (unconsumed) code.
CREATE INDEX idx_email_verification_codes_user ON email_verification_codes(user_id) WHERE consumed_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS email_verification_codes;
