package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/admin"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/postgres/db"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
	"github.com/thumbnailiq/thumbnailiq/pkg/pagination"
)

type AdminRepo struct {
	q *db.Queries
}

func NewAdminRepo(pool *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{q: db.New(pool)}
}

func toDomainAdmin(a db.AdminUser) *admin.Admin {
	return &admin.Admin{
		ID:           a.ID,
		Email:        a.Email,
		PasswordHash: a.PasswordHash,
		FullName:     a.FullName,
		Role:         a.Role,
		IsActive:     a.IsActive,
		LastLoginAt:  tsPtr(a.LastLoginAt),
		CreatedAt:    tsVal(a.CreatedAt),
	}
}

func toDomainRefreshToken(rt db.AdminRefreshToken) *admin.RefreshToken {
	return &admin.RefreshToken{
		ID:        rt.ID,
		AdminID:   rt.AdminID,
		TokenHash: rt.TokenHash,
		ExpiresAt: tsVal(rt.ExpiresAt),
		IsRevoked: rt.IsRevoked,
	}
}

func toUserSummary(u db.User, plan string, analysesCount int, storageUsed int64) *admin.UserSummary {
	return &admin.UserSummary{
		ID: u.ID, Email: u.Email, FullName: u.FullName, Status: u.Status,
		Plan: plan, AnalysesCount: analysesCount, StorageUsedBytes: storageUsed,
		CreatedAt: tsVal(u.CreatedAt), DeletedAt: tsPtr(u.DeletedAt),
	}
}

func toUploadSummary(a db.Analysis) *admin.UploadSummary {
	return &admin.UploadSummary{
		ID: a.ID, WorkspaceID: a.WorkspaceID, UserID: a.UserID, Keyword: a.Keyword,
		ThumbnailS3Key: a.ThumbnailS3Key, Status: a.Status,
		Score: int4Ptr(a.Score), FileSizeBytes: int8Ptr(a.FileSizeBytes),
		CreatedAt: tsVal(a.CreatedAt), DeletedAt: tsPtr(a.DeletedAt),
	}
}

func toDomainAuditLog(l db.AdminAuditLog) *admin.AuditLog {
	return &admin.AuditLog{
		ID:         l.ID,
		AdminID:    l.AdminID,
		Action:     l.Action,
		TargetType: l.TargetType,
		TargetID:   l.TargetID,
		Metadata:   json.RawMessage(l.Metadata),
		CreatedAt:  tsVal(l.CreatedAt),
	}
}

func (r *AdminRepo) CreateAdmin(ctx context.Context, email, passwordHash, fullName, role string) (*admin.Admin, error) {
	a, err := r.q.CreateAdmin(ctx, db.CreateAdminParams{
		Email: email, PasswordHash: passwordHash, FullName: fullName, Role: role,
	})
	if err != nil {
		return nil, err
	}
	return toDomainAdmin(a), nil
}

func (r *AdminRepo) GetAdminByID(ctx context.Context, id uuid.UUID) (*admin.Admin, error) {
	a, err := r.q.GetAdminByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toDomainAdmin(a), nil
}

func (r *AdminRepo) GetAdminByEmail(ctx context.Context, email string) (*admin.Admin, error) {
	a, err := r.q.GetAdminByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toDomainAdmin(a), nil
}

func (r *AdminRepo) UpdateAdminLastLogin(ctx context.Context, id uuid.UUID) error {
	return r.q.UpdateAdminLastLogin(ctx, id)
}

func (r *AdminRepo) UpdateAdminPassword(ctx context.Context, id uuid.UUID, newPasswordHash string) error {
	return r.q.UpdateAdminPassword(ctx, db.UpdateAdminPasswordParams{ID: id, PasswordHash: newPasswordHash})
}

func (r *AdminRepo) CreateRefreshToken(ctx context.Context, adminID uuid.UUID, tokenHash, deviceInfo string, expiresAt time.Time) (*admin.RefreshToken, error) {
	rt, err := r.q.CreateAdminRefreshToken(ctx, db.CreateAdminRefreshTokenParams{
		AdminID: adminID, TokenHash: tokenHash, DeviceInfo: textOrNil(deviceInfo), ExpiresAt: tsNow(expiresAt),
	})
	if err != nil {
		return nil, err
	}
	return toDomainRefreshToken(rt), nil
}

