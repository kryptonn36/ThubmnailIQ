-- name: UpsertViralThumbnail :one
INSERT INTO viral_thumbnails (video_id, channel_id, channel_name, video_title, thumbnail_url, niche, tags, view_count, score, has_face, cv_results)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (video_id) DO UPDATE
SET score = EXCLUDED.score,
    view_count = EXCLUDED.view_count,
    cv_results = EXCLUDED.cv_results,
    tags = (SELECT array_agg(DISTINCT tag) FROM unnest(viral_thumbnails.tags || EXCLUDED.tags) AS tag),
    niche = COALESCE(viral_thumbnails.niche, EXCLUDED.niche),
    has_face = EXCLUDED.has_face
RETURNING *;

-- name: SearchViralThumbnails :many
SELECT * FROM viral_thumbnails
WHERE (sqlc.narg('keyword')::text IS NULL OR tags @> ARRAY[sqlc.narg('keyword')::text])
  AND (sqlc.narg('niche')::varchar IS NULL OR niche = sqlc.narg('niche'))
  AND (sqlc.narg('min_score')::int IS NULL OR score >= sqlc.narg('min_score'))
  AND (sqlc.narg('has_face')::bool IS NULL OR has_face = sqlc.narg('has_face'))
ORDER BY score DESC
LIMIT $1 OFFSET $2;
