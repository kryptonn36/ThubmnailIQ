package payment

import (
	"context"

	"github.com/google/uuid"
)

// PendingOrder records what a checkout order was created for. It is persisted
// when an order is opened and read back at confirmation time, so the plan and
// amount that get activated come from the server's record of the order — never
// from client-supplied values that could be tampered with.
type PendingOrder struct {
	WorkspaceID uuid.UUID
	Plan        string
	AmountMinor int64
}

// PendingOrderStore persists PendingOrders between checkout creation and
// confirmation. Consume must atomically fetch-and-delete the record so a single
// payment cannot be replayed to activate a plan more than once.
type PendingOrderStore interface {
	Save(ctx context.Context, orderID string, order PendingOrder) error
	Consume(ctx context.Context, orderID string) (*PendingOrder, error)
}

// CreateOrderParams describes a one-time charge for a billing cycle.
type CreateOrderParams struct {
	Plan        string
	AmountMinor int64 // amount in the smallest currency unit (e.g. paise, cents)
	Currency    string
	Receipt     string
}

// Order is the provider-side checkout object the frontend needs to open
// the provider's payment widget.
type Order struct {
	ID          string
	AmountMinor int64
	Currency    string
	Provider    string
	KeyID       string // public/publishable key the frontend needs to open checkout
}

// VerifyParams carries whatever the provider's checkout widget hands back
// after a successful payment, so the backend can confirm it's genuine.
type VerifyParams struct {
	OrderID   string
	PaymentID string
	Signature string
}

// Gateway is implemented by each payment provider. ThumbnailIQ runs on
// Razorpay today (internal/infra/payment/razorpay) because Stripe access is
// invite-gated; once Stripe approves the account, switching is a config
// change (PAYMENT_PROVIDER=stripe) to the Gateway implementation in
// internal/infra/payment/stripe — no usecase/handler changes needed.
type Gateway interface {
	Name() string
	CreateOrder(ctx context.Context, p CreateOrderParams) (*Order, error)
	VerifyPayment(ctx context.Context, v VerifyParams) error
}
