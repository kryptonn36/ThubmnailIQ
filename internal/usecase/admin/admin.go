package admin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/admin"
	"github.com/thumbnailiq/thumbnailiq/pkg/errors"
	"github.com/thumbnailiq/thumbnailiq/pkg/hash"
	"github.com/thumbnailiq/thumbnailiq/pkg/jwt"
	"github.com/thumbnailiq/thumbnailiq/pkg/validator"
)

type Usecase struct {
	admins admin.Repository
	jwt    *jwt.Service
	health admin.HealthChecker
}

func NewUsecase(admins admin.Repository, jwtSvc *jwt.Service, health admin.HealthChecker) *Usecase {
	return &Usecase{admins: admins, jwt: jwtSvc, health: health}
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	Admin        *admin.Admin
}

// Login is the only entry point into an admin session — there is no
// self-service registration; admin accounts are created out-of-band via the
// admin-seed CLI.
func (u *Usecase) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	email = validator.NormalizeEmail(email)
	a, err := u.admins.GetAdminByEmail(ctx, email)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}
	if !a.IsActive {
		return nil, errors.ErrForbidden
	}
	if !hash.CheckPassword(a.PasswordHash, password) {
		return nil, errors.ErrUnauthorized
	}
	_ = u.admins.UpdateAdminLastLogin(ctx, a.ID)
	return u.issueTokens(ctx, a)
}

func (u *Usecase) Refresh(ctx context.Context, refreshToken string) (*AuthResult, error) {
	tokenHash := hash.SHA256Hex(refreshToken)
	rt, err := u.admins.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}
	a, err := u.admins.GetAdminByID(ctx, rt.AdminID)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}
	if !a.IsActive {
		return nil, errors.ErrForbidden
	}
	_ = u.admins.RevokeRefreshToken(ctx, tokenHash)
	return u.issueTokens(ctx, a)
}

func (u *Usecase) issueTokens(ctx context.Context, a *admin.Admin) (*AuthResult, error) {
	access, ttl, err := u.jwt.GenerateAccessToken(a.ID, a.Email)
	if err != nil {
		return nil, err
	}
	refreshRaw, err := hash.GenerateRandomToken()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(u.jwt.RefreshTTL())
	if _, err := u.admins.CreateRefreshToken(ctx, a.ID, hash.SHA256Hex(refreshRaw), "", expiresAt); err != nil {
		return nil, err
	}

	return &AuthResult{
		AccessToken:  access,
		RefreshToken: refreshRaw,
		ExpiresIn:    int(ttl.Seconds()),
		Admin:        a,
	}, nil
}

// LogAction records what an admin did for the audit trail. Every mutating
// admin usecase method (suspend/activate/delete user, delete/restore
// upload, update settings, ...) calls this after its write succeeds, and
// propagates any error from it the same way as the write itself — an audit
// entry that silently failed to record would defeat the point of having an
// audit trail at all.
func (u *Usecase) LogAction(ctx context.Context, adminID uuid.UUID, action, targetType, targetID string, metadata map[string]any) error {
	var raw []byte
	if metadata != nil {
		encoded, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		raw = encoded
	}
	_, err := u.admins.CreateAuditLog(ctx, adminID, action, targetType, targetID, raw)
	return err
}

func (u *Usecase) ListAuditLogs(ctx context.Context, p admin.Pagination) ([]*admin.AuditLog, int, error) {
	return u.admins.ListAuditLogs(ctx, p)
}

func (u *Usecase) Dashboard(ctx context.Context) (*admin.DashboardStats, error) {
	stats, err := u.admins.DashboardStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.SystemHealth = admin.SystemHealth{
		Database:  u.health.CheckDatabase(ctx) == nil,
		Redis:     u.health.CheckRedis(ctx) == nil,
		CVService: u.health.CheckCVService(ctx) == nil,
	}
	return stats, nil
}

func (u *Usecase) ListUsers(ctx context.Context, filter admin.UserFilter, p admin.Pagination) ([]*admin.UserSummary, int, error) {
	return u.admins.ListUsers(ctx, filter, p)
}

func (u *Usecase) GetUserDetail(ctx context.Context, id uuid.UUID) (*admin.UserDetail, error) {
	return u.admins.GetUserDetail(ctx, id)
}

func (u *Usecase) SuspendUser(ctx context.Context, adminID, userID uuid.UUID) error {
	if err := u.admins.SuspendUser(ctx, userID); err != nil {
		return err
	}
	return u.LogAction(ctx, adminID, "user.suspend", "user", userID.String(), nil)
}

func (u *Usecase) ActivateUser(ctx context.Context, adminID, userID uuid.UUID) error {
	if err := u.admins.ActivateUser(ctx, userID); err != nil {
		return err
	}
	return u.LogAction(ctx, adminID, "user.activate", "user", userID.String(), nil)
}

