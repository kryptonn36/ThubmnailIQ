package analysis

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusComplete   = "complete"
	StatusFailed     = "failed"
)

type Analysis struct {
	ID                uuid.UUID
	WorkspaceID       uuid.UUID
	ProjectID         *uuid.UUID
	UserID            uuid.UUID
	Keyword           string
	ThumbnailS3Key    string
	Status            string
	ErrorMessage      string
	Score             *int
	VisibilityScore   *int
	ContrastScore     *int
	AttentionScore    *int
	MobileScore       *int
	BrandingScore     *int
	CuriosityScore    *int
	CVResults         []byte
	CompetitorAvg     []byte
	Suggestions       []byte
	CompetitorCount   int
	RankInCompetitors *int
	RelevanceScore    *int
	RelevanceReasoning string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ThumbnailVersion struct {
	ID            uuid.UUID
	AnalysisID    uuid.UUID
	VersionNumber int
	S3Key         string
	Score         *int
	CVResults     []byte
	CreatedAt     time.Time
}

type ListFilter struct {
	WorkspaceID uuid.UUID
	Status      *string
	Page        int
	PerPage     int
}

type ListResult struct {
	Items   []*Analysis
	Total   int
	Page    int
	PerPage int
}

type Repository interface {
	Create(ctx context.Context, a *Analysis) (*Analysis, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Analysis, error)
	List(ctx context.Context, f ListFilter) (*ListResult, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status, errMsg string) error
	UpdateResults(ctx context.Context, a *Analysis) error
	CreateVersion(ctx context.Context, v *ThumbnailVersion) (*ThumbnailVersion, error)
	ListVersions(ctx context.Context, analysisID uuid.UUID) ([]*ThumbnailVersion, error)
	NextVersionNumber(ctx context.Context, analysisID uuid.UUID) (int, error)
}
