package billing

import (
	"context"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/billing"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/payment"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/user"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/workspace"
	"github.com/thumbnailiq/thumbnailiq/pkg/errors"
	"github.com/thumbnailiq/thumbnailiq/pkg/hash"
)

// Mailer is the narrow email dependency this usecase needs for subscription
// confirmation emails. The concrete implementation lives in internal/infra/email.
type Mailer interface {
	Send(ctx context.Context, to, subject, textBody, htmlBody string) error
}

type Usecase struct {
	billing    billing.Repository
	workspaces workspace.Repository
	users      user.Repository
	gateway    payment.Gateway
	orders     payment.PendingOrderStore
	currency   string
	mailer     Mailer
	log        zerolog.Logger
}

func NewUsecase(billingRepo billing.Repository, workspaces workspace.Repository, users user.Repository, gateway payment.Gateway, orders payment.PendingOrderStore, currency string, mailer Mailer, log zerolog.Logger) *Usecase {
	return &Usecase{billing: billingRepo, workspaces: workspaces, users: users, gateway: gateway, orders: orders, currency: currency, mailer: mailer, log: log}
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

	amountMinor := int64(amountDue * 100)
	order, err := u.gateway.CreateOrder(ctx, payment.CreateOrderParams{
		Plan:        plan,
		AmountMinor: amountMinor,
		Currency:    u.currency,
		Receipt:     fmt.Sprintf("ws_%s_%d", strings.ReplaceAll(workspaceID.String(), "-", "")[:8], time.Now().UnixNano()),
	})
	if err != nil {
		return nil, err
	}
	// Record what this order was actually for. ConfirmCheckout reads this back
	// so the activated plan is bound to the order the user paid, not to whatever
	// plan the client echoes at confirmation time. If we can't persist it, fail
	// now rather than let the user pay for something we can't later verify.
	if err := u.orders.Save(ctx, order.ID, payment.PendingOrder{
		WorkspaceID: workspaceID,
		Plan:        plan,
		AmountMinor: amountMinor,
	}); err != nil {
		return nil, err
	}
	return &Checkout{RequiresPayment: true, Plan: plan, Order: order}, nil
}

// ConfirmCheckout verifies a payment the frontend reports as successful and,
// once verified, activates the plan it paid for. The client-supplied `plan` is
// deliberately NOT trusted: the plan and workspace that get activated are read
// back from the server-side order record created in CreateCheckout, so a caller
// can't pay for a cheap plan and confirm an expensive one (or target another
// workspace) by tampering with the confirmation request.
func (u *Usecase) ConfirmCheckout(ctx context.Context, workspaceID uuid.UUID, plan, orderID, paymentID, signature string) (*billing.Subscription, error) {
	if err := u.gateway.VerifyPayment(ctx, payment.VerifyParams{
		OrderID: orderID, PaymentID: paymentID, Signature: signature,
	}); err != nil {
		return nil, err
	}

	// Consume the order we recorded at checkout time. A missing record means the
	// order was never created here, already confirmed, or expired — reject it.
	pending, err := u.orders.Consume(ctx, orderID)
	if err != nil {
		return nil, errors.ErrInvalidInput
	}
	// The confirmation must target the same workspace the order was opened for.
	if pending.WorkspaceID != workspaceID {
		return nil, errors.ErrForbidden
	}
	planDef := findPlan(pending.Plan)
	if planDef == nil {
		return nil, errors.ErrInvalidInput
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
	// The plan is live at this point; the confirmation email is notification
	// only and must never fail or delay the activation.
	u.notifySubscriptionActivated(workspaceID, planDef, periodEnd)
	return sub, nil
}

// notifySubscriptionActivated emails the workspace owner that their plan is
// active. It runs in the background and only logs failures: a mail outage
// must never make a paid activation look failed to the customer.
func (u *Usecase) notifySubscriptionActivated(workspaceID uuid.UUID, planDef *billing.Plan, periodEnd time.Time) {
	if u.mailer == nil {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		ws, err := u.workspaces.GetByID(ctx, workspaceID)
		if err != nil {
			u.log.Error().Err(err).Str("workspace_id", workspaceID.String()).Msg("loading workspace for subscription email")
			return
		}
		owner, err := u.users.GetByID(ctx, ws.OwnerID)
		if err != nil {
			u.log.Error().Err(err).Str("workspace_id", workspaceID.String()).Msg("loading workspace owner for subscription email")
			return
		}

		price := "Free"
		if planDef.PriceMonthly > 0 {
			price = fmt.Sprintf("%.2f %s/month", planDef.PriceMonthly, u.currency)
		}
		until := periodEnd.Format("Jan 2, 2006")

		subject := fmt.Sprintf("Your ThumbnailIQ %s plan is now active", planDef.Name)
		text := fmt.Sprintf("Hi %s,\n\nYour %q workspace is now on the %s plan (%s).\n\nWhat you get:\n- %s\n\nYour current period runs until %s.\n\nThanks for using ThumbnailIQ!",
			owner.FullName, ws.Name, planDef.Name, price, strings.Join(planDef.Features, "\n- "), until)

		features := make([]string, 0, len(planDef.Features))
		for _, f := range planDef.Features {
			features = append(features, "<li>"+html.EscapeString(f)+"</li>")
		}
		htmlBody := fmt.Sprintf(`<p>Hi %s,</p><p>Your <strong>%s</strong> workspace is now on the <strong>%s</strong> plan (%s).</p>`+
			`<p>What you get:</p><ul>%s</ul>`+
			`<p>Your current period runs until <strong>%s</strong>.</p><p>Thanks for using ThumbnailIQ!</p>`,
			html.EscapeString(owner.FullName), html.EscapeString(ws.Name), html.EscapeString(planDef.Name),
			html.EscapeString(price), strings.Join(features, ""), until)

		if err := u.mailer.Send(ctx, owner.Email, subject, text, htmlBody); err != nil {
			u.log.Error().Err(err).Str("workspace_id", workspaceID.String()).Str("plan", planDef.ID).Msg("sending subscription confirmation email")
		}
	}()
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
