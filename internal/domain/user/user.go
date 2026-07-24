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

// PendingRegistration holds a not-yet-verified signup: enough to create the
// real account once its OTP is confirmed, and nothing more. It is never
// written to the users table — see PendingRegistrationStore — so an email
// that's never verified leaves no permanent record behind.
type PendingRegistration struct {
	Email        string
	PasswordHash string
	FullName     string
	CodeHash     string
	Attempts     int
	ExpiresAt    time.Time
}

// PendingRegistrationStore persists PendingRegistrations between Register and
// VerifyEmail. Entries carry a TTL so an abandoned signup simply expires
// instead of leaving unverified data behind; the concrete implementation
// lives in internal/infra/redis, mirroring payment.PendingOrderStore.
type PendingRegistrationStore interface {
	Save(ctx context.Context, reg PendingRegistration) error
	Get(ctx context.Context, email string) (*PendingRegistration, error)
	IncrementAttempts(ctx context.Context, email string) error
	// Consume atomically reads and deletes the record, so a verification
	// code can't be replayed to create the account twice.
	Consume(ctx context.Context, email string) (*PendingRegistration, error)
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

	MarkEmailVerified(ctx context.Context, userID uuid.UUID) error

	CreatePasswordResetCode(ctx context.Context, userID uuid.UUID, codeHash string, expiresAt time.Time) (*PasswordResetCode, error)
	GetLatestPasswordResetCode(ctx context.Context, userID uuid.UUID) (*PasswordResetCode, error)
	IncrementPasswordResetAttempts(ctx context.Context, id uuid.UUID) error
	ConsumePasswordResetCode(ctx context.Context, id uuid.UUID) error
	InvalidatePasswordResetCodes(ctx context.Context, userID uuid.UUID) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	RevokeAllRefreshTokens(ctx context.Context, userID uuid.UUID) error
}
