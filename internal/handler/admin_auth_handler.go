package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	adminuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/admin"
)

type AdminAuthHandler struct {
	uc *adminuc.Usecase
}

func NewAdminAuthHandler(uc *adminuc.Usecase) *AdminAuthHandler {
	return &AdminAuthHandler{uc: uc}
}

type adminLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type adminRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func adminAuthResponse(c *gin.Context, status int, res *adminuc.AuthResult) {
	c.JSON(status, gin.H{
		"access_token":  res.AccessToken,
		"refresh_token": res.RefreshToken,
		"expires_in":    res.ExpiresIn,
		"admin": gin.H{
			"id":        res.Admin.ID,
			"email":     res.Admin.Email,
			"full_name": res.Admin.FullName,
			"role":      res.Admin.Role,
		},
	})
}

func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.uc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		respondError(c, err)
		return
	}
	adminAuthResponse(c, http.StatusOK, res)
}

func (h *AdminAuthHandler) Refresh(c *gin.Context) {
	var req adminRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.uc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		respondError(c, err)
		return
	}
	adminAuthResponse(c, http.StatusOK, res)
}
