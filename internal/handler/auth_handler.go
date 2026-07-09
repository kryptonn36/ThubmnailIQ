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
	usr, err := h.uc.Register(c.Request.Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		respondError(c, err)
		return
	}
	// No tokens yet — the account must verify its email first. The client
	// should send the user to the verification screen with this email.
	c.JSON(http.StatusCreated, gin.H{
		"requires_verification": true,
		"email":                 usr.Email,
	})
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

type verifyEmailRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// VerifyEmail confirms an account's email via the OTP code emailed at signup.
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req verifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.uc.VerifyEmail(c.Request.Context(), req.Email, req.Code)
	if err != nil {
		respondError(c, err)
		return
	}
	// Verification succeeded: the user is now logged in (tokens issued).
	authResponse(c, http.StatusOK, res)
}

type resendVerificationRequest struct {
	Email string `json:"email" binding:"required"`
}

// ResendVerification emails a fresh OTP. The response is always 200 regardless
// of whether the email exists or is already verified, to avoid account
// enumeration.
func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req resendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.ResendVerification(c.Request.Context(), req.Email); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "if the account exists and is unverified, a code has been sent"})
}

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

// ForgotPassword emails a password-reset OTP. Always 200 (no account
// enumeration), whether or not the email is registered.
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req forgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.RequestPasswordReset(c.Request.Context(), req.Email); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "if the account exists, a reset code has been sent"})
}

type resetPasswordRequest struct {
	Email       string `json:"email" binding:"required"`
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ResetPassword validates the OTP and sets a new password. On success the user
// is logged out everywhere and must sign in again with the new password.
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.uc.ResetPassword(c.Request.Context(), req.Email, req.Code, req.NewPassword); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password updated; please log in with your new password"})
}
