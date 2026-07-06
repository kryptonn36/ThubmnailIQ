package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hibiken/asynq"
	goredis "github.com/redis/go-redis/v9"

	"github.com/thumbnailiq/thumbnailiq/internal/config"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/payment"
	"github.com/thumbnailiq/thumbnailiq/internal/handler"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/cdn"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/gemini"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/health"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/payment/razorpay"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/payment/stripe"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/postgres"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/redis"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/s3"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/youtube"
	"github.com/thumbnailiq/thumbnailiq/internal/server"
	adminuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/admin"
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
	if cfg.Server.Env == "production"{
		if err = productionCfgCheck(cfg); err != nil {
			panic(fmt.Errorf("Production env: %v", err))
		}
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
		Bucket: cfg.S3.UploadBucket, UsePathStyle: cfg.S3.UsePathStyle,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("configuring s3 storage")
	}
	// S3_PUBLIC_READ should only be true for local/MinIO dev where nothing
	// else fronts the bucket. In production, once a CloudFront distribution
	// with Origin Access Control sits in front of it, set this to false —
	// the bucket must stay private and grant access to CloudFront's OAC
	// principal only (done in AWS, not here).
	if cfg.S3.PublicRead {
		if err := storage.EnsurePublicReadBucket(ctx); err != nil {
			log.Warn().Err(err).Msg("could not ensure upload bucket exists/public (continuing)")
		}
	} else if err := storage.EnsureBucket(ctx); err != nil {
		log.Warn().Err(err).Msg("could not ensure upload bucket exists (continuing)")
	}

	cdnBuilder := cdn.NewBuilder(cfg.CDN.Domain)

	cvClient := cv.NewClient(cfg.CVService.URL)

	var ytFetcher youtube.Fetcher
	if cfg.YouTube.APIKey != "" {
		ytFetcher = youtube.NewClient(cfg.YouTube.APIKey, cvClient)
		log.Info().Msg("using real YouTube Data API for competitor data")
	} else {
		ytFetcher = youtube.NewMockClient()
		log.Info().Msg("YOUTUBE_API_KEY not set, using synthetic competitor data")
	}

	geminiClient := gemini.NewClient(cfg.Gemini.APIKey, cfg.Gemini.Model)

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

	redisClient := goredis.NewClient(&goredis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password, DB: cfg.Redis.DB})
	defer redisClient.Close()
	rateLimiter := redis.NewRateLimiter(redisClient)
	healthChecker := health.NewChecker(pool, redisClient, cvClient)

	userRepo := postgres.NewUserRepo(pool)
	workspaceRepo := postgres.NewWorkspaceRepo(pool)
	analysisRepo := postgres.NewAnalysisRepo(pool)
	competitorRepo := postgres.NewCompetitorRepo(pool)
	billingRepo := postgres.NewBillingRepo(pool)
	viralDBRepo := postgres.NewViralDBRepo(pool)
	adminRepo := postgres.NewAdminRepo(pool)

	jwtSvc := jwt.NewService(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	// A distinct *jwt.Service instance (own secrets) so admin and customer
	// tokens can never be parsed as one another.
	adminJWTSvc := jwt.NewService(cfg.AdminJWT.AccessSecret, cfg.AdminJWT.RefreshSecret, cfg.AdminJWT.AccessTTL, cfg.AdminJWT.RefreshTTL)

	userUC := useruc.NewUsecase(userRepo, workspaceRepo, jwtSvc)
	workspaceUC := workspaceuc.NewUsecase(workspaceRepo, userRepo)
	analysisUC := analysisuc.NewUsecase(analysisRepo, workspaceRepo, storage, cdnBuilder, cvClient, queueClient)
	billingUC := billinguc.NewUsecase(billingRepo, workspaceRepo, gateway, cfg.Payment.Currency)
	trackingUC := trackinguc.NewUsecase(competitorRepo)
	viralDBUC := viraldbuc.NewUsecase(viralDBRepo, ytFetcher)
	adminUC := adminuc.NewUsecase(adminRepo, adminJWTSvc, healthChecker)

	handlers := &server.Handlers{
		Auth:           handler.NewAuthHandler(userUC),
		Workspace:      handler.NewWorkspaceHandler(workspaceUC),
		Analysis:       handler.NewAnalysisHandler(analysisUC, competitorRepo, workspaceUC, cdnBuilder),
		Competitor:     handler.NewCompetitorHandler(ytFetcher, geminiClient),
		Tracking:       handler.NewTrackingHandler(trackingUC, workspaceUC),
		Billing:        handler.NewBillingHandler(billingUC, workspaceUC),
		ViralDB:        handler.NewViralDBHandler(viralDBUC),
		QuickAnalyze:   handler.NewQuickAnalyzeHandler(cvClient),
		AdminAuth:      handler.NewAdminAuthHandler(adminUC),
		AdminDashboard: handler.NewAdminDashboardHandler(adminUC),
		AdminUsers:     handler.NewAdminUsersHandler(adminUC),
		AdminUploads:   handler.NewAdminUploadsHandler(adminUC, cdnBuilder),
		AdminAnalytics: handler.NewAdminAnalyticsHandler(adminUC),
		AdminSettings:  handler.NewAdminSettingsHandler(adminUC),
		AdminLogs:      handler.NewAdminLogsHandler(adminUC),
		AdminProfile:   handler.NewAdminProfileHandler(adminUC),
	}

	router := server.NewRouter(handlers, jwtSvc, adminJWTSvc, billingRepo, rateLimiter, log)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info().Str("addr", addr).Msg("starting ThumbnailIQ API server")
	if err := router.Run(addr); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}


func productionCfgCheck(cfg *config.Config) error{
	if strings.Contains(cfg.Database.URL, "sslmode=disable"){
		return errors.New("database SSL must be enabled in production")
	}
	return nil 
}