func (r *AdminRepo) GetRefreshToken(ctx context.Context, tokenHash string) (*admin.RefreshToken, error) {
	rt, err := r.q.GetAdminRefreshToken(ctx, tokenHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toDomainRefreshToken(rt), nil
}

func (r *AdminRepo) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	return r.q.RevokeAdminRefreshToken(ctx, tokenHash)
}

func (r *AdminRepo) CreateAuditLog(ctx context.Context, adminID uuid.UUID, action, targetType, targetID string, metadata []byte) (*admin.AuditLog, error) {
	l, err := r.q.CreateAuditLog(ctx, db.CreateAuditLogParams{
		AdminID: adminID, Action: action, TargetType: targetType, TargetID: targetID, Metadata: metadata,
	})
	if err != nil {
		return nil, err
	}
	return toDomainAuditLog(l), nil
}

func (r *AdminRepo) ListAuditLogs(ctx context.Context, p admin.Pagination) ([]*admin.AuditLog, int, error) {
	page, perPage := pagination.Normalize(p.Page, p.PerPage)
	rows, err := r.q.ListAuditLogs(ctx, db.ListAuditLogsParams{
		Limit: int32(perPage), Offset: int32(pagination.Offset(page, perPage)),
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := r.q.CountAuditLogs(ctx)
	if err != nil {
		return nil, 0, err
	}
	logs := make([]*admin.AuditLog, len(rows))
	for i, row := range rows {
		logs[i] = toDomainAuditLog(row)
	}
	return logs, int(total), nil
}

func (r *AdminRepo) DashboardStats(ctx context.Context) (*admin.DashboardStats, error) {
	totalUsers, err := r.q.CountUsers(ctx)
	if err != nil {
		return nil, err
	}
	activeUsers, err := r.q.CountActiveUsers(ctx)
	if err != nil {
		return nil, err
	}
	totalUploads, err := r.q.CountUploads(ctx)
	if err != nil {
		return nil, err
	}
	storageUsed, err := r.q.SumStorageUsed(ctx)
	if err != nil {
		return nil, err
	}
	dailyUploads, err := r.q.CountUploadsToday(ctx)
	if err != nil {
		return nil, err
	}
	monthlyUploads, err := r.q.CountUploadsThisMonth(ctx)
	if err != nil {
		return nil, err
	}
	recentUsersRows, err := r.q.ListRecentUsers(ctx, 5)
	if err != nil {
		return nil, err
	}
	recentUploadsRows, err := r.q.ListRecentUploads(ctx, 5)
	if err != nil {
		return nil, err
	}

	recentUsers := make([]*admin.UserSummary, len(recentUsersRows))
	for i, u := range recentUsersRows {
		// The recent-activity widget only needs identity fields, so plan/
		// analyses-count/storage are left at zero here rather than paying
		// for the same per-user joins ListUsers uses for its full table.
		recentUsers[i] = toUserSummary(u, "", 0, 0)
	}
	recentUploads := make([]*admin.UploadSummary, len(recentUploadsRows))
	for i, a := range recentUploadsRows {
		recentUploads[i] = toUploadSummary(a)
	}

	return &admin.DashboardStats{
		TotalUsers: totalUsers, ActiveUsers: activeUsers, TotalUploads: totalUploads,
		StorageUsedBytes: storageUsed, DailyUploads: dailyUploads, MonthlyUploads: monthlyUploads,
		RecentUsers: recentUsers, RecentUploads: recentUploads,
	}, nil
}

func (r *AdminRepo) ListUsers(ctx context.Context, filter admin.UserFilter, p admin.Pagination) ([]*admin.UserSummary, int, error) {
	page, perPage := pagination.Normalize(p.Page, p.PerPage)
	rows, err := r.q.ListUsersAdmin(ctx, db.ListUsersAdminParams{
		Limit: int32(perPage), Offset: int32(pagination.Offset(page, perPage)),
		Search: filter.Search, StatusFilter: filter.Status,
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := r.q.CountUsersAdmin(ctx, db.CountUsersAdminParams{Search: filter.Search, StatusFilter: filter.Status})
	if err != nil {
		return nil, 0, err
	}
	users := make([]*admin.UserSummary, len(rows))
	for i, row := range rows {
		users[i] = &admin.UserSummary{
			ID: row.ID, Email: row.Email, FullName: row.FullName, Status: row.Status,
			Plan: row.Plan, AnalysesCount: int(row.AnalysesCount), StorageUsedBytes: row.StorageUsed,
			CreatedAt: tsVal(row.CreatedAt), DeletedAt: tsPtr(row.DeletedAt),
		}
	}
	return users, int(total), nil
}

func (r *AdminRepo) GetUserDetail(ctx context.Context, id uuid.UUID) (*admin.UserDetail, error) {
	row, err := r.q.GetUserDetailAdmin(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	workspaceID := row.WorkspaceID
	return &admin.UserDetail{
		UserSummary: admin.UserSummary{
			ID: row.ID, Email: row.Email, FullName: row.FullName, Status: row.Status,
			Plan: row.Plan, AnalysesCount: int(row.AnalysesCount), StorageUsedBytes: row.StorageUsed,
			CreatedAt: tsVal(row.CreatedAt), DeletedAt: tsPtr(row.DeletedAt),
		},
		WorkspaceID: &workspaceID,
	}, nil
}

func (r *AdminRepo) SuspendUser(ctx context.Context, id uuid.UUID) error {
	return r.q.SuspendUserAdmin(ctx, id)
}

func (r *AdminRepo) ActivateUser(ctx context.Context, id uuid.UUID) error {
	return r.q.ActivateUserAdmin(ctx, id)
}

func (r *AdminRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.q.SoftDeleteUserAdmin(ctx, id)
}

func (r *AdminRepo) ResetUserPassword(ctx context.Context, id uuid.UUID, newPasswordHash string) error {
	return r.q.ResetUserPasswordAdmin(ctx, db.ResetUserPasswordAdminParams{ID: id, PasswordHash: textOrNil(newPasswordHash)})
}

func (r *AdminRepo) ChangeUserWorkspaceRole(ctx context.Context, userID, workspaceID uuid.UUID, role string) error {
	return r.q.ChangeUserWorkspaceRoleAdmin(ctx, db.ChangeUserWorkspaceRoleAdminParams{
		UserID: userID, WorkspaceID: workspaceID, Role: role,
	})
}

func (r *AdminRepo) ListUserUploads(ctx context.Context, userID uuid.UUID, p admin.Pagination) ([]*admin.UploadSummary, int, error) {
	page, perPage := pagination.Normalize(p.Page, p.PerPage)
	rows, err := r.q.ListUserUploadsAdmin(ctx, db.ListUserUploadsAdminParams{
		UserID: userID, Limit: int32(perPage), Offset: int32(pagination.Offset(page, perPage)),
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := r.q.CountUserUploads(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	uploads := make([]*admin.UploadSummary, len(rows))
	for i, row := range rows {
		uploads[i] = toUploadSummary(row)
	}
	return uploads, int(total), nil
}

func (r *AdminRepo) ListUploads(ctx context.Context, filter admin.UploadFilter, p admin.Pagination) ([]*admin.UploadSummary, int, error) {
	page, perPage := pagination.Normalize(p.Page, p.PerPage)
	rows, err := r.q.ListUploadsAdmin(ctx, db.ListUploadsAdminParams{
		Limit: int32(perPage), Offset: int32(pagination.Offset(page, perPage)),
		IncludeDeleted: filter.IncludeDeleted, Search: filter.Search, StatusFilter: filter.Status,
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := r.q.CountUploadsAdmin(ctx, db.CountUploadsAdminParams{
		IncludeDeleted: filter.IncludeDeleted, Search: filter.Search, StatusFilter: filter.Status,
	})
	if err != nil {
		return nil, 0, err
	}
	uploads := make([]*admin.UploadSummary, len(rows))
	for i, row := range rows {
		uploads[i] = toUploadSummary(row)
	}
	return uploads, int(total), nil
}

func (r *AdminRepo) GetUpload(ctx context.Context, id uuid.UUID) (*admin.UploadSummary, error) {
	row, err := r.q.GetUploadAdmin(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toUploadSummary(row), nil
}

func (r *AdminRepo) DeleteUpload(ctx context.Context, id uuid.UUID) error {
	return r.q.SoftDeleteUploadAdmin(ctx, id)
}

func (r *AdminRepo) RestoreUpload(ctx context.Context, id uuid.UUID) error {
	return r.q.RestoreUploadAdmin(ctx, id)
}

func (r *AdminRepo) GetAnalytics(ctx context.Context) (*admin.Analytics, error) {
	dailySignups, err := r.q.DailySignupTrend(ctx)
	if err != nil {
		return nil, err
	}
	monthlySignups, err := r.q.MonthlySignupTrend(ctx)
	if err != nil {
		return nil, err
	}
	dailyUploads, err := r.q.DailyUploadTrend(ctx)
	if err != nil {
		return nil, err
	}
	storageUsed, err := r.q.SumStorageUsed(ctx)
	if err != nil {
		return nil, err
	}
	topUsersRows, err := r.q.TopActiveUsers(ctx, 10)
	if err != nil {
		return nil, err
	}
	fileTypes, err := r.q.FileTypeBreakdown(ctx)
	if err != nil {
		return nil, err
	}
	apiUsage, err := r.q.APIUsageStats(ctx)
	if err != nil {
		return nil, err
	}

	toTrend := func(bucket pgtype.Date, count int64) admin.TrendPoint {
		return admin.TrendPoint{Date: dateStr(bucket), Count: count}
	}
	dailyUserTrend := make([]admin.TrendPoint, len(dailySignups))
	for i, row := range dailySignups {
		dailyUserTrend[i] = toTrend(row.Bucket, row.Count)
	}
	monthlyUserTrend := make([]admin.TrendPoint, len(monthlySignups))
	for i, row := range monthlySignups {
		monthlyUserTrend[i] = toTrend(row.Bucket, row.Count)
	}
	uploadTrend := make([]admin.TrendPoint, len(dailyUploads))
	for i, row := range dailyUploads {
		uploadTrend[i] = toTrend(row.Bucket, row.Count)
	}

	topUsers := make([]*admin.UserSummary, len(topUsersRows))
	for i, row := range topUsersRows {
		topUsers[i] = &admin.UserSummary{
			ID: row.ID, Email: row.Email, FullName: row.FullName, Status: row.Status,
			Plan: row.Plan, AnalysesCount: int(row.AnalysesCount), StorageUsedBytes: row.StorageUsed,
			CreatedAt: tsVal(row.CreatedAt), DeletedAt: tsPtr(row.DeletedAt),
		}
	}

	fileTypeBreakdown := make(map[string]int64, len(fileTypes))
	for _, row := range fileTypes {
		fileTypeBreakdown[row.Extension] = row.Count
	}

	return &admin.Analytics{
		DailyUserSignups: dailyUserTrend, MonthlyUserSignups: monthlyUserTrend,
		DailyUploadTrend: uploadTrend, StorageUsedBytes: storageUsed,
		TopActiveUsers: topUsers, FileTypeBreakdown: fileTypeBreakdown,
		APIUsage: admin.APIUsageStats{
			TotalRequestsThisMonth: apiUsage.TotalRequestsThisMonth,
			ActiveKeys:             apiUsage.ActiveKeys,
		},
	}, nil
}

func toDomainSettings(s db.AppSetting) (*admin.Settings, error) {
	flags := map[string]bool{}
	if len(s.FeatureFlags) > 0 {
		if err := json.Unmarshal(s.FeatureFlags, &flags); err != nil {
			return nil, err
		}
	}
	return &admin.Settings{
		MaxUploadSizeBytes: s.MaxUploadSizeBytes,
		AllowedExtensions:  s.AllowedExtensions,
		FeatureFlags:       flags,
		StorageProvider:    s.StorageProvider,
		EmailProvider:      s.EmailProvider,
		EmailFromAddress:   s.EmailFromAddress,
		UpdatedAt:          tsVal(s.UpdatedAt),
	}, nil
}

func (r *AdminRepo) GetSettings(ctx context.Context) (*admin.Settings, error) {
	s, err := r.q.GetAppSettings(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainSettings(s)
}

func (r *AdminRepo) UpdateSettings(ctx context.Context, settings *admin.Settings) (*admin.Settings, error) {
	flagsJSON, err := json.Marshal(settings.FeatureFlags)
	if err != nil {
		return nil, err
	}
	s, err := r.q.UpdateAppSettings(ctx, db.UpdateAppSettingsParams{
		MaxUploadSizeBytes: settings.MaxUploadSizeBytes,
		AllowedExtensions:  settings.AllowedExtensions,
		FeatureFlags:       flagsJSON,
		StorageProvider:    settings.StorageProvider,
		EmailProvider:      settings.EmailProvider,
		EmailFromAddress:   settings.EmailFromAddress,
	})
	if err != nil {
		return nil, err
	}
	return toDomainSettings(s)
}
