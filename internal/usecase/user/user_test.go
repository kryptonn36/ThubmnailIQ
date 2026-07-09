package user

import (
	"context"
	"testing"
	"time"

	"regexp"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	domainuser "github.com/thumbnailiq/thumbnailiq/internal/domain/user"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/workspace"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
	"github.com/thumbnailiq/thumbnailiq/pkg/hash"
	"github.com/thumbnailiq/thumbnailiq/pkg/jwt"
)

func TestRegisterNormalizesEmailCreatesWorkspaceAndQueuesCode(t *testing.T) {
	users := newFakeUserRepo(t)
	workspaces := &fakeWorkspaceRepo{}
	uc := newTestUsecase(users, workspaces)

	usr, err := uc.Register(context.Background(), "  USER@Example.COM  ", "password123", "  Jane Creator  ")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if usr.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", usr.Email)
	}
	if usr.FullName != "Jane Creator" {
		t.Fatalf("expected trimmed full name, got %q", usr.FullName)
	}
	// Registration must NOT log the user in: no tokens until email is verified.
	if usr.EmailVerified {
		t.Fatal("new account should start unverified")
	}
	if len(users.refreshByHash) != 0 {
		t.Fatal("Register must not issue a refresh token before verification")
	}
	// A verification code should have been generated for the new user.
	if users.verifications[usr.ID] == nil {
		t.Fatal("expected a verification code to be queued on register")
	}

	if len(workspaces.created) != 1 {
		t.Fatalf("expected one workspace to be created, got %d", len(workspaces.created))
	}
	if len(workspaces.members) != 1 {
		t.Fatalf("expected owner membership to be created, got %d", len(workspaces.members))
	}
}

