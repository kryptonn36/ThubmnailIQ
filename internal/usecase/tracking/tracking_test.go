package tracking

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/competitor"
)

func TestCreateDefaultsIntervalTo24Hours(t *testing.T) {
	repo := &fakeCompetitorRepo{}
	uc := NewUsecase(repo)

	job, err := uc.Create(context.Background(), CreateParams{
		WorkspaceID: uuid.New(),
		UserID:      uuid.New(),
		Type:        "keyword",
		Keyword:     "thumbnail ideas",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if job.CheckIntervalHours != 24 {
		t.Fatalf("interval = %d, want 24", job.CheckIntervalHours)
	}
	if repo.created == nil {
		t.Fatal("expected repository CreateTrackingJob to be called")
	}
}

func TestCreateUsesProvidedInterval(t *testing.T) {
	repo := &fakeCompetitorRepo{}
	uc := NewUsecase(repo)

	job, err := uc.Create(context.Background(), CreateParams{
		WorkspaceID:   uuid.New(),
		UserID:        uuid.New(),
		Type:          "channel",
		ChannelID:     "UC123",
		IntervalHours: 6,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if job.CheckIntervalHours != 6 {
		t.Fatalf("interval = %d, want 6", job.CheckIntervalHours)
	}
}

func TestListReturnsWorkspaceJobs(t *testing.T) {
	workspaceID := uuid.New()
	repo := &fakeCompetitorRepo{
		jobs: []*competitor.TrackingJob{{ID: uuid.New(), WorkspaceID: workspaceID}},
	}
	uc := NewUsecase(repo)

	jobs, err := uc.List(context.Background(), workspaceID)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("got %d jobs, want 1", len(jobs))
	}
}

type fakeCompetitorRepo struct {
	created *competitor.TrackingJob
	jobs    []*competitor.TrackingJob
}

func (r *fakeCompetitorRepo) CreateSnapshot(_ context.Context, s *competitor.Snapshot) (*competitor.Snapshot, error) {
	return s, nil
}

func (r *fakeCompetitorRepo) ListByAnalysis(context.Context, uuid.UUID) ([]*competitor.Snapshot, error) {
	return nil, nil
}

func (r *fakeCompetitorRepo) ListLatestByKeyword(context.Context, string, int) ([]*competitor.Snapshot, error) {
	return nil, nil
}

func (r *fakeCompetitorRepo) CreateTrackingJob(_ context.Context, t *competitor.TrackingJob) (*competitor.TrackingJob, error) {
	t.ID = uuid.New()
	t.CreatedAt = time.Now()
	r.created = t
	r.jobs = append(r.jobs, t)
	return t, nil
}

func (r *fakeCompetitorRepo) ListTrackingJobsByWorkspace(_ context.Context, workspaceID uuid.UUID) ([]*competitor.TrackingJob, error) {
	var out []*competitor.TrackingJob
	for _, job := range r.jobs {
		if job.WorkspaceID == workspaceID {
			out = append(out, job)
		}
	}
	return out, nil
}

func (r *fakeCompetitorRepo) ListDueTrackingJobs(context.Context) ([]*competitor.TrackingJob, error) {
	return nil, nil
}

func (r *fakeCompetitorRepo) MarkTrackingJobChecked(context.Context, uuid.UUID) error { return nil }
