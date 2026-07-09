package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// StatusActive / StatusSuspended mirror the users.status column
// (see migration 20240101000014). A suspended user must not be able to log in
// or refresh a session.
const (
	StatusActive    = "active"
	StatusSuspended = "suspended"
)

type User struct {
	ID            uuid.UUID
	Email         string
	PasswordHash  string
	FullName      string
	AvatarURL     string
	EmailVerified bool
	Status        string
	CreatedAt     time.Time
}

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	IsRevoked bool
}

// EmailVerificationCode is a short-lived OTP emailed to the user. Only the
// SHA-256 hash of the code is stored; the plaintext lives only in the email.
type EmailVerificationCode struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CodeHash  string
	Attempts  int
	ExpiresAt time.Time
	Consumed  bool
}

// PasswordResetCode is a short-lived OTP emailed to a user who requested a
// password reset. Same storage posture as EmailVerificationCode.
type PasswordResetCode struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	CodeHash  string
	Attempts  int
	ExpiresAt time.Time
	Consumed  bool
}

type Repository interface {
	Create(ctx context.Context, email, passwordHash, fullName string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash, deviceInfo string, expiresAt time.Time) (*RefreshToken, error)
	GetRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error

	CreateEmailVerificationCode(ctx context.Context, userID uuid.UUID, codeHash string, expiresAt time.Time) (*EmailVerificationCode, error)
	GetLatestEmailVerificationCode(ctx context.Context, userID uuid.UUID) (*EmailVerificationCode, error)
	IncrementEmailVerificationAttempts(ctx context.Context, id uuid.UUID) error
	ConsumeEmailVerificationCode(ctx context.Context, id uuid.UUID) error
	InvalidateEmailVerificationCodes(ctx context.Context, userID uuid.UUID) error
	MarkEmailVerified(ctx context.Context, userID uuid.UUID) error

	CreatePasswordResetCode(ctx context.Context, userID uuid.UUID, codeHash string, expiresAt time.Time) (*PasswordResetCode, error)
	GetLatestPasswordResetCode(ctx context.Context, userID uuid.UUID) (*PasswordResetCode, error)
	IncrementPasswordResetAttempts(ctx context.Context, id uuid.UUID) error
	ConsumePasswordResetCode(ctx context.Context, id uuid.UUID) error
	InvalidatePasswordResetCodes(ctx context.Context, userID uuid.UUID) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	RevokeAllRefreshTokens(ctx context.Context, userID uuid.UUID) error
}
