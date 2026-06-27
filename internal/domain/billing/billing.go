package billing

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	PriceMonthly     float64  `json:"price_monthly"`
	Features         []string `json:"features"`
	AnalysesLimit    int      `json:"analyses_limit"`
	APIRequestsLimit int      `json:"api_requests_limit"`
}

var Plans = []Plan{
	{ID: "free", Name: "Free", PriceMonthly: 0, AnalysesLimit: 5, APIRequestsLimit: 100,
		Features: []string{"5 analyses/month", "Basic scoring", "Community support"}},
	{ID: "starter", Name: "Starter", PriceMonthly: 19, AnalysesLimit: 50, APIRequestsLimit: 1000,
		Features: []string{"50 analyses/month", "Competitor comparison", "Email support"}},
	{ID: "pro", Name: "Pro", PriceMonthly: 49, AnalysesLimit: 250, APIRequestsLimit: 10000,
		Features: []string{"250 analyses/month", "Historical tracking", "Priority support"}},
	{ID: "agency", Name: "Agency", PriceMonthly: 149, AnalysesLimit: 1000, APIRequestsLimit: 100000,
		Features: []string{"1000 analyses/month", "White-label reports", "Team collaboration"}},
}

type Subscription struct {
	ID                     uuid.UUID `json:"id"`
	WorkspaceID            uuid.UUID `json:"workspace_id"`
	Provider               string    `json:"provider"`
	ProviderSubscriptionID string    `json:"provider_subscription_id"`
	Plan                   string    `json:"plan"`
	Status                 string    `json:"status"`
	CurrentPeriodStart     time.Time `json:"current_period_start"`
	CurrentPeriodEnd       time.Time `json:"current_period_end"`
}

type APIKey struct {
	ID                uuid.UUID `json:"id"`
	WorkspaceID       uuid.UUID `json:"workspace_id"`
	Name              string    `json:"name"`
	KeyHash           string    `json:"-"`
	KeyPrefix         string    `json:"key_prefix"`
	RequestsThisMonth int       `json:"requests_this_month"`
	RequestsLimit     int       `json:"requests_limit"`
	CreatedAt         time.Time `json:"created_at"`
}

type Repository interface {
	UpsertSubscription(ctx context.Context, s *Subscription) (*Subscription, error)
	GetSubscriptionByWorkspace(ctx context.Context, workspaceID uuid.UUID) (*Subscription, error)

	CreateAPIKey(ctx context.Context, workspaceID uuid.UUID, name, keyHash, keyPrefix string, createdBy uuid.UUID) (*APIKey, error)
	ListAPIKeys(ctx context.Context, workspaceID uuid.UUID) ([]*APIKey, error)
	GetAPIKeyByHash(ctx context.Context, hash string) (*APIKey, error)
	RevokeAPIKey(ctx context.Context, id, workspaceID uuid.UUID) error
}
