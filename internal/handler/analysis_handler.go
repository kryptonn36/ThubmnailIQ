package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/analysis"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/competitor"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/cdn"
	"github.com/thumbnailiq/thumbnailiq/internal/middleware"
	analysisuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/analysis"
	workspaceuc "github.com/thumbnailiq/thumbnailiq/internal/usecase/workspace"
)

type AnalysisHandler struct {
	uc          *analysisuc.Usecase
	competitors competitor.Repository
	workspaces  *workspaceuc.Usecase
	cdn         *cdn.Builder
}

func NewAnalysisHandler(uc *analysisuc.Usecase, competitors competitor.Repository, workspaces *workspaceuc.Usecase, cdnBuilder *cdn.Builder) *AnalysisHandler {
	return &AnalysisHandler{uc: uc, competitors: competitors, workspaces: workspaces, cdn: cdnBuilder}
}

// thumbnailURL renders the CDN URL for an object key, logging nothing and
// falling back to an empty string on failure so one bad row can't break an
// entire list response.
func (h *AnalysisHandler) thumbnailURL(key string) string {
	url, err := h.cdn.URL(key)
	if err != nil {
		return ""
	}
	return url
}

// allowedImageTypes is the set of content types we accept for thumbnail
// uploads, keyed by what http.DetectContentType reports for that format.
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

// detectImageContentType sniffs the real MIME type from the file's leading
// bytes (never the client-supplied header) and returns it only if it's an
// allowed image format.
func detectImageContentType(data []byte) (string, bool) {
	ct := http.DetectContentType(data)
	if !allowedImageTypes[ct] {
		return "", false
	}
	return ct, true
}

func (h *AnalysisHandler) Create(c *gin.Context) {
	const maxUploadSize = 10 << 20
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)


	workspaceID, err := resolveWorkspaceID(c, c.PostForm("workspace_id"), h.workspaces)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not resolve workspace_id: " + err.Error()})
		return
	}
	keyword := c.PostForm("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is required"})
		return
	}

	fileHeader, err := c.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "thumbnail file is required"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not open uploaded file"})
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		// Triggered by MaxBytesReader when the body exceeds maxUploadSize, or
		// on a genuine read error. Either way the upload is rejected.
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "uploaded file is too large or could not be read"})
		return
	}

	// Never trust the client-supplied Content-Type. Sniff the real type from
	// the file's magic bytes and only accept a small allowlist of image
	// formats, so non-image (or active) content can't be stored and later
	// served from our CDN domain.
	contentType, ok := detectImageContentType(data)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "thumbnail must be a JPEG, PNG or WebP image"})
		return
	}

	var projectID *uuid.UUID
	if pid := c.PostForm("project_id"); pid != "" {
		if parsed, err := uuid.Parse(pid); err == nil {
			projectID = &parsed
		}
	}

	created, err := h.uc.Create(c.Request.Context(), analysisuc.CreateParams{
		WorkspaceID: workspaceID,
		ProjectID:   projectID,
		UserID:      middleware.UserID(c),
		Keyword:     keyword,
		FileBytes:   data,
		ContentType: contentType,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"id": created.ID, "status": created.Status, "keyword": created.Keyword,
		"thumbnail_url": h.thumbnailURL(created.ThumbnailS3Key), "created_at": created.CreatedAt,
	})
}

