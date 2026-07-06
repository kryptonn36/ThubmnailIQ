package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/user"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/workspace"
	"github.com/thumbnailiq/thumbnailiq/pkg/errors"
	"github.com/thumbnailiq/thumbnailiq/pkg/hash"
	"github.com/thumbnailiq/thumbnailiq/pkg/jwt"
	"github.com/thumbnailiq/thumbnailiq/pkg/validator"
)

const (
	otpLength      = 6
	otpTTL         = 10 * time.Minute
	maxOTPAttempts = 5
)

// Mailer is the narrow email dependency this usecase needs. The concrete
// implementation lives in internal/infra/email.
type Mailer interface {
	Send(ctx context.Context, to, subject, textBody, htmlBody string) error
}

type Usecase struct {
	users      user.Repository
	workspaces workspace.Repository
	jwt        *jwt.Service
	mailer     Mailer
	log        zerolog.Logger
}

func NewUsecase(users user.Repository, workspaces workspace.Repository, jwtSvc *jwt.Service, mailer Mailer, log zerolog.Logger) *Usecase {
	return &Usecase{users: users, workspaces: workspaces, jwt: jwtSvc, mailer: mailer, log: log}
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	User         *user.User
}

// Register creates the account and its workspace and emails a verification
// code, but deliberately does NOT log the user in: no tokens are issued until
// the email is verified (see VerifyEmail). This keeps unverified accounts out
// of the app entirely, consistent with Login rejecting unverified users.
func (u *Usecase) Register(ctx context.Context, email, password, fullName string) (*user.User, error) {
	email = validator.NormalizeEmail(email)
	fullName = strings.TrimSpace(fullName)
	if !validator.IsValidEmail(email) || !validator.IsValidPassword(password) || fullName == "" {
		return nil, errors.ErrInvalidInput
	}
	if existing, _ := u.users.GetByEmail(ctx, email); existing != nil {
		return nil, errors.ErrAlreadyExists
	}

	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		return nil, err
	}
	created, err := u.users.Create(ctx, email, passwordHash, fullName)
	if err != nil {
		return nil, err
	}

	slug := fmt.Sprintf("%s-%s", validator.Slugify(fullName), created.ID.String()[:8])
	ws, err := u.workspaces.Create(ctx, fmt.Sprintf("%s's Workspace", fullName), slug, created.ID)
	if err != nil {
		return nil, err
	}
	if _, err := u.workspaces.AddMember(ctx, ws.ID, created.ID, "owner"); err != nil {
		return nil, err
	}

	// Send the verification code best-effort: a mail/SMTP failure must not lose
	// a just-created account. The user can request a fresh code via
	// ResendVerification if this send doesn't land.
	if err := u.SendVerificationCode(ctx, created); err != nil {
		u.log.Error().Err(err).Str("user_id", created.ID.String()).Msg("sending verification code on register")
	}

	return created, nil
}

// SendVerificationCode generates a fresh OTP, invalidates any earlier codes for
// the user, stores only its hash, and emails the plaintext code. It's a no-op
// (with a warning) when no mailer is configured, so local dev works without SMTP.
func (u *Usecase) SendVerificationCode(ctx context.Context, usr *user.User) error {
	code, err := hash.GenerateNumericOTP(otpLength)
	if err != nil {
		return err
	}
	if err := u.users.InvalidateEmailVerificationCodes(ctx, usr.ID); err != nil {
		return err
	}
	if _, err := u.users.CreateEmailVerificationCode(ctx, usr.ID, hash.SHA256Hex(code), time.Now().Add(otpTTL)); err != nil {
		return err
	}

	if u.mailer == nil {
		u.log.Warn().Str("user_id", usr.ID.String()).Msg("no mailer configured; verification code generated but not emailed")
		return nil
	}

	subject := "Your ThumbnailIQ verification code"
	minutes := int(otpTTL.Minutes())
	text := fmt.Sprintf("Hi %s,\n\nYour ThumbnailIQ verification code is: %s\n\nIt expires in %d minutes. If you didn't create an account, you can ignore this email.",
		usr.FullName, code, minutes)
	html := fmt.Sprintf(`<p>Hi %s,</p><p>Your ThumbnailIQ verification code is:</p>`+
		`<p style="font-size:28px;font-weight:bold;letter-spacing:4px">%s</p>`+
		`<p>It expires in %d minutes. If you didn't create an account, you can ignore this email.</p>`,
		usr.FullName, code, minutes)
	return u.mailer.Send(ctx, usr.Email, subject, text, html)
}

