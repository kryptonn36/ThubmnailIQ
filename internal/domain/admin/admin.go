package admin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Admin struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FullName     string     `json:"full_name"`
	Role         string     `json:"role"`
	IsActive     bool       `json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	AdminID   uuid.UUID `json:"admin_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	IsRevoked bool      `json:"is_revoked"`
}

type AuditLog struct {
	ID         uuid.UUID       `json:"id"`
	AdminID    uuid.UUID       `json:"admin_id"`
	Action     string          `json:"action"`
	TargetType string          `json:"target_type"`
	TargetID   string          `json:"target_id"`
	Metadata   json.RawMessage `json:"metadata"`
	CreatedAt  time.Time       `json:"created_at"`
}

// UserSummary is the admin-facing view of a customer user — a superset of
// what user.User exposes to the customer themselves (adds Status, plan,
// upload/storage counters).
type UserSummary struct {
	ID               uuid.UUID  `json:"id"`
	Email            string     `json:"email"`
	FullName         string     `json:"full_name"`
	Status           string     `json:"status"`
	Plan             string     `json:"plan"`
	AnalysesCount    int        `json:"analyses_count"`
	StorageUsedBytes int64      `json:"storage_used_bytes"`
	CreatedAt        time.Time  `json:"created_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}

type UserDetail struct {
	UserSummary
	WorkspaceID *uuid.UUID `json:"workspace_id"`
}

type UploadSummary struct {
	ID             uuid.UUID  `json:"id"`
	WorkspaceID    uuid.UUID  `json:"workspace_id"`
	UserID         uuid.UUID  `json:"user_id"`
	Keyword        string     `json:"keyword"`
	ThumbnailS3Key string     `json:"thumbnail_s3_key"`
	Status         string     `json:"status"`
	Score          *int       `json:"score"`
	FileSizeBytes  *int64     `json:"file_size_bytes"`
	CreatedAt      time.Time  `json:"created_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

type SystemHealth struct {
	Database  bool `json:"database"`
	Redis     bool `json:"redis"`
	CVService bool `json:"cv_service"`
}

type DashboardStats struct {
	TotalUsers       int64            `json:"total_users"`
	ActiveUsers      int64            `json:"active_users"`
	TotalUploads     int64            `json:"total_uploads"`
	StorageUsedBytes int64            `json:"storage_used_bytes"`
	DailyUploads     int64            `json:"daily_uploads"`
	MonthlyUploads   int64            `json:"monthly_uploads"`
	RecentUsers      []*UserSummary   `json:"recent_users"`
	RecentUploads    []*UploadSummary `json:"recent_uploads"`
	SystemHealth     SystemHealth     `json:"system_health"`
}

// HealthChecker is implemented by internal/infra/health — kept as a small
// interface here (rather than importing pgxpool/redis/cv directly into the
// domain package) so the admin bounded context stays free of infra
// dependencies, matching every other domain package in this codebase.
type HealthChecker interface {
	CheckDatabase(ctx context.Context) error
	CheckRedis(ctx context.Context) error
	CheckCVService(ctx context.Context) error
}

type UserFilter struct {
	Search string
	Status string
}

type UploadFilter struct {
	Search         string
	Status         string
	IncludeDeleted bool
}

type TrendPoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type APIUsageStats struct {
	TotalRequestsThisMonth int64 `json:"total_requests_this_month"`
	ActiveKeys             int64 `json:"active_keys"`
}

type Settings struct {
	MaxUploadSizeBytes int64           `json:"max_upload_size_bytes"`
	AllowedExtensions  []string        `json:"allowed_extensions"`
	FeatureFlags       map[string]bool `json:"feature_flags"`
	StorageProvider    string          `json:"storage_provider"`
	EmailProvider      string          `json:"email_provider"`
	EmailFromAddress   string          `json:"email_from_address"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type Analytics struct {
	DailyUserSignups   []TrendPoint     `json:"daily_user_signups"`
	MonthlyUserSignups []TrendPoint     `json:"monthly_user_signups"`
	DailyUploadTrend   []TrendPoint     `json:"daily_upload_trend"`
	StorageUsedBytes   int64            `json:"storage_used_bytes"`
	TopActiveUsers     []*UserSummary   `json:"top_active_users"`
	FileTypeBreakdown  map[string]int64 `json:"file_type_breakdown"`
	APIUsage           APIUsageStats    `json:"api_usage"`
}

type Pagination struct {
	Page    int
	PerPage int
}

// Repository is the admin bounded context's single repository — it owns the
// admin identity/session/audit tables plus every admin-facing read/write
// against customer tables (users, analyses, workspaces). Keeping these off
// the customer-facing user.Repository/analysis.Repository interfaces means
// no existing repository contract has to change for the admin panel to
// exist.
type Repository interface {
	// Identity
	CreateAdmin(ctx context.Context, email, passwordHash, fullName, role string) (*Admin, error)
	GetAdminByID(ctx context.Context, id uuid.UUID) (*Admin, error)
	GetAdminByEmail(ctx context.Context, email string) (*Admin, error)
	UpdateAdminLastLogin(ctx context.Context, id uuid.UUID) error
	UpdateAdminPassword(ctx context.Context, id uuid.UUID, newPasswordHash string) error

	// Sessions
	CreateRefreshToken(ctx context.Context, adminID uuid.UUID, tokenHash, deviceInfo string, expiresAt time.Time) (*RefreshToken, error)
	GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error

	// Audit
	CreateAuditLog(ctx context.Context, adminID uuid.UUID, action, targetType, targetID string, metadata []byte) (*AuditLog, error)
	ListAuditLogs(ctx context.Context, p Pagination) ([]*AuditLog, int, error)

	// Dashboard
	DashboardStats(ctx context.Context) (*DashboardStats, error)

	// Users
	ListUsers(ctx context.Context, filter UserFilter, p Pagination) ([]*UserSummary, int, error)
	GetUserDetail(ctx context.Context, id uuid.UUID) (*UserDetail, error)
	SuspendUser(ctx context.Context, id uuid.UUID) error
	ActivateUser(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ResetUserPassword(ctx context.Context, id uuid.UUID, newPasswordHash string) error
	ChangeUserWorkspaceRole(ctx context.Context, userID, workspaceID uuid.UUID, role string) error
	ListUserUploads(ctx context.Context, userID uuid.UUID, p Pagination) ([]*UploadSummary, int, error)

	// Uploads
	ListUploads(ctx context.Context, filter UploadFilter, p Pagination) ([]*UploadSummary, int, error)
	GetUpload(ctx context.Context, id uuid.UUID) (*UploadSummary, error)
	DeleteUpload(ctx context.Context, id uuid.UUID) error
	RestoreUpload(ctx context.Context, id uuid.UUID) error

	// Analytics
	GetAnalytics(ctx context.Context) (*Analytics, error)

	// Settings
	GetSettings(ctx context.Context) (*Settings, error)
	UpdateSettings(ctx context.Context, s *Settings) (*Settings, error)
}
