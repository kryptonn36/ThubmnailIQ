package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID
	Email         string
	PasswordHash  string
	FullName      string
	AvatarURL     string
	EmailVerified bool
	CreatedAt     time.Time
}

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	IsRevoked bool
}

type Repository interface {
	Create(ctx context.Context, email, passwordHash, fullName string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash, deviceInfo string, expiresAt time.Time) (*RefreshToken, error)
	GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
}
