-- +goose Up
CREATE TABLE tracking_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    channel_id VARCHAR(50),
    keyword VARCHAR(500),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    check_interval_hours INT NOT NULL DEFAULT 24,
    last_checked_at TIMESTAMPTZ,
    next_check_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tracking_jobs_workspace ON tracking_jobs(workspace_id);
CREATE INDEX idx_tracking_jobs_next_check ON tracking_jobs(next_check_at) WHERE status = 'active';

CREATE TABLE competitor_snapshots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    analysis_id UUID REFERENCES analyses(id) ON DELETE SET NULL,
    tracking_job_id UUID REFERENCES tracking_jobs(id) ON DELETE SET NULL,
    video_id VARCHAR(20) NOT NULL,
    video_title TEXT NOT NULL,
    channel_id VARCHAR(50) NOT NULL,
    channel_name VARCHAR(255) NOT NULL,
    thumbnail_url TEXT NOT NULL,
    thumbnail_s3_key TEXT,
    view_count BIGINT,
    like_count BIGINT,
    subscriber_count BIGINT,
    rank_position INT,
    keyword VARCHAR(500),
    cv_results JSONB,
    score INT,
    snapshot_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_competitor_snapshots_analysis ON competitor_snapshots(analysis_id);
CREATE INDEX idx_competitor_snapshots_video ON competitor_snapshots(video_id);
CREATE INDEX idx_competitor_snapshots_channel ON competitor_snapshots(channel_id);
CREATE INDEX idx_competitor_snapshots_keyword ON competitor_snapshots(keyword, snapshot_date);
CREATE INDEX idx_competitor_snapshots_date ON competitor_snapshots(snapshot_date DESC);

-- +goose Down
DROP TABLE IF EXISTS competitor_snapshots;
DROP TABLE IF EXISTS tracking_jobs;
