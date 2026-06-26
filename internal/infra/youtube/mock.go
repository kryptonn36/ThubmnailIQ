package youtube

import (
	"context"
	"fmt"
	"hash/fnv"
	"math/rand"

	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
)

// MockClient generates deterministic synthetic competitor data without
// calling any external API. Used when YOUTUBE_API_KEY is not configured,
// so the full analysis flow (including competitor comparison) is testable
// without real credentials or network access.
type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

var emotions = []string{"happy", "neutral", "surprised"}
var channelNames = []string{"CreatorHub", "TopTier Media", "NextLevel Channel", "Viral Studio", "PrimeContent", "StreamCraft"}

func (m *MockClient) FetchCompetitors(ctx context.Context, keyword string, count int) ([]CompetitorResult, error) {
	seed := hashSeed(keyword)
	rng := rand.New(rand.NewSource(seed))

	results := make([]CompetitorResult, 0, count)
	for i := 0; i < count; i++ {
		faceCount := rng.Intn(2)
		textDensity := 10 + rng.Float64()*35
		clutter := 15 + rng.Float64()*55
		contrast := 1.5 + rng.Float64()*12
		saturation := 30 + rng.Float64()*60

		var result cv.Result
		result.Face.FaceCount = faceCount
		if faceCount > 0 {
			result.Face.PrimaryEmotion = emotions[rng.Intn(len(emotions))]
		}
		result.OCR.TextDetected = textDensity > 5
		result.OCR.TextDensityPct = textDensity
		result.OCR.AvgTextHeightPct = 6 + rng.Float64()*10
		result.Colors.ContrastScore = contrast
		result.Colors.SaturationScore = saturation
		result.Colors.BrightnessScore = 60 + rng.Float64()*120
		result.Clutter.ClutterScore = clutter
		result.VisualComplexity = (clutter + textDensity) / 2

		videoID := fmt.Sprintf("mock%s%d", hashShort(keyword), i)
		results = append(results, CompetitorResult{
			VideoID:      videoID,
			Title:        fmt.Sprintf("%s — Top Result #%d", titleCase(keyword), i+1),
			ChannelID:    fmt.Sprintf("UCmock%d", i),
			ChannelName:  channelNames[i%len(channelNames)],
			ThumbnailURL: fmt.Sprintf("https://placehold.co/480x270?text=Competitor+%d", i+1),
			ViewCount:    int64(5000 + rng.Intn(2000000)),
			CV:           &result,
		})
	}
	return results, nil
}

func hashSeed(s string) int64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())
}

func hashShort(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum32())[:6]
}

func titleCase(s string) string {
	if s == "" {
		return s
	}
	b := []byte(s)
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 32
	}
	return string(b)
}
