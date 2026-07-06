package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/viraldb"
	viraldbuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/viraldb"
)

type ViralDBHandler struct {
	uc *viraldbuc.Usecase
}

func NewViralDBHandler(uc *viraldbuc.Usecase) *ViralDBHandler {
	return &ViralDBHandler{uc: uc}
}

func (h *ViralDBHandler) Search(c *gin.Context) {
	f := viraldb.SearchFilter{Limit: 50}
	if keyword := c.Query("keyword"); keyword != "" {
		f.Keyword = &keyword
	}
	if niche := c.Query("niche"); niche != "" {
		f.Niche = &niche
	}
	if minScoreStr := c.Query("min_score"); minScoreStr != "" {
		if v, err := strconv.Atoi(minScoreStr); err == nil {
			f.MinScore = &v
		}
	}
	if hasFaceStr := c.Query("has_face"); hasFaceStr != "" {
		if v, err := strconv.ParseBool(hasFaceStr); err == nil {
			f.HasFace = &v
		}
	}

	results, err := h.uc.Search(c.Request.Context(), f)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}
