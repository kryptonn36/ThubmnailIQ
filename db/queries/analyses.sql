-- name: CreateAnalysis :one
INSERT INTO analyses (
    workspace_id, project_id, user_id, keyword, keyword_normalized, thumbnail_s3_key, status
) VALUES (
    $1, $2, $3, $4, $5, $6, 'pending'
) RETURNING *;

-- name: GetAnalysisByID :one
SELECT * FROM analyses WHERE id = $1 AND deleted_at IS NULL;

-- name: ListAnalysesByWorkspace :many
SELECT * FROM analyses
WHERE workspace_id = $1 AND deleted_at IS NULL
  AND (sqlc.narg('status')::varchar IS NULL OR status = sqlc.narg('status'))
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountAnalysesByWorkspace :one
SELECT COUNT(*) FROM analyses
WHERE workspace_id = $1 AND deleted_at IS NULL
  AND (sqlc.narg('status')::varchar IS NULL OR status = sqlc.narg('status'));

-- name: UpdateAnalysisStatus :one
UPDATE analyses
SET status = $2, error_message = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateAnalysisResults :one
UPDATE analyses
SET
    status = 'complete',
    score = $2,
    visibility_score = $3,
    contrast_score = $4,
    attention_score = $5,
    mobile_score = $6,
    branding_score = $7,
    curiosity_score = $8,
    cv_results = $9,
    competitor_avg = $10,
    suggestions = $11,
    competitor_count = $12,
    rank_in_competitors = $13,
    relevance_score = $14,
    relevance_reasoning = $15,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateThumbnailVersion :one
INSERT INTO thumbnail_versions (analysis_id, version_number, s3_key, score, cv_results)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListThumbnailVersions :many
SELECT * FROM thumbnail_versions WHERE analysis_id = $1 ORDER BY version_number;

-- name: NextThumbnailVersionNumber :one
SELECT COALESCE(MAX(version_number), 0) + 1 FROM thumbnail_versions WHERE analysis_id = $1;
