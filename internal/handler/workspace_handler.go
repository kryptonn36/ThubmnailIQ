package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/middleware"
	workspaceuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/workspace"
)

type WorkspaceHandler struct {
	uc *workspaceuc.Usecase
}

func NewWorkspaceHandler(uc *workspaceuc.Usecase) *WorkspaceHandler {
	return &WorkspaceHandler{uc: uc}
}

type createWorkspaceRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *WorkspaceHandler) Create(c *gin.Context) {
	var req createWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ws, err := h.uc.Create(c.Request.Context(), middleware.UserID(c), req.Name)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, ws)
}

func (h *WorkspaceHandler) List(c *gin.Context) {
	list, err := h.uc.List(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

type inviteMemberRequest struct {
	Email string `json:"email" binding:"required"`
	Role  string `json:"role" binding:"required"`
}

func (h *WorkspaceHandler) InviteMember(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
		return
	}
	var req inviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	member, err := h.uc.InviteMember(c.Request.Context(), workspaceID, req.Email, req.Role)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, member)
}

func (h *WorkspaceHandler) ListMembers(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
		return
	}
	members, err := h.uc.ListMembers(c.Request.Context(), workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, members)
}
