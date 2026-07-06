package viraldb

import (
	"context"
	"reflect"
	"testing"

	domainviraldb "github.com/thumbnailiq/thumbnailiq/internal/domain/viraldb"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/youtube"
)

func TestSearchForwardsFilterToRepository(t *testing.T) {
	niche := "gaming"
	minScore := 80
	hasFace := true
	filter := domainviraldb.SearchFilter{
		Niche:    &niche,
		MinScore: &minScore,
		HasFace:  &hasFace,
		Limit:    20,
		Offset:   40,
	}
	repo := &fakeViralRepo{
		results: []*domainviraldb.Thumbnail{{ID: "thumb_1", Score: 90}},
	}
	uc := NewUsecase(repo, nil)

	results, err := uc.Search(context.Background(), filter)
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if !reflect.DeepEqual(repo.lastFilter, filter) {
		t.Fatalf("repository filter = %+v, want %+v", repo.lastFilter, filter)
	}
	if len(results) != 1 || results[0].ID != "thumb_1" {
		t.Fatalf("unexpected search results: %+v", results)
	}
}

func TestSearchWithKeywordIngestsYouTubeResultsWhenCacheIsEmpty(t *testing.T) {
	keyword := "Gaming Thumbnails"
	repo := &fakeViralRepo{}
	fetcher := &fakeFetcher{
		results: []youtube.CompetitorResult{
			{
				VideoID:      "video_1",
				Title:        "Gaming Thumbnail Breakdown",
				ChannelID:    "UC1",
				ChannelName:  "Creator Lab",
				ThumbnailURL: "https://example.com/thumb.jpg",
				ViewCount:    100000,
				CV:           sampleCV(),
			},
		},
	}
	uc := NewUsecase(repo, fetcher)

	results, err := uc.Search(context.Background(), domainviraldb.SearchFilter{Keyword: &keyword, Limit: 50})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if fetcher.calls != 1 {
		t.Fatalf("fetcher calls = %d, want 1", fetcher.calls)
	}
	if len(repo.upserts) != 1 {
		t.Fatalf("upserts = %d, want 1", len(repo.upserts))
	}
	if repo.upserts[0].Tags[0] != "gaming thumbnails" {
		t.Fatalf("stored keyword tag = %q, want normalized keyword", repo.upserts[0].Tags[0])
	}
	if repo.upserts[0].Score == 0 {
		t.Fatal("expected ingested thumbnail to be scored")
	}
	if !repo.upserts[0].HasFace {
		t.Fatal("expected has_face to be derived from CV")
	}
	if len(results) != 1 || results[0].VideoID != "video_1" {
		t.Fatalf("unexpected results: %+v", results)
	}
}

type fakeViralRepo struct {
	lastFilter domainviraldb.SearchFilter
	results    []*domainviraldb.Thumbnail
	upserts    []*domainviraldb.ThumbnailInput
}

func (r *fakeViralRepo) Search(_ context.Context, f domainviraldb.SearchFilter) ([]*domainviraldb.Thumbnail, error) {
	r.lastFilter = f
	if f.Keyword != nil {
		var out []*domainviraldb.Thumbnail
		for _, result := range r.results {
			if result.Niche == *f.Keyword {
				out = append(out, result)
			}
		}
		return out, nil
	}
	return r.results, nil
}

func (r *fakeViralRepo) Upsert(_ context.Context, input *domainviraldb.ThumbnailInput) (*domainviraldb.Thumbnail, error) {
	r.upserts = append(r.upserts, input)
	thumb := &domainviraldb.Thumbnail{
		ID:           input.VideoID,
		VideoID:      input.VideoID,
		ChannelName:  input.ChannelName,
		VideoTitle:   input.VideoTitle,
		ThumbnailURL: input.ThumbnailURL,
		Niche:        input.Niche,
		Score:        input.Score,
		ViewCount:    input.ViewCount,
		HasFace:      input.HasFace,
	}
	r.results = append(r.results, thumb)
	return thumb, nil
}

type fakeFetcher struct {
	calls   int
	results []youtube.CompetitorResult
}

func (f *fakeFetcher) FetchCompetitors(_ context.Context, _ string, _ int) ([]youtube.CompetitorResult, error) {
	f.calls++
	return f.results, nil
}

func sampleCV() *cv.Result {
	result := &cv.Result{}
	result.Face.FaceCount = 1
	result.Face.HasEyeContact = true
	result.OCR.TextDetected = true
	result.OCR.TextDensityPct = 20
	result.OCR.AvgTextHeightPct = 12
	result.Colors.ContrastScore = 5
	result.Colors.SaturationScore = 70
	result.Colors.DominantColors = []cv.DominantColor{{Hex: "#ff0000", Saturation: 85}}
	result.Clutter.ClutterScore = 20
	return result
}
