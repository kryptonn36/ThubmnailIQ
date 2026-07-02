"use client";

import { useCallback, useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Check, Copy, CreditCard, KeyRound } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import type {
  BillingPlan,
  CheckoutResponse,
  ConfirmCheckoutResponse,
  CurrentSubscriptionResponse,
  APIKey,
  CreatedAPIKey,
} from "@/types";
import { Alert } from "@/components/ui/Alert";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { EmptyState } from "@/components/ui/EmptyState";
import { Input } from "@/components/ui/Input";
import { PageHeader } from "@/components/ui/PageHeader";
import { Skeleton } from "@/components/ui/Skeleton";
import { useMotionVariants } from "@/lib/motion";

declare global {
  interface Window {
    Razorpay: new (options: Record<string, unknown>) => { open: () => void };
  }
}

const RAZORPAY_CHECKOUT_SRC = "https://checkout.razorpay.com/v1/checkout.js";

function loadRazorpayScript(): Promise<void> {
  return new Promise((resolve, reject) => {
    if (window.Razorpay) {
      resolve();
      return;
    }
    const script = document.createElement("script");
    script.src = RAZORPAY_CHECKOUT_SRC;
    script.onload = () => resolve();
    script.onerror = () => reject(new Error("Failed to load Razorpay checkout"));
    document.body.appendChild(script);
  });
}

