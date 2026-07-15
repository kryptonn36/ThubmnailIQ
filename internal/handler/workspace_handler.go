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

// List returns the caller's workspaces enriched with owner name/email, the
// caller's role, and member count, so the client can always show whose
// workspace each one is. The plain workspace fields are unchanged, so older
// clients reading only those keep working.
func (h *WorkspaceHandler) List(c *gin.Context) {
	list, err := h.uc.ListWithContext(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

type renameWorkspaceRequest struct {
	Name string `json:"name" binding:"required"`
}

// Rename updates the workspace's display name (owner/admin only).
func (h *WorkspaceHandler) Rename(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
		return
	}
	if err := h.uc.EnsureRole(c.Request.Context(), middleware.UserID(c), workspaceID, "owner", "admin"); err != nil {
		respondError(c, err)
		return
	}
	var req renameWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ws, err := h.uc.Rename(c.Request.Context(), workspaceID, req.Name)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, ws)
}

// RemoveMember removes a member from the workspace. The usecase enforces the
// rules (owner is irremovable; self-removal = leaving; removing others needs
// owner/admin), so membership here is only checked implicitly through them.
func (h *WorkspaceHandler) RemoveMember(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
		return
	}
	memberUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member user id"})
		return
	}
	if err := h.uc.RemoveMember(c.Request.Context(), workspaceID, middleware.UserID(c), memberUserID); err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
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
	// Inviting members is an owner/admin-only action, not something any member
	// of the workspace can do.
	if err := h.uc.EnsureRole(c.Request.Context(), middleware.UserID(c), workspaceID, "owner", "admin"); err != nil {
		respondError(c, err)
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

type updateBrandRequest struct {
	PrimaryColor   string `json:"primary_color" binding:"required"`
	SecondaryColor string `json:"secondary_color" binding:"required"`
	Font           string `json:"font" binding:"required"`
}

func (h *WorkspaceHandler) UpdateBrand(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
		return
	}
	if err := h.uc.EnsureMember(c.Request.Context(), middleware.UserID(c), workspaceID); err != nil {
		respondError(c, err)
		return
	}
	var req updateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ws, err := h.uc.UpdateBrand(c.Request.Context(), workspaceID, req.PrimaryColor, req.SecondaryColor, req.Font)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, ws)
}

func (h *WorkspaceHandler) ListMembers(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
		return
	}
	if err := h.uc.EnsureMember(c.Request.Context(), middleware.UserID(c), workspaceID); err != nil {
		respondError(c, err)
		return
	}
	members, err := h.uc.ListMembers(c.Request.Context(), workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, members)
}
