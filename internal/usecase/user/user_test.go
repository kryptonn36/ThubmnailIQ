package user

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	domainuser "github.com/thumbnailiq/thumbnailiq/internal/domain/user"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/workspace"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
	"github.com/thumbnailiq/thumbnailiq/pkg/hash"
	"github.com/thumbnailiq/thumbnailiq/pkg/jwt"
)

func TestRegisterNormalizesEmailAndIssuesRefreshToken(t *testing.T) {
	users := newFakeUserRepo(t)
	workspaces := &fakeWorkspaceRepo{}
	uc := newTestUsecase(users, workspaces)

	res, err := uc.Register(context.Background(), "  USER@Example.COM  ", "password123", "  Jane Creator  ")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if res.User.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", res.User.Email)
	}
	if res.User.FullName != "Jane Creator" {
		t.Fatalf("expected trimmed full name, got %q", res.User.FullName)
	}
	if res.AccessToken == "" {
		t.Fatal("expected access token")
	}
	if res.RefreshToken == "" {
		t.Fatal("expected refresh token")
	}
	if res.ExpiresIn != 900 {
		t.Fatalf("expected 900 second access token ttl, got %d", res.ExpiresIn)
	}

	stored := users.refreshByHash[hash.SHA256Hex(res.RefreshToken)]
	if stored == nil {
		t.Fatal("expected refresh token hash to be stored")
	}
	if stored.TokenHash == res.RefreshToken {
		t.Fatal("refresh token was stored in raw form instead of hashed form")
	}
	if stored.UserID != res.User.ID {
		t.Fatalf("refresh token stored for wrong user: got %s want %s", stored.UserID, res.User.ID)
	}
	if time.Until(stored.ExpiresAt) <= 6*24*time.Hour {
		t.Fatalf("refresh token expiry is too short: %s", stored.ExpiresAt)
	}

	if len(workspaces.created) != 1 {
		t.Fatalf("expected one workspace to be created, got %d", len(workspaces.created))
	}
	if len(workspaces.members) != 1 {
		t.Fatalf("expected owner membership to be created, got %d", len(workspaces.members))
	}
}

func TestRegisterRejectsBlankTrimmedFullName(t *testing.T) {
	uc := newTestUsecase(newFakeUserRepo(t), &fakeWorkspaceRepo{})

	_, err := uc.Register(context.Background(), "user@example.com", "password123", "   ")
	if err != apperrors.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestLoginNormalizesEmailAndIssuesRefreshToken(t *testing.T) {
	users := newFakeUserRepo(t)
	uc := newTestUsecase(users, &fakeWorkspaceRepo{})
	users.seedUser("user@example.com", "password123", "Jane Creator")

	res, err := uc.Login(context.Background(), "  USER@Example.COM  ", "password123")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}

	if res.AccessToken == "" {
		t.Fatal("expected access token")
	}
	if res.RefreshToken == "" {
		t.Fatal("expected refresh token")
	}
	if users.refreshByHash[hash.SHA256Hex(res.RefreshToken)] == nil {
		t.Fatal("expected refresh token hash to be stored")
	}
}

func TestRefreshRotatesTokenAndRejectsReuse(t *testing.T) {
	users := newFakeUserRepo(t)
	uc := newTestUsecase(users, &fakeWorkspaceRepo{})
	users.seedUser("user@example.com", "password123", "Jane Creator")

	loginRes, err := uc.Login(context.Background(), "user@example.com", "password123")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	oldRefreshHash := hash.SHA256Hex(loginRes.RefreshToken)

	refreshRes, err := uc.Refresh(context.Background(), loginRes.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh returned error: %v", err)
	}

	if refreshRes.AccessToken == "" {
		t.Fatal("expected rotated access token")
	}
	if refreshRes.RefreshToken == "" {
		t.Fatal("expected rotated refresh token")
	}
	if refreshRes.RefreshToken == loginRes.RefreshToken {
		t.Fatal("expected refresh token rotation to issue a different token")
	}
	if !users.refreshByHash[oldRefreshHash].IsRevoked {
		t.Fatal("expected old refresh token to be revoked")
	}
	if users.refreshByHash[hash.SHA256Hex(refreshRes.RefreshToken)] == nil {
		t.Fatal("expected new refresh token hash to be stored")
	}

	_, err = uc.Refresh(context.Background(), loginRes.RefreshToken)
	if err != apperrors.ErrUnauthorized {
		t.Fatalf("expected old refresh token reuse to be unauthorized, got %v", err)
	}
}

func newTestUsecase(users *fakeUserRepo, workspaces *fakeWorkspaceRepo) *Usecase {
	return NewUsecase(
		users,
		workspaces,
		jwt.NewService("test-access-secret", "test-refresh-secret", 15*time.Minute, 7*24*time.Hour),
	)
}

