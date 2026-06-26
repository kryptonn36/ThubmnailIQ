package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	useruc "github.com/thumbnailiq/thumbnailiq/internal/usecase/user"
)

type AuthHandler struct {
	uc *useruc.Usecase
}

func NewAuthHandler(uc *useruc.Usecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func authResponse(c *gin.Context, status int, res *useruc.AuthResult) {
	c.JSON(status, gin.H{
		"access_token":  res.AccessToken,
		"refresh_token": res.RefreshToken,
		"expires_in":    res.ExpiresIn,
		"user": gin.H{
			"id":         res.User.ID,
			"email":      res.User.Email,
			"full_name":  res.User.FullName,
			"avatar_url": res.User.AvatarURL,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.uc.Register(c.Request.Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		respondError(c, err)
		return
	}
	authResponse(c, http.StatusCreated, res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.uc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		respondError(c, err)
		return
	}
	authResponse(c, http.StatusOK, res)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.uc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		respondError(c, err)
		return
	}
	authResponse(c, http.StatusOK, res)
}
