package errors

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrInvalidInput     = errors.New("invalid input")
	ErrRateLimited      = errors.New("rate limited")
	ErrQuotaExceeded    = errors.New("quota exceeded")
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}
