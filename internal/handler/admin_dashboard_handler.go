package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	adminuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/admin"
)

type AdminDashboardHandler struct {
	uc *adminuc.Usecase
}

func NewAdminDashboardHandler(uc *adminuc.Usecase) *AdminDashboardHandler {
	return &AdminDashboardHandler{uc: uc}
}

func (h *AdminDashboardHandler) Get(c *gin.Context) {
	stats, err := h.uc.Dashboard(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, stats)
}
