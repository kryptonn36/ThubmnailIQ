package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/thumbnailiq/thumbnailiq/internal/infra/gemini"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/youtube"
)

type CompetitorHandler struct {
	fetcher youtube.Fetcher
	gemini  *gemini.Client
}

func NewCompetitorHandler(fetcher youtube.Fetcher, geminiClient *gemini.Client) *CompetitorHandler {
	return &CompetitorHandler{fetcher: fetcher, gemini: geminiClient}
}

func (h *CompetitorHandler) ListForKeyword(c *gin.Context) {
	keyword := c.Param("keyword")
	count, err := strconv.Atoi(c.DefaultQuery("count", "20"))
	if err != nil || count <= 0 || count > 50 {
		count = 20
	}

	results, err := h.fetcher.FetchCompetitors(c.Request.Context(), keyword, count)
	if err != nil {
		respondError(c, err)
		return
	}

	out := make([]gin.H, 0, len(results))
	for i, r := range results {
		item := gin.H{
			"video_id": r.VideoID, "video_title": r.Title, "channel_name": r.ChannelName,
			"thumbnail_url": r.ThumbnailURL, "view_count": r.ViewCount, "rank_position": i + 1,
		}
		out = append(out, item)
	}
	c.JSON(http.StatusOK, gin.H{"competitors": out})
}

// VideoIdeas generates fresh video ideas for a keyword by studying what
// top-performing competitors are already doing and either asking Gemini to
// suggest differentiated angles, or falling back to a template-based
// heuristic when Gemini is unavailable.
func (h *CompetitorHandler) VideoIdeas(c *gin.Context) {
	keyword := c.Param("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is required"})
		return
	}

	count := 15
	results, err := h.fetcher.FetchCompetitors(c.Request.Context(), keyword, count)
	if err != nil {
		results = nil
	}

	titles := make([]string, 0, len(results))
	for _, r := range results {
		titles = append(titles, r.Title)
	}

	ideas, err := h.gemini.GenerateVideoIdeas(c.Request.Context(), keyword, titles)
	if err != nil || ideas == nil {
		if err != nil {
			c.Error(err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, ideas)
}
