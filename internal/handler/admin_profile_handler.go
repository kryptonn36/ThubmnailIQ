package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thumbnailiq/thumbnailiq/internal/middleware"
	adminuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/admin"
)

type AdminProfileHandler struct {
	uc *adminuc.Usecase
}

func NewAdminProfileHandler(uc *adminuc.Usecase) *AdminProfileHandler {
	return &AdminProfileHandler{uc: uc}
}

func (h *AdminProfileHandler) Get(c *gin.Context) {
	profile, err := h.uc.GetProfile(c.Request.Context(), middleware.AdminID(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

func (h *AdminProfileHandler) ChangePassword(c *gin.Context) {
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.ChangeOwnPassword(c.Request.Context(), middleware.AdminID(c), req.CurrentPassword, req.NewPassword); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "password updated"})
}
