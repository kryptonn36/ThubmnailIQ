"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api, ApiError } from "@/lib/api";
import { getScoreBand } from "@/components/ScoreGauge";
import type { AnalysisSummary, PaginatedResponse, Workspace } from "@/types";

const QUICK_ACTIONS = [
  {
    href: "/dashboard/analyses/new",
    icon: "🖼️",
    title: "Analyze a Thumbnail",
    desc: "Upload and score a new thumbnail against live competitors",
    cta: "Start analysis",
    accent: "brand",
  },
  {
    href: "/dashboard/ideas",
    icon: "💡",
    title: "Generate Video Ideas",
    desc: "AI-powered ideas based on what's working in your niche",
    cta: "Get ideas",
    accent: "emerald",
  },
  {
    href: "/dashboard/tracking",
    icon: "📡",
    title: "Track Competitors",
    desc: "Monitor how rival thumbnails change over time",
    cta: "Set up tracking",
    accent: "amber",
  },
  {
    href: "/dashboard/database",
    icon: "🔥",
    title: "Viral Thumbnail DB",
    desc: "Browse top-performing thumbnails across all niches",
    cta: "Browse database",
    accent: "rose",
  },
];

const ACCENT_CLASSES: Record<string, string> = {
  brand: "border-brand-500/30 hover:border-brand-500/60",
  emerald: "border-emerald-500/30 hover:border-emerald-500/60",
  amber: "border-amber-500/30 hover:border-amber-500/60",
  rose: "border-rose-500/30 hover:border-rose-500/60",
};

