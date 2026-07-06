package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/workspace"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/postgres/db"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
)

type WorkspaceRepo struct {
	q *db.Queries
}

func NewWorkspaceRepo(pool *pgxpool.Pool) *WorkspaceRepo {
	return &WorkspaceRepo{q: db.New(pool)}
}

func toDomainWorkspace(w db.Workspace) *workspace.Workspace {
	return &workspace.Workspace{
		ID:                  w.ID,
		Name:                w.Name,
		Slug:                w.Slug,
		Plan:                w.Plan,
		OwnerID:             w.OwnerID,
		AnalysesThisMonth:   int(w.AnalysesThisMonth),
		AnalysesLimit:       int(w.AnalysesLimit),
		BrandPrimaryColor:   textVal(w.BrandPrimaryColor),
		BrandSecondaryColor: textVal(w.BrandSecondaryColor),
		BrandFont:           textVal(w.BrandFont),
		CreatedAt:           tsVal(w.CreatedAt),
	}
}

func (r *WorkspaceRepo) Create(ctx context.Context, name, slug string, ownerID uuid.UUID) (*workspace.Workspace, error) {
	w, err := r.q.CreateWorkspace(ctx, db.CreateWorkspaceParams{Name: name, Slug: slug, OwnerID: ownerID})
	if err != nil {
		return nil, err
	}
	return toDomainWorkspace(w), nil
}

func (r *WorkspaceRepo) GetByID(ctx context.Context, id uuid.UUID) (*workspace.Workspace, error) {
	w, err := r.q.GetWorkspaceByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toDomainWorkspace(w), nil
}

func (r *WorkspaceRepo) ListForUser(ctx context.Context, userID uuid.UUID) ([]*workspace.Workspace, error) {
	rows, err := r.q.ListWorkspacesForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]*workspace.Workspace, 0, len(rows))
	for _, w := range rows {
		out = append(out, toDomainWorkspace(w))
	}
	return out, nil
}

func (r *WorkspaceRepo) AddMember(ctx context.Context, workspaceID, userID uuid.UUID, role string) (*workspace.Member, error) {
	m, err := r.q.AddWorkspaceMember(ctx, db.AddWorkspaceMemberParams{WorkspaceID: workspaceID, UserID: userID, Role: role})
	if err != nil {
		return nil, err
	}
	return &workspace.Member{ID: m.ID, UserID: m.UserID, Role: m.Role, JoinedAt: tsVal(m.JoinedAt)}, nil
}

func (r *WorkspaceRepo) ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]*workspace.Member, error) {
	rows, err := r.q.ListWorkspaceMembers(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	out := make([]*workspace.Member, 0, len(rows))
	for _, m := range rows {
		out = append(out, &workspace.Member{
			ID: m.ID, UserID: m.UserID, Email: m.Email, FullName: m.FullName,
			Role: m.Role, JoinedAt: tsVal(m.JoinedAt),
		})
	}
	return out, nil
}

func (r *WorkspaceRepo) IsMember(ctx context.Context, workspaceID, userID uuid.UUID) (bool, error) {
	return r.q.IsWorkspaceMember(ctx, db.IsWorkspaceMemberParams{WorkspaceID: workspaceID, UserID: userID})
}

func (r *WorkspaceRepo) IncrementAnalysesUsage(ctx context.Context, workspaceID uuid.UUID) error {
	return r.q.IncrementWorkspaceAnalysesUsage(ctx, workspaceID)
}

func (r *WorkspaceRepo) UpdatePlan(ctx context.Context, workspaceID uuid.UUID, plan string, analysesLimit int) (*workspace.Workspace, error) {
	w, err := r.q.UpdateWorkspacePlan(ctx, db.UpdateWorkspacePlanParams{ID: workspaceID, Plan: plan, AnalysesLimit: int32(analysesLimit)})
	if err != nil {
		return nil, err
	}
	return toDomainWorkspace(w), nil
}

func (r *WorkspaceRepo) UpdateBrand(ctx context.Context, workspaceID uuid.UUID, primary, secondary, font string) (*workspace.Workspace, error) {
	w, err := r.q.UpdateWorkspaceBrand(ctx, db.UpdateWorkspaceBrandParams{
		ID: workspaceID, BrandPrimaryColor: textOrNil(primary), BrandSecondaryColor: textOrNil(secondary), BrandFont: textOrNil(font),
	})
	if err != nil {
		return nil, err
	}
	return toDomainWorkspace(w), nil
}
