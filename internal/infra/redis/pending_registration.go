package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/user"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
)

// PendingRegistrationStore persists PendingRegistrations in Redis for the
// short window between signup and email verification, mirroring how
// PendingOrderStore holds checkout state between order creation and payment
// confirmation. The fixed ttl below is a generous backstop; the actual
// verification deadline is the ExpiresAt the usecase stores on the record
// (see internal/usecase/user), which is checked explicitly on every verify
// attempt — so Redis eventually cleaning up an abandoned signup on its own
// is a bonus, not the enforcement mechanism.
type PendingRegistrationStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewPendingRegistrationStore(client *redis.Client) *PendingRegistrationStore {
	return &PendingRegistrationStore{client: client, ttl: 30 * time.Minute}
}

type pendingRegistrationRecord struct {
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	CodeHash     string    `json:"code_hash"`
	Attempts     int       `json:"attempts"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func pendingRegistrationKey(email string) string {
	return "pending_registration:" + email
}

func (s *PendingRegistrationStore) Save(ctx context.Context, reg user.PendingRegistration) error {
	data, err := json.Marshal(pendingRegistrationRecord{
		PasswordHash: reg.PasswordHash,
		FullName:     reg.FullName,
		CodeHash:     reg.CodeHash,
		Attempts:     reg.Attempts,
		ExpiresAt:    reg.ExpiresAt,
	})
	if err != nil {
		return err
	}
	return s.client.Set(ctx, pendingRegistrationKey(reg.Email), data, s.ttl).Err()
}

func (s *PendingRegistrationStore) Get(ctx context.Context, email string) (*user.PendingRegistration, error) {
	val, err := s.client.Get(ctx, pendingRegistrationKey(email)).Result()
	if err == redis.Nil {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return decodePendingRegistration(email, val)
}

// IncrementAttempts records a failed OTP guess. It preserves the record's
// original ExpiresAt (and derives the Redis TTL from it) rather than
// resetting to the full backstop TTL, so repeated wrong guesses can't extend
// a pending signup's lifetime.
func (s *PendingRegistrationStore) IncrementAttempts(ctx context.Context, email string) error {
	reg, err := s.Get(ctx, email)
	if err != nil {
		return err
	}
	reg.Attempts++
	remaining := time.Until(reg.ExpiresAt)
	if remaining <= 0 {
		return apperrors.ErrNotFound
	}
	data, err := json.Marshal(pendingRegistrationRecord{
		PasswordHash: reg.PasswordHash,
		FullName:     reg.FullName,
		CodeHash:     reg.CodeHash,
		Attempts:     reg.Attempts,
		ExpiresAt:    reg.ExpiresAt,
	})
	if err != nil {
		return err
	}
	return s.client.Set(ctx, pendingRegistrationKey(email), data, remaining).Err()
}

// Consume atomically reads and deletes the pending registration (GetDel), so
// a verification code can't be replayed to create the account twice.
func (s *PendingRegistrationStore) Consume(ctx context.Context, email string) (*user.PendingRegistration, error) {
	val, err := s.client.GetDel(ctx, pendingRegistrationKey(email)).Result()
	if err == redis.Nil {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return decodePendingRegistration(email, val)
}

func decodePendingRegistration(email, val string) (*user.PendingRegistration, error) {
	var rec pendingRegistrationRecord
	if err := json.Unmarshal([]byte(val), &rec); err != nil {
		return nil, err
	}
	return &user.PendingRegistration{
		Email:        email,
		PasswordHash: rec.PasswordHash,
		FullName:     rec.FullName,
		CodeHash:     rec.CodeHash,
		Attempts:     rec.Attempts,
		ExpiresAt:    rec.ExpiresAt,
	}, nil
}
