package billing

import (
	"context"

	"github.com/google/uuid"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/billing"
	"github.com/thumbnailiq/thumbnailiq/internal/domain/workspace"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/stripe"
	"github.com/thumbnailiq/thumbnailiq/pkg/errors"
	"github.com/thumbnailiq/thumbnailiq/pkg/hash"
)

type Usecase struct {
	billing    billing.Repository
	workspaces workspace.Repository
	stripe     *stripe.Client
}

func NewUsecase(billingRepo billing.Repository, workspaces workspace.Repository, stripeClient *stripe.Client) *Usecase {
	return &Usecase{billing: billingRepo, workspaces: workspaces, stripe: stripeClient}
}

func (u *Usecase) Plans() []billing.Plan {
	return billing.Plans
}

func (u *Usecase) Subscribe(ctx context.Context, workspaceID uuid.UUID, plan string) (*billing.Subscription, error) {
	var planDef *billing.Plan
	for i := range billing.Plans {
		if billing.Plans[i].ID == plan {
			planDef = &billing.Plans[i]
			break
		}
	}
	if planDef == nil {
		return nil, errors.ErrInvalidInput
	}

	subID, periodStart, periodEnd, err := u.stripe.CreateSubscription(plan)
	if err != nil {
		return nil, err
	}

	sub, err := u.billing.UpsertSubscription(ctx, &billing.Subscription{
		WorkspaceID: workspaceID, StripeSubscriptionID: subID, Plan: plan, Status: "active",
		CurrentPeriodStart: periodStart, CurrentPeriodEnd: periodEnd,
	})
	if err != nil {
		return nil, err
	}

	if _, err := u.workspaces.UpdatePlan(ctx, workspaceID, plan, planDef.AnalysesLimit); err != nil {
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
