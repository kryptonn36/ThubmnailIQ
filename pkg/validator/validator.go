package validator

import (
	"regexp"
	"strings"
)

var emailRe = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func IsValidEmail(email string) bool {
	return emailRe.MatchString(email)
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func IsValidPassword(password string) bool {
	return len(password) >= 8
}

func NormalizeKeyword(keyword string) string {
	return strings.ToLower(strings.TrimSpace(keyword))
}

func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}
