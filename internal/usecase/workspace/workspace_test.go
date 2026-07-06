package workspace

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	domainuser "github.com/thumbnailiq/thumbnailiq/internal/domain/user"
	domainworkspace "github.com/thumbnailiq/thumbnailiq/internal/domain/workspace"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
)

func TestCreateBuildsSlugAndAddsOwnerMember(t *testing.T) {
	workspaces := &fakeWorkspaceRepo{
		byID: make(map[uuid.UUID]*domainworkspace.Workspace),
	}
	uc := NewUsecase(workspaces, &fakeUserRepo{})
	ownerID := uuid.New()

	ws, err := uc.Create(context.Background(), ownerID, "My First Workspace!")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	wantSlug := "my-first-workspace-" + ownerID.String()[:8]
	if ws.Slug != wantSlug {
		t.Fatalf("slug = %q, want %q", ws.Slug, wantSlug)
	}
	if len(workspaces.members) != 1 {
		t.Fatalf("expected owner member to be added, got %d members", len(workspaces.members))
	}
	if workspaces.members[0].Role != "owner" {
		t.Fatalf("member role = %q, want owner", workspaces.members[0].Role)
	}
}

func TestInviteMemberFindsExistingUserByEmail(t *testing.T) {
	invited := &domainuser.User{ID: uuid.New(), Email: "member@example.com"}
	users := &fakeUserRepo{byEmail: map[string]*domainuser.User{invited.Email: invited}}
	workspaces := &fakeWorkspaceRepo{byID: make(map[uuid.UUID]*domainworkspace.Workspace)}
	uc := NewUsecase(workspaces, users)

	member, err := uc.InviteMember(context.Background(), uuid.New(), "member@example.com", "editor")
	if err != nil {
		t.Fatalf("InviteMember returned error: %v", err)
	}
	if member.UserID != invited.ID {
		t.Fatalf("member user id = %s, want %s", member.UserID, invited.ID)
	}
	if member.Role != "editor" {
		t.Fatalf("member role = %q, want editor", member.Role)
	}
}

func TestInviteMemberReturnsNotFoundForUnknownEmail(t *testing.T) {
	uc := NewUsecase(&fakeWorkspaceRepo{}, &fakeUserRepo{})

	_, err := uc.InviteMember(context.Background(), uuid.New(), "missing@example.com", "viewer")
	if err != apperrors.ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

type fakeWorkspaceRepo struct {
	byID    map[uuid.UUID]*domainworkspace.Workspace
	members []*domainworkspace.Member
}

func (r *fakeWorkspaceRepo) Create(_ context.Context, name, slug string, ownerID uuid.UUID) (*domainworkspace.Workspace, error) {
	ws := &domainworkspace.Workspace{
		ID:        uuid.New(),
		Name:      name,
		Slug:      slug,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
	}
	if r.byID == nil {
		r.byID = make(map[uuid.UUID]*domainworkspace.Workspace)
	}
	r.byID[ws.ID] = ws
	return ws, nil
}

func (r *fakeWorkspaceRepo) GetByID(_ context.Context, id uuid.UUID) (*domainworkspace.Workspace, error) {
	ws := r.byID[id]
	if ws == nil {
		return nil, apperrors.ErrNotFound
	}
	return ws, nil
}

func (r *fakeWorkspaceRepo) ListForUser(_ context.Context, userID uuid.UUID) ([]*domainworkspace.Workspace, error) {
	var out []*domainworkspace.Workspace
	for _, ws := range r.byID {
		if ws.OwnerID == userID {
			out = append(out, ws)
		}
	}
	return out, nil
}

func (r *fakeWorkspaceRepo) AddMember(_ context.Context, workspaceID, userID uuid.UUID, role string) (*domainworkspace.Member, error) {
	member := &domainworkspace.Member{
		ID:     uuid.New(),
		UserID: userID,
		Role:   role,
	}
	r.members = append(r.members, member)
	return member, nil
}

func (r *fakeWorkspaceRepo) ListMembers(_ context.Context, _ uuid.UUID) ([]*domainworkspace.Member, error) {
	return r.members, nil
}

func (r *fakeWorkspaceRepo) IsMember(_ context.Context, _, userID uuid.UUID) (bool, error) {
	for _, m := range r.members {
		if m.UserID == userID {
			return true, nil
		}
	}
	return false, nil
}

func (r *fakeWorkspaceRepo) IncrementAnalysesUsage(context.Context, uuid.UUID) error { return nil }

func (r *fakeWorkspaceRepo) UpdatePlan(_ context.Context, id uuid.UUID, plan string, limit int) (*domainworkspace.Workspace, error) {
	ws := r.byID[id]
	if ws == nil {
		return nil, apperrors.ErrNotFound
	}
	ws.Plan = plan
	ws.AnalysesLimit = limit
	return ws, nil
}

func (r *fakeWorkspaceRepo) UpdateBrand(_ context.Context, id uuid.UUID, primary, secondary, font string) (*domainworkspace.Workspace, error) {
	ws := r.byID[id]
	if ws == nil {
		return nil, apperrors.ErrNotFound
	}
	ws.BrandPrimaryColor = primary
	ws.BrandSecondaryColor = secondary
	ws.BrandFont = font
	return ws, nil
}

type fakeUserRepo struct {
	byEmail map[string]*domainuser.User
}

func (r *fakeUserRepo) Create(_ context.Context, email, passwordHash, fullName string) (*domainuser.User, error) {
	usr := &domainuser.User{ID: uuid.New(), Email: email, PasswordHash: passwordHash, FullName: fullName}
	if r.byEmail == nil {
		r.byEmail = make(map[string]*domainuser.User)
	}
	r.byEmail[email] = usr
	return usr, nil
}

func (r *fakeUserRepo) GetByID(_ context.Context, _ uuid.UUID) (*domainuser.User, error) {
	return nil, apperrors.ErrNotFound
}

func (r *fakeUserRepo) GetByEmail(_ context.Context, email string) (*domainuser.User, error) {
	usr := r.byEmail[email]
	if usr == nil {
		return nil, apperrors.ErrNotFound
	}
	return usr, nil
}

func (r *fakeUserRepo) CreateRefreshToken(context.Context, uuid.UUID, string, string, time.Time) (*domainuser.RefreshToken, error) {
	return nil, nil
}

func (r *fakeUserRepo) GetRefreshToken(context.Context, string) (*domainuser.RefreshToken, error) {
	return nil, apperrors.ErrNotFound
}

func (r *fakeUserRepo) RevokeRefreshToken(context.Context, string) error { return nil }
