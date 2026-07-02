// Package health implements admin.HealthChecker by wrapping the same
// clients already constructed in cmd/api/main.go (pool, redis, cv) — no new
// connections, just a thin liveness probe over each one for the admin
// dashboard's system health widget.
package health

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"

	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
)

type Checker struct {
	pool     *pgxpool.Pool
	redis    *goredis.Client
	cvClient *cv.Client
}

func NewChecker(pool *pgxpool.Pool, redisClient *goredis.Client, cvClient *cv.Client) *Checker {
	return &Checker{pool: pool, redis: redisClient, cvClient: cvClient}
}

func (c *Checker) CheckDatabase(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

func (c *Checker) CheckRedis(ctx context.Context) error {
	return c.redis.Ping(ctx).Err()
}

func (c *Checker) CheckCVService(ctx context.Context) error {
	return c.cvClient.Health(ctx)
}
