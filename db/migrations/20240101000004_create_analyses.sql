-- +goose Up
CREATE TABLE analyses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    keyword VARCHAR(500) NOT NULL,
    keyword_normalized VARCHAR(500) NOT NULL,
    thumbnail_url TEXT NOT NULL,
    thumbnail_s3_key TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    score INT,
    visibility_score INT,
    contrast_score INT,
    attention_score INT,
    mobile_score INT,
    branding_score INT,
    curiosity_score INT,
    cv_results JSONB,
    competitor_avg JSONB,
    suggestions JSONB,
    competitor_count INT DEFAULT 0,
    rank_in_competitors INT,
    error_message TEXT,
    retry_count INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_analyses_workspace ON analyses(workspace_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_analyses_project ON analyses(project_id) WHERE project_id IS NOT NULL;
CREATE INDEX idx_analyses_keyword ON analyses USING gin(to_tsvector('english', keyword));
CREATE INDEX idx_analyses_status ON analyses(status) WHERE status != 'complete';
CREATE INDEX idx_analyses_created ON analyses(created_at DESC);
CREATE INDEX idx_analyses_workspace_created ON analyses(workspace_id, created_at DESC) WHERE deleted_at IS NULL;

CREATE TABLE thumbnail_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    version_number INT NOT NULL DEFAULT 1,
    s3_key TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,
    score INT,
    cv_results JSONB,
    is_selected_winner BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(analysis_id, version_number)
);

-- +goose Down
DROP TABLE IF EXISTS thumbnail_versions;
DROP TABLE IF EXISTS analyses;
