// admin-seed bootstraps the very first admin account. There is no
// self-service admin registration endpoint by design (see
// internal/usecase/admin), so this is the only way an admin_users row gets
// created. Safe to re-run: it's a no-op if the email already exists.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/thumbnailiq/thumbnailiq/internal/config"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/postgres"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
	"github.com/thumbnailiq/thumbnailiq/pkg/hash"
	"github.com/thumbnailiq/thumbnailiq/pkg/logger"
	"github.com/thumbnailiq/thumbnailiq/pkg/validator"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("loading config: %w", err))
	}
	log := logger.New(cfg.Server.Env)

	email := os.Getenv("ADMIN_SEED_EMAIL")
	password := os.Getenv("ADMIN_SEED_PASSWORD")
	fullName := os.Getenv("ADMIN_SEED_FULL_NAME")
	if fullName == "" {
		fullName = "Admin"
	}
	if email == "" || password == "" {
		log.Fatal().Msg("ADMIN_SEED_EMAIL and ADMIN_SEED_PASSWORD must be set")
	}
	if !validator.IsValidEmail(email) || !validator.IsValidPassword(password) {
		log.Fatal().Msg("ADMIN_SEED_EMAIL must be a valid email and ADMIN_SEED_PASSWORD must be at least 8 characters")
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("connecting to database")
	}
	defer pool.Close()

	adminRepo := postgres.NewAdminRepo(pool)

	if existing, err := adminRepo.GetAdminByEmail(ctx, email); err == nil && existing != nil {
		log.Info().Str("email", email).Msg("admin already exists, nothing to do")
		return
	} else if err != nil && !apperrors.Is(err, apperrors.ErrNotFound) {
		log.Fatal().Err(err).Msg("checking for existing admin")
	}

	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		log.Fatal().Err(err).Msg("hashing password")
	}

	created, err := adminRepo.CreateAdmin(ctx, email, passwordHash, fullName, "admin")
	if err != nil {
		log.Fatal().Err(err).Msg("creating admin")
	}

	log.Info().Str("email", created.Email).Str("id", created.ID.String()).Msg("admin account created")
}
