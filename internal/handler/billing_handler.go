package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/middleware"
	billinguc "github.com/thumbnailiq/thumbnailiq/internal/usecase/billing"
	workspaceuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/workspace"
)

type BillingHandler struct {
	uc         *billinguc.Usecase
	workspaces *workspaceuc.Usecase
}

func NewBillingHandler(uc *billinguc.Usecase, workspaces *workspaceuc.Usecase) *BillingHandler {
	return &BillingHandler{uc: uc, workspaces: workspaces}
}

func (h *BillingHandler) Plans(c *gin.Context) {
	c.JSON(http.StatusOK, h.uc.Plans())
}

type subscribeRequest struct {
	WorkspaceID string `json:"workspace_id"`
	Plan        string `json:"plan" binding:"required"`
}

func (h *BillingHandler) Subscribe(c *gin.Context) {
	var req subscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspaceID, err := resolveWorkspaceID(c, req.WorkspaceID, h.workspaces)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not resolve workspace_id: " + err.Error()})
		return
	}
	sub, err := h.uc.Subscribe(c.Request.Context(), workspaceID, req.Plan)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"plan": sub.Plan, "status": sub.Status})
}

type createAPIKeyRequest struct {
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name" binding:"required"`
}

func (h *BillingHandler) CreateAPIKey(c *gin.Context) {
	var req createAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspaceID, err := resolveWorkspaceID(c, req.WorkspaceID, h.workspaces)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not resolve workspace_id: " + err.Error()})
		return
	}
	raw, key, err := h.uc.CreateAPIKey(c.Request.Context(), workspaceID, middleware.UserID(c), req.Name)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": key.ID, "name": key.Name, "key_prefix": key.KeyPrefix, "key": raw})
}

func (h *BillingHandler) ListAPIKeys(c *gin.Context) {
	workspaceID, err := resolveWorkspaceID(c, c.Query("workspace_id"), h.workspaces)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not resolve workspace_id: " + err.Error()})
		return
	}
	keys, err := h.uc.ListAPIKeys(c.Request.Context(), workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, keys)
}

func (h *BillingHandler) RevokeAPIKey(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid api key id"})
		return
	}
	workspaceID, err := resolveWorkspaceID(c, c.Query("workspace_id"), h.workspaces)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not resolve workspace_id: " + err.Error()})
		return
	}
	if err := h.uc.RevokeAPIKey(c.Request.Context(), id, workspaceID); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
