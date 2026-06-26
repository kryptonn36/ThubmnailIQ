package scoring

import (
	"sort"

	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
)

func impactRank(level string) int {
	switch level {
	case "high":
		return 0
	case "medium":
		return 1
	default:
		return 2
	}
}

// BuildSuggestions mirrors the catalog in blueprint section 3.4: a ranked
// list of improvement actions derived from threshold checks against the
// computed CV metrics and sub-scores.
func BuildSuggestions(own *cv.Result, sub SubScores, avg CompetitorAvg) []Suggestion {
	var out []Suggestion

	if own.Face.FaceCount == 0 {
		out = append(out, Suggestion{
			Type: "add_human_face", ImpactLevel: "high",
			Headline:        "Add a human face showing strong emotion",
			Explanation:     "Faces with visible emotion draw viewer attention within the first 300ms and consistently outperform faceless thumbnails in this niche.",
			EstimatedCTRMin: 0.5, EstimatedCTRMax: 1.2,
		})
	}
	if own.OCR.TextDetected && own.OCR.AvgTextHeightPct < 8 {
		out = append(out, Suggestion{
			Type: "increase_text_size", ImpactLevel: "high",
			Headline:        "Your text is too small for mobile viewers",
			Explanation:     "Text under roughly 8% of the thumbnail height becomes unreadable at mobile feed sizes, where most impressions happen.",
			EstimatedCTRMin: 0.3, EstimatedCTRMax: 0.9,
		})
	}
	if own.Colors.ContrastScore < 4.5 {
		out = append(out, Suggestion{
			Type: "improve_contrast", ImpactLevel: "high",
			Headline:        "Low contrast makes text and subject hard to read",
			Explanation:     "Your dominant-color contrast ratio is below the WCAG AA threshold of 4.5, which hurts legibility at thumbnail scale.",
			EstimatedCTRMin: 0.4, EstimatedCTRMax: 1.0,
		})
	}
	if own.Clutter.ClutterScore > 50 {
		out = append(out, Suggestion{
			Type: "reduce_clutter", ImpactLevel: "high",
			Headline:        "Too many elements compete for attention",
			Explanation:     "High edge density indicates a busy frame; viewers decide to click in under 300ms and clutter slows that decision.",
			EstimatedCTRMin: 0.4, EstimatedCTRMax: 1.1,
		})
	}
	if sub.Visibility < 50 {
		out = append(out, Suggestion{
			Type: "use_complementary_color", ImpactLevel: "medium",
			Headline:        "Your palette blends in with competitor averages",
			Explanation:     "A more saturated or distinct color palette than the competitor norm helps your thumbnail stand out in a crowded search results page.",
			EstimatedCTRMin: 0.2, EstimatedCTRMax: 0.6,
		})
	}
	if sub.Curiosity < 50 {
		out = append(out, Suggestion{
			Type: "add_curiosity_gap", ImpactLevel: "medium",
			Headline:        "No clear promise or payoff visible",
			Explanation:     "Thumbnails that hint at a reward, transformation, or unresolved question tend to convert better than purely descriptive ones.",
			EstimatedCTRMin: 0.3, EstimatedCTRMax: 0.8,
		})
	}
	if sub.Attention < 50 {
		out = append(out, Suggestion{
			Type: "add_directional_element", ImpactLevel: "medium",
			Headline:        "Add an arrow or visual cue pointing to the key subject",
			Explanation:     "Directional cues guide the eye to the most important part of the frame within the first glance.",
			EstimatedCTRMin: 0.2, EstimatedCTRMax: 0.5,
		})
	}
	if own.Colors.SaturationScore < avg.SaturationScore-10 {
		out = append(out, Suggestion{
			Type: "increase_saturation", ImpactLevel: "medium",
			Headline:        "Increase color vibrancy to stand out",
			Explanation:     "Your saturation is noticeably below the competitor average for this keyword, which can make the thumbnail look washed out next to ranking results.",
			EstimatedCTRMin: 0.1, EstimatedCTRMax: 0.4,
		})
	}
	out = append(out, Suggestion{
		Type: "add_channel_branding", ImpactLevel: "low",
		Headline:        "Consistent branding builds recognition",
		Explanation:     "Using consistent colors, fonts, or a watermark across thumbnails helps repeat viewers recognize your content in their feed.",
		EstimatedCTRMin: 0.1, EstimatedCTRMax: 0.3,
	})
	out = append(out, Suggestion{
		Type: "optimize_for_dark_mode", ImpactLevel: "low",
		Headline:        "Check thumbnail in YouTube dark mode",
		Explanation:     "Dark backgrounds can lose contrast against YouTube's dark theme; verify legibility in both light and dark contexts.",
		EstimatedCTRMin: 0.05, EstimatedCTRMax: 0.2,
	})

	sort.SliceStable(out, func(i, j int) bool {
		return impactRank(out[i].ImpactLevel) < impactRank(out[j].ImpactLevel)
	})
	return out
}
