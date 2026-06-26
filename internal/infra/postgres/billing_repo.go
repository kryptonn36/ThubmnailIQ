package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/billing"
	"github.com/thumbnailiq/thumbnailiq/internal/infra/postgres/db"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
)

type BillingRepo struct {
	q *db.Queries
}

func NewBillingRepo(pool *pgxpool.Pool) *BillingRepo {
	return &BillingRepo{q: db.New(pool)}
}

func (r *BillingRepo) UpsertSubscription(ctx context.Context, s *billing.Subscription) (*billing.Subscription, error) {
	created, err := r.q.UpsertSubscription(ctx, db.UpsertSubscriptionParams{
		WorkspaceID:          s.WorkspaceID,
		StripeSubscriptionID: s.StripeSubscriptionID,
		StripePriceID:        "price_" + s.Plan,
		Plan:                 s.Plan,
		Status:               s.Status,
		CurrentPeriodStart:   tsNow(s.CurrentPeriodStart),
		CurrentPeriodEnd:     tsNow(s.CurrentPeriodEnd),
	})
	if err != nil {
		return nil, err
	}
	return toDomainSubscription(created), nil
}

func toDomainSubscription(s db.Subscription) *billing.Subscription {
	return &billing.Subscription{
		ID: s.ID, WorkspaceID: s.WorkspaceID, StripeSubscriptionID: s.StripeSubscriptionID,
		Plan: s.Plan, Status: s.Status,
		CurrentPeriodStart: tsVal(s.CurrentPeriodStart), CurrentPeriodEnd: tsVal(s.CurrentPeriodEnd),
	}
}

func (r *BillingRepo) GetSubscriptionByWorkspace(ctx context.Context, workspaceID uuid.UUID) (*billing.Subscription, error) {
	s, err := r.q.GetSubscriptionByWorkspace(ctx, workspaceID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toDomainSubscription(s), nil
}

func (r *BillingRepo) CreateAPIKey(ctx context.Context, workspaceID uuid.UUID, name, keyHash, keyPrefix string, createdBy uuid.UUID) (*billing.APIKey, error) {
	k, err := r.q.CreateApiKey(ctx, db.CreateApiKeyParams{
		WorkspaceID: workspaceID, Name: name, KeyHash: keyHash, KeyPrefix: keyPrefix, CreatedBy: createdBy,
	})
	if err != nil {
		return nil, err
	}
	return toDomainAPIKey(k), nil
}

func toDomainAPIKey(k db.ApiKey) *billing.APIKey {
	return &billing.APIKey{
		ID: k.ID, WorkspaceID: k.WorkspaceID, Name: k.Name, KeyHash: k.KeyHash, KeyPrefix: k.KeyPrefix,
		RequestsThisMonth: int(k.RequestsThisMonth), RequestsLimit: int(k.RequestsLimit), CreatedAt: tsVal(k.CreatedAt),
	}
}

func (r *BillingRepo) ListAPIKeys(ctx context.Context, workspaceID uuid.UUID) ([]*billing.APIKey, error) {
	rows, err := r.q.ListApiKeysByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	out := make([]*billing.APIKey, 0, len(rows))
	for _, k := range rows {
		out = append(out, toDomainAPIKey(k))
	}
	return out, nil
}

func (r *BillingRepo) GetAPIKeyByHash(ctx context.Context, hash string) (*billing.APIKey, error) {
	k, err := r.q.GetApiKeyByHash(ctx, hash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return toDomainAPIKey(k), nil
}

func (r *BillingRepo) RevokeAPIKey(ctx context.Context, id, workspaceID uuid.UUID) error {
	return r.q.RevokeApiKey(ctx, db.RevokeApiKeyParams{ID: id, WorkspaceID: workspaceID})
}
