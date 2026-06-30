-- name: CreateWorkspace :one
INSERT INTO workspaces (name, slug, owner_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetWorkspaceByID :one
SELECT * FROM workspaces WHERE id = $1 AND deleted_at IS NULL;

-- name: ListWorkspacesForUser :many
SELECT w.* FROM workspaces w
JOIN workspace_members wm ON wm.workspace_id = w.id
WHERE wm.user_id = $1 AND w.deleted_at IS NULL
ORDER BY w.created_at DESC;

-- name: AddWorkspaceMember :one
INSERT INTO workspace_members (workspace_id, user_id, role)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListWorkspaceMembers :many
SELECT wm.*, u.email, u.full_name FROM workspace_members wm
JOIN users u ON u.id = wm.user_id
WHERE wm.workspace_id = $1
ORDER BY wm.joined_at;

-- name: IncrementWorkspaceAnalysesUsage :exec
UPDATE workspaces SET analyses_this_month = analyses_this_month + 1, updated_at = NOW()
WHERE id = $1;

-- name: UpdateWorkspaceBrand :one
UPDATE workspaces
SET brand_primary_color = $2, brand_secondary_color = $3, brand_font = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateWorkspacePlan :one
UPDATE workspaces SET plan = $2, analyses_limit = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;
