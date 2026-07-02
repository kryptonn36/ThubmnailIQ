package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/admin"
	"github.com/thumbnailiq/thumbnailiq/internal/middleware"
	adminuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/admin"
)

type AdminSettingsHandler struct {
	uc *adminuc.Usecase
}

func NewAdminSettingsHandler(uc *adminuc.Usecase) *AdminSettingsHandler {
	return &AdminSettingsHandler{uc: uc}
}

func (h *AdminSettingsHandler) Get(c *gin.Context) {
	settings, err := h.uc.GetSettings(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, settings)
}

type updateSettingsRequest struct {
	MaxUploadSizeBytes int64           `json:"max_upload_size_bytes" binding:"required"`
	AllowedExtensions  []string        `json:"allowed_extensions" binding:"required"`
	FeatureFlags       map[string]bool `json:"feature_flags"`
	StorageProvider    string          `json:"storage_provider" binding:"required"`
	EmailProvider      string          `json:"email_provider"`
	EmailFromAddress   string          `json:"email_from_address"`
}

func (h *AdminSettingsHandler) Update(c *gin.Context) {
	var req updateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.FeatureFlags == nil {
		req.FeatureFlags = map[string]bool{}
	}
	updated, err := h.uc.UpdateSettings(c.Request.Context(), middleware.AdminID(c), &admin.Settings{
		MaxUploadSizeBytes: req.MaxUploadSizeBytes,
		AllowedExtensions:  req.AllowedExtensions,
		FeatureFlags:       req.FeatureFlags,
		StorageProvider:    req.StorageProvider,
		EmailProvider:      req.EmailProvider,
		EmailFromAddress:   req.EmailFromAddress,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}
