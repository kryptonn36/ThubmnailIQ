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

// CurrentSubscription reports the workspace's current plan and whether it's
// still valid, so the frontend can show it and decide what's purchasable.
func (h *BillingHandler) CurrentSubscription(c *gin.Context) {
	workspaceID, err := resolveWorkspaceID(c, c.Query("workspace_id"), h.workspaces)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not resolve workspace_id: " + err.Error()})
		return
	}
	info, err := h.uc.CurrentSubscription(c.Request.Context(), workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

type checkoutRequest struct {
	WorkspaceID string `json:"workspace_id"`
	Plan        string `json:"plan" binding:"required"`
}

// Checkout starts a plan change. Free plans are activated immediately; paid
// plans return an order for the frontend to open the payment gateway's
// checkout widget against (Razorpay today; Stripe once approved).
func (h *BillingHandler) Checkout(c *gin.Context) {
	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspaceID, err := resolveWorkspaceID(c, req.WorkspaceID, h.workspaces)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not resolve workspace_id: " + err.Error()})
		return
	}
	checkout, err := h.uc.CreateCheckout(c.Request.Context(), workspaceID, req.Plan)
	if err != nil {
		respondError(c, err)
		return
	}
	if !checkout.RequiresPayment {
		c.JSON(http.StatusOK, gin.H{"requires_payment": false, "plan": checkout.Plan, "status": checkout.Status})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"requires_payment": true,
		"plan":             checkout.Plan,
		"provider":         checkout.Order.Provider,
		"order_id":         checkout.Order.ID,
		"amount":           checkout.Order.AmountMinor,
		"currency":         checkout.Order.Currency,
		"key_id":           checkout.Order.KeyID,
	})
}

type confirmCheckoutRequest struct {
	WorkspaceID string `json:"workspace_id"`
	Plan        string `json:"plan" binding:"required"`
	OrderID     string `json:"order_id" binding:"required"`
	PaymentID   string `json:"payment_id" binding:"required"`
	Signature   string `json:"signature" binding:"required"`
}

// ConfirmCheckout verifies a payment the gateway's checkout widget reports as
// successful, then activates the plan it paid for.
func (h *BillingHandler) ConfirmCheckout(c *gin.Context) {
	var req confirmCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspaceID, err := resolveWorkspaceID(c, req.WorkspaceID, h.workspaces)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not resolve workspace_id: " + err.Error()})
		return
	}
	sub, err := h.uc.ConfirmCheckout(c.Request.Context(), workspaceID, req.Plan, req.OrderID, req.PaymentID, req.Signature)
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
