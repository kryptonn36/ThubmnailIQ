package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateAndParseAccessToken(t *testing.T) {
	svc := NewService("access-secret", "refresh-secret", 15*time.Minute, 7*24*time.Hour)
	userID := uuid.New()

	token, ttl, err := svc.GenerateAccessToken(userID, "user@example.com")
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}
	if token == "" {
		t.Fatal("expected signed token")
	}
	if ttl != 15*time.Minute {
		t.Fatalf("ttl = %s, want 15m", ttl)
	}

	claims, err := svc.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("ParseAccessToken returned error: %v", err)
	}
	if claims.UserID != userID {
		t.Fatalf("claims user id = %s, want %s", claims.UserID, userID)
	}
	if claims.Email != "user@example.com" {
		t.Fatalf("claims email = %q", claims.Email)
	}
}

func TestParseAccessTokenRejectsWrongSecret(t *testing.T) {
	svc := NewService("access-secret", "refresh-secret", 15*time.Minute, 7*24*time.Hour)
	token, _, err := svc.GenerateAccessToken(uuid.New(), "user@example.com")
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	otherSvc := NewService("other-secret", "refresh-secret", 15*time.Minute, 7*24*time.Hour)
	if _, err := otherSvc.ParseAccessToken(token); err == nil {
		t.Fatal("expected token signed with different secret to be rejected")
	}
}

func TestRefreshTTL(t *testing.T) {
	svc := NewService("access-secret", "refresh-secret", 15*time.Minute, 14*24*time.Hour)
	if got := svc.RefreshTTL(); got != 14*24*time.Hour {
		t.Fatalf("RefreshTTL() = %s", got)
	}
}
