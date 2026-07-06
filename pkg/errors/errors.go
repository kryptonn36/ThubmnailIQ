package errors

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrInvalidInput     = errors.New("invalid input")
	ErrEmailNotVerified = errors.New("email not verified")
	ErrRateLimited      = errors.New("rate limited")
	ErrQuotaExceeded    = errors.New("quota exceeded")
	ErrPlanAlreadyActive = errors.New("plan already active")
	ErrDowngradeNotAllowed = errors.New("downgrade not allowed")
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}