// VerifyEmail checks an OTP the user submits and, on success, marks their email
// verified and logs them in (issues tokens) — the validated code is proof of
// email ownership, like clicking a magic link. Errors are deliberately uniform
// (ErrInvalidInput) so the endpoint can't be used to probe which emails exist
// or which codes are close. Because a matched code is consumed, a replay of the
// same code fails generically, so there's no token issuance without a fresh
// validated code.
func (u *Usecase) VerifyEmail(ctx context.Context, email, code string) (*AuthResult, error) {
	email = validator.NormalizeEmail(email)
	usr, err := u.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.ErrInvalidInput
	}

	vc, err := u.users.GetLatestEmailVerificationCode(ctx, usr.ID)
	if err != nil {
		return nil, errors.ErrInvalidInput
	}
	if time.Now().After(vc.ExpiresAt) || vc.Attempts >= maxOTPAttempts {
		return nil, errors.ErrInvalidInput
	}
	if hash.SHA256Hex(code) != vc.CodeHash {
		// Count the failed try so a code can't be brute-forced indefinitely.
		if incErr := u.users.IncrementEmailVerificationAttempts(ctx, vc.ID); incErr != nil {
			u.log.Error().Err(incErr).Msg("incrementing email verification attempts")
		}
		return nil, errors.ErrInvalidInput
	}

	if err := u.users.ConsumeEmailVerificationCode(ctx, vc.ID); err != nil {
		return nil, err
	}
	if err := u.users.MarkEmailVerified(ctx, usr.ID); err != nil {
		return nil, err
	}
	usr.EmailVerified = true
	return u.issueTokens(ctx, usr)
}

// ResendVerification issues a new code for an unverified account. It never
// reveals whether the email exists or is already verified — callers always see
// success — so it can't be used for account enumeration.
func (u *Usecase) ResendVerification(ctx context.Context, email string) error {
	email = validator.NormalizeEmail(email)
	usr, err := u.users.GetByEmail(ctx, email)
	if err != nil || usr.EmailVerified {
		return nil
	}
	if err := u.SendVerificationCode(ctx, usr); err != nil {
		u.log.Error().Err(err).Str("user_id", usr.ID.String()).Msg("resending verification code")
	}
	return nil
}

func (u *Usecase) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	email = validator.NormalizeEmail(email)
	usr, err := u.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}
	if !hash.CheckPassword(usr.PasswordHash, password) {
		return nil, errors.ErrUnauthorized
	}
	if usr.Status == user.StatusSuspended {
		return nil, errors.ErrForbidden
	}
	// Gate the app behind email verification. Returned only after a correct
	// password so it doesn't reveal verification status to an attacker who
	// doesn't know the password.
	if !usr.EmailVerified {
		return nil, errors.ErrEmailNotVerified
	}
	return u.issueTokens(ctx, usr)
}

func (u *Usecase) Refresh(ctx context.Context, refreshToken string) (*AuthResult, error) {
	tokenHash := hash.SHA256Hex(refreshToken)
	rt, err := u.users.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}
	usr, err := u.users.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}
	// A user suspended after they logged in still holds a valid refresh token;
	// block the rotation here so suspension takes effect within one access-token
	// TTL instead of lasting until the refresh token itself expires.
	if usr.Status == user.StatusSuspended {
		_ = u.users.RevokeRefreshToken(ctx, tokenHash)
		return nil, errors.ErrForbidden
	}
	_ = u.users.RevokeRefreshToken(ctx, tokenHash)
	return u.issueTokens(ctx, usr)
}

func (u *Usecase) issueTokens(ctx context.Context, usr *user.User) (*AuthResult, error) {
	access, ttl, err := u.jwt.GenerateAccessToken(usr.ID, usr.Email)
	if err != nil {
		return nil, err
	}
	refreshRaw, err := hash.GenerateRandomToken()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(u.jwt.RefreshTTL())
	if _, err := u.users.CreateRefreshToken(ctx, usr.ID, hash.SHA256Hex(refreshRaw), "", expiresAt); err != nil {
		return nil, err
	}

	return &AuthResult{
		AccessToken:  access,
		RefreshToken: refreshRaw,
		ExpiresIn:    int(ttl.Seconds()),
		User:         usr,
	}, nil
}
