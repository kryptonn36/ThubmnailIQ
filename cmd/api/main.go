package main

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"

	"github.com/thumbnailiq/thumbnailiq/internal/config"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/payment"
	"github.com/thumbnailiq/thumbnailiq/internal/handler"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/payment/razorpay"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/payment/stripe"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/postgres"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/s3"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/youtube"
	"github.com/thumbnailiq/thumbnailiq/internal/server"
	analysisuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/analysis"
	billinguc "github.com/thumbnailiq/thumbnailiq/internal/usecase/billing"
	trackinguc "github.com/thumbnailiq/thumbnailiq/internal/usecase/tracking"
	useruc "github.com/thumbnailiq/thumbnailiq/internal/usecase/user"
	viraldbuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/viraldb"
	workspaceuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/workspace"
	"github.com/thumbnailiq/thumbnailiq/pkg/jwt"
	"github.com/thumbnailiq/thumbnailiq/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("loading config: %w", err))
	}
	log := logger.New(cfg.Server.Env)

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("connecting to database")
	}
	defer pool.Close()

	storage, err := s3.NewStorage(ctx, s3.Config{
		Endpoint: cfg.S3.Endpoint, Region: cfg.S3.Region,
		AccessKeyID: cfg.S3.AccessKeyID, SecretAccessKey: cfg.S3.SecretAccessKey,
		Bucket: cfg.S3.UploadBucket, PublicBaseURL: cfg.S3.PublicBaseURL, UsePathStyle: cfg.S3.UsePathStyle,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("configuring s3 storage")
	}
	if err := storage.EnsurePublicReadBucket(ctx); err != nil {
		log.Warn().Err(err).Msg("could not ensure upload bucket exists/public (continuing)")
	}

	cvClient := cv.NewClient(cfg.CVService.URL)

	var ytFetcher youtube.Fetcher
	if cfg.YouTube.APIKey != "" {
		ytFetcher = youtube.NewClient(cfg.YouTube.APIKey, cvClient)
		log.Info().Msg("using real YouTube Data API for competitor data")
	} else {
		ytFetcher = youtube.NewMockClient()
		log.Info().Msg("YOUTUBE_API_KEY not set, using synthetic competitor data")
	}

	var gateway payment.Gateway
	switch cfg.Payment.Provider {
	case "stripe":
		gateway = stripe.NewClient(cfg.Stripe.SecretKey)
	default:
		gateway = razorpay.NewClient(cfg.Razorpay.KeyID, cfg.Razorpay.KeySecret)
	}
	log.Info().Str("provider", cfg.Payment.Provider).Msg("payment gateway configured")

	queueClient := asynq.NewClient(asynq.RedisClientOpt{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password, DB: cfg.Redis.DB})
	defer queueClient.Close()

	userRepo := postgres.NewUserRepo(pool)
	workspaceRepo := postgres.NewWorkspaceRepo(pool)
	analysisRepo := postgres.NewAnalysisRepo(pool)
	competitorRepo := postgres.NewCompetitorRepo(pool)
	billingRepo := postgres.NewBillingRepo(pool)
	viralDBRepo := postgres.NewViralDBRepo(pool)

	jwtSvc := jwt.NewService(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	userUC := useruc.NewUsecase(userRepo, workspaceRepo, jwtSvc)
	workspaceUC := workspaceuc.NewUsecase(workspaceRepo, userRepo)
	analysisUC := analysisuc.NewUsecase(analysisRepo, workspaceRepo, storage, cvClient, queueClient)
	billingUC := billinguc.NewUsecase(billingRepo, workspaceRepo, gateway, cfg.Payment.Currency)
	trackingUC := trackinguc.NewUsecase(competitorRepo)
	viralDBUC := viraldbuc.NewUsecase(viralDBRepo)

	handlers := &server.Handlers{
		Auth:       handler.NewAuthHandler(userUC),
		Workspace:  handler.NewWorkspaceHandler(workspaceUC),
		Analysis:   handler.NewAnalysisHandler(analysisUC, competitorRepo, workspaceUC),
		Competitor: handler.NewCompetitorHandler(ytFetcher),
		Tracking:   handler.NewTrackingHandler(trackingUC, workspaceUC),
		Billing:    handler.NewBillingHandler(billingUC, workspaceUC),
		ViralDB:    handler.NewViralDBHandler(viralDBUC),
	}

	router := server.NewRouter(handlers, jwtSvc, log)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info().Str("addr", addr).Msg("starting ThumbnailIQ API server")
	if err := router.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}
