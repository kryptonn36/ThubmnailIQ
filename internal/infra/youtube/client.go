package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
)

// Client calls the real YouTube Data API v3 and analyzes each competitor
// thumbnail through the CV service. Used when YOUTUBE_API_KEY is configured.
type Client struct {
	apiKey   string
	cvClient *cv.Client
	http     *http.Client
}

func NewClient(apiKey string, cvClient *cv.Client) *Client {
	return &Client{apiKey: apiKey, cvClient: cvClient, http: &http.Client{Timeout: 15 * time.Second}}
}

type searchResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title      string `json:"title"`
			ChannelID  string `json:"channelId"`
			ChannelTitle string `json:"channelTitle"`
			Thumbnails struct {
				High struct {
					URL string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
}

type videoStatsResponse struct {
	Items []struct {
		Statistics struct {
			ViewCount string `json:"viewCount"`
		} `json:"statistics"`
	} `json:"items"`
}

func (c *Client) FetchCompetitors(ctx context.Context, keyword string, count int) ([]CompetitorResult, error) {
	searchURL := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/search?part=snippet&type=video&maxResults=%d&q=%s&key=%s",
		count, url.QueryEscape(keyword), c.apiKey,
	)
	var sr searchResponse
	if err := c.getJSON(ctx, searchURL, &sr); err != nil {
		return nil, fmt.Errorf("youtube search: %w", err)
	}

	results := make([]CompetitorResult, 0, len(sr.Items))
	for _, item := range sr.Items {
		viewCount := c.fetchViewCount(ctx, item.ID.VideoID)
		var cvResult *cv.Result
		if c.cvClient != nil {
			cvResult, _ = c.cvClient.Analyze(ctx, item.Snippet.Thumbnails.High.URL)
		}
		results = append(results, CompetitorResult{
			VideoID:      item.ID.VideoID,
			Title:        item.Snippet.Title,
			ChannelID:    item.Snippet.ChannelID,
			ChannelName:  item.Snippet.ChannelTitle,
			ThumbnailURL: item.Snippet.Thumbnails.High.URL,
			ViewCount:    viewCount,
			CV:           cvResult,
		})
	}
	return results, nil
}

func (c *Client) fetchViewCount(ctx context.Context, videoID string) int64 {
	statsURL := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/videos?part=statistics&id=%s&key=%s",
		videoID, c.apiKey,
	)
	var vr videoStatsResponse
	if err := c.getJSON(ctx, statsURL, &vr); err != nil || len(vr.Items) == 0 {
		return 0
	}
	v, _ := strconv.ParseInt(vr.Items[0].Statistics.ViewCount, 10, 64)
	return v
}

func (c *Client) getJSON(ctx context.Context, u string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("youtube api returned status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
