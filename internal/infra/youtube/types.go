package youtube

import (
	"context"

	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
)

type CompetitorResult struct {
	VideoID         string
	Title           string
	ChannelID       string
	ChannelName     string
	ThumbnailURL    string
	ViewCount       int64
	SubscriberCount int64
	CV              *cv.Result
}

type Fetcher interface {
	FetchCompetitors(ctx context.Context, keyword string, count int) ([]CompetitorResult, error)
}