func (h *AnalysisHandler) List(c *gin.Context) {
	workspaceID, err := resolveWorkspaceID(c, c.Query("workspace_id"), h.workspaces)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not resolve workspace_id: " + err.Error()})
		return
	}
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))
	var status *string
	if s := c.Query("status"); s != "" {
		status = &s
	}

	result, err := h.uc.List(c.Request.Context(), analysis.ListFilter{
		WorkspaceID: workspaceID, Status: status, Page: page, PerPage: perPage,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	items := make([]gin.H, 0, len(result.Items))
	for _, a := range result.Items {
		items = append(items, gin.H{
			"id": a.ID, "keyword": a.Keyword, "thumbnail_url": h.thumbnailURL(a.ThumbnailS3Key),
			"status": a.Status, "score": a.Score, "created_at": a.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "page": result.Page, "per_page": result.PerPage, "total": result.Total})
}

func (h *AnalysisHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid analysis id"})
		return
	}
	a, err := h.uc.Get(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}
	// Analyses are addressed by their own UUID, so authorize by the analysis's
	// workspace — otherwise any authenticated user could read another tenant's
	// analysis just by knowing its ID.
	if err := h.workspaces.EnsureMember(c.Request.Context(), middleware.UserID(c), a.WorkspaceID); err != nil {
		respondError(c, err)
		return
	}

	resp := gin.H{
		"id": a.ID, "status": a.Status, "keyword": a.Keyword, "thumbnail_url": h.thumbnailURL(a.ThumbnailS3Key),
		"created_at": a.CreatedAt, "score": a.Score,
		"visibility_score": a.VisibilityScore, "contrast_score": a.ContrastScore,
		"attention_score": a.AttentionScore, "mobile_score": a.MobileScore,
		"branding_score": a.BrandingScore, "curiosity_score": a.CuriosityScore,
		"competitor_count": a.CompetitorCount, "rank_in_competitors": a.RankInCompetitors,
		"relevance_score": a.RelevanceScore, "relevance_reasoning": a.RelevanceReasoning,
		"actual_ctr": a.ActualCTR, "published_at": a.PublishedAt,
		"error_message": a.ErrorMessage,
	}
	resp["cv_results"] = rawJSONOrNil(a.CVResults)
	resp["competitor_avg"] = rawJSONOrNil(a.CompetitorAvg)
	resp["suggestions"] = rawJSONOrNil(a.Suggestions)

	if comps, err := h.competitors.ListByAnalysis(c.Request.Context(), id); err == nil {
		compList := make([]gin.H, 0, len(comps))
		for _, comp := range comps {
			compList = append(compList, gin.H{
				"video_id": comp.VideoID, "video_title": comp.VideoTitle, "channel_name": comp.ChannelName,
				"thumbnail_url": comp.ThumbnailURL, "view_count": comp.ViewCount,
				"rank_position": comp.RankPosition, "score": comp.Score,
			})
		}
		resp["competitors"] = compList
	}

	if versions, err := h.uc.ListVersions(c.Request.Context(), id); err == nil && len(versions) > 0 {
		versionList := make([]gin.H, 0, len(versions))
		for _, v := range versions {
			versionList = append(versionList, gin.H{
				"id": v.ID, "version_number": v.VersionNumber, "thumbnail_url": h.thumbnailURL(v.S3Key),
				"score": v.Score, "cv_results": rawJSONOrNil(v.CVResults),
			})
		}
		resp["versions"] = versionList
	}

	c.JSON(http.StatusOK, resp)
}

// authorizeAnalysis loads an analysis and confirms the caller is a member of
// its workspace, returning ErrForbidden otherwise. Every analysis-by-ID action
// that isn't a plain read routes through here.
func (h *AnalysisHandler) authorizeAnalysis(c *gin.Context, id uuid.UUID) error {
	a, err := h.uc.Get(c.Request.Context(), id)
	if err != nil {
		return err
	}
	return h.workspaces.EnsureMember(c.Request.Context(), middleware.UserID(c), a.WorkspaceID)
}

func (h *AnalysisHandler) AddCompareVersion(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid analysis id"})
		return
	}
	if err := h.authorizeAnalysis(c, id); err != nil {
		respondError(c, err)
		return
	}
	fileHeader, err := c.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "thumbnail file is required"})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not open uploaded file"})
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read uploaded file"})
		return
	}

	version, err := h.uc.AddCompareVersion(c.Request.Context(), id, data, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id": version.ID, "version_number": version.VersionNumber,
		"thumbnail_url": h.thumbnailURL(version.S3Key), "score": version.Score,
		"cv_results": rawJSONOrNil(version.CVResults),
	})
}

type updateCTRRequest struct {
	ActualCTR   float64 `json:"actual_ctr" binding:"required,min=0,max=100"`
	PublishedAt string  `json:"published_at"`
}

func (h *AnalysisHandler) UpdateCTR(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid analysis id"})
		return
	}
	if err := h.authorizeAnalysis(c, id); err != nil {
		respondError(c, err)
		return
	}
	var req updateCTRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	publishedAt := time.Now()
	if req.PublishedAt != "" {
		if t, err := time.Parse(time.RFC3339, req.PublishedAt); err == nil {
			publishedAt = t
		}
	}
	if err := h.uc.UpdateCTR(c.Request.Context(), id, req.ActualCTR, publishedAt); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"actual_ctr": req.ActualCTR, "published_at": publishedAt})
}

func rawJSONOrNil(b []byte) json.RawMessage {
	if len(b) == 0 {
		return json.RawMessage("null")
	}
	return json.RawMessage(b)
}