func TestRegisterRejectsBlankTrimmedFullName(t *testing.T) {
	uc := newTestUsecase(newFakeUserRepo(t), &fakeWorkspaceRepo{})

	_, err := uc.Register(context.Background(), "user@example.com", "password123", "   ")
	if err != apperrors.ErrInvalidInput {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestLoginNormalizesEmailAndIssuesRefreshToken(t *testing.T) {
	users := newFakeUserRepo(t)
	uc := newTestUsecase(users, &fakeWorkspaceRepo{})
	users.seedUser("user@example.com", "password123", "Jane Creator")

	res, err := uc.Login(context.Background(), "  USER@Example.COM  ", "password123")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}

	if res.AccessToken == "" {
		t.Fatal("expected access token")
	}
	if res.RefreshToken == "" {
		t.Fatal("expected refresh token")
	}
	if users.refreshByHash[hash.SHA256Hex(res.RefreshToken)] == nil {
		t.Fatal("expected refresh token hash to be stored")
	}
}

func TestRefreshRotatesTokenAndRejectsReuse(t *testing.T) {
	users := newFakeUserRepo(t)
	uc := newTestUsecase(users, &fakeWorkspaceRepo{})
	users.seedUser("user@example.com", "password123", "Jane Creator")

	loginRes, err := uc.Login(context.Background(), "user@example.com", "password123")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	oldRefreshHash := hash.SHA256Hex(loginRes.RefreshToken)

	refreshRes, err := uc.Refresh(context.Background(), loginRes.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh returned error: %v", err)
	}

	if refreshRes.AccessToken == "" {
		t.Fatal("expected rotated access token")
	}
	if refreshRes.RefreshToken == "" {
		t.Fatal("expected rotated refresh token")
	}
	if refreshRes.RefreshToken == loginRes.RefreshToken {
		t.Fatal("expected refresh token rotation to issue a different token")
	}
	if !users.refreshByHash[oldRefreshHash].IsRevoked {
		t.Fatal("expected old refresh token to be revoked")
	}
	if users.refreshByHash[hash.SHA256Hex(refreshRes.RefreshToken)] == nil {
		t.Fatal("expected new refresh token hash to be stored")
	}

	_, err = uc.Refresh(context.Background(), loginRes.RefreshToken)
	if err != apperrors.ErrUnauthorized {
		t.Fatalf("expected old refresh token reuse to be unauthorized, got %v", err)
	}
}

func newTestUsecase(users *fakeUserRepo, workspaces *fakeWorkspaceRepo) *Usecase {
	return NewUsecase(
		users,
		workspaces,
		jwt.NewService("test-access-secret", "test-refresh-secret", 15*time.Minute, 7*24*time.Hour),
		nil,
		zerolog.Nop(),
	)
}

type captureMailer struct {
	calls    int
	lastText string
}

func (m *captureMailer) Send(_ context.Context, _, _, textBody, _ string) error {
	m.calls++
	m.lastText = textBody
	return nil
}

var otpPattern = regexp.MustCompile(`\b\d{6}\b`)

func TestEmailVerificationFlow(t *testing.T) {
	users := newFakeUserRepo(t)
	mailer := &captureMailer{}
	uc := NewUsecase(
		users,
		&fakeWorkspaceRepo{},
		jwt.NewService("a", "b", 15*time.Minute, 7*24*time.Hour),
		mailer,
		zerolog.Nop(),
	)

	if _, err := uc.Register(context.Background(), "user@example.com", "password123", "Jane"); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if mailer.calls != 1 {
		t.Fatalf("expected 1 verification email on register, got %d", mailer.calls)
	}
	code := otpPattern.FindString(mailer.lastText)
	if code == "" {
		t.Fatalf("no 6-digit code found in verification email: %q", mailer.lastText)
	}

	// A wrong code is rejected and must not verify the account.
	wrong := "000000"
	if wrong == code {
		wrong = "999999"
	}
	if _, err := uc.VerifyEmail(context.Background(), "user@example.com", wrong); err != apperrors.ErrInvalidInput {
		t.Fatalf("wrong code: expected ErrInvalidInput, got %v", err)
	}
	if usr, _ := users.GetByEmail(context.Background(), "user@example.com"); usr.EmailVerified {
		t.Fatal("account should not be verified after a wrong code")
	}

	// The correct code verifies (email match is case-insensitive via
	// normalization) and logs the user in by issuing tokens.
	res, err := uc.VerifyEmail(context.Background(), "USER@example.com", code)
	if err != nil {
		t.Fatalf("correct code: expected success, got %v", err)
	}
	if res.AccessToken == "" || res.RefreshToken == "" {
		t.Fatal("expected tokens to be issued on successful verification")
	}
	usr, _ := users.GetByEmail(context.Background(), "user@example.com")
	if !usr.EmailVerified {
		t.Fatal("account should be verified after the correct code")
	}

	// The code is single-use: replaying it (now consumed) is rejected generically.
	if _, err := uc.VerifyEmail(context.Background(), "user@example.com", code); err != apperrors.ErrInvalidInput {
		t.Fatalf("consumed code replay: expected ErrInvalidInput, got %v", err)
	}
}

func TestVerifyEmailRejectsAfterMaxAttempts(t *testing.T) {
	users := newFakeUserRepo(t)
	mailer := &captureMailer{}
	uc := NewUsecase(
		users,
		&fakeWorkspaceRepo{},
		jwt.NewService("a", "b", 15*time.Minute, 7*24*time.Hour),
		mailer,
		zerolog.Nop(),
	)
	if _, err := uc.Register(context.Background(), "user@example.com", "password123", "Jane"); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	code := otpPattern.FindString(mailer.lastText)
	wrong := "111111"
	if wrong == code {
		wrong = "222222"
	}

	// Exhaust the attempt budget with wrong guesses.
	for i := 0; i < maxOTPAttempts; i++ {
		if _, err := uc.VerifyEmail(context.Background(), "user@example.com", wrong); err != apperrors.ErrInvalidInput {
			t.Fatalf("attempt %d: expected ErrInvalidInput, got %v", i, err)
		}
	}
	// Even the correct code is now refused because the code is locked.
	if _, err := uc.VerifyEmail(context.Background(), "user@example.com", code); err != apperrors.ErrInvalidInput {
		t.Fatalf("expected lockout after max attempts, got %v", err)
	}
}

func TestPasswordResetFlow(t *testing.T) {
	users := newFakeUserRepo(t)
	mailer := &captureMailer{}
	uc := NewUsecase(
		users,
		&fakeWorkspaceRepo{},
		jwt.NewService("a", "b", 15*time.Minute, 7*24*time.Hour),
		mailer,
		zerolog.Nop(),
	)
	// An established, verified user with an active session.
	usr := users.seedUser("user@example.com", "oldpassword1", "Jane")
	users.refreshByHash["existing"] = &domainuser.RefreshToken{UserID: usr.ID, TokenHash: "existing"}

	if err := uc.RequestPasswordReset(context.Background(), "  USER@Example.com "); err != nil {
		t.Fatalf("RequestPasswordReset returned error: %v", err)
	}
	if mailer.calls != 1 {
		t.Fatalf("expected one reset email, got %d", mailer.calls)
	}
	code := otpPattern.FindString(mailer.lastText)
	if code == "" {
		t.Fatalf("no code in reset email: %q", mailer.lastText)
	}

	// Wrong code is rejected and doesn't change the password.
	wrong := "000000"
	if wrong == code {
		wrong = "999999"
	}
	if err := uc.ResetPassword(context.Background(), "user@example.com", wrong, "newpassword1"); err != apperrors.ErrInvalidInput {
		t.Fatalf("wrong code: expected ErrInvalidInput, got %v", err)
	}

	// A too-short password is rejected before any code check.
	if err := uc.ResetPassword(context.Background(), "user@example.com", code, "short"); err != apperrors.ErrInvalidInput {
		t.Fatalf("weak password: expected ErrInvalidInput, got %v", err)
	}

	// Correct code + valid password succeeds.
	if err := uc.ResetPassword(context.Background(), "user@example.com", code, "newpassword1"); err != nil {
		t.Fatalf("reset: expected success, got %v", err)
	}
	// New password now works; old one doesn't.
	if _, err := uc.Login(context.Background(), "user@example.com", "newpassword1"); err != nil {
		t.Fatalf("login with new password failed: %v", err)
	}
	if _, err := uc.Login(context.Background(), "user@example.com", "oldpassword1"); err != apperrors.ErrUnauthorized {
		t.Fatalf("old password should be rejected, got %v", err)
	}
	// Existing sessions were revoked.
	if !users.refreshByHash["existing"].IsRevoked {
		t.Fatal("expected existing refresh tokens to be revoked after reset")
	}
	// The consumed code can't be reused.
	if err := uc.ResetPassword(context.Background(), "user@example.com", code, "another1234"); err != apperrors.ErrInvalidInput {
		t.Fatalf("consumed code replay: expected ErrInvalidInput, got %v", err)
	}
}

type fakeUserRepo struct {
	t             *testing.T
	byID          map[uuid.UUID]*domainuser.User
	byEmail       map[string]*domainuser.User
	refreshByHash map[string]*domainuser.RefreshToken
	verifications map[uuid.UUID]*domainuser.EmailVerificationCode
	resets        map[uuid.UUID]*domainuser.PasswordResetCode
}

func newFakeUserRepo(t *testing.T) *fakeUserRepo {
	t.Helper()
	return &fakeUserRepo{
		t:             t,
		byID:          make(map[uuid.UUID]*domainuser.User),
		byEmail:       make(map[string]*domainuser.User),
		refreshByHash: make(map[string]*domainuser.RefreshToken),
	}
}

func (r *fakeUserRepo) seedUser(email, password, fullName string) *domainuser.User {
	r.t.Helper()
	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		r.t.Fatalf("hash password: %v", err)
	}
	usr := &domainuser.User{
		ID:            uuid.New(),
		Email:         email,
		PasswordHash:  passwordHash,
		FullName:      fullName,
		EmailVerified: true, // seeded users represent established, verified accounts
		CreatedAt:     time.Now(),
	}
	r.byID[usr.ID] = usr
	r.byEmail[email] = usr
	return usr
}

