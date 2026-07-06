package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/user"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/workspace"
	"github.com/thumbnailiq/thumbnailiq/pkg/errors"
	"github.com/thumbnailiq/thumbnailiq/pkg/hash"
	"github.com/thumbnailiq/thumbnailiq/pkg/jwt"
	"github.com/thumbnailiq/thumbnailiq/pkg/validator"
)

type Usecase struct {
	users      user.Repository
	workspaces workspace.Repository
	jwt        *jwt.Service
}

func NewUsecase(users user.Repository, workspaces workspace.Repository, jwtSvc *jwt.Service) *Usecase {
	return &Usecase{users: users, workspaces: workspaces, jwt: jwtSvc}
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	User         *user.User
}

func (u *Usecase) Register(ctx context.Context, email, password, fullName string) (*AuthResult, error) {
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

	return u.issueTokens(ctx, created)
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
