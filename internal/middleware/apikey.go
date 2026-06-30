package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/billing"
	"github.com/thumbnailiq/thumbnailiq/pkg/hash"
)

// APIKeyAuth authenticates requests coming from the browser extension (and
// any other API key consumers) without requiring a JWT browser session.
// The caller passes their raw API key as `Authorization: Bearer <key>`.
// The middleware hashes it and looks it up in the billing repository to
// resolve the workspace ID, then stores it in the request context so
// downstream handlers can access it.
func APIKeyAuth(repo billing.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		raw := strings.TrimPrefix(header, "Bearer ")
		key, err := repo.GetAPIKeyByHash(context.Background(), hash.SHA256Hex(raw))
		if err != nil {
			c.Error(fmt.Errorf("api key auth failed for %s %s: %w", c.Request.Method, c.Request.URL.Path, err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			return
		}
		c.Set("workspace_id", key.WorkspaceID)
		c.Next()
	}
}
