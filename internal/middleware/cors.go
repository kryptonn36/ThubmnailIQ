package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS echoes the request Origin back only when it is on the configured
// allowlist. It never responds with a wildcard `*`, so an arbitrary third-party
// site can't script authenticated cross-origin calls against the API.
func CORS(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		if o != "" {
			allowed[o] = true
		}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && allowed[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-API-Key")
		}
		if c.Request.Method == http.MethodOptions {
			// Only short-circuit as a successful preflight for allowed origins;
			// a disallowed cross-origin preflight gets 403 so the browser blocks it.
			if origin != "" && !allowed[origin] {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
