package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/analysis"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/postgres/db"
	"github.com/thumbnailiq/thumbnailiq/pkg/pagination"

	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
)

type AnalysisRepo struct {
	q *db.Queries
}

func NewAnalysisRepo(pool *pgxpool.Pool) *AnalysisRepo {
	return &AnalysisRepo{q: db.New(pool)}
}

func toDomainAnalysis(a db.Analysis) *analysis.Analysis {
	return &analysis.Analysis{
		ID:                a.ID,
		WorkspaceID:       a.WorkspaceID,
		ProjectID:         uuidPtr(a.ProjectID),
		UserID:            a.UserID,
		Keyword:           a.Keyword,
		ThumbnailS3Key:    a.ThumbnailS3Key,
		Status:            a.Status,
		ErrorMessage:      textVal(a.ErrorMessage),
		Score:             int4Ptr(a.Score),
		VisibilityScore:   int4Ptr(a.VisibilityScore),
		ContrastScore:     int4Ptr(a.ContrastScore),
		AttentionScore:    int4Ptr(a.AttentionScore),
		MobileScore:       int4Ptr(a.MobileScore),
		BrandingScore:     int4Ptr(a.BrandingScore),
		CuriosityScore:    int4Ptr(a.CuriosityScore),
		CVResults:         a.CvResults,
		CompetitorAvg:     a.CompetitorAvg,
		Suggestions:       a.Suggestions,
		CompetitorCount:   int4Val(a.CompetitorCount),
		RankInCompetitors: int4Ptr(a.RankInCompetitors),
		RelevanceScore:     int4Ptr(a.RelevanceScore),
		RelevanceReasoning: textVal(a.RelevanceReasoning),
		ActualCTR:          float8Ptr(a.ActualCtr),
		PublishedAt:        tsPtr(a.PublishedAt),
		CreatedAt:          tsVal(a.CreatedAt),
		UpdatedAt:         tsVal(a.UpdatedAt),
	}
}

func (r *AnalysisRepo) Create(ctx context.Context, a *analysis.Analysis) (*analysis.Analysis, error) {
	created, err := r.q.CreateAnalysis(ctx, db.CreateAnalysisParams{
		WorkspaceID:       a.WorkspaceID,
		ProjectID:         uuidOrNil(a.ProjectID),
		UserID:            a.UserID,
		Keyword:           a.Keyword,
		KeywordNormalized: a.Keyword,
		ThumbnailS3Key:    a.ThumbnailS3Key,
	})
	if err != nil {
		return nil, err
	}
	return toDomainAnalysis(created), nil
}

func (r *AnalysisRepo) GetByID(ctx context.Context, id uuid.UUID) (*analysis.Analysis, error) {
	a, err := r.q.GetAnalysisByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toDomainAnalysis(a), nil
}

func (r *AnalysisRepo) List(ctx context.Context, f analysis.ListFilter) (*analysis.ListResult, error) {
	page, perPage := pagination.Normalize(f.Page, f.PerPage)
	status := pgtypeText(f.Status)

	rows, err := r.q.ListAnalysesByWorkspace(ctx, db.ListAnalysesByWorkspaceParams{
		WorkspaceID: f.WorkspaceID,
		Limit:       int32(perPage),
		Offset:      int32(pagination.Offset(page, perPage)),
		Status:      status,
	})
	if err != nil {
		return nil, err
	}
	total, err := r.q.CountAnalysesByWorkspace(ctx, db.CountAnalysesByWorkspaceParams{
		WorkspaceID: f.WorkspaceID,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*analysis.Analysis, 0, len(rows))
	for _, a := range rows {
		items = append(items, toDomainAnalysis(a))
	}
	return &analysis.ListResult{Items: items, Total: int(total), Page: page, PerPage: perPage}, nil
}

func (r *AnalysisRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status, errMsg string) error {
	_, err := r.q.UpdateAnalysisStatus(ctx, db.UpdateAnalysisStatusParams{
		ID: id, Status: status, ErrorMessage: textOrNil(errMsg),
	})
	return err
}

func (r *AnalysisRepo) UpdateResults(ctx context.Context, a *analysis.Analysis) error {
	_, err := r.q.UpdateAnalysisResults(ctx, db.UpdateAnalysisResultsParams{
		ID:                a.ID,
		Score:             int4OrNil(a.Score),
		VisibilityScore:   int4OrNil(a.VisibilityScore),
		ContrastScore:     int4OrNil(a.ContrastScore),
		AttentionScore:    int4OrNil(a.AttentionScore),
		MobileScore:       int4OrNil(a.MobileScore),
		BrandingScore:     int4OrNil(a.BrandingScore),
		CuriosityScore:    int4OrNil(a.CuriosityScore),
		CvResults:         a.CVResults,
		CompetitorAvg:     a.CompetitorAvg,
		Suggestions:       a.Suggestions,
		CompetitorCount:   int4OrNil(&a.CompetitorCount),
		RankInCompetitors: int4OrNil(a.RankInCompetitors),
		RelevanceScore:     int4OrNil(a.RelevanceScore),
		RelevanceReasoning: textOrNil(a.RelevanceReasoning),
	})
	return err
}

func (r *AnalysisRepo) CreateVersion(ctx context.Context, v *analysis.ThumbnailVersion) (*analysis.ThumbnailVersion, error) {
	created, err := r.q.CreateThumbnailVersion(ctx, db.CreateThumbnailVersionParams{
		AnalysisID:    v.AnalysisID,
		VersionNumber: int32(v.VersionNumber),
		S3Key:         v.S3Key,
		Score:         int4OrNil(v.Score),
		CvResults:     v.CVResults,
	})
	if err != nil {
		return nil, err
	}
	return &analysis.ThumbnailVersion{
		ID: created.ID, AnalysisID: created.AnalysisID, VersionNumber: int(created.VersionNumber),
		S3Key: created.S3Key, Score: int4Ptr(created.Score),
		CVResults: created.CvResults, CreatedAt: tsVal(created.CreatedAt),
	}, nil
}

func (r *AnalysisRepo) ListVersions(ctx context.Context, analysisID uuid.UUID) ([]*analysis.ThumbnailVersion, error) {
	rows, err := r.q.ListThumbnailVersions(ctx, analysisID)
	if err != nil {
		return nil, err
	}
	out := make([]*analysis.ThumbnailVersion, 0, len(rows))
	for _, v := range rows {
		out = append(out, &analysis.ThumbnailVersion{
			ID: v.ID, AnalysisID: v.AnalysisID, VersionNumber: int(v.VersionNumber),
			S3Key: v.S3Key, Score: int4Ptr(v.Score),
			CVResults: v.CvResults, CreatedAt: tsVal(v.CreatedAt),
		})
	}
	return out, nil
}

func (r *AnalysisRepo) NextVersionNumber(ctx context.Context, analysisID uuid.UUID) (int, error) {
	n, err := r.q.NextThumbnailVersionNumber(ctx, analysisID)
	return int(n), err
}

func (r *AnalysisRepo) UpdateCTR(ctx context.Context, id uuid.UUID, ctr float64, publishedAt time.Time) error {
	_, err := r.q.UpdateAnalysisCTR(ctx, db.UpdateAnalysisCTRParams{
		ID:          id,
		ActualCtr:   float8OrNil(&ctr),
		PublishedAt: tsNow(publishedAt),
	})
	return err
}
