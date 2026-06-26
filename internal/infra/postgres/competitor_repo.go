package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/competitor"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/postgres/db"
)

type CompetitorRepo struct {
	q *db.Queries
}

func NewCompetitorRepo(pool *pgxpool.Pool) *CompetitorRepo {
	return &CompetitorRepo{q: db.New(pool)}
}

func toDomainSnapshot(s db.CompetitorSnapshot) *competitor.Snapshot {
	return &competitor.Snapshot{
		ID:              s.ID,
		AnalysisID:      uuidPtr(s.AnalysisID),
		TrackingJobID:   uuidPtr(s.TrackingJobID),
		VideoID:         s.VideoID,
		VideoTitle:      s.VideoTitle,
		ChannelID:       s.ChannelID,
		ChannelName:     s.ChannelName,
		ThumbnailURL:    s.ThumbnailUrl,
		ViewCount:       int8Val(s.ViewCount),
		LikeCount:       int8Val(s.LikeCount),
		SubscriberCount: int8Val(s.SubscriberCount),
		RankPosition:    int4Val(s.RankPosition),
		Keyword:         textVal(s.Keyword),
		CVResults:       s.CvResults,
		Score:           int4Ptr(s.Score),
		SnapshotDate:    s.SnapshotDate.Time,
		CreatedAt:       tsVal(s.CreatedAt),
	}
}

func (r *CompetitorRepo) CreateSnapshot(ctx context.Context, s *competitor.Snapshot) (*competitor.Snapshot, error) {
	created, err := r.q.CreateCompetitorSnapshot(ctx, db.CreateCompetitorSnapshotParams{
		AnalysisID:      uuidOrNil(s.AnalysisID),
		TrackingJobID:   uuidOrNil(s.TrackingJobID),
		VideoID:         s.VideoID,
		VideoTitle:      s.VideoTitle,
		ChannelID:       s.ChannelID,
		ChannelName:     s.ChannelName,
		ThumbnailUrl:    s.ThumbnailURL,
		ViewCount:       int8OrNil(s.ViewCount),
		LikeCount:       int8OrNil(s.LikeCount),
		SubscriberCount: int8OrNil(s.SubscriberCount),
		RankPosition:    int4OrNil(&s.RankPosition),
		Keyword:         textOrNil(s.Keyword),
		CvResults:       s.CVResults,
		Score:           int4OrNil(s.Score),
	})
	if err != nil {
		return nil, err
	}
	return toDomainSnapshot(created), nil
}

func (r *CompetitorRepo) ListByAnalysis(ctx context.Context, analysisID uuid.UUID) ([]*competitor.Snapshot, error) {
	rows, err := r.q.ListCompetitorsByAnalysis(ctx, uuidOrNil(&analysisID))
	if err != nil {
		return nil, err
	}
	out := make([]*competitor.Snapshot, 0, len(rows))
	for _, s := range rows {
		out = append(out, toDomainSnapshot(s))
	}
	return out, nil
}

func (r *CompetitorRepo) ListLatestByKeyword(ctx context.Context, keyword string, limit int) ([]*competitor.Snapshot, error) {
	rows, err := r.q.ListLatestCompetitorsByKeyword(ctx, db.ListLatestCompetitorsByKeywordParams{
		Keyword: textOrNil(keyword), Limit: int32(limit),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*competitor.Snapshot, 0, len(rows))
	for _, s := range rows {
		out = append(out, toDomainSnapshot(s))
	}
	return out, nil
}

func (r *CompetitorRepo) CreateTrackingJob(ctx context.Context, t *competitor.TrackingJob) (*competitor.TrackingJob, error) {
	created, err := r.q.CreateTrackingJob(ctx, db.CreateTrackingJobParams{
		WorkspaceID:        t.WorkspaceID,
		Type:               t.Type,
		ChannelID:          textOrNil(t.ChannelID),
		Keyword:            textOrNil(t.Keyword),
		CheckIntervalHours: int32(t.CheckIntervalHours),
		CreatedBy:          t.CreatedBy,
	})
	if err != nil {
		return nil, err
	}
	return toDomainTrackingJob(created), nil
}

func toDomainTrackingJob(t db.TrackingJob) *competitor.TrackingJob {
	return &competitor.TrackingJob{
		ID:                 t.ID,
		WorkspaceID:        t.WorkspaceID,
		Type:               t.Type,
		ChannelID:          textVal(t.ChannelID),
		Keyword:            textVal(t.Keyword),
		Status:             t.Status,
		CheckIntervalHours: int(t.CheckIntervalHours),
		LastCheckedAt:      tsPtr(t.LastCheckedAt),
		NextCheckAt:        tsVal(t.NextCheckAt),
		CreatedBy:          t.CreatedBy,
		CreatedAt:          tsVal(t.CreatedAt),
	}
}

func (r *CompetitorRepo) ListTrackingJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*competitor.TrackingJob, error) {
	rows, err := r.q.ListTrackingJobsByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	out := make([]*competitor.TrackingJob, 0, len(rows))
	for _, t := range rows {
		out = append(out, toDomainTrackingJob(t))
	}
	return out, nil
}

func (r *CompetitorRepo) ListDueTrackingJobs(ctx context.Context) ([]*competitor.TrackingJob, error) {
	rows, err := r.q.ListDueTrackingJobs(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*competitor.TrackingJob, 0, len(rows))
	for _, t := range rows {
		out = append(out, toDomainTrackingJob(t))
	}
	return out, nil
}

func (r *CompetitorRepo) MarkTrackingJobChecked(ctx context.Context, id uuid.UUID) error {
	return r.q.MarkTrackingJobChecked(ctx, id)
}
