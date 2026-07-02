package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/billing"
	"github.com/thumbnailiq/thumbnailiq/internal/handler"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/redis"
	"github.com/thumbnailiq/thumbnailiq/internal/middleware"
	"github.com/thumbnailiq/thumbnailiq/pkg/jwt"
)

type Handlers struct {
	Auth           *handler.AuthHandler
	Workspace      *handler.WorkspaceHandler
	Analysis       *handler.AnalysisHandler
	Competitor     *handler.CompetitorHandler
	Tracking       *handler.TrackingHandler
	Billing        *handler.BillingHandler
	ViralDB        *handler.ViralDBHandler
	QuickAnalyze   *handler.QuickAnalyzeHandler
	AdminAuth      *handler.AdminAuthHandler
	AdminDashboard *handler.AdminDashboardHandler
	AdminUsers     *handler.AdminUsersHandler
	AdminUploads   *handler.AdminUploadsHandler
	AdminAnalytics *handler.AdminAnalyticsHandler
	AdminSettings  *handler.AdminSettingsHandler
	AdminLogs      *handler.AdminLogsHandler
	AdminProfile   *handler.AdminProfileHandler
}

func NewRouter(h *Handlers, jwtSvc, adminJWTSvc *jwt.Service, billingRepo billing.Repository, rateLimiter *redis.RateLimiter, log zerolog.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.RequestLogger(log), middleware.CORS())

	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.POST("/register", h.Auth.Register)
	auth.POST("/login", h.Auth.Login)
	auth.POST("/refresh", h.Auth.Refresh)

	// Admin panel — fully separate identity/session space from customer auth
	// above. Login is rate-limited by client IP (no admin JWT exists yet at
	// this point) to slow down credential-stuffing against admin accounts.
	adminAuth := v1.Group("/admin/auth")
	adminAuth.Use(middleware.RateLimit(rateLimiter, 10, time.Minute))
	adminAuth.POST("/login", h.AdminAuth.Login)
	adminAuth.POST("/refresh", h.AdminAuth.Refresh)

	adminAuthed := v1.Group("/admin")
	adminAuthed.Use(middleware.AdminAuth(adminJWTSvc))

	adminAuthed.GET("/dashboard", h.AdminDashboard.Get)

	adminAuthed.GET("/users", h.AdminUsers.List)
	adminAuthed.GET("/users/:id", h.AdminUsers.Get)
	adminAuthed.PATCH("/users/:id/suspend", h.AdminUsers.Suspend)
	adminAuthed.PATCH("/users/:id/activate", h.AdminUsers.Activate)
	adminAuthed.DELETE("/users/:id", h.AdminUsers.Delete)
	adminAuthed.POST("/users/:id/reset-password", h.AdminUsers.ResetPassword)
	adminAuthed.PATCH("/users/:id/role", h.AdminUsers.ChangeRole)
	adminAuthed.GET("/users/:id/uploads", h.AdminUsers.Uploads)

	adminAuthed.GET("/uploads", h.AdminUploads.List)
	adminAuthed.GET("/uploads/:id", h.AdminUploads.Get)
	adminAuthed.DELETE("/uploads/:id", h.AdminUploads.Delete)
	adminAuthed.POST("/uploads/:id/restore", h.AdminUploads.Restore)
	adminAuthed.GET("/uploads/:id/download", h.AdminUploads.Download)

	adminAuthed.GET("/analytics", h.AdminAnalytics.Get)

	adminAuthed.GET("/settings", h.AdminSettings.Get)
	adminAuthed.PATCH("/settings", h.AdminSettings.Update)

	adminAuthed.GET("/logs/audit", h.AdminLogs.List)

	adminAuthed.GET("/profile", h.AdminProfile.Get)
	adminAuthed.PATCH("/profile/password", h.AdminProfile.ChangePassword)

	authed := v1.Group("")
	authed.Use(middleware.Auth(jwtSvc))

	authed.POST("/workspaces", h.Workspace.Create)
	authed.GET("/workspaces", h.Workspace.List)
	authed.PATCH("/workspaces/:id/brand", h.Workspace.UpdateBrand)
	authed.GET("/workspaces/:id/members", h.Workspace.ListMembers)
	authed.POST("/workspaces/:id/members", h.Workspace.InviteMember)

	authed.POST("/analyses", h.Analysis.Create)
	authed.GET("/analyses", h.Analysis.List)
	authed.GET("/analyses/:id", h.Analysis.Get)
	authed.POST("/analyses/:id/compare", h.Analysis.AddCompareVersion)
	authed.PATCH("/analyses/:id/ctr", h.Analysis.UpdateCTR)

	authed.GET("/keywords/:keyword/competitors", h.Competitor.ListForKeyword)
	authed.GET("/keywords/:keyword/ideas", h.Competitor.VideoIdeas)

	authed.POST("/tracking", h.Tracking.Create)
	authed.GET("/tracking", h.Tracking.List)

	authed.GET("/billing/plans", h.Billing.Plans)
	authed.GET("/billing/subscription", h.Billing.CurrentSubscription)
	authed.POST("/billing/checkout", h.Billing.Checkout)
	authed.POST("/billing/checkout/verify", h.Billing.ConfirmCheckout)
	authed.POST("/api-keys", h.Billing.CreateAPIKey)
	authed.GET("/api-keys", h.Billing.ListAPIKeys)
	authed.DELETE("/api-keys/:id", h.Billing.RevokeAPIKey)

	authed.GET("/viral-db", h.ViralDB.Search)

	// API-key-authenticated routes — used by the browser extension and
	// any other non-browser API consumer that holds a ThumbnailIQ API key.
	apiKeyed := v1.Group("/public")
	apiKeyed.Use(middleware.APIKeyAuth(billingRepo))
	apiKeyed.POST("/quick-analyze", h.QuickAnalyze.Analyze)

	return r
}