func (r *fakeUserRepo) Create(_ context.Context, email, passwordHash, fullName string) (*domainuser.User, error) {
	if _, exists := r.byEmail[email]; exists {
		return nil, apperrors.ErrAlreadyExists
	}
	usr := &domainuser.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		CreatedAt:    time.Now(),
	}
	r.byID[usr.ID] = usr
	r.byEmail[email] = usr
	return usr, nil
}

func (r *fakeUserRepo) GetByID(_ context.Context, id uuid.UUID) (*domainuser.User, error) {
	usr := r.byID[id]
	if usr == nil {
		return nil, apperrors.ErrNotFound
	}
	return usr, nil
}

func (r *fakeUserRepo) GetByEmail(_ context.Context, email string) (*domainuser.User, error) {
	usr := r.byEmail[email]
	if usr == nil {
		return nil, apperrors.ErrNotFound
	}
	return usr, nil
}

func (r *fakeUserRepo) CreateRefreshToken(_ context.Context, userID uuid.UUID, tokenHash, _ string, expiresAt time.Time) (*domainuser.RefreshToken, error) {
	rt := &domainuser.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	r.refreshByHash[tokenHash] = rt
	return rt, nil
}

func (r *fakeUserRepo) GetRefreshToken(_ context.Context, tokenHash string) (*domainuser.RefreshToken, error) {
	rt := r.refreshByHash[tokenHash]
	if rt == nil || rt.IsRevoked || time.Now().After(rt.ExpiresAt) {
		return nil, apperrors.ErrNotFound
	}
	return rt, nil
}

