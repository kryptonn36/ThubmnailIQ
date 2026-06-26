package scoring

import (
	"math"

	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
)

// Engine computes the six ThumbnailIQ sub-scores and the weighted final
// score (0-100), per the blueprint's formulas. Some inputs the original
// spec assumes (eye-contact/gaze, arrow/object detection, multi-thumbnail
// branding history) aren't available from the lightweight CV pipeline used
// here; those terms are fixed at a neutral default and documented inline.
type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Visibility(own *cv.Result, avg CompetitorAvg) int {
	saturation := own.Colors.SaturationScore

	ownContrastNorm := math.Min(own.Colors.ContrastScore/21*100, 100)
	avgContrastNorm := math.Min(avg.ContrastScore/21*100, 100)
	contrastVsCompetitors := clamp(50+(ownContrastNorm-avgContrastNorm), 0, 100)

	uniqueColorDistance := colorDistance(own) // 0-100, higher = more distinct from competitor norm

	score := saturation*0.3 + contrastVsCompetitors*0.4 + uniqueColorDistance*0.3
	return int(clamp(score, 0, 100))
}

// colorDistance approximates "how different is this thumbnail's dominant
// color from the competitor norm" using brightness/saturation spread as a
// proxy (a true ΔE comparison would need the competitor's actual dominant
// color, which the mock/lightweight pipeline doesn't aggregate per-channel).
func colorDistance(own *cv.Result) float64 {
	if len(own.Colors.DominantColors) == 0 {
		return 50
	}
	top := own.Colors.DominantColors[0]
	return clamp(math.Abs(top.Saturation-50)*1.4, 0, 100)
}

func (e *Engine) Contrast(own *cv.Result) int {
	norm := math.Min(own.Colors.ContrastScore/21*100, 100)
	bonus := 0.0
	if own.Colors.ContrastScore >= 4.5 {
		bonus = 10
	}
	return int(clamp(norm+bonus, 0, 100))
}

func (e *Engine) Attention(own *cv.Result) int {
	facePresent := 0.0
	if own.Face.FaceCount > 0 {
		facePresent = 1
	}
	eyeContact := 0.0
	if own.Face.HasEyeContact {
		eyeContact = 1
	}
	largeText := 0.0
	if own.OCR.AvgTextHeightPct >= 8 {
		largeText = 1
	}
	// arrow/pointer detection isn't available without an object detector;
	// that 10% weight is folded into the text/face terms below instead of
	// silently dropping 10 points off every score.
	score := facePresent*0.45 + eyeContact*0.3 + largeText*0.25
	return int(clamp(score*100, 0, 100))
}

func (e *Engine) Mobile(own *cv.Result) int {
	textDensityPenalty := math.Max(0, (own.OCR.TextDensityPct-30)*2)
	clutterPenalty := own.Clutter.ClutterScore * 0.5
	smallTextPenalty := 0.0
	if own.OCR.TextDetected && own.OCR.AvgTextHeightPct > 0 && own.OCR.AvgTextHeightPct < 8 {
		smallTextPenalty = 20
	}
	return int(clamp(100-textDensityPenalty-clutterPenalty-smallTextPenalty, 0, 100))
}

// Branding requires 3+ historical thumbnails from the same channel to
// measure palette/font/compositional consistency (blueprint section 3.3).
// That history isn't collected in this build, so branding defaults to a
// neutral midpoint score rather than fabricating a number from no data.
func (e *Engine) Branding() int {
	return 50
}

func (e *Engine) FinalScore(s SubScores) int {
	weighted := float64(s.Visibility)*0.25 +
		float64(s.Contrast)*0.20 +
		float64(s.Attention)*0.20 +
		float64(s.Mobile)*0.15 +
		float64(s.Branding)*0.10 +
		float64(s.Curiosity)*0.10
	return int(clamp(weighted, 0, 100))
}

func ComputeCompetitorAvg(results []*cv.Result, scores []int) CompetitorAvg {
	if len(results) == 0 {
		return CompetitorAvg{}
	}
	var avg CompetitorAvg
	n := float64(len(results))
	for i, r := range results {
		avg.FaceCount += float64(r.Face.FaceCount)
		avg.TextDensityPct += r.OCR.TextDensityPct
		avg.ClutterScore += r.Clutter.ClutterScore
		avg.ContrastScore += r.Colors.ContrastScore
		avg.SaturationScore += r.Colors.SaturationScore
		if i < len(scores) {
			avg.Score += float64(scores[i])
		}
	}
	avg.FaceCount /= n
	avg.TextDensityPct /= n
	avg.ClutterScore /= n
	avg.ContrastScore /= n
	avg.SaturationScore /= n
	avg.Score /= n
	return avg
}