export default function BillingPage() {
  const { fadeInUp, staggerContainer } = useMotionVariants();
  const [plans, setPlans] = useState<BillingPlan[]>([]);
  const [currentSub, setCurrentSub] = useState<CurrentSubscriptionResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [subscribing, setSubscribing] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  // API Keys
  const [apiKeys, setApiKeys] = useState<APIKey[]>([]);
  const [newKeyName, setNewKeyName] = useState("");
  const [creatingKey, setCreatingKey] = useState(false);
  const [createdKey, setCreatedKey] = useState<CreatedAPIKey | null>(null);
  const [keyCopied, setKeyCopied] = useState(false);
  const [keyError, setKeyError] = useState<string | null>(null);

  const loadApiKeys = useCallback(async () => {
    try {
      const keys = await api.get<APIKey[]>("/api-keys");
      setApiKeys(keys);
    } catch {
      // non-fatal, don't block the whole page
    }
  }, []);

  async function handleCreateKey(e: React.FormEvent) {
    e.preventDefault();
    if (!newKeyName.trim()) return;
    setCreatingKey(true);
    setKeyError(null);
    setCreatedKey(null);
    try {
      const result = await api.post<CreatedAPIKey>("/api-keys", { name: newKeyName.trim() });
      setCreatedKey(result);
      setNewKeyName("");
      await loadApiKeys();
    } catch (err) {
      setKeyError(err instanceof ApiError ? err.message : "Failed to create API key.");
    } finally {
      setCreatingKey(false);
    }
  }

  async function handleRevokeKey(id: string) {
    try {
      await api.del(`/api-keys/${id}`);
      setApiKeys((prev) => prev.filter((k) => k.id !== id));
      if (createdKey?.id === id) setCreatedKey(null);
    } catch {
      // ignore
    }
  }

  function copyKey(key: string) {
    navigator.clipboard.writeText(key);
    setKeyCopied(true);
    setTimeout(() => setKeyCopied(false), 2000);
  }

  useEffect(() => {
    let cancelled = false;
    async function load() {
      try {
        const [plansData, subData] = await Promise.all([
          api.get<BillingPlan[]>("/billing/plans"),
          api.get<CurrentSubscriptionResponse>("/billing/subscription"),
        ]);
        if (cancelled) return;
        setPlans(plansData);
        setCurrentSub(subData);
        loadApiKeys();
      } catch (err) {
        if (cancelled) return;
        setError(err instanceof ApiError ? err.message : "Failed to load billing info.");
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    load();
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const currentPlanDef = plans.find((p) => p.id === currentSub?.plan);

  async function handleSubscribe(planId: string) {
    setSubscribing(planId);
    setMessage(null);
    setError(null);
    try {
      const checkout = await api.post<CheckoutResponse>("/billing/checkout", { plan: planId });

      if (!checkout.requires_payment) {
        setCurrentSub({ plan: checkout.plan, status: checkout.status ?? "active", is_active: true });
        setMessage(`Successfully switched to ${checkout.plan}.`);
        setSubscribing(null);
        return;
      }

      await loadRazorpayScript();
      const razorpay = new window.Razorpay({
        key: checkout.key_id,
        amount: checkout.amount,
        currency: checkout.currency,
        order_id: checkout.order_id,
        name: "ThumbnailIQ",
        description: `${checkout.plan} plan subscription`,
        handler: async (response: {
          razorpay_order_id: string;
          razorpay_payment_id: string;
          razorpay_signature: string;
        }) => {
          try {
            const confirmed = await api.post<ConfirmCheckoutResponse>("/billing/checkout/verify", {
              plan: planId,
              order_id: response.razorpay_order_id,
              payment_id: response.razorpay_payment_id,
              signature: response.razorpay_signature,
            });
            setCurrentSub({ plan: confirmed.plan, status: confirmed.status, is_active: true });
            setMessage(`Successfully subscribed to ${confirmed.plan}.`);
          } catch (err) {
            setError(err instanceof ApiError ? err.message : "Payment succeeded but activation failed.");
          } finally {
            setSubscribing(null);
          }
        },
        modal: {
          ondismiss: () => setSubscribing(null),
        },
      });
      razorpay.open();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to subscribe.");
      setSubscribing(null);
    }
  }

  return (
    <div className="space-y-6">
      <PageHeader title="Billing" description="Choose the plan that fits your workspace." />

      {currentSub && (
        <Card className="text-sm text-gray-300">
          You&apos;re on the <span className="font-semibold text-white">{currentPlanDef?.name ?? currentSub.plan}</span> plan
          {currentSub.is_active && currentSub.current_period_end && currentSub.plan !== "free" && (
            <> &middot; valid until {new Date(currentSub.current_period_end).toLocaleDateString()}</>
          )}
        </Card>
      )}

      {message && <Alert variant="success">{message}</Alert>}
      {error && <Alert variant="danger">{error}</Alert>}

      {loading ? (
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {[...Array(4)].map((_, i) => (
            <Card key={i} className="space-y-3">
              <Skeleton className="h-5 w-1/2" />
              <Skeleton className="h-8 w-1/3" />
              <Skeleton className="h-3 w-full" />
              <Skeleton className="h-3 w-full" />
              <Skeleton className="h-10 w-full" />
            </Card>
          ))}
        </div>
      ) : (
        <motion.div
          initial="hidden"
          animate="visible"
          variants={staggerContainer}
          className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4"
        >
          {plans.map((plan) => {
            const isCurrent = currentSub?.is_active && currentSub.plan === plan.id;
            const isUpgrade =
              currentSub?.is_active &&
              !isCurrent &&
              currentPlanDef &&
              plan.price_monthly > currentPlanDef.price_monthly;
            const isDowngrade =
              currentSub?.is_active &&
              !isCurrent &&
              currentPlanDef &&
              plan.price_monthly <= currentPlanDef.price_monthly;
            const upgradePrice = isUpgrade ? plan.price_monthly - currentPlanDef!.price_monthly : null;

            return (
              <motion.div key={plan.id} variants={fadeInUp}>
                <Card
                  accent={isCurrent}
                  className={`flex h-full flex-col ${isCurrent ? "border-brand-500 shadow-glow" : ""}`}
                >
                  <h2 className="text-lg font-semibold text-white">{plan.name}</h2>
                  <p className="mt-2 text-3xl font-bold tracking-tight text-white">
                    ${plan.price_monthly}
                    <span className="text-sm font-normal text-gray-500">/mo</span>
                  </p>
                  <ul className="mt-4 flex-1 space-y-2 text-sm text-gray-400">
                    <li>{plan.analyses_limit} analyses / month</li>
                    <li>{plan.api_requests_limit} API requests / month</li>
                    {plan.features.map((feature) => (
                      <li key={feature}>{feature}</li>
                    ))}
                  </ul>
                  <Button
                    variant={isCurrent || isDowngrade ? "secondary" : "primary"}
                    onClick={() => handleSubscribe(plan.id)}
                    disabled={isCurrent || isDowngrade}
                    loading={subscribing === plan.id}
                    title={isDowngrade ? "Downgrades aren't available while your current plan is active" : undefined}
                    className="mt-6 w-full"
                  >
                    {isCurrent
                      ? "Current Plan"
                      : isDowngrade
                        ? "Unavailable"
                        : isUpgrade
                          ? `Upgrade for $${upgradePrice}/mo`
                          : "Subscribe"}
                  </Button>
                </Card>
              </motion.div>
            );
          })}
        </motion.div>
      )}

      {/* ── API Keys ── */}
      <div className="space-y-4">
        <div>
          <h2 className="text-xl font-bold tracking-tight text-white">API Keys</h2>
          <p className="mt-1 text-sm text-gray-400">
            Use an API key to authenticate the browser extension or any other integration. The raw key is shown only once — copy it immediately.
          </p>
        </div>

        <Card as="form" onSubmit={handleCreateKey}>
          <h3 className="mb-3 text-sm font-semibold text-white">Create a new key</h3>
          <div className="flex gap-3">
            <Input
              type="text"
              value={newKeyName}
              onChange={(e) => setNewKeyName(e.target.value)}
              placeholder="e.g. Chrome Extension"
              className="flex-1"
            />
            <Button type="submit" loading={creatingKey} disabled={!newKeyName.trim()} icon={<KeyRound className="h-4 w-4" aria-hidden="true" />}>
              Create Key
            </Button>
          </div>
          {keyError && <p className="mt-2 text-xs text-danger">{keyError}</p>}
        </Card>

        {createdKey && (
          <Card className="border-success/40 bg-success/10">
            <p className="mb-1 flex items-center gap-1.5 text-sm font-semibold text-success">
              <Check className="h-4 w-4" aria-hidden="true" />
              Key created — copy it now, it won&apos;t be shown again
            </p>
            <div className="mt-2 flex items-center gap-2">
              <code className="flex-1 break-all rounded-lg border border-success/30 bg-black/30 px-3 py-2 font-mono text-xs text-success">
                {createdKey.key}
              </code>
              <Button
                type="button"
                variant="secondary"
                size="sm"
                onClick={() => copyKey(createdKey.key)}
                icon={<Copy className="h-3.5 w-3.5" aria-hidden="true" />}
              >
                {keyCopied ? "Copied" : "Copy"}
              </Button>
            </div>
          </Card>
        )}

        {apiKeys.length > 0 ? (
          <div className="space-y-2">
            {apiKeys.map((k) => (
              <div
                key={k.id}
                className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-surface-300 bg-surface-100 px-4 py-3"
              >
                <div>
                  <p className="text-sm font-medium text-white">{k.name}</p>
                  <p className="mt-0.5 font-mono text-xs text-gray-500">{k.key_prefix}••••••••</p>
                </div>
                <div className="flex items-center gap-4">
                  <div className="text-right text-xs text-gray-500">
                    <p>{k.requests_this_month} / {k.requests_limit} requests</p>
                    <p>Created {new Date(k.created_at).toLocaleDateString()}</p>
                  </div>
                  <Button variant="danger" size="sm" onClick={() => handleRevokeKey(k.id)}>
                    Revoke
                  </Button>
                </div>
              </div>
            ))}
          </div>
        ) : (
          !loading && (
            <EmptyState
              icon={<CreditCard className="h-5 w-5" aria-hidden="true" />}
              title="No API keys yet."
              description="Create one above to connect the browser extension or external tools."
            />
          )
        )}
      </div>
    </div>
  );
}
