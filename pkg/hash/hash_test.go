package hash

import "testing"

func TestHashPasswordAndCheckPassword(t *testing.T) {
	hashed, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if hashed == "password123" {
		t.Fatal("password hash should not equal raw password")
	}
	if !CheckPassword(hashed, "password123") {
		t.Fatal("expected correct password to match hash")
	}
	if CheckPassword(hashed, "wrong-password") {
		t.Fatal("expected wrong password to fail hash check")
	}
}

func TestSHA256Hex(t *testing.T) {
	got := SHA256Hex("thumbnailiq")
	want := "549e16a2bf8301e1bdde8aaec88c1431e91ec15c4a4453a95a2a678fb39f79ae"
	if got != want {
		t.Fatalf("SHA256Hex() = %q, want %q", got, want)
	}
}

func TestGenerateTokens(t *testing.T) {
	apiKey, prefix, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("GenerateAPIKey returned error: %v", err)
	}
	if len(apiKey) != 52 {
		t.Fatalf("GenerateAPIKey raw length = %d, want 52", len(apiKey))
	}
	if len(prefix) != 12 {
		t.Fatalf("GenerateAPIKey prefix length = %d, want 12", len(prefix))
	}
	if prefix != apiKey[:12] {
		t.Fatal("GenerateAPIKey prefix does not match raw key prefix")
	}

	refresh, err := GenerateRandomToken()
	if err != nil {
		t.Fatalf("GenerateRandomToken returned error: %v", err)
	}
	if len(refresh) != 64 {
		t.Fatalf("GenerateRandomToken length = %d, want 64", len(refresh))
	}
}
