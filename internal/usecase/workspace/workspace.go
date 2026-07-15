package workspace

import (
	"context"
	"fmt"
	"strings"

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

// ListWithContext returns the user's workspaces enriched with owner identity,
// the user's own role, and member count — everything the UI needs to show
// whose workspace each one is and what the viewer may do in it.
func (u *Usecase) ListWithContext(ctx context.Context, userID uuid.UUID) ([]*workspace.Summary, error) {
	return u.workspaces.ListForUserWithContext(ctx, userID)
}

// Rename changes the workspace's display name. The slug is deliberately left
// untouched: it's already used in external references and stays stable for
// the workspace's lifetime, like most SaaS products handle renames.
func (u *Usecase) Rename(ctx context.Context, workspaceID uuid.UUID, name string) (*workspace.Workspace, error) {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 100 {
		return nil, errors.ErrInvalidInput
	}
	return u.workspaces.Rename(ctx, workspaceID, name)
}

// RemoveMember removes a member from the workspace. Rules mirror common team
// products: the owner can never be removed (ownership transfer is not
// supported here), any non-owner may remove themselves (leave), and removing
// someone else requires the actor to be owner or admin.
func (u *Usecase) RemoveMember(ctx context.Context, workspaceID, actorID, memberUserID uuid.UUID) error {
	ws, err := u.workspaces.GetByID(ctx, workspaceID)
	if err != nil {
		return err
	}
	if memberUserID == ws.OwnerID {
		return errors.ErrForbidden
	}
	if actorID != memberUserID {
		if err := u.EnsureRole(ctx, actorID, workspaceID, "owner", "admin"); err != nil {
			return err
		}
	}
	return u.workspaces.RemoveMember(ctx, workspaceID, memberUserID)
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

// invitableRoles are the roles a member can be granted. "owner" is
// deliberately absent: there is exactly one owner (the creator) and ownership
// can't be handed out via invites.
var invitableRoles = map[string]bool{"admin": true, "editor": true, "viewer": true}

// InviteMember requires the invited user to already have an account; this
// MVP doesn't implement invite-by-email for not-yet-registered users.
func (u *Usecase) InviteMember(ctx context.Context, workspaceID uuid.UUID, email, role string) (*workspace.Member, error) {
	if !invitableRoles[role] {
		return nil, errors.ErrInvalidInput
	}
	invited, err := u.users.GetByEmail(ctx, validator.NormalizeEmail(email))
	if err != nil {
		return nil, errors.ErrNotFound
	}
	if already, err := u.workspaces.IsMember(ctx, workspaceID, invited.ID); err != nil {
		return nil, err
	} else if already {
		return nil, errors.ErrAlreadyExists
	}
	return u.workspaces.AddMember(ctx, workspaceID, invited.ID, role)
}

func (u *Usecase) ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]*workspace.Member, error) {
	return u.workspaces.ListMembers(ctx, workspaceID)
}

func (u *Usecase) UpdateBrand(ctx context.Context, workspaceID uuid.UUID, primary, secondary, font string) (*workspace.Workspace, error) {
	return u.workspaces.UpdateBrand(ctx, workspaceID, primary, secondary, font)
}
