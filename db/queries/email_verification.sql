-- name: CreateEmailVerificationCode :one
INSERT INTO email_verification_codes (user_id, code_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetLatestEmailVerificationCode :one
SELECT * FROM email_verification_codes
WHERE user_id = $1 AND consumed_at IS NULL
ORDER BY created_at DESC
LIMIT 1;

-- name: IncrementEmailVerificationAttempts :exec
UPDATE email_verification_codes SET attempts = attempts + 1 WHERE id = $1;

-- name: ConsumeEmailVerificationCode :exec
UPDATE email_verification_codes SET consumed_at = NOW() WHERE id = $1;

-- name: InvalidateEmailVerificationCodes :exec
UPDATE email_verification_codes SET consumed_at = NOW()
WHERE user_id = $1 AND consumed_at IS NULL;