func (r *fakeUserRepo) RevokeRefreshToken(_ context.Context, tokenHash string) error {
	rt := r.refreshByHash[tokenHash]
	if rt == nil {
		return apperrors.ErrNotFound
	}
	rt.IsRevoked = true
	return nil
}

func (r *fakeUserRepo) CreateEmailVerificationCode(_ context.Context, userID uuid.UUID, codeHash string, expiresAt time.Time) (*domainuser.EmailVerificationCode, error) {
	if r.verifications == nil {
		r.verifications = make(map[uuid.UUID]*domainuser.EmailVerificationCode)
	}
	vc := &domainuser.EmailVerificationCode{ID: uuid.New(), UserID: userID, CodeHash: codeHash, ExpiresAt: expiresAt}
	r.verifications[userID] = vc
	return vc, nil
}

func (r *fakeUserRepo) GetLatestEmailVerificationCode(_ context.Context, userID uuid.UUID) (*domainuser.EmailVerificationCode, error) {
	vc := r.verifications[userID]
	if vc == nil || vc.Consumed {
		return nil, apperrors.ErrNotFound
	}
	return vc, nil
}

func (r *fakeUserRepo) IncrementEmailVerificationAttempts(_ context.Context, id uuid.UUID) error {
	for _, vc := range r.verifications {
		if vc.ID == id {
			vc.Attempts++
		}
	}
	return nil
}

func (r *fakeUserRepo) ConsumeEmailVerificationCode(_ context.Context, id uuid.UUID) error {
	for _, vc := range r.verifications {
		if vc.ID == id {
			vc.Consumed = true
		}
	}
	return nil
}

func (r *fakeUserRepo) InvalidateEmailVerificationCodes(_ context.Context, userID uuid.UUID) error {
	if vc := r.verifications[userID]; vc != nil {
		vc.Consumed = true
	}
	return nil
}

func (r *fakeUserRepo) MarkEmailVerified(_ context.Context, userID uuid.UUID) error {
	if usr := r.byID[userID]; usr != nil {
		usr.EmailVerified = true
	}
	return nil
}

