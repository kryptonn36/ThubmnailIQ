package workspace

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Workspace struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Slug                string    `json:"slug"`
	Plan                string    `json:"plan"`
	OwnerID             uuid.UUID `json:"owner_id"`
	AnalysesThisMonth   int       `json:"analyses_this_month"`
	AnalysesLimit       int       `json:"analyses_limit"`
	BrandPrimaryColor   string    `json:"brand_primary_color"`
	BrandSecondaryColor string    `json:"brand_secondary_color"`
	BrandFont           string    `json:"brand_font"`
	CreatedAt           time.Time `json:"created_at"`
}

type Member struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

// Summary is a workspace as seen by one specific member: the workspace itself
// plus who owns it and what the viewing user's role in it is. It exists so
// the UI can always show whose workspace is being operated on (e.g. "Owned by
// Jane Doe · you're an editor") without extra lookups per workspace.
type Summary struct {
	Workspace
	OwnerName   string `json:"owner_name"`
	OwnerEmail  string `json:"owner_email"`
	Role        string `json:"role"`
	MemberCount int    `json:"member_count"`
}

type Repository interface {
	Create(ctx context.Context, name, slug string, ownerID uuid.UUID) (*Workspace, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Workspace, error)
	ListForUser(ctx context.Context, userID uuid.UUID) ([]*Workspace, error)
	ListForUserWithContext(ctx context.Context, userID uuid.UUID) ([]*Summary, error)
	Rename(ctx context.Context, workspaceID uuid.UUID, name string) (*Workspace, error)
	AddMember(ctx context.Context, workspaceID, userID uuid.UUID, role string) (*Member, error)
	RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error
	ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]*Member, error)
	IsMember(ctx context.Context, workspaceID, userID uuid.UUID) (bool, error)
	IncrementAnalysesUsage(ctx context.Context, workspaceID uuid.UUID) error
	UpdatePlan(ctx context.Context, workspaceID uuid.UUID, plan string, analysesLimit int) (*Workspace, error)
	UpdateBrand(ctx context.Context, workspaceID uuid.UUID, primary, secondary, font string) (*Workspace, error)
}
