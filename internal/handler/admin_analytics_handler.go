package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	adminuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/admin"
)

type AdminAnalyticsHandler struct {
	uc *adminuc.Usecase
}

func NewAdminAnalyticsHandler(uc *adminuc.Usecase) *AdminAnalyticsHandler {
	return &AdminAnalyticsHandler{uc: uc}
}

func (h *AdminAnalyticsHandler) Get(c *gin.Context) {
	stats, err := h.uc.GetAnalytics(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, stats)
}