func (r *fakeUserRepo) CreatePasswordResetCode(_ context.Context, userID uuid.UUID, codeHash string, expiresAt time.Time) (*domainuser.PasswordResetCode, error) {
	if r.resets == nil {
		r.resets = make(map[uuid.UUID]*domainuser.PasswordResetCode)
	}
	rc := &domainuser.PasswordResetCode{ID: uuid.New(), UserID: userID, CodeHash: codeHash, ExpiresAt: expiresAt}
	r.resets[userID] = rc
	return rc, nil
}

func (r *fakeUserRepo) GetLatestPasswordResetCode(_ context.Context, userID uuid.UUID) (*domainuser.PasswordResetCode, error) {
	rc := r.resets[userID]
	if rc == nil || rc.Consumed {
		return nil, apperrors.ErrNotFound
	}
	return rc, nil
}

func (r *fakeUserRepo) IncrementPasswordResetAttempts(_ context.Context, id uuid.UUID) error {
	for _, rc := range r.resets {
		if rc.ID == id {
			rc.Attempts++
		}
	}
	return nil
}

func (r *fakeUserRepo) ConsumePasswordResetCode(_ context.Context, id uuid.UUID) error {
	for _, rc := range r.resets {
		if rc.ID == id {
			rc.Consumed = true
		}
	}
	return nil
}

func (r *fakeUserRepo) InvalidatePasswordResetCodes(_ context.Context, userID uuid.UUID) error {
	if rc := r.resets[userID]; rc != nil {
		rc.Consumed = true
	}
	return nil
}

func (r *fakeUserRepo) UpdatePassword(_ context.Context, userID uuid.UUID, passwordHash string) error {
	if usr := r.byID[userID]; usr != nil {
		usr.PasswordHash = passwordHash
	}
	return nil
}

func (r *fakeUserRepo) RevokeAllRefreshTokens(_ context.Context, userID uuid.UUID) error {
	for _, rt := range r.refreshByHash {
		if rt.UserID == userID {
			rt.IsRevoked = true
		}
	}
	return nil
}

type fakeWorkspaceRepo struct {
	created []*workspace.Workspace
	members []*workspace.Member
}

func (r *fakeWorkspaceRepo) Create(_ context.Context, name, slug string, ownerID uuid.UUID) (*workspace.Workspace, error) {
	ws := &workspace.Workspace{
		ID:      uuid.New(),
		Name:    name,
		Slug:    slug,
		OwnerID: ownerID,
	}
	r.created = append(r.created, ws)
	return ws, nil
}

func (r *fakeWorkspaceRepo) GetByID(_ context.Context, _ uuid.UUID) (*workspace.Workspace, error) {
	return nil, apperrors.ErrNotFound
}

func (r *fakeWorkspaceRepo) ListForUser(_ context.Context, _ uuid.UUID) ([]*workspace.Workspace, error) {
	return r.created, nil
}

func (r *fakeWorkspaceRepo) AddMember(_ context.Context, workspaceID, userID uuid.UUID, role string) (*workspace.Member, error) {
	member := &workspace.Member{
		ID:     uuid.New(),
		UserID: userID,
		Role:   role,
	}
	r.members = append(r.members, member)
	return member, nil
}

func (r *fakeWorkspaceRepo) ListMembers(_ context.Context, _ uuid.UUID) ([]*workspace.Member, error) {
	return r.members, nil
}

func (r *fakeWorkspaceRepo) IsMember(_ context.Context, _, _ uuid.UUID) (bool, error) {
	return true, nil
}

func (r *fakeWorkspaceRepo) IncrementAnalysesUsage(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (r *fakeWorkspaceRepo) UpdatePlan(_ context.Context, _ uuid.UUID, _ string, _ int) (*workspace.Workspace, error) {
	return nil, apperrors.ErrNotFound
}

func (r *fakeWorkspaceRepo) UpdateBrand(_ context.Context, _ uuid.UUID, _, _, _ string) (*workspace.Workspace, error) {
	return nil, apperrors.ErrNotFound
}
