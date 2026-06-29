package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/thumbnailiq/thumbnailiq/internal/handler"
	"github.com/thumbnailiq/thumbnailiq/internal/middleware"
	"github.com/thumbnailiq/thumbnailiq/pkg/jwt"
)

type Handlers struct {
	Auth       *handler.AuthHandler
	Workspace  *handler.WorkspaceHandler
	Analysis   *handler.AnalysisHandler
	Competitor *handler.CompetitorHandler
	Tracking   *handler.TrackingHandler
	Billing    *handler.BillingHandler
	ViralDB    *handler.ViralDBHandler
}

func NewRouter(h *Handlers, jwtSvc *jwt.Service, log zerolog.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.RequestLogger(log), middleware.CORS())

	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	v1 := r.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.POST("/register", h.Auth.Register)
	auth.POST("/login", h.Auth.Login)
	auth.POST("/refresh", h.Auth.Refresh)

	authed := v1.Group("")
	authed.Use(middleware.Auth(jwtSvc))

	authed.POST("/workspaces", h.Workspace.Create)
	authed.GET("/workspaces", h.Workspace.List)
	authed.GET("/workspaces/:id/members", h.Workspace.ListMembers)
	authed.POST("/workspaces/:id/members", h.Workspace.InviteMember)

	authed.POST("/analyses", h.Analysis.Create)
	authed.GET("/analyses", h.Analysis.List)
	authed.GET("/analyses/:id", h.Analysis.Get)
	authed.POST("/analyses/:id/compare", h.Analysis.AddCompareVersion)

	authed.GET("/keywords/:keyword/competitors", h.Competitor.ListForKeyword)

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

	return r
}