type fakeUserRepo struct {
	t             *testing.T
	byID          map[uuid.UUID]*domainuser.User
	byEmail       map[string]*domainuser.User
	refreshByHash map[string]*domainuser.RefreshToken
}

func newFakeUserRepo(t *testing.T) *fakeUserRepo {
	t.Helper()
	return &fakeUserRepo{
		t:             t,
		byID:          make(map[uuid.UUID]*domainuser.User),
		byEmail:       make(map[string]*domainuser.User),
		refreshByHash: make(map[string]*domainuser.RefreshToken),
	}
}

func (r *fakeUserRepo) seedUser(email, password, fullName string) *domainuser.User {
	r.t.Helper()
	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		r.t.Fatalf("hash password: %v", err)
	}
	usr := &domainuser.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		CreatedAt:    time.Now(),
	}
	r.byID[usr.ID] = usr
	r.byEmail[email] = usr
	return usr
}

func (r *fakeUserRepo) Create(_ context.Context, email, passwordHash, fullName string) (*domainuser.User, error) {
	if _, exists := r.byEmail[email]; exists {
		return nil, apperrors.ErrAlreadyExists
	}
	usr := &domainuser.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		CreatedAt:    time.Now(),
	}
	r.byID[usr.ID] = usr
	r.byEmail[email] = usr
	return usr, nil
}

func (r *fakeUserRepo) GetByID(_ context.Context, id uuid.UUID) (*domainuser.User, error) {
	usr := r.byID[id]
	if usr == nil {
		return nil, apperrors.ErrNotFound
	}
	return usr, nil
}

func (r *fakeUserRepo) GetByEmail(_ context.Context, email string) (*domainuser.User, error) {
	usr := r.byEmail[email]
	if usr == nil {
		return nil, apperrors.ErrNotFound
	}
	return usr, nil
}

func (r *fakeUserRepo) CreateRefreshToken(_ context.Context, userID uuid.UUID, tokenHash, _ string, expiresAt time.Time) (*domainuser.RefreshToken, error) {
	rt := &domainuser.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	r.refreshByHash[tokenHash] = rt
	return rt, nil
}

func (r *fakeUserRepo) GetRefreshToken(_ context.Context, tokenHash string) (*domainuser.RefreshToken, error) {
	rt := r.refreshByHash[tokenHash]
	if rt == nil || rt.IsRevoked || time.Now().After(rt.ExpiresAt) {
		return nil, apperrors.ErrNotFound
	}
	return rt, nil
}

func (r *fakeUserRepo) RevokeRefreshToken(_ context.Context, tokenHash string) error {
	rt := r.refreshByHash[tokenHash]
	if rt == nil {
		return apperrors.ErrNotFound
	}
	rt.IsRevoked = true
	return nil
}

type fakeWorkspaceRepo struct {
	created []*workspace.Workspace
	members []*workspace.Member
}

func (r *fakeWorkspaceRepo) Create(_ context.Context, name, slug string, ownerID uuid.UUID) (*workspace.Workspace, error) {
	ws := &workspace.Workspace{
		ID:      uuid.New(),
		Name:    name,
		Slug:    slug,
		OwnerID: ownerID,
	}
	r.created = append(r.created, ws)
	return ws, nil
}

func (r *fakeWorkspaceRepo) GetByID(_ context.Context, _ uuid.UUID) (*workspace.Workspace, error) {
	return nil, apperrors.ErrNotFound
}

func (r *fakeWorkspaceRepo) ListForUser(_ context.Context, _ uuid.UUID) ([]*workspace.Workspace, error) {
	return r.created, nil
}

func (r *fakeWorkspaceRepo) AddMember(_ context.Context, workspaceID, userID uuid.UUID, role string) (*workspace.Member, error) {
	member := &workspace.Member{
		ID:     uuid.New(),
		UserID: userID,
		Role:   role,
	}
	r.members = append(r.members, member)
	return member, nil
}

func (r *fakeWorkspaceRepo) ListMembers(_ context.Context, _ uuid.UUID) ([]*workspace.Member, error) {
	return r.members, nil
}

func (r *fakeWorkspaceRepo) IsMember(_ context.Context, _, _ uuid.UUID) (bool, error) {
	return true, nil
}

func (r *fakeWorkspaceRepo) IncrementAnalysesUsage(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (r *fakeWorkspaceRepo) UpdatePlan(_ context.Context, _ uuid.UUID, _ string, _ int) (*workspace.Workspace, error) {
	return nil, apperrors.ErrNotFound
}

func (r *fakeWorkspaceRepo) UpdateBrand(_ context.Context, _ uuid.UUID, _, _, _ string) (*workspace.Workspace, error) {
	return nil, apperrors.ErrNotFound
}
