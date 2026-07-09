-- +goose Up
-- password_reset_codes holds short-lived OTP codes emailed to users who
-- request a password reset. Same posture as email_verification_codes: only the
-- SHA-256 hash of the code is stored, with an expiry and an attempt counter.
CREATE TABLE password_reset_codes (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    code_hash   VARCHAR(64) NOT NULL,
    attempts    INT NOT NULL DEFAULT 0,
    expires_at  TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_password_reset_codes_user ON password_reset_codes(user_id) WHERE consumed_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS password_reset_codes;
