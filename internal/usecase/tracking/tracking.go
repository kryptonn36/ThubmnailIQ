package tracking

import (
	"context"

	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/competitor"
)

type Usecase struct {
	repo competitor.Repository
}

func NewUsecase(repo competitor.Repository) *Usecase {
	return &Usecase{repo: repo}
}

type CreateParams struct {
	WorkspaceID   uuid.UUID
	UserID        uuid.UUID
	Type          string
	ChannelID     string
	Keyword       string
	IntervalHours int
}

func (u *Usecase) Create(ctx context.Context, p CreateParams) (*competitor.TrackingJob, error) {
	interval := p.IntervalHours
	if interval <= 0 {
		interval = 24
	}
	return u.repo.CreateTrackingJob(ctx, &competitor.TrackingJob{
		WorkspaceID: p.WorkspaceID, Type: p.Type, ChannelID: p.ChannelID,
		Keyword: p.Keyword, CheckIntervalHours: interval, CreatedBy: p.UserID,
	})
}

func (u *Usecase) List(ctx context.Context, workspaceID uuid.UUID) ([]*competitor.TrackingJob, error) {
	return u.repo.ListTrackingJobsByWorkspace(ctx, workspaceID)
}
