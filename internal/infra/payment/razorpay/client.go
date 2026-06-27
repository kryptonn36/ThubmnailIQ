// Package razorpay implements the payment.Gateway interface against
// Razorpay's Orders API. It's the active gateway while ThumbnailIQ waits on
// Stripe's invitation-only account approval (see internal/infra/payment/stripe).
package razorpay

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	razorpaygo "github.com/razorpay/razorpay-go"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/payment"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
)

type Client struct {
	rzp       *razorpaygo.Client
	keyID     string
	keySecret string
}

func NewClient(keyID, keySecret string) *Client {
	return &Client{rzp: razorpaygo.NewClient(keyID, keySecret), keyID: keyID, keySecret: keySecret}
}

func (c *Client) Name() string {
	return "razorpay"
}

// IsLive reports whether real Razorpay API keys are configured.
func (c *Client) IsLive() bool {
	return c.keyID != "" && c.keySecret != ""
}

func (c *Client) CreateOrder(ctx context.Context, p payment.CreateOrderParams) (*payment.Order, error) {
	resp, err := c.rzp.Order.Create(map[string]interface{}{
		"amount":   p.AmountMinor,
		"currency": p.Currency,
		"receipt":  p.Receipt,
		"notes":    map[string]interface{}{"plan": p.Plan},
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("razorpay: create order: %w", err)
	}
	orderID, _ := resp["id"].(string)
	if orderID == "" {
		return nil, fmt.Errorf("razorpay: create order: missing id in response")
	}
	return &payment.Order{
		ID:          orderID,
		AmountMinor: p.AmountMinor,
		Currency:    p.Currency,
		Provider:    c.Name(),
		KeyID:       c.keyID,
	}, nil
}

// VerifyPayment checks the HMAC-SHA256 signature Razorpay Checkout returns
// alongside a successful payment: signature = HMAC(order_id|payment_id, key_secret).
// See https://razorpay.com/docs/payments/payment-gateway/web-integration/standard/build-integration/#step-5-verify-payment-signature
func (c *Client) VerifyPayment(ctx context.Context, v payment.VerifyParams) error {
	if v.OrderID == "" || v.PaymentID == "" || v.Signature == "" {
		return apperrors.ErrInvalidInput
	}
	payload := v.OrderID + "|" + v.PaymentID
	mac := hmac.New(sha256.New, []byte(c.keySecret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))

	sig, err := hex.DecodeString(v.Signature)
	if err != nil {
		return apperrors.ErrInvalidInput
	}
	expectedBytes, err := hex.DecodeString(expected)
	if err != nil {
		return err
	}
	if !hmac.Equal(sig, expectedBytes) {
		return apperrors.ErrInvalidInput
	}
	return nil
}
