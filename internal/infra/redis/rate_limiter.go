package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
}

func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{client: client}
}

// Allow implements a fixed-window limiter: limit requests per window per key.
func (r *RateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	cacheKey := fmt.Sprintf("ratelimit:%s", key)
	pipe := r.client.Pipeline()
	incr := pipe.Incr(ctx, cacheKey)
	pipe.Expire(ctx, cacheKey, window)
	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}
	return incr.Val() <= int64(limit), nil
}
