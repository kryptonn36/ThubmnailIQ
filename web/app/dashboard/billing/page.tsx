"use client";

import { useEffect, useState } from "react";
import { api, ApiError } from "@/lib/api";
import type { BillingPlan, CheckoutResponse, ConfirmCheckoutResponse } from "@/types";

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
  const [plans, setPlans] = useState<BillingPlan[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [subscribing, setSubscribing] = useState<string | null>(null);
  const [activePlan, setActivePlan] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    async function load() {
      try {
        const data = await api.get<BillingPlan[]>("/billing/plans");
        if (cancelled) return;
        setPlans(data);
      } catch (err) {
        if (cancelled) return;
        setError(err instanceof ApiError ? err.message : "Failed to load plans.");
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    load();
    return () => {
      cancelled = true;
    };
  }, []);

  async function handleSubscribe(planId: string) {
    setSubscribing(planId);
    setMessage(null);
    setError(null);
    try {
      const checkout = await api.post<CheckoutResponse>("/billing/checkout", { plan: planId });

      if (!checkout.requires_payment) {
        setActivePlan(checkout.plan);
        setMessage(`Successfully subscribed to ${checkout.plan}.`);
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
            setActivePlan(confirmed.plan);
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
      <div>
        <h1 className="text-2xl font-bold text-white">Billing</h1>
        <p className="mt-1 text-sm text-gray-400">
          Choose the plan that fits your workspace.
        </p>
      </div>

      {message && (
        <p className="rounded-lg bg-emerald-500/10 px-4 py-3 text-sm text-emerald-400">
          {message}
        </p>
      )}
      {error && (
        <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">{error}</p>
      )}

      {loading ? (
        <p className="text-sm text-gray-500">Loading plans...</p>
      ) : (
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {plans.map((plan) => {
            const isActive = activePlan === plan.id;
            return (
              <div
                key={plan.id}
                className={`flex flex-col rounded-2xl border p-6 ${
                  isActive
                    ? "border-brand-500 shadow-glow"
                    : "border-surface-300 bg-surface-100"
                }`}
              >
                <h2 className="text-lg font-semibold text-white">{plan.name}</h2>
                <p className="mt-2 text-3xl font-bold text-white">
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
                <button
                  onClick={() => handleSubscribe(plan.id)}
                  disabled={subscribing === plan.id}
                  className={`mt-6 rounded-lg px-4 py-2.5 text-sm font-semibold transition disabled:opacity-60 ${
                    isActive
                      ? "bg-surface-300 text-gray-300"
                      : "bg-brand-gradient text-white shadow-glow hover:opacity-90"
                  }`}
                >
                  {subscribing === plan.id
                    ? "Subscribing..."
                    : isActive
                    ? "Current Plan"
                    : "Subscribe"}
                </button>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
