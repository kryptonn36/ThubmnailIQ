-- +goose Up
CREATE TABLE viral_thumbnails (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    video_id VARCHAR(20) UNIQUE NOT NULL,
    channel_id VARCHAR(50) NOT NULL,
    channel_name VARCHAR(255) NOT NULL,
    video_title TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,
    thumbnail_s3_key TEXT,
    niche VARCHAR(100),
    tags TEXT[],
    view_count BIGINT,
    view_count_when_captured BIGINT,
    score INT,
    has_face BOOLEAN NOT NULL DEFAULT FALSE,
    cv_results JSONB,
    captured_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_viral_thumbnails_niche ON viral_thumbnails(niche);
CREATE INDEX idx_viral_thumbnails_score ON viral_thumbnails(score DESC);
CREATE INDEX idx_viral_thumbnails_tags ON viral_thumbnails USING gin(tags);
CREATE INDEX idx_viral_thumbnails_cv ON viral_thumbnails USING gin(cv_results);

-- +goose Down
DROP TABLE IF EXISTS viral_thumbnails;
