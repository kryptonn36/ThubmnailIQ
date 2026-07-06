package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/user"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/postgres/db"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
)

type UserRepo struct {
	q *db.Queries
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{q: db.New(pool)}
}

func toDomainUser(u db.User) *user.User {
	return &user.User{
		ID:            u.ID,
		Email:         u.Email,
		PasswordHash:  textVal(u.PasswordHash),
		FullName:      u.FullName,
		AvatarURL:     textVal(u.AvatarUrl),
		EmailVerified: u.EmailVerified,
		Status:        u.Status,
		CreatedAt:     tsVal(u.CreatedAt),
	}
}

func (r *UserRepo) Create(ctx context.Context, email, passwordHash, fullName string) (*user.User, error) {
	u, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: textOrNil(passwordHash),
		FullName:     fullName,
	})
	if err != nil {
		return nil, err
	}
	return toDomainUser(u), nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	u, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toDomainUser(u), nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	u, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toDomainUser(u), nil
}

func (r *UserRepo) CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash, deviceInfo string, expiresAt time.Time) (*user.RefreshToken, error) {
	rt, err := r.q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:     userID,
		TokenHash:  tokenHash,
		DeviceInfo: textOrNil(deviceInfo),
		ExpiresAt:  tsNow(expiresAt),
	})
	if err != nil {
		return nil, err
	}
	return &user.RefreshToken{
		ID:        rt.ID,
		UserID:    rt.UserID,
		TokenHash: rt.TokenHash,
		ExpiresAt: tsVal(rt.ExpiresAt),
		IsRevoked: rt.IsRevoked,
	}, nil
}

func (r *UserRepo) GetRefreshToken(ctx context.Context, tokenHash string) (*user.RefreshToken, error) {
	rt, err := r.q.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &user.RefreshToken{
		ID:        rt.ID,
		UserID:    rt.UserID,
		TokenHash: rt.TokenHash,
		ExpiresAt: tsVal(rt.ExpiresAt),
		IsRevoked: rt.IsRevoked,
	}, nil
}

func (r *UserRepo) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	return r.q.RevokeRefreshToken(ctx, tokenHash)
}

func toDomainVerificationCode(c db.EmailVerificationCode) *user.EmailVerificationCode {
	return &user.EmailVerificationCode{
		ID:        c.ID,
		UserID:    c.UserID,
		CodeHash:  c.CodeHash,
		Attempts:  int(c.Attempts),
		ExpiresAt: tsVal(c.ExpiresAt),
		Consumed:  c.ConsumedAt.Valid,
	}
}

func (r *UserRepo) CreateEmailVerificationCode(ctx context.Context, userID uuid.UUID, codeHash string, expiresAt time.Time) (*user.EmailVerificationCode, error) {
	c, err := r.q.CreateEmailVerificationCode(ctx, db.CreateEmailVerificationCodeParams{
		UserID:    userID,
		CodeHash:  codeHash,
		ExpiresAt: tsNow(expiresAt),
	})
	if err != nil {
		return nil, err
	}
	return toDomainVerificationCode(c), nil
}

func (r *UserRepo) GetLatestEmailVerificationCode(ctx context.Context, userID uuid.UUID) (*user.EmailVerificationCode, error) {
	c, err := r.q.GetLatestEmailVerificationCode(ctx, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toDomainVerificationCode(c), nil
}

func (r *UserRepo) IncrementEmailVerificationAttempts(ctx context.Context, id uuid.UUID) error {
	return r.q.IncrementEmailVerificationAttempts(ctx, id)
}

func (r *UserRepo) ConsumeEmailVerificationCode(ctx context.Context, id uuid.UUID) error {
	return r.q.ConsumeEmailVerificationCode(ctx, id)
}

func (r *UserRepo) InvalidateEmailVerificationCodes(ctx context.Context, userID uuid.UUID) error {
	return r.q.InvalidateEmailVerificationCodes(ctx, userID)
}

func (r *UserRepo) MarkEmailVerified(ctx context.Context, userID uuid.UUID) error {
	return r.q.MarkUserEmailVerified(ctx, userID)
}
