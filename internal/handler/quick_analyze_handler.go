package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/infra/cv"
	"github.com/thumbnailiq/thumbnailiq/internal/scoring"
)

// QuickAnalyzeHandler is the synchronous analysis endpoint used by the
// browser extension (and any API key consumer) — it runs the CV pipeline
// inline without the async queue, so results come back in a single HTTP
// request. Only used for overlay scoring; full pipeline analyses still
// go through POST /analyses.
type QuickAnalyzeHandler struct {
	cvClient *cv.Client
	engine   *scoring.Engine
}

func NewQuickAnalyzeHandler(cvClient *cv.Client) *QuickAnalyzeHandler {
	return &QuickAnalyzeHandler{cvClient: cvClient, engine: scoring.NewEngine()}
}

type quickAnalyzeRequest struct {
	ThumbnailURL string `json:"thumbnail_url" binding:"required"`
	Keyword      string `json:"keyword"`
}

func (h *QuickAnalyzeHandler) Analyze(c *gin.Context) {
	var req quickAnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cvResult, err := h.cvClient.Analyze(c.Request.Context(), req.ThumbnailURL)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Workspace ID from API-key auth middleware
	var workspaceIDStr string
	if v, ok := c.Get("workspace_id"); ok {
		if id, ok := v.(uuid.UUID); ok {
			workspaceIDStr = id.String()
		}
	}

	sub := scoring.SubScores{
		Visibility: h.engine.Visibility(cvResult, scoring.CompetitorAvg{}),
		Contrast:   h.engine.Contrast(cvResult),
		Attention:  h.engine.Attention(cvResult),
		Mobile:     h.engine.Mobile(cvResult),
		Branding:   h.engine.Branding(),
		Curiosity:  50,
	}
	finalScore := h.engine.FinalScore(sub)

	cvJSON, _ := json.Marshal(cvResult)

	c.JSON(http.StatusOK, gin.H{
		"score":            finalScore,
		"visibility_score": sub.Visibility,
		"contrast_score":   sub.Contrast,
		"attention_score":  sub.Attention,
		"mobile_score":     sub.Mobile,
		"branding_score":   sub.Branding,
		"cv_results":       json.RawMessage(cvJSON),
		"workspace_id":     workspaceIDStr,
		"keyword":          req.Keyword,
	})
}
