package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/admin"
	"github.com/thumbnailiq/thumbnailiq/internal/middleware"
	adminuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/admin"
	"github.com/thumbnailiq/thumbnailiq/pkg/pagination"
)

type AdminUsersHandler struct {
	uc *adminuc.Usecase
}

func NewAdminUsersHandler(uc *adminuc.Usecase) *AdminUsersHandler {
	return &AdminUsersHandler{uc: uc}
}

func (h *AdminUsersHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	displayPage, displayPerPage := pagination.Normalize(page, perPage)

	users, total, err := h.uc.ListUsers(c.Request.Context(), admin.UserFilter{
		Search: c.Query("search"), Status: c.Query("status"),
	}, admin.Pagination{Page: page, PerPage: perPage})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users, "page": displayPage, "per_page": displayPerPage, "total": total})
}

func (h *AdminUsersHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	detail, err := h.uc.GetUserDetail(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, detail)
}

func (h *AdminUsersHandler) Suspend(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	if err := h.uc.SuspendUser(c.Request.Context(), middleware.AdminID(c), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "suspended"})
}

func (h *AdminUsersHandler) Activate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	if err := h.uc.ActivateUser(c.Request.Context(), middleware.AdminID(c), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "active"})
}

func (h *AdminUsersHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	if err := h.uc.DeleteUser(c.Request.Context(), middleware.AdminID(c), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *AdminUsersHandler) ResetPassword(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	tempPassword, err := h.uc.ResetUserPassword(c.Request.Context(), middleware.AdminID(c), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"temporary_password": tempPassword})
}

type changeUserRoleRequest struct {
	WorkspaceID string `json:"workspace_id" binding:"required"`
	Role        string `json:"role" binding:"required"`
}

func (h *AdminUsersHandler) ChangeRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var req changeUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspaceID, err := uuid.Parse(req.WorkspaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace_id"})
		return
	}
	if err := h.uc.ChangeUserWorkspaceRole(c.Request.Context(), middleware.AdminID(c), id, workspaceID, req.Role); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *AdminUsersHandler) Uploads(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	displayPage, displayPerPage := pagination.Normalize(page, perPage)

	uploads, total, err := h.uc.ListUserUploads(c.Request.Context(), id, admin.Pagination{Page: page, PerPage: perPage})
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": uploads, "page": displayPage, "per_page": displayPerPage, "total": total})
}
