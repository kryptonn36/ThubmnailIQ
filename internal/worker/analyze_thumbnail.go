package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/analysis"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/competitor"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/gemini"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/youtube"
	"github.com/thumbnailiq/thumbnailiq/internal/scoring"
)

const CompetitorCount = 20

type AnalysisHandler struct {
	analyses    analysis.Repository
	competitors competitor.Repository
	cvClient    *cv.Client
	youtube     youtube.Fetcher
	gemini      *gemini.Client
	engine      *scoring.Engine
	log         zerolog.Logger
}

func NewAnalysisHandler(
	analyses analysis.Repository,
	competitors competitor.Repository,
	cvClient *cv.Client,
	yt youtube.Fetcher,
	geminiClient *gemini.Client,
	log zerolog.Logger,
) *AnalysisHandler {
	return &AnalysisHandler{
		analyses: analyses, competitors: competitors, cvClient: cvClient,
		youtube: yt, gemini: geminiClient, engine: scoring.NewEngine(), log: log,
	}
}

func (h *AnalysisHandler) HandleAnalyzeThumbnail(ctx context.Context, t *asynq.Task) error {
	var payload AnalyzeThumbnailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	a, err := h.analyses.GetByID(ctx, payload.AnalysisID)
	if err != nil {
		return fmt.Errorf("load analysis: %w", err)
	}

	if err := h.analyses.UpdateStatus(ctx, a.ID, analysis.StatusProcessing, ""); err != nil {
		return err
	}

	result, err := h.run(ctx, a)
	if err != nil {
		h.log.Error().Err(err).Str("analysis_id", a.ID.String()).Msg("analysis failed")
		_ = h.analyses.UpdateStatus(ctx, a.ID, analysis.StatusFailed, err.Error())
		return err
	}

	return h.analyses.UpdateResults(ctx, result)
}

func (h *AnalysisHandler) run(ctx context.Context, a *analysis.Analysis) (*analysis.Analysis, error) {
	ownCV, err := h.cvClient.Analyze(ctx, a.ThumbnailURL)
	if err != nil {
		return nil, fmt.Errorf("cv analyze own thumbnail: %w", err)
	}

	competitors, err := h.youtube.FetchCompetitors(ctx, a.Keyword, CompetitorCount)
	if err != nil {
		h.log.Warn().Err(err).Msg("competitor fetch failed, continuing with no competitor context")
		competitors = nil
	}

	var competitorResults []*cv.Result
	var competitorScores []int
	for i, c := range competitors {
		if c.CV == nil {
			continue
		}
		sub := scoring.SubScores{
			Visibility: h.engine.Visibility(c.CV, scoring.CompetitorAvg{}),
			Contrast:   h.engine.Contrast(c.CV),
			Attention:  h.engine.Attention(c.CV),
			Mobile:     h.engine.Mobile(c.CV),
			Branding:   h.engine.Branding(),
		}
		score := h.engine.FinalScore(sub)
		competitorResults = append(competitorResults, c.CV)
		competitorScores = append(competitorScores, score)

		cvJSON, _ := json.Marshal(c.CV)
		rank := i + 1
		_, _ = h.competitors.CreateSnapshot(ctx, &competitor.Snapshot{
			AnalysisID: &a.ID, VideoID: c.VideoID, VideoTitle: c.Title,
			ChannelID: c.ChannelID, ChannelName: c.ChannelName, ThumbnailURL: c.ThumbnailURL,
			ViewCount: c.ViewCount, RankPosition: rank, Keyword: a.Keyword,
			CVResults: cvJSON, Score: &score,
		})
	}

	avg := scoring.ComputeCompetitorAvg(competitorResults, competitorScores)

	curiosity, err := h.gemini.ScoreCuriosity(ctx, ownCV.OCR.TextStrings, ownCV.Face.PrimaryEmotion)
	if err != nil {
		curiosity = &gemini.CuriosityResult{Score: 50, Reasoning: "scoring unavailable"}
	}

	sub := scoring.SubScores{
		Visibility: h.engine.Visibility(ownCV, avg),
		Contrast:   h.engine.Contrast(ownCV),
		Attention:  h.engine.Attention(ownCV),
		Mobile:     h.engine.Mobile(ownCV),
		Branding:   h.engine.Branding(),
		Curiosity:  curiosity.Score,
	}
	final := h.engine.FinalScore(sub)
	suggestions := scoring.BuildSuggestions(ownCV, sub, avg)

	rank := rankAmongCompetitors(final, competitorScores)

	cvJSON, _ := json.Marshal(ownCV)
	avgJSON, _ := json.Marshal(avg)
	suggestionsJSON, _ := json.Marshal(suggestions)

	a.Score = &final
	a.VisibilityScore = &sub.Visibility
	a.ContrastScore = &sub.Contrast
	a.AttentionScore = &sub.Attention
	a.MobileScore = &sub.Mobile
	a.BrandingScore = &sub.Branding
	a.CuriosityScore = &sub.Curiosity
	a.CVResults = cvJSON
	a.CompetitorAvg = avgJSON
	a.Suggestions = suggestionsJSON
	a.CompetitorCount = len(competitors)
	a.RankInCompetitors = &rank

	return a, nil
}

func rankAmongCompetitors(score int, competitorScores []int) int {
	rank := 1
	for _, s := range competitorScores {
		if s > score {
			rank++
		}
	}
	return rank
}
