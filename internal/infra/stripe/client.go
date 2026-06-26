package stripe

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Client is a stub billing client. Real Stripe checkout/webhook handling
// requires a live secret key and webhook secret which aren't available in
// this dev environment, so subscriptions are simulated: a fake Stripe
// subscription ID is minted and the workspace plan is updated immediately.
// Swap this implementation for github.com/stripe/stripe-go once real keys
// are available; the interface (CreateSubscription) would stay the same.
type Client struct {
	secretKey string
}

func NewClient(secretKey string) *Client {
	return &Client{secretKey: secretKey}
}

func (c *Client) IsLive() bool {
	return c.secretKey != ""
}

func (c *Client) CreateSubscription(plan string) (subscriptionID string, periodStart, periodEnd time.Time, err error) {
	subscriptionID = fmt.Sprintf("sub_mock_%s", uuid.New().String()[:12])
	periodStart = time.Now()
	periodEnd = periodStart.AddDate(0, 1, 0)
	return subscriptionID, periodStart, periodEnd, nil
}
