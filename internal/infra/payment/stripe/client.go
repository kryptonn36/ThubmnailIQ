// Package stripe implements the payment.Gateway interface for Stripe.
//
// Stripe's account is currently invitation-gated, so this is a stub: it
// mints a fake order/payment ID and verification always succeeds, mirroring
// the shape a real implementation would have. Once the Stripe account is
// approved, swap the body of CreateOrder/VerifyPayment for real calls
// against github.com/stripe/stripe-go (Checkout Session + webhook signature
// verification) — the Gateway interface and call sites in the billing
// usecase don't need to change. Until then, set PAYMENT_PROVIDER=razorpay.
package stripe

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/payment"
)

type Client struct {
	secretKey string
}

func NewClient(secretKey string) *Client {
	return &Client{secretKey: secretKey}
}

func (c *Client) Name() string {
	return "stripe"
}

// IsLive reports whether a real Stripe secret key is configured.
func (c *Client) IsLive() bool {
	return c.secretKey != ""
}

func (c *Client) CreateOrder(ctx context.Context, p payment.CreateOrderParams) (*payment.Order, error) {
	return &payment.Order{
		ID:          fmt.Sprintf("cs_mock_%s", uuid.New().String()[:12]),
		AmountMinor: p.AmountMinor,
		Currency:    p.Currency,
		Provider:    c.Name(),
		KeyID:       c.secretKey,
	}, nil
}

func (c *Client) VerifyPayment(ctx context.Context, v payment.VerifyParams) error {
	return nil
}
