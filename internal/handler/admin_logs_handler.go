package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/admin"
	adminuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/admin"
	"github.com/thumbnailiq/thumbnailiq/pkg/pagination"
)

type AdminLogsHandler struct {
	uc *adminuc.Usecase
}

func NewAdminLogsHandler(uc *adminuc.Usecase) *AdminLogsHandler {
	return &AdminLogsHandler{uc: uc}
}

func (h *AdminLogsHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	displayPage, displayPerPage := pagination.Normalize(page, perPage)

	logs, total, err := h.uc.ListAuditLogs(c.Request.Context(), admin.Pagination{Page: page, PerPage: perPage})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": logs, "page": displayPage, "per_page": displayPerPage, "total": total})
}
