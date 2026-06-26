package competitor

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Snapshot struct {
	ID               uuid.UUID
	AnalysisID       *uuid.UUID
	TrackingJobID    *uuid.UUID
	VideoID          string
	VideoTitle       string
	ChannelID        string
	ChannelName      string
	ThumbnailURL     string
	ViewCount        int64
	LikeCount        int64
	SubscriberCount  int64
	RankPosition     int
	Keyword          string
	CVResults        []byte
	Score            *int
	SnapshotDate     time.Time
	CreatedAt        time.Time
}

type TrackingJob struct {
	ID                 uuid.UUID  `json:"id"`
	WorkspaceID        uuid.UUID  `json:"workspace_id"`
	Type               string     `json:"type"`
	ChannelID          string     `json:"channel_id"`
	Keyword            string     `json:"keyword"`
	Status             string     `json:"status"`
	CheckIntervalHours int        `json:"check_interval_hours"`
	LastCheckedAt      *time.Time `json:"last_checked_at"`
	NextCheckAt        time.Time `json:"next_check_at"`
	CreatedBy          uuid.UUID `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
}

type Repository interface {
	CreateSnapshot(ctx context.Context, s *Snapshot) (*Snapshot, error)
	ListByAnalysis(ctx context.Context, analysisID uuid.UUID) ([]*Snapshot, error)
	ListLatestByKeyword(ctx context.Context, keyword string, limit int) ([]*Snapshot, error)

	CreateTrackingJob(ctx context.Context, t *TrackingJob) (*TrackingJob, error)
	ListTrackingJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*TrackingJob, error)
	ListDueTrackingJobs(ctx context.Context) ([]*TrackingJob, error)
	MarkTrackingJobChecked(ctx context.Context, id uuid.UUID) error
}
