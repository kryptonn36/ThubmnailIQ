package admin

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	domainadmin "github.com/thumbnailiq/thumbnailiq/internal/domain/admin"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
	"github.com/thumbnailiq/thumbnailiq/pkg/hash"
	"github.com/thumbnailiq/thumbnailiq/pkg/jwt"
)

func TestLoginNormalizesEmailAndIssuesRefreshToken(t *testing.T) {
	admins := newFakeAdminRepo(t)
	uc := newTestUsecase(admins)
	admins.seedAdmin("admin@example.com", "password123", true)

	res, err := uc.Login(context.Background(), "  ADMIN@Example.COM  ", "password123")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
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
	if admins.refreshByHash[hash.SHA256Hex(res.RefreshToken)] == nil {
		t.Fatal("expected refresh token hash to be stored")
	}
	if admins.lastLoginID == uuid.Nil {
		t.Fatal("expected last login to be updated")
	}
}

func TestLoginRejectsInactiveAdmin(t *testing.T) {
	admins := newFakeAdminRepo(t)
	uc := newTestUsecase(admins)
	admins.seedAdmin("admin@example.com", "password123", false)

	_, err := uc.Login(context.Background(), "admin@example.com", "password123")
	if err != apperrors.ErrForbidden {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestRefreshRotatesTokenAndRejectsReuse(t *testing.T) {
	admins := newFakeAdminRepo(t)
	uc := newTestUsecase(admins)
	admins.seedAdmin("admin@example.com", "password123", true)

	loginRes, err := uc.Login(context.Background(), "admin@example.com", "password123")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	oldRefreshHash := hash.SHA256Hex(loginRes.RefreshToken)

	refreshRes, err := uc.Refresh(context.Background(), loginRes.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh returned error: %v", err)
	}

	if refreshRes.RefreshToken == "" {
		t.Fatal("expected rotated refresh token")
	}
	if refreshRes.RefreshToken == loginRes.RefreshToken {
		t.Fatal("expected refresh token rotation to issue a different token")
	}
	if !admins.refreshByHash[oldRefreshHash].IsRevoked {
		t.Fatal("expected old refresh token to be revoked")
	}
	if admins.refreshByHash[hash.SHA256Hex(refreshRes.RefreshToken)] == nil {
		t.Fatal("expected new refresh token hash to be stored")
	}

	_, err = uc.Refresh(context.Background(), loginRes.RefreshToken)
	if err != apperrors.ErrUnauthorized {
		t.Fatalf("expected old refresh token reuse to be unauthorized, got %v", err)
	}
}

func newTestUsecase(admins *fakeAdminRepo) *Usecase {
	return NewUsecase(
		admins,
		jwt.NewService("test-admin-access-secret", "test-admin-refresh-secret", 15*time.Minute, 7*24*time.Hour),
		fakeHealthChecker{},
	)
}

type fakeHealthChecker struct{}

func (fakeHealthChecker) CheckDatabase(context.Context) error  { return nil }
func (fakeHealthChecker) CheckRedis(context.Context) error     { return nil }
func (fakeHealthChecker) CheckCVService(context.Context) error { return nil }

type fakeAdminRepo struct {
	t             *testing.T
	byID          map[uuid.UUID]*domainadmin.Admin
	byEmail       map[string]*domainadmin.Admin
	refreshByHash map[string]*domainadmin.RefreshToken
	lastLoginID   uuid.UUID
}

func newFakeAdminRepo(t *testing.T) *fakeAdminRepo {
	t.Helper()
	return &fakeAdminRepo{
		t:             t,
		byID:          make(map[uuid.UUID]*domainadmin.Admin),
		byEmail:       make(map[string]*domainadmin.Admin),
		refreshByHash: make(map[string]*domainadmin.RefreshToken),
	}
}

func (r *fakeAdminRepo) seedAdmin(email, password string, active bool) *domainadmin.Admin {
	r.t.Helper()
	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		r.t.Fatalf("hash password: %v", err)
	}
	a := &domainadmin.Admin{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     "Admin User",
		Role:         "admin",
		IsActive:     active,
		CreatedAt:    time.Now(),
	}
	r.byID[a.ID] = a
	r.byEmail[email] = a
	return a
}

func (r *fakeAdminRepo) CreateAdmin(_ context.Context, email, passwordHash, fullName, role string) (*domainadmin.Admin, error) {
	a := &domainadmin.Admin{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		Role:         role,
		IsActive:     true,
		CreatedAt:    time.Now(),
	}
	r.byID[a.ID] = a
	r.byEmail[email] = a
	return a, nil
}

func (r *fakeAdminRepo) GetAdminByID(_ context.Context, id uuid.UUID) (*domainadmin.Admin, error) {
	a := r.byID[id]
	if a == nil {
		return nil, apperrors.ErrNotFound
	}
	return a, nil
}

func (r *fakeAdminRepo) GetAdminByEmail(_ context.Context, email string) (*domainadmin.Admin, error) {
	a := r.byEmail[email]
	if a == nil {
		return nil, apperrors.ErrNotFound
	}
	return a, nil
}

func (r *fakeAdminRepo) UpdateAdminLastLogin(_ context.Context, id uuid.UUID) error {
	r.lastLoginID = id
	return nil
}

func (r *fakeAdminRepo) UpdateAdminPassword(context.Context, uuid.UUID, string) error {
	return nil
}

func (r *fakeAdminRepo) CreateRefreshToken(_ context.Context, adminID uuid.UUID, tokenHash, _ string, expiresAt time.Time) (*domainadmin.RefreshToken, error) {
	rt := &domainadmin.RefreshToken{
		ID:        uuid.New(),
		AdminID:   adminID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	r.refreshByHash[tokenHash] = rt
	return rt, nil
}

func (r *fakeAdminRepo) GetRefreshToken(_ context.Context, tokenHash string) (*domainadmin.RefreshToken, error) {
	rt := r.refreshByHash[tokenHash]
	if rt == nil || rt.IsRevoked || time.Now().After(rt.ExpiresAt) {
		return nil, apperrors.ErrNotFound
	}
	return rt, nil
}

func (r *fakeAdminRepo) RevokeRefreshToken(_ context.Context, tokenHash string) error {
	rt := r.refreshByHash[tokenHash]
	if rt == nil {
		return apperrors.ErrNotFound
	}
	rt.IsRevoked = true
	return nil
}

func (r *fakeAdminRepo) CreateAuditLog(context.Context, uuid.UUID, string, string, string, []byte) (*domainadmin.AuditLog, error) {
	return &domainadmin.AuditLog{ID: uuid.New()}, nil
}

func (r *fakeAdminRepo) ListAuditLogs(context.Context, domainadmin.Pagination) ([]*domainadmin.AuditLog, int, error) {
	return nil, 0, nil
}

func (r *fakeAdminRepo) DashboardStats(context.Context) (*domainadmin.DashboardStats, error) {
	return &domainadmin.DashboardStats{}, nil
}

func (r *fakeAdminRepo) ListUsers(context.Context, domainadmin.UserFilter, domainadmin.Pagination) ([]*domainadmin.UserSummary, int, error) {
	return nil, 0, nil
}

func (r *fakeAdminRepo) GetUserDetail(context.Context, uuid.UUID) (*domainadmin.UserDetail, error) {
	return nil, apperrors.ErrNotFound
}

func (r *fakeAdminRepo) SuspendUser(context.Context, uuid.UUID) error { return nil }

func (r *fakeAdminRepo) ActivateUser(context.Context, uuid.UUID) error { return nil }

func (r *fakeAdminRepo) DeleteUser(context.Context, uuid.UUID) error { return nil }

func (r *fakeAdminRepo) ResetUserPassword(context.Context, uuid.UUID, string) error { return nil }

func (r *fakeAdminRepo) ChangeUserWorkspaceRole(context.Context, uuid.UUID, uuid.UUID, string) error {
	return nil
}

func (r *fakeAdminRepo) ListUserUploads(context.Context, uuid.UUID, domainadmin.Pagination) ([]*domainadmin.UploadSummary, int, error) {
	return nil, 0, nil
}

func (r *fakeAdminRepo) ListUploads(context.Context, domainadmin.UploadFilter, domainadmin.Pagination) ([]*domainadmin.UploadSummary, int, error) {
	return nil, 0, nil
}

func (r *fakeAdminRepo) GetUpload(context.Context, uuid.UUID) (*domainadmin.UploadSummary, error) {
	return nil, apperrors.ErrNotFound
}

func (r *fakeAdminRepo) DeleteUpload(context.Context, uuid.UUID) error { return nil }

func (r *fakeAdminRepo) RestoreUpload(context.Context, uuid.UUID) error { return nil }

func (r *fakeAdminRepo) GetAnalytics(context.Context) (*domainadmin.Analytics, error) {
	return &domainadmin.Analytics{}, nil
}

func (r *fakeAdminRepo) GetSettings(context.Context) (*domainadmin.Settings, error) {
	return &domainadmin.Settings{}, nil
}

func (r *fakeAdminRepo) UpdateSettings(_ context.Context, s *domainadmin.Settings) (*domainadmin.Settings, error) {
	return s, nil
}
