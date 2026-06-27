-- +goose Up
ALTER TABLE subscriptions RENAME COLUMN stripe_subscription_id TO provider_subscription_id;
ALTER TABLE subscriptions RENAME COLUMN stripe_price_id TO provider_plan_id;
ALTER TABLE subscriptions ADD COLUMN provider VARCHAR(20) NOT NULL DEFAULT 'razorpay';

ALTER INDEX idx_subscriptions_stripe RENAME TO idx_subscriptions_provider_subscription;

-- +goose Down
ALTER INDEX idx_subscriptions_provider_subscription RENAME TO idx_subscriptions_stripe;

ALTER TABLE subscriptions DROP COLUMN provider;
ALTER TABLE subscriptions RENAME COLUMN provider_plan_id TO stripe_price_id;
ALTER TABLE subscriptions RENAME COLUMN provider_subscription_id TO stripe_subscription_id;
