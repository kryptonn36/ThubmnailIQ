package workspace

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	domainuser "github.com/thumbnailiq/thumbnailiq/internal/domain/user"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/workspace"
	"github.com/thumbnailiq/thumbnailiq/pkg/errors"
	"github.com/thumbnailiq/thumbnailiq/pkg/validator"
)

type Usecase struct {
	workspaces workspace.Repository
	users      domainuser.Repository
}

func NewUsecase(workspaces workspace.Repository, users domainuser.Repository) *Usecase {
	return &Usecase{workspaces: workspaces, users: users}
}

func (u *Usecase) Create(ctx context.Context, ownerID uuid.UUID, name string) (*workspace.Workspace, error) {
	slug := fmt.Sprintf("%s-%s", validator.Slugify(name), ownerID.String()[:8])
	ws, err := u.workspaces.Create(ctx, name, slug, ownerID)
	if err != nil {
		return nil, err
	}
	_, _ = u.workspaces.AddMember(ctx, ws.ID, ownerID, "owner")
	return ws, nil
}

func (u *Usecase) List(ctx context.Context, userID uuid.UUID) ([]*workspace.Workspace, error) {
	return u.workspaces.ListForUser(ctx, userID)
}

func (u *Usecase) Get(ctx context.Context, id uuid.UUID) (*workspace.Workspace, error) {
	return u.workspaces.GetByID(ctx, id)
}

// EnsureMember authorizes a user against a workspace. It returns ErrForbidden
// when the user is not a member, so callers can pass a client-supplied
// workspace ID straight through without trusting it. This is the single
// object-level authorization check for every workspace-scoped action.
func (u *Usecase) EnsureMember(ctx context.Context, userID, workspaceID uuid.UUID) error {
	member, err := u.workspaces.IsMember(ctx, workspaceID, userID)
	if err != nil {
		return err
	}
	if !member {
		return errors.ErrForbidden
	}
	return nil
}

// EnsureRole authorizes a user against a workspace and additionally requires
// their membership role to be one of `allowed` (e.g. "owner", "admin"). Used
// for privileged actions like inviting members or changing the billing plan,
// where plain membership isn't enough.
func (u *Usecase) EnsureRole(ctx context.Context, userID, workspaceID uuid.UUID, allowed ...string) error {
	members, err := u.workspaces.ListMembers(ctx, workspaceID)
	if err != nil {
		return err
	}
	for _, m := range members {
		if m.UserID != userID {
			continue
		}
		for _, role := range allowed {
			if m.Role == role {
				return nil
			}
		}
		return errors.ErrForbidden
	}
	return errors.ErrForbidden
}

// InviteMember requires the invited user to already have an account; this
// MVP doesn't implement invite-by-email for not-yet-registered users.
func (u *Usecase) InviteMember(ctx context.Context, workspaceID uuid.UUID, email, role string) (*workspace.Member, error) {
	invited, err := u.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return u.workspaces.AddMember(ctx, workspaceID, invited.ID, role)
}

func (u *Usecase) ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]*workspace.Member, error) {
	return u.workspaces.ListMembers(ctx, workspaceID)
}

func (u *Usecase) UpdateBrand(ctx context.Context, workspaceID uuid.UUID, primary, secondary, font string) (*workspace.Workspace, error) {
	return u.workspaces.UpdateBrand(ctx, workspaceID, primary, secondary, font)
}