export default function DashboardPage() {
  const [workspace, setWorkspace] = useState<Workspace | null>(null);
  const [analyses, setAnalyses] = useState<AnalysisSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    async function load() {
      try {
        const workspaces = await api.get<Workspace[]>("/workspaces");
        const ws = workspaces[0] ?? null;
        if (cancelled) return;
        setWorkspace(ws);

        const res = await api.get<PaginatedResponse<AnalysisSummary>>("/analyses", {
          workspace_id: ws?.id,
          page: 1,
          per_page: 8,
        });
        if (cancelled) return;
        setAnalyses(res.data);
      } catch (err) {
        if (cancelled) return;
        setError(err instanceof ApiError ? err.message : "Failed to load dashboard data");
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    load();
    return () => { cancelled = true; };
  }, []);

  const completedAnalyses = analyses.filter((a) => a.status === "complete" && a.score !== null);
  const avgScore = completedAnalyses.length
    ? Math.round(completedAnalyses.reduce((s, a) => s + (a.score ?? 0), 0) / completedAnalyses.length)
    : null;
  const bestScore = completedAnalyses.length
    ? Math.max(...completedAnalyses.map((a) => a.score ?? 0))
    : null;

  const usagePct = workspace
    ? Math.min(100, (workspace.analyses_this_month / Math.max(workspace.analyses_limit, 1)) * 100)
    : 0;

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white">
            {workspace ? `${workspace.name}` : "Overview"}
          </h1>
          <p className="mt-1 text-sm text-gray-400">
            {workspace
              ? `${workspace.plan} plan · ${workspace.analyses_this_month} of ${workspace.analyses_limit} analyses used this month`
              : "Your ThumbnailIQ workspace"}
          </p>
        </div>
        <Link
          href="/dashboard/analyses/new"
          className="rounded-lg bg-brand-gradient px-4 py-2.5 text-sm font-semibold text-white shadow-glow transition hover:opacity-90"
        >
          + New Analysis
        </Link>
      </div>

      {error && (
        <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">{error}</p>
      )}

      {/* Stats row */}
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
        {[
          { label: "Analyses this month", value: workspace?.analyses_this_month ?? "—", sub: `of ${workspace?.analyses_limit ?? "—"} limit` },
          { label: "Average score", value: avgScore !== null ? avgScore : "—", sub: "across completed analyses" },
          { label: "Best score", value: bestScore !== null ? bestScore : "—", sub: "highest scoring thumbnail" },
          { label: "Current plan", value: workspace ? workspace.plan.charAt(0).toUpperCase() + workspace.plan.slice(1) : "—", sub: <Link href="/dashboard/billing" className="text-brand-300 hover:underline">Upgrade plan</Link> },
        ].map((s, i) => (
          <div key={i} className="rounded-2xl border border-surface-300 bg-surface-100 p-5">
            <p className="text-xs font-medium uppercase tracking-wide text-gray-500">{s.label}</p>
            <p className="mt-2 text-3xl font-bold text-white">{s.value}</p>
            <p className="mt-1 text-xs text-gray-500">{s.sub}</p>
          </div>
        ))}
      </div>

      {/* Usage bar */}
      {workspace && (
        <div className="rounded-2xl border border-surface-300 bg-surface-100 p-5">
          <div className="mb-2 flex items-center justify-between text-sm">
            <span className="text-gray-300">Monthly usage</span>
            <span className={`font-medium ${usagePct > 80 ? "text-amber-400" : "text-white"}`}>
              {workspace.analyses_this_month} / {workspace.analyses_limit}
            </span>
          </div>
          <div className="h-2 overflow-hidden rounded-full bg-surface-300">
            <div
              className={`h-full transition-all ${usagePct > 80 ? "bg-amber-500" : "bg-brand-gradient"}`}
              style={{ width: `${usagePct}%` }}
            />
          </div>
          {usagePct > 80 && (
            <p className="mt-2 text-xs text-amber-400">
              Running low — <Link href="/dashboard/billing" className="underline">upgrade your plan</Link>
            </p>
          )}
        </div>
      )}

      {/* Quick actions */}
      <div>
        <h2 className="mb-4 text-base font-semibold text-white">Tools</h2>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {QUICK_ACTIONS.map((action) => (
            <Link
              key={action.href}
              href={action.href}
              className={`group flex flex-col rounded-2xl border bg-surface-100 p-5 transition ${ACCENT_CLASSES[action.accent]}`}
            >
              <span className="mb-3 text-2xl">{action.icon}</span>
              <p className="text-sm font-semibold text-white">{action.title}</p>
              <p className="mt-1 flex-1 text-xs text-gray-500">{action.desc}</p>
              <p className="mt-4 text-xs font-medium text-brand-300 group-hover:text-brand-200">
                {action.cta} →
              </p>
            </Link>
          ))}
        </div>
      </div>

      {/* Recent analyses */}
      <div>
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-base font-semibold text-white">Recent Analyses</h2>
          <Link href="/dashboard/analyses" className="text-sm text-brand-300 hover:text-brand-200">
            View all →
          </Link>
        </div>

        {loading ? (
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="animate-pulse rounded-xl border border-surface-300 bg-surface-100 p-3">
                <div className="aspect-video w-full rounded-lg bg-surface-300 mb-2" />
                <div className="h-3 w-3/4 rounded bg-surface-300 mb-1" />
                <div className="h-3 w-1/2 rounded bg-surface-300" />
              </div>
            ))}
          </div>
        ) : analyses.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-surface-300 p-14 text-center">
            <p className="text-4xl mb-3">🖼️</p>
            <p className="font-medium text-gray-300">No analyses yet</p>
            <p className="mt-1 text-sm text-gray-500">Upload a thumbnail to see how it compares to your niche.</p>
            <Link
              href="/dashboard/analyses/new"
              className="mt-4 inline-block rounded-lg bg-brand-gradient px-5 py-2.5 text-sm font-semibold text-white shadow-glow hover:opacity-90"
            >
              Create First Analysis
            </Link>
          </div>
        ) : (
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
            {analyses.map((a) => {
              const band = typeof a.score === "number" ? getScoreBand(a.score) : null;
              return (
                <Link
                  key={a.id}
                  href={`/dashboard/analyses/${a.id}`}
                  className="group rounded-xl border border-surface-300 bg-surface-100 p-3 transition hover:border-brand-500/50"
                >
                  <div className="relative mb-2 overflow-hidden rounded-lg bg-surface-300">
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                      src={a.thumbnail_url}
                      alt={a.keyword}
                      className="aspect-video w-full object-cover transition group-hover:scale-105"
                    />
                    {band && (
                      <span className={`absolute right-1.5 top-1.5 rounded-md px-1.5 py-0.5 text-xs font-bold ${band.bgClass} ${band.textClass}`}>
                        {a.score}
                      </span>
                    )}
                    {a.status === "pending" || a.status === "processing" ? (
                      <span className="absolute right-1.5 top-1.5 rounded-md bg-amber-500/90 px-1.5 py-0.5 text-xs font-bold text-white">
                        ⏳
                      </span>
                    ) : null}
                  </div>
                  <p className="truncate text-sm font-medium text-gray-200">{a.keyword}</p>
                  <p className="mt-0.5 text-xs capitalize text-gray-500">{a.status}</p>
                </Link>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
