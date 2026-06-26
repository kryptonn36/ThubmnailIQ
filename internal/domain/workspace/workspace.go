package workspace

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Workspace struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Slug              string    `json:"slug"`
	Plan              string    `json:"plan"`
	OwnerID           uuid.UUID `json:"owner_id"`
	AnalysesThisMonth int       `json:"analyses_this_month"`
	AnalysesLimit     int       `json:"analyses_limit"`
	CreatedAt         time.Time `json:"created_at"`
}

type Member struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type Repository interface {
	Create(ctx context.Context, name, slug string, ownerID uuid.UUID) (*Workspace, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Workspace, error)
	ListForUser(ctx context.Context, userID uuid.UUID) ([]*Workspace, error)
	AddMember(ctx context.Context, workspaceID, userID uuid.UUID, role string) (*Member, error)
	ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]*Member, error)
	IncrementAnalysesUsage(ctx context.Context, workspaceID uuid.UUID) error
	UpdatePlan(ctx context.Context, workspaceID uuid.UUID, plan string, analysesLimit int) (*Workspace, error)
}
