package analysis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/analysis"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/workspace"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/s3"
	"github.com/thumbnailiq/thumbnailiq/internal/scoring"
	"github.com/thumbnailiq/thumbnailiq/internal/worker"
	"github.com/thumbnailiq/thumbnailiq/pkg/errors"
	"github.com/thumbnailiq/thumbnailiq/pkg/validator"
)

type Usecase struct {
	analyses   analysis.Repository
	workspaces workspace.Repository
	storage    *s3.Storage
	cvClient   *cv.Client
	queue      *asynq.Client
	engine     *scoring.Engine
}

func NewUsecase(analyses analysis.Repository, workspaces workspace.Repository, storage *s3.Storage, cvClient *cv.Client, queue *asynq.Client) *Usecase {
	return &Usecase{
		analyses: analyses, workspaces: workspaces, storage: storage,
		cvClient: cvClient, queue: queue, engine: scoring.NewEngine(),
	}
}

type CreateParams struct {
	WorkspaceID uuid.UUID
	ProjectID   *uuid.UUID
	UserID      uuid.UUID
	Keyword     string
	FileBytes   []byte
	ContentType string
}

func (u *Usecase) Create(ctx context.Context, p CreateParams) (*analysis.Analysis, error) {
	ws, err := u.workspaces.GetByID(ctx, p.WorkspaceID)
	if err != nil {
		return nil, err
	}
	if ws.AnalysesThisMonth >= ws.AnalysesLimit {
		return nil, errors.ErrQuotaExceeded
	}

	key := fmt.Sprintf("%s/%s/original.jpg", p.WorkspaceID, uuid.New().String())
	thumbnailURL, err := u.storage.Upload(ctx, key, p.FileBytes, p.ContentType)
	if err != nil {
		return nil, err
	}

	created, err := u.analyses.Create(ctx, &analysis.Analysis{
		WorkspaceID:    p.WorkspaceID,
		ProjectID:      p.ProjectID,
		UserID:         p.UserID,
		Keyword:        validator.NormalizeKeyword(p.Keyword),
		ThumbnailURL:   thumbnailURL,
		ThumbnailS3Key: key,
	})
	if err != nil {
		return nil, err
	}

	if err := u.workspaces.IncrementAnalysesUsage(ctx, p.WorkspaceID); err != nil {
		return nil, err
	}

	task, err := worker.NewAnalyzeThumbnailTask(created.ID)
	if err != nil {
		return nil, err
	}
	if _, err := u.queue.EnqueueContext(ctx, task); err != nil {
		return nil, fmt.Errorf("enqueue analyze task: %w", err)
	}

	return created, nil
}

func (u *Usecase) Get(ctx context.Context, id uuid.UUID) (*analysis.Analysis, error) {
	return u.analyses.GetByID(ctx, id)
}

func (u *Usecase) List(ctx context.Context, f analysis.ListFilter) (*analysis.ListResult, error) {
	return u.analyses.List(ctx, f)
}

func (u *Usecase) ListVersions(ctx context.Context, analysisID uuid.UUID) ([]*analysis.ThumbnailVersion, error) {
	return u.analyses.ListVersions(ctx, analysisID)
}

// AddCompareVersion scores a new variant synchronously (a single CV call,
// unlike the full async pipeline) reusing the parent analysis's already
// computed competitor average so variants are still ranked against the
// same niche benchmark.
func (u *Usecase) AddCompareVersion(ctx context.Context, analysisID uuid.UUID, fileBytes []byte, contentType string) (*analysis.ThumbnailVersion, error) {
	parent, err := u.analyses.GetByID(ctx, analysisID)
	if err != nil {
		return nil, err
	}

	versionNum, err := u.analyses.NextVersionNumber(ctx, analysisID)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%s/%s/version-%d.jpg", parent.WorkspaceID, analysisID, versionNum)
	thumbnailURL, err := u.storage.Upload(ctx, key, fileBytes, contentType)
	if err != nil {
		return nil, err
	}

	cvResult, err := u.cvClient.Analyze(ctx, thumbnailURL)
	if err != nil {
		return nil, fmt.Errorf("cv analyze: %w", err)
	}

	var avg scoring.CompetitorAvg
	if len(parent.CompetitorAvg) > 0 {
		_ = json.Unmarshal(parent.CompetitorAvg, &avg)
	}

	sub := scoring.SubScores{
		Visibility: u.engine.Visibility(cvResult, avg),
		Contrast:   u.engine.Contrast(cvResult),
		Attention:  u.engine.Attention(cvResult),
		Mobile:     u.engine.Mobile(cvResult),
		Branding:   u.engine.Branding(),
	}
	final := u.engine.FinalScore(sub)

	cvJSON, _ := json.Marshal(cvResult)
	version, err := u.analyses.CreateVersion(ctx, &analysis.ThumbnailVersion{
		AnalysisID: analysisID, VersionNumber: versionNum, S3Key: key,
		ThumbnailURL: thumbnailURL, Score: &final, CVResults: cvJSON,
	})
	if err != nil {
		return nil, err
	}
	return version, nil
}
