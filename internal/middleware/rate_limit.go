package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/thumbnailiq/thumbnailiq/internal/infra/redis"
)

func RateLimit(limiter *redis.RateLimiter, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if uid := UserID(c); uid.String() != "00000000-0000-0000-0000-000000000000" {
			key = uid.String()
		}
		allowed, err := limiter.Allow(c.Request.Context(), key+":"+c.FullPath(), limit, window)
		if err != nil {
			c.Next()
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
