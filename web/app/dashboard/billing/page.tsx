"use client";

import { useEffect, useState } from "react";
import { api, ApiError } from "@/lib/api";
import type { BillingPlan } from "@/types";

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
      const res = await api.post<{ plan: string; status: string }>("/billing/subscribe", {
        plan: planId,
      });
      setActivePlan(res.plan);
      setMessage(`Successfully subscribed to ${res.plan}.`);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to subscribe.");
    } finally {
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
