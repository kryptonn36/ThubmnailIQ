package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/admin"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/cdn"
	"github.com/thumbnailiq/thumbnailiq/internal/middleware"
	adminuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/admin"
	"github.com/thumbnailiq/thumbnailiq/pkg/pagination"
)

type AdminUploadsHandler struct {
	uc  *adminuc.Usecase
	cdn *cdn.Builder
}

func NewAdminUploadsHandler(uc *adminuc.Usecase, cdnBuilder *cdn.Builder) *AdminUploadsHandler {
	return &AdminUploadsHandler{uc: uc, cdn: cdnBuilder}
}

// thumbnailURL mirrors AnalysisHandler.thumbnailURL: falls back to an empty
// string on failure so one bad row can't break an entire list response.
func (h *AdminUploadsHandler) thumbnailURL(key string) string {
	url, err := h.cdn.URL(key)
	if err != nil {
		return ""
	}
	return url
}

func (h *AdminUploadsHandler) uploadJSON(u *admin.UploadSummary) gin.H {
	return gin.H{
		"id": u.ID, "workspace_id": u.WorkspaceID, "user_id": u.UserID, "keyword": u.Keyword,
		"thumbnail_url": h.thumbnailURL(u.ThumbnailS3Key), "status": u.Status, "score": u.Score,
		"file_size_bytes": u.FileSizeBytes, "created_at": u.CreatedAt, "deleted_at": u.DeletedAt,
	}
}

func (h *AdminUploadsHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	displayPage, displayPerPage := pagination.Normalize(page, perPage)
	includeDeleted, _ := strconv.ParseBool(c.Query("include_deleted"))

	uploads, total, err := h.uc.ListUploads(c.Request.Context(), admin.UploadFilter{
		Search: c.Query("search"), Status: c.Query("status"), IncludeDeleted: includeDeleted,
	}, admin.Pagination{Page: page, PerPage: perPage})
	if err != nil {
		respondError(c, err)
		return
	}
	data := make([]gin.H, len(uploads))
	for i, u := range uploads {
		data[i] = h.uploadJSON(u)
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "page": displayPage, "per_page": displayPerPage, "total": total})
}

func (h *AdminUploadsHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid upload id"})
		return
	}
	upload, err := h.uc.GetUpload(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, h.uploadJSON(upload))
}

func (h *AdminUploadsHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid upload id"})
		return
	}
	if err := h.uc.DeleteUpload(c.Request.Context(), middleware.AdminID(c), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *AdminUploadsHandler) Restore(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid upload id"})
		return
	}
	if err := h.uc.RestoreUpload(c.Request.Context(), middleware.AdminID(c), id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "restored"})
}

// Download redirects to the same public CDN URL every other read path in
// this codebase already uses (see cdn.Builder) — the bucket is public-read
// behind CloudFront, so there is no separate presigned-URL flow to build.
func (h *AdminUploadsHandler) Download(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid upload id"})
		return
	}
	upload, err := h.uc.GetUpload(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	url, err := h.cdn.URL(upload.ThumbnailS3Key)
	if err != nil {
		respondError(c, err)
		return
	}
	c.Redirect(http.StatusFound, url)
}
