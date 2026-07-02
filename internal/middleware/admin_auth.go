package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/pkg/jwt"
)

const ContextAdminID = "admin_id"

// AdminAuth is structurally identical to Auth, but is constructed with a
// separate *jwt.Service (its own secrets, see config.AdminJWT) so admin and
// customer tokens are cryptographically isolated from each other — a
// customer's access token can never pass this check, and vice versa.
func AdminAuth(jwtSvc *jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtSvc.ParseAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		c.Set(ContextAdminID, claims.UserID)
		c.Next()
	}
}

func AdminID(c *gin.Context) uuid.UUID {
	v, _ := c.Get(ContextAdminID)
	id, _ := v.(uuid.UUID)
	return id
}
