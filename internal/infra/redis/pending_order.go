package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/thumbnailiq/thumbnailiq/internal/domain/payment"
	apperrors "github.com/thumbnailiq/thumbnailiq/pkg/errors"
)

// PendingOrderStore persists checkout orders in Redis for the short window
// between creating the order and confirming payment. Entries carry a TTL so an
// abandoned checkout simply expires; if the record is gone at confirmation time
// the confirmation fails closed rather than trusting client input.
type PendingOrderStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewPendingOrderStore(client *redis.Client) *PendingOrderStore {
	return &PendingOrderStore{client: client, ttl: time.Hour}
}

type pendingOrderRecord struct {
	WorkspaceID string `json:"workspace_id"`
	Plan        string `json:"plan"`
	AmountMinor int64  `json:"amount_minor"`
}

func pendingOrderKey(orderID string) string {
	return "pending_order:" + orderID
}

func (s *PendingOrderStore) Save(ctx context.Context, orderID string, order payment.PendingOrder) error {
	data, err := json.Marshal(pendingOrderRecord{
		WorkspaceID: order.WorkspaceID.String(),
		Plan:        order.Plan,
		AmountMinor: order.AmountMinor,
	})
	if err != nil {
		return err
	}
	return s.client.Set(ctx, pendingOrderKey(orderID), data, s.ttl).Err()
}

// Consume atomically reads and deletes the stored order (GetDel), so a payment
// can't be replayed against the same order to activate a plan twice. A missing
// record returns ErrNotFound.
func (s *PendingOrderStore) Consume(ctx context.Context, orderID string) (*payment.PendingOrder, error) {
	val, err := s.client.GetDel(ctx, pendingOrderKey(orderID)).Result()
	if err == redis.Nil {
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var rec pendingOrderRecord
	if err := json.Unmarshal([]byte(val), &rec); err != nil {
		return nil, err
	}
	workspaceID, err := uuid.Parse(rec.WorkspaceID)
	if err != nil {
		return nil, err
	}
	return &payment.PendingOrder{
		WorkspaceID: workspaceID,
		Plan:        rec.Plan,
		AmountMinor: rec.AmountMinor,
	}, nil
}
