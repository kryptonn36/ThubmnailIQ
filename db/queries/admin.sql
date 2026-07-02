-- name: CreateAdmin :one
INSERT INTO admin_users (email, password_hash, full_name, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAdminByID :one
SELECT * FROM admin_users WHERE id = $1;

-- name: GetAdminByEmail :one
SELECT * FROM admin_users WHERE email = $1;

-- name: UpdateAdminLastLogin :exec
UPDATE admin_users SET last_login_at = NOW(), updated_at = NOW() WHERE id = $1;

-- name: CreateAdminRefreshToken :one
INSERT INTO admin_refresh_tokens (admin_id, token_hash, device_info, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAdminRefreshToken :one
SELECT * FROM admin_refresh_tokens WHERE token_hash = $1 AND is_revoked = FALSE AND expires_at > NOW();

-- name: RevokeAdminRefreshToken :exec
UPDATE admin_refresh_tokens SET is_revoked = TRUE WHERE token_hash = $1;

-- name: CreateAuditLog :one
INSERT INTO admin_audit_logs (admin_id, action, target_type, target_id, metadata)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListAuditLogs :many
SELECT * FROM admin_audit_logs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAuditLogs :one
SELECT COUNT(*) FROM admin_audit_logs;

-- Dashboard -------------------------------------------------------------

-- name: CountUsers :one
SELECT COUNT(*) FROM users WHERE deleted_at IS NULL;

-- name: CountActiveUsers :one
SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND status = 'active';

-- name: CountUploads :one
SELECT COUNT(*) FROM analyses WHERE deleted_at IS NULL;

-- name: SumStorageUsed :one
SELECT COALESCE(SUM(file_size_bytes), 0)::bigint FROM analyses WHERE deleted_at IS NULL;

-- name: CountUploadsToday :one
SELECT COUNT(*) FROM analyses WHERE deleted_at IS NULL AND created_at >= CURRENT_DATE;

-- name: CountUploadsThisMonth :one
SELECT COUNT(*) FROM analyses WHERE deleted_at IS NULL AND created_at >= date_trunc('month', NOW());

-- name: ListRecentUsers :many
SELECT * FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1;

-- name: ListRecentUploads :many
SELECT * FROM analyses WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1;

-- User management ---------------------------------------------------------
--
-- ListUsersAdmin/GetUserDetailAdmin fold each user's primary-workspace plan
-- and analyses/storage totals into the same query via LATERAL joins, instead
-- of issuing 2-3 extra round-trips per user (avoids N+1 on the paginated
-- list).

-- name: ListUsersAdmin :many
SELECT
    u.*,
    COALESCE(ws.plan, '') AS plan,
    COALESCE(agg.analyses_count, 0)::bigint AS analyses_count,
    COALESCE(agg.storage_used, 0)::bigint AS storage_used
FROM users u
LEFT JOIN LATERAL (
    SELECT w.plan FROM workspace_members wm
    JOIN workspaces w ON w.id = wm.workspace_id AND w.deleted_at IS NULL
    WHERE wm.user_id = u.id
    ORDER BY wm.joined_at ASC
    LIMIT 1
) ws ON true
LEFT JOIN LATERAL (
    SELECT COUNT(*) AS analyses_count, COALESCE(SUM(file_size_bytes), 0) AS storage_used
    FROM analyses
    WHERE analyses.user_id = u.id AND analyses.deleted_at IS NULL
) agg ON true
WHERE u.deleted_at IS NULL
  AND (sqlc.arg(search)::text = '' OR u.email ILIKE '%' || sqlc.arg(search) || '%' OR u.full_name ILIKE '%' || sqlc.arg(search) || '%')
  AND (sqlc.arg(status_filter)::text = '' OR u.status = sqlc.arg(status_filter))
ORDER BY u.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsersAdmin :one
SELECT COUNT(*) FROM users
WHERE deleted_at IS NULL
  AND (sqlc.arg(search)::text = '' OR email ILIKE '%' || sqlc.arg(search) || '%' OR full_name ILIKE '%' || sqlc.arg(search) || '%')
  AND (sqlc.arg(status_filter)::text = '' OR status = sqlc.arg(status_filter));

-- name: GetUserDetailAdmin :one
SELECT
    u.*,
    COALESCE(ws.plan, '') AS plan,
    ws.id AS workspace_id,
    COALESCE(agg.analyses_count, 0)::bigint AS analyses_count,
    COALESCE(agg.storage_used, 0)::bigint AS storage_used
FROM users u
LEFT JOIN LATERAL (
    SELECT w.id, w.plan FROM workspace_members wm
    JOIN workspaces w ON w.id = wm.workspace_id AND w.deleted_at IS NULL
    WHERE wm.user_id = u.id
    ORDER BY wm.joined_at ASC
    LIMIT 1
) ws ON true
LEFT JOIN LATERAL (
    SELECT COUNT(*) AS analyses_count, COALESCE(SUM(file_size_bytes), 0) AS storage_used
    FROM analyses
    WHERE analyses.user_id = u.id AND analyses.deleted_at IS NULL
) agg ON true
WHERE u.id = $1;

-- name: CountUserUploads :one
SELECT COUNT(*) FROM analyses WHERE user_id = $1 AND deleted_at IS NULL;

-- name: SuspendUserAdmin :exec
UPDATE users SET status = 'suspended', updated_at = NOW() WHERE id = $1;

-- name: ActivateUserAdmin :exec
UPDATE users SET status = 'active', updated_at = NOW() WHERE id = $1;

-- name: SoftDeleteUserAdmin :exec
UPDATE users SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1;

-- name: ResetUserPasswordAdmin :exec
UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1;

-- name: ChangeUserWorkspaceRoleAdmin :exec
UPDATE workspace_members SET role = $3 WHERE user_id = $1 AND workspace_id = $2;

-- name: ListUserUploadsAdmin :many
SELECT * FROM analyses WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- Upload management ---------------------------------------------------------
--
-- include_deleted lets the admin list soft-deleted uploads too (needed to
-- find something to restore); the customer-facing analyses list never does.

-- name: ListUploadsAdmin :many
SELECT * FROM analyses
WHERE (sqlc.arg(include_deleted)::bool OR deleted_at IS NULL)
  AND (sqlc.arg(search)::text = '' OR keyword ILIKE '%' || sqlc.arg(search) || '%')
  AND (sqlc.arg(status_filter)::text = '' OR status = sqlc.arg(status_filter))
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUploadsAdmin :one
SELECT COUNT(*) FROM analyses
WHERE (sqlc.arg(include_deleted)::bool OR deleted_at IS NULL)
  AND (sqlc.arg(search)::text = '' OR keyword ILIKE '%' || sqlc.arg(search) || '%')
  AND (sqlc.arg(status_filter)::text = '' OR status = sqlc.arg(status_filter));

-- name: GetUploadAdmin :one
SELECT * FROM analyses WHERE id = $1;

-- name: SoftDeleteUploadAdmin :exec
UPDATE analyses SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1;

-- name: RestoreUploadAdmin :exec
UPDATE analyses SET deleted_at = NULL, updated_at = NOW() WHERE id = $1;

-- Analytics -----------------------------------------------------------------

-- name: DailySignupTrend :many
SELECT date_trunc('day', created_at)::date AS bucket, COUNT(*) AS count
FROM users
WHERE deleted_at IS NULL AND created_at >= NOW() - interval '30 days'
GROUP BY bucket
ORDER BY bucket;

-- name: MonthlySignupTrend :many
SELECT date_trunc('month', created_at)::date AS bucket, COUNT(*) AS count
FROM users
WHERE deleted_at IS NULL AND created_at >= NOW() - interval '12 months'
GROUP BY bucket
ORDER BY bucket;

-- name: DailyUploadTrend :many
SELECT date_trunc('day', created_at)::date AS bucket, COUNT(*) AS count
FROM analyses
WHERE deleted_at IS NULL AND created_at >= NOW() - interval '30 days'
GROUP BY bucket
ORDER BY bucket;

-- name: TopActiveUsers :many
SELECT
    u.*,
    COALESCE(ws.plan, '') AS plan,
    COALESCE(agg.analyses_count, 0)::bigint AS analyses_count,
    COALESCE(agg.storage_used, 0)::bigint AS storage_used
FROM users u
LEFT JOIN LATERAL (
    SELECT w.plan FROM workspace_members wm
    JOIN workspaces w ON w.id = wm.workspace_id AND w.deleted_at IS NULL
    WHERE wm.user_id = u.id
    ORDER BY wm.joined_at ASC
    LIMIT 1
) ws ON true
JOIN LATERAL (
    SELECT COUNT(*) AS analyses_count, COALESCE(SUM(file_size_bytes), 0) AS storage_used
    FROM analyses
    WHERE analyses.user_id = u.id AND analyses.deleted_at IS NULL
) agg ON true
WHERE u.deleted_at IS NULL AND agg.analyses_count > 0
ORDER BY agg.analyses_count DESC
LIMIT $1;

-- name: FileTypeBreakdown :many
-- The upload pipeline currently always writes .jpg regardless of the
-- original content type (see usecase/analysis.Create), so this will read as
-- entirely "jpg" today — grouping by the real stored extension keeps this
-- query correct if that ever changes, rather than hardcoding "jpg".
SELECT lower(substring(thumbnail_s3_key from '\.([a-zA-Z0-9]+)$')) AS extension, COUNT(*) AS count
FROM analyses
WHERE deleted_at IS NULL
GROUP BY extension
ORDER BY count DESC;

-- name: APIUsageStats :one
SELECT
    COALESCE(SUM(requests_this_month), 0)::bigint AS total_requests_this_month,
    COUNT(*) FILTER (WHERE revoked_at IS NULL)::bigint AS active_keys
FROM api_keys;

-- Settings --------------------------------------------------------------

-- name: GetAppSettings :one
SELECT * FROM app_settings WHERE id = 1;

-- name: UpdateAppSettings :one
UPDATE app_settings
SET
    max_upload_size_bytes = $1,
    allowed_extensions = $2,
    feature_flags = $3,
    storage_provider = $4,
    email_provider = $5,
    email_from_address = $6,
    updated_at = NOW()
WHERE id = 1
RETURNING *;

-- name: UpdateAdminPassword :exec
UPDATE admin_users SET password_hash = $2, updated_at = NOW() WHERE id = $1;
