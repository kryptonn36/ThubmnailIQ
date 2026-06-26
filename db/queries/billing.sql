-- name: CreateApiKey :one
INSERT INTO api_keys (workspace_id, name, key_hash, key_prefix, created_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListApiKeysByWorkspace :many
SELECT * FROM api_keys WHERE workspace_id = $1 AND revoked_at IS NULL ORDER BY created_at DESC;

-- name: GetApiKeyByHash :one
SELECT * FROM api_keys WHERE key_hash = $1 AND revoked_at IS NULL;

-- name: RevokeApiKey :exec
UPDATE api_keys SET revoked_at = NOW() WHERE id = $1 AND workspace_id = $2;

-- name: UpsertSubscription :one
INSERT INTO subscriptions (workspace_id, stripe_subscription_id, stripe_price_id, plan, status, current_period_start, current_period_end)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (stripe_subscription_id) DO UPDATE
SET plan = EXCLUDED.plan, status = EXCLUDED.status,
    current_period_start = EXCLUDED.current_period_start,
    current_period_end = EXCLUDED.current_period_end,
    updated_at = NOW()
RETURNING *;

-- name: GetSubscriptionByWorkspace :one
SELECT * FROM subscriptions WHERE workspace_id = $1 ORDER BY created_at DESC LIMIT 1;
