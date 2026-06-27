package payment

import "context"

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
