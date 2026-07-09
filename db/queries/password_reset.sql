-- name: CreatePasswordResetCode :one
INSERT INTO password_reset_codes (user_id, code_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetLatestPasswordResetCode :one
SELECT * FROM password_reset_codes
WHERE user_id = $1 AND consumed_at IS NULL
ORDER BY created_at DESC
LIMIT 1;

-- name: IncrementPasswordResetAttempts :exec
UPDATE password_reset_codes SET attempts = attempts + 1 WHERE id = $1;

-- name: ConsumePasswordResetCode :exec
UPDATE password_reset_codes SET consumed_at = NOW() WHERE id = $1;

-- name: InvalidatePasswordResetCodes :exec
UPDATE password_reset_codes SET consumed_at = NOW()
WHERE user_id = $1 AND consumed_at IS NULL;
