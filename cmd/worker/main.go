package main

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"

	"github.com/thumbnailiq/thumbnailiq/internal/config"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/gemini"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/postgres"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/youtube"
	"github.com/thumbnailiq/thumbnailiq/internal/worker"
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
	if cfg.Gemini.APIKey == "" {
		log.Info().Msg("GEMINI_API_KEY not set, using heuristic curiosity scoring")
	}

	analysisRepo := postgres.NewAnalysisRepo(pool)
	competitorRepo := postgres.NewCompetitorRepo(pool)

	analysisHandler := worker.NewAnalysisHandler(analysisRepo, competitorRepo, cvClient, ytFetcher, geminiClient, log)

	mux := asynq.NewServeMux()
	mux.HandleFunc(worker.TypeAnalyzeThumbnail, analysisHandler.HandleAnalyzeThumbnail)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password, DB: cfg.Redis.DB},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	log.Info().Msg("starting ThumbnailIQ worker")
	if err := srv.Run(mux); err != nil {
		log.Fatal().Err(err).Msg("worker stopped")
	}
}
