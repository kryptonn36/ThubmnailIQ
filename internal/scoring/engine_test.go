package scoring

import (
	"testing"

	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
)

func TestEngineScoresExpectedRanges(t *testing.T) {
	engine := NewEngine()
	result := sampleCVResult()

	if got := engine.Contrast(result); got != 31 {
		t.Fatalf("Contrast() = %d, want 31", got)
	}
	if got := engine.Attention(result); got != 100 {
		t.Fatalf("Attention() = %d, want 100", got)
	}
	if got := engine.Mobile(result); got != 70 {
		t.Fatalf("Mobile() = %d, want 70", got)
	}
	if got := engine.Branding(); got != 50 {
		t.Fatalf("Branding() = %d, want neutral 50", got)
	}
}

func TestFinalScoreUsesWeights(t *testing.T) {
	engine := NewEngine()
	got := engine.FinalScore(SubScores{
		Visibility: 80,
		Contrast:   70,
		Attention:  60,
		Mobile:     50,
		Branding:   40,
		Curiosity:  30,
	})
	if got != 60 {
		t.Fatalf("FinalScore() = %d, want 60", got)
	}
}

func TestComputeCompetitorAvg(t *testing.T) {
	one := sampleCVResult()
	two := sampleCVResult()
	two.Face.FaceCount = 0
	two.OCR.TextDensityPct = 40
	two.Clutter.ClutterScore = 20
	two.Colors.ContrastScore = 6
	two.Colors.SaturationScore = 80

	avg := ComputeCompetitorAvg([]*cv.Result{one, two}, []int{70, 90})
	if avg.FaceCount != 0.5 {
		t.Fatalf("FaceCount avg = %.2f, want 0.5", avg.FaceCount)
	}
	if avg.TextDensityPct != 35 {
		t.Fatalf("TextDensityPct avg = %.2f, want 35", avg.TextDensityPct)
	}
	if avg.ClutterScore != 40 {
		t.Fatalf("ClutterScore avg = %.2f, want 40", avg.ClutterScore)
	}
	if avg.ContrastScore != 5.25 {
		t.Fatalf("ContrastScore avg = %.2f, want 5.25", avg.ContrastScore)
	}
	if avg.SaturationScore != 70 {
		t.Fatalf("SaturationScore avg = %.2f, want 70", avg.SaturationScore)
	}
	if avg.Score != 80 {
		t.Fatalf("Score avg = %.2f, want 80", avg.Score)
	}
}

func TestBuildSuggestionsRanksHighImpactFirst(t *testing.T) {
	result := sampleCVResult()
	result.Face.FaceCount = 0
	result.OCR.TextDetected = true
	result.OCR.AvgTextHeightPct = 4
	result.Colors.ContrastScore = 3
	result.Clutter.ClutterScore = 70

	suggestions := BuildSuggestions(result, SubScores{Visibility: 40, Attention: 40, Curiosity: 40}, CompetitorAvg{
		SaturationScore: 80,
	})
	if len(suggestions) == 0 {
		t.Fatal("expected suggestions")
	}
	if suggestions[0].ImpactLevel != "high" {
		t.Fatalf("first suggestion impact = %q, want high", suggestions[0].ImpactLevel)
	}

	seenBranding := false
	for _, suggestion := range suggestions {
		if suggestion.Type == "add_channel_branding" {
			seenBranding = true
		}
	}
	if !seenBranding {
		t.Fatal("expected baseline channel branding suggestion")
	}
}

func sampleCVResult() *cv.Result {
	result := &cv.Result{}
	result.Face.FaceCount = 1
	result.Face.HasEyeContact = true
	result.OCR.TextDetected = true
	result.OCR.TextDensityPct = 30
	result.OCR.AvgTextHeightPct = 10
	result.Clutter.ClutterScore = 60
	result.Colors.ContrastScore = 4.5
	result.Colors.SaturationScore = 60
	result.Colors.DominantColors = []cv.DominantColor{{Hex: "#ff0000", Saturation: 90}}
	return result
}