func (u *Usecase) DeleteUser(ctx context.Context, adminID, userID uuid.UUID) error {
	if err := u.admins.DeleteUser(ctx, userID); err != nil {
		return err
	}
	return u.LogAction(ctx, adminID, "user.delete", "user", userID.String(), nil)
}

// ResetUserPassword generates a random temporary password, stores its hash,
// and returns the plaintext once so the admin can relay it to the user —
// the same "shown once, never again" pattern the existing API-key creation
// flow already uses (see BillingHandler.CreateAPIKey), since this codebase
// has no outbound email delivery to send it automatically.
func (u *Usecase) ResetUserPassword(ctx context.Context, adminID, userID uuid.UUID) (string, error) {
	rawToken, err := hash.GenerateRandomToken()
	if err != nil {
		return "", err
	}
	tempPassword := rawToken[:16]
	passwordHash, err := hash.HashPassword(tempPassword)
	if err != nil {
		return "", err
	}
	if err := u.admins.ResetUserPassword(ctx, userID, passwordHash); err != nil {
		return "", err
	}
	if err := u.LogAction(ctx, adminID, "user.reset_password", "user", userID.String(), nil); err != nil {
		return "", err
	}
	return tempPassword, nil
}

var validWorkspaceRoles = map[string]bool{"editor": true, "viewer": true}

func (u *Usecase) ChangeUserWorkspaceRole(ctx context.Context, adminID, userID, workspaceID uuid.UUID, role string) error {
	// "owner" is deliberately excluded — it's assigned exactly once, to the
	// workspace creator, at registration time (see usecase/user.Register),
	// and reassigning it isn't a case this endpoint handles.
	if !validWorkspaceRoles[role] {
		return errors.ErrInvalidInput
	}
	if err := u.admins.ChangeUserWorkspaceRole(ctx, userID, workspaceID, role); err != nil {
		return err
	}
	return u.LogAction(ctx, adminID, "user.change_role", "user", userID.String(), map[string]any{
		"workspace_id": workspaceID.String(), "role": role,
	})
}

func (u *Usecase) ListUserUploads(ctx context.Context, userID uuid.UUID, p admin.Pagination) ([]*admin.UploadSummary, int, error) {
	return u.admins.ListUserUploads(ctx, userID, p)
}

func (u *Usecase) ListUploads(ctx context.Context, filter admin.UploadFilter, p admin.Pagination) ([]*admin.UploadSummary, int, error) {
	return u.admins.ListUploads(ctx, filter, p)
}

func (u *Usecase) GetUpload(ctx context.Context, id uuid.UUID) (*admin.UploadSummary, error) {
	return u.admins.GetUpload(ctx, id)
}

func (u *Usecase) DeleteUpload(ctx context.Context, adminID, uploadID uuid.UUID) error {
	if err := u.admins.DeleteUpload(ctx, uploadID); err != nil {
		return err
	}
	return u.LogAction(ctx, adminID, "upload.delete", "upload", uploadID.String(), nil)
}

func (u *Usecase) RestoreUpload(ctx context.Context, adminID, uploadID uuid.UUID) error {
	if err := u.admins.RestoreUpload(ctx, uploadID); err != nil {
		return err
	}
	return u.LogAction(ctx, adminID, "upload.restore", "upload", uploadID.String(), nil)
}

func (u *Usecase) GetAnalytics(ctx context.Context) (*admin.Analytics, error) {
	return u.admins.GetAnalytics(ctx)
}

func (u *Usecase) GetSettings(ctx context.Context) (*admin.Settings, error) {
	return u.admins.GetSettings(ctx)
}

func (u *Usecase) UpdateSettings(ctx context.Context, adminID uuid.UUID, s *admin.Settings) (*admin.Settings, error) {
	if s.MaxUploadSizeBytes <= 0 {
		return nil, errors.ErrInvalidInput
	}
	updated, err := u.admins.UpdateSettings(ctx, s)
	if err != nil {
		return nil, err
	}
	if err := u.LogAction(ctx, adminID, "settings.update", "settings", "app", nil); err != nil {
		return nil, err
	}
	return updated, nil
}

func (u *Usecase) GetProfile(ctx context.Context, adminID uuid.UUID) (*admin.Admin, error) {
	return u.admins.GetAdminByID(ctx, adminID)
}

func (u *Usecase) ChangeOwnPassword(ctx context.Context, adminID uuid.UUID, currentPassword, newPassword string) error {
	a, err := u.admins.GetAdminByID(ctx, adminID)
	if err != nil {
		return err
	}
	if !hash.CheckPassword(a.PasswordHash, currentPassword) {
		return errors.ErrUnauthorized
	}
	if !validator.IsValidPassword(newPassword) {
		return errors.ErrInvalidInput
	}
	newHash, err := hash.HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := u.admins.UpdateAdminPassword(ctx, adminID, newHash); err != nil {
		return err
	}
	return u.LogAction(ctx, adminID, "profile.change_password", "admin", adminID.String(), nil)
}
