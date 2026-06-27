package billing

import (
	"context"
	"fmt"
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

// CreateCheckout starts a plan change. Free plans are activated immediately;
// paid plans get an order placed with the active payment gateway, which the
// frontend uses to open that gateway's checkout widget.
func (u *Usecase) CreateCheckout(ctx context.Context, workspaceID uuid.UUID, plan string) (*Checkout, error) {
	planDef := findPlan(plan)
	if planDef == nil {
		return nil, errors.ErrInvalidInput
	}

	if planDef.PriceMonthly == 0 {
		sub, err := u.activate(ctx, workspaceID, planDef, fmt.Sprintf("free_%s", uuid.New().String()))
		if err != nil {
			return nil, err
		}
		return &Checkout{RequiresPayment: false, Plan: sub.Plan, Status: sub.Status}, nil
	}

	order, err := u.gateway.CreateOrder(ctx, payment.CreateOrderParams{
		Plan:        plan,
		AmountMinor: int64(planDef.PriceMonthly * 100),
		Currency:    u.currency,
		Receipt:     fmt.Sprintf("ws_%s_%s", workspaceID, plan),
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
