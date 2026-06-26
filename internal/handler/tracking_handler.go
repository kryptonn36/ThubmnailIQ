package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/thumbnailiq/thumbnailiq/internal/middleware"
	trackinguc "github.com/thumbnailiq/thumbnailiq/internal/usecase/tracking"
	workspaceuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/workspace"
)

type TrackingHandler struct {
	uc         *trackinguc.Usecase
	workspaces *workspaceuc.Usecase
}

func NewTrackingHandler(uc *trackinguc.Usecase, workspaces *workspaceuc.Usecase) *TrackingHandler {
	return &TrackingHandler{uc: uc, workspaces: workspaces}
}

type createTrackingRequest struct {
	WorkspaceID   string `json:"workspace_id"`
	Type          string `json:"type" binding:"required"`
	ChannelID     string `json:"channel_id"`
	Keyword       string `json:"keyword"`
	IntervalHours int    `json:"interval_hours"`
}

func (h *TrackingHandler) Create(c *gin.Context) {
	var req createTrackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspaceID, err := resolveWorkspaceID(c, req.WorkspaceID, h.workspaces)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not resolve workspace_id: " + err.Error()})
		return
	}

	job, err := h.uc.Create(c.Request.Context(), trackinguc.CreateParams{
		WorkspaceID: workspaceID, UserID: middleware.UserID(c), Type: req.Type,
		ChannelID: req.ChannelID, Keyword: req.Keyword, IntervalHours: req.IntervalHours,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, job)
}

func (h *TrackingHandler) List(c *gin.Context) {
	workspaceID, err := resolveWorkspaceID(c, c.Query("workspace_id"), h.workspaces)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not resolve workspace_id: " + err.Error()})
		return
	}
	jobs, err := h.uc.List(c.Request.Context(), workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, jobs)
}
