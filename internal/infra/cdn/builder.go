// Package cdn is the single place in the codebase allowed to turn an S3
// object key into a public URL. Every upload flow stores just the object
// key; every read path that needs a public URL goes through Builder.URL.
package cdn

import (
	"fmt"
	"net/url"
	"strings"
)

type Builder struct {
	domain string // e.g. "https://cdn.example.com" — no trailing slash
}

func NewBuilder(domain string) *Builder {
	return &Builder{domain: strings.TrimRight(domain, "/")}
}

// URL builds the public URL for an object key, validating the key isn't
// empty and won't produce a malformed path (leading or duplicate slashes)
// once joined to the CDN domain.
func (b *Builder) URL(key string) (string, error) {
	if b.domain == "" {
		return "", fmt.Errorf("cdn: domain is not configured")
	}
	trimmed := strings.TrimLeft(key, "/")
	if trimmed == "" {
		return "", fmt.Errorf("cdn: object key is empty")
	}
	if strings.Contains(trimmed, "//") {
		return "", fmt.Errorf("cdn: object key %q contains a duplicate slash", key)
	}

	full := b.domain + "/" + trimmed
	if _, err := url.ParseRequestURI(full); err != nil {
		return "", fmt.Errorf("cdn: generated URL is invalid: %w", err)
	}
	return full, nil
}
