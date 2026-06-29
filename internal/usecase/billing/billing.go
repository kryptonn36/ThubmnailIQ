package billing

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/billing"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/payment"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/workspace"
	"github.com/thumbnailiq/thumbnailiq/pkg/errors"
	"github.com/thumbnailiq/thumbnailiq/pkg/hash"
)

type Usecase struct {
	billing    billing.Repository
	workspaces workspace.Repository
	gateway    payment.Gateway
	currency   string
}

func NewUsecase(billingRepo billing.Repository, workspaces workspace.Repository, gateway payment.Gateway, currency string) *Usecase {
	return &Usecase{billing: billingRepo, workspaces: workspaces, gateway: gateway, currency: currency}
}

func (u *Usecase) Plans() []billing.Plan {
	return billing.Plans
}

func findPlan(id string) *billing.Plan {
	for i := range billing.Plans {
		if billing.Plans[i].ID == id {
			return &billing.Plans[i]
		}
	}
	return nil
}

// Checkout is what the frontend needs to either confirm an already-activated
// free plan, or open the payment provider's checkout widget for a paid plan.
type Checkout struct {
	RequiresPayment bool
	Plan            string
	Status          string
	Order           *payment.Order
}

// activeSubscription returns the workspace's current subscription and
// whether it's still within its paid period. A workspace with no
// subscription row yet (never been through checkout) is treated as having
// no active subscription, not an error.
func (u *Usecase) activeSubscription(ctx context.Context, workspaceID uuid.UUID) (*billing.Subscription, bool, error) {
	sub, err := u.billing.GetSubscriptionByWorkspace(ctx, workspaceID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return sub, sub.Status == "active" && sub.CurrentPeriodEnd.After(time.Now()), nil
}

// CurrentPlanInfo is what the frontend needs to show a workspace's current
// plan and whether it's still valid.
type CurrentPlanInfo struct {
	Plan             string     `json:"plan"`
	Status           string     `json:"status"`
	CurrentPeriodEnd *time.Time `json:"current_period_end,omitempty"`
	IsActive         bool       `json:"is_active"`
}

// CurrentSubscription reports the workspace's current plan. Workspaces that
// have never been through checkout default to the free plan.
func (u *Usecase) CurrentSubscription(ctx context.Context, workspaceID uuid.UUID) (*CurrentPlanInfo, error) {
	sub, isActive, err := u.activeSubscription(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return &CurrentPlanInfo{Plan: "free", Status: "active", IsActive: true}, nil
	}
	periodEnd := sub.CurrentPeriodEnd
	return &CurrentPlanInfo{Plan: sub.Plan, Status: sub.Status, CurrentPeriodEnd: &periodEnd, IsActive: isActive}, nil
}

// CreateCheckout starts a plan change. Free plans are activated immediately
// (only possible when there's no active paid subscription, since downgrades
// are blocked below). Re-buying the plan a workspace is already on while
// it's still valid is rejected. Upgrading to a pricier plan only charges the
// difference against what's already been paid for the current period.
// Switching to a same-or-cheaper plan is rejected outright, the same way
// most subscription products only let you upgrade self-serve, never
// downgrade, while a paid period is still running.
func (u *Usecase) CreateCheckout(ctx context.Context, workspaceID uuid.UUID, plan string) (*Checkout, error) {
	planDef := findPlan(plan)
	if planDef == nil {
		return nil, errors.ErrInvalidInput
	}

	current, isActive, err := u.activeSubscription(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	if isActive && current.Plan == plan {
		return nil, fmt.Errorf("you're already on the %s plan until %s: %w",
			planDef.Name, current.CurrentPeriodEnd.Format("Jan 2, 2006"), errors.ErrPlanAlreadyActive)
	}

	var currentPlanDef *billing.Plan
	if isActive {
		currentPlanDef = findPlan(current.Plan)
	}
	if isActive && currentPlanDef != nil && planDef.PriceMonthly <= currentPlanDef.PriceMonthly {
		return nil, fmt.Errorf("you're on the %s plan until %s; downgrades aren't available while a plan is active, only upgrades: %w",
			currentPlanDef.Name, current.CurrentPeriodEnd.Format("Jan 2, 2006"), errors.ErrDowngradeNotAllowed)
	}

	if planDef.PriceMonthly == 0 {
		sub, err := u.activate(ctx, workspaceID, planDef, fmt.Sprintf("free_%s", uuid.New().String()))
		if err != nil {
			return nil, err
		}
		return &Checkout{RequiresPayment: false, Plan: sub.Plan, Status: sub.Status}, nil
	}

	amountDue := planDef.PriceMonthly
	if currentPlanDef != nil {
		amountDue = planDef.PriceMonthly - currentPlanDef.PriceMonthly
	}

	order, err := u.gateway.CreateOrder(ctx, payment.CreateOrderParams{
		Plan:        plan,
		AmountMinor: int64(amountDue * 100),
		Currency:    u.currency,
		Receipt:     fmt.Sprintf("ws_%s_%d", strings.ReplaceAll(workspaceID.String(), "-", "")[:8], time.Now().UnixNano()),
	})
	if err != nil {
		return nil, err
	}
	return &Checkout{RequiresPayment: true, Plan: plan, Order: order}, nil
}

// ConfirmCheckout verifies a payment the frontend reports as successful and,
// once verified, activates the plan it paid for.
func (u *Usecase) ConfirmCheckout(ctx context.Context, workspaceID uuid.UUID, plan, orderID, paymentID, signature string) (*billing.Subscription, error) {
	planDef := findPlan(plan)
	if planDef == nil {
		return nil, errors.ErrInvalidInput
	}
	if err := u.gateway.VerifyPayment(ctx, payment.VerifyParams{
		OrderID: orderID, PaymentID: paymentID, Signature: signature,
	}); err != nil {
		return nil, err
	}
	return u.activate(ctx, workspaceID, planDef, paymentID)
}

func (u *Usecase) activate(ctx context.Context, workspaceID uuid.UUID, planDef *billing.Plan, providerSubID string) (*billing.Subscription, error) {
	periodStart := time.Now()
	periodEnd := periodStart.AddDate(0, 1, 0)

	sub, err := u.billing.UpsertSubscription(ctx, &billing.Subscription{
		WorkspaceID: workspaceID, Provider: u.gateway.Name(), ProviderSubscriptionID: providerSubID,
		Plan: planDef.ID, Status: "active", CurrentPeriodStart: periodStart, CurrentPeriodEnd: periodEnd,
	})
	if err != nil {
		return nil, err
	}
	if _, err := u.workspaces.UpdatePlan(ctx, workspaceID, planDef.ID, planDef.AnalysesLimit); err != nil {
		return nil, err
	}
	return sub, nil
}

func (u *Usecase) CreateAPIKey(ctx context.Context, workspaceID, userID uuid.UUID, name string) (string, *billing.APIKey, error) {
	raw, prefix, err := hash.GenerateAPIKey()
	if err != nil {
		return "", nil, err
	}
	key, err := u.billing.CreateAPIKey(ctx, workspaceID, name, hash.SHA256Hex(raw), prefix, userID)
	if err != nil {
		return "", nil, err
	}
	return raw, key, nil
}

func (u *Usecase) ListAPIKeys(ctx context.Context, workspaceID uuid.UUID) ([]*billing.APIKey, error) {
	return u.billing.ListAPIKeys(ctx, workspaceID)
}

func (u *Usecase) RevokeAPIKey(ctx context.Context, id, workspaceID uuid.UUID) error {
	return u.billing.RevokeAPIKey(ctx, id, workspaceID)
}
