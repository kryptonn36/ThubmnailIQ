package validator

import "testing"

func TestNormalizeEmail(t *testing.T) {
	got := NormalizeEmail("  USER.Name+tag@Example.COM  ")
	want := "user.name+tag@example.com"
	if got != want {
		t.Fatalf("NormalizeEmail() = %q, want %q", got, want)
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{name: "valid", email: "user@example.com", want: true},
		{name: "missing domain", email: "user@", want: false},
		{name: "missing at", email: "user.example.com", want: false},
		{name: "contains spaces", email: "user @example.com", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEmail(tt.email); got != tt.want {
				t.Fatalf("IsValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	if IsValidPassword("short") {
		t.Fatal("expected short password to be invalid")
	}
	if !IsValidPassword("long-enough") {
		t.Fatal("expected password with at least 8 characters to be valid")
	}
}

func TestNormalizeKeywordAndSlugify(t *testing.T) {
	if got := NormalizeKeyword("  Viral Thumbnail Ideas  "); got != "viral thumbnail ideas" {
		t.Fatalf("NormalizeKeyword() = %q", got)
	}
	if got := Slugify(" Jane's Best Thumbnail! "); got != "jane-s-best-thumbnail" {
		t.Fatalf("Slugify() = %q", got)
	}
}
