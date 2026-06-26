package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/thumbnailiq/thumbnailiq/internal/infra/youtube"
)

type CompetitorHandler struct {
	fetcher youtube.Fetcher
}

func NewCompetitorHandler(fetcher youtube.Fetcher) *CompetitorHandler {
	return &CompetitorHandler{fetcher: fetcher}
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
