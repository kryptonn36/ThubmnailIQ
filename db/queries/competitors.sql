-- name: CreateCompetitorSnapshot :one
INSERT INTO competitor_snapshots (
    analysis_id, tracking_job_id, video_id, video_title, channel_id, channel_name,
    thumbnail_url, view_count, like_count, subscriber_count, rank_position, keyword, cv_results, score
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;

-- name: ListCompetitorsByAnalysis :many
SELECT * FROM competitor_snapshots WHERE analysis_id = $1 ORDER BY rank_position;

-- name: ListLatestCompetitorsByKeyword :many
SELECT * FROM competitor_snapshots
WHERE keyword = $1
ORDER BY snapshot_date DESC, rank_position
LIMIT $2;

-- name: CreateTrackingJob :one
INSERT INTO tracking_jobs (workspace_id, type, channel_id, keyword, check_interval_hours, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListTrackingJobsByWorkspace :many
SELECT * FROM tracking_jobs WHERE workspace_id = $1 ORDER BY created_at DESC;

-- name: ListDueTrackingJobs :many
SELECT * FROM tracking_jobs WHERE status = 'active' AND next_check_at <= NOW();

-- name: MarkTrackingJobChecked :exec
UPDATE tracking_jobs
SET last_checked_at = NOW(), next_check_at = NOW() + (check_interval_hours || ' hours')::interval, updated_at = NOW()
WHERE id = $1;
