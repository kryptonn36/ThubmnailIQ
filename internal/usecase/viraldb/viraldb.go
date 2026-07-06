package viraldb

import (
	"context"
	"encoding/json"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/viraldb"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/youtube"
	"github.com/thumbnailiq/thumbnailiq/internal/scoring"
	"github.com/thumbnailiq/thumbnailiq/pkg/validator"
)

type Usecase struct {
	repo    viraldb.Repository
	fetcher youtube.Fetcher
	scorer  *scoring.Engine
}

func NewUsecase(repo viraldb.Repository, fetcher youtube.Fetcher) *Usecase {
	return &Usecase{repo: repo, fetcher: fetcher, scorer: scoring.NewEngine()}
}

func (u *Usecase) Search(ctx context.Context, f viraldb.SearchFilter) ([]*viraldb.Thumbnail, error) {
	f.Keyword = normalizeOptional(f.Keyword)

	if f.Keyword != nil && u.fetcher != nil {
		if err := u.ensureKeywordCached(ctx, *f.Keyword); err != nil {
			return nil, err
		}
	}
	return u.repo.Search(ctx, f)
}

func (u *Usecase) ensureKeywordCached(ctx context.Context, keyword string) error {
	cached, err := u.repo.Search(ctx, viraldb.SearchFilter{Keyword: &keyword, Limit: 1})
	if err != nil {
		return err
	}
	if len(cached) > 0 {
		return nil
	}

	results, err := u.fetcher.FetchCompetitors(ctx, keyword, 20)
	if err != nil {
		return err
	}
	avg := scoring.ComputeCompetitorAvg(cvResults(results), nil)
	for _, result := range results {
		score, hasFace, rawCV := u.scoreResult(result.CV, avg)
		if _, err := u.repo.Upsert(ctx, &viraldb.ThumbnailInput{
			VideoID:      result.VideoID,
			ChannelID:    result.ChannelID,
			ChannelName:  result.ChannelName,
			VideoTitle:   result.Title,
			ThumbnailURL: result.ThumbnailURL,
			Niche:        keyword,
			Tags:         []string{keyword},
			ViewCount:    result.ViewCount,
			Score:        score,
			HasFace:      hasFace,
			CVResults:    rawCV,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (u *Usecase) scoreResult(result *cv.Result, avg scoring.CompetitorAvg) (int, bool, []byte) {
	if result == nil {
		return 0, false, nil
	}
	sub := scoring.SubScores{
		Visibility: u.scorer.Visibility(result, avg),
		Contrast:   u.scorer.Contrast(result),
		Attention:  u.scorer.Attention(result),
		Mobile:     u.scorer.Mobile(result),
		Branding:   u.scorer.Branding(),
		Curiosity:  50,
	}
	raw, _ := json.Marshal(result)
	return u.scorer.FinalScore(sub), result.Face.FaceCount > 0, raw
}

func cvResults(results []youtube.CompetitorResult) []*cv.Result {
	out := make([]*cv.Result, 0, len(results))
	for _, result := range results {
		if result.CV != nil {
			out = append(out, result.CV)
		}
	}
	return out
}

func normalizeOptional(value *string) *string {
	if value == nil {
		return nil
	}
	normalized := validator.NormalizeKeyword(*value)
	if normalized == "" {
		return nil
	}
	return &normalized
}
