"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api, ApiError } from "@/lib/api";
import { getScoreBand } from "@/components/ScoreGauge";
import type { AnalysisSummary, PaginatedResponse, Workspace } from "@/types";

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
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <div className="space-y-8">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white">Overview</h1>
          <p className="mt-1 text-sm text-gray-400">
            Your recent thumbnail analyses and workspace usage.
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

      <div className="rounded-2xl border border-surface-300 bg-surface-100 p-6">
        <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-gray-400">
          Workspace Usage
        </h2>
        {workspace ? (
          <>
            <div className="flex items-end justify-between">
              <div>
                <p className="text-2xl font-bold text-white">
                  {workspace.analyses_this_month}
                  <span className="text-base font-normal text-gray-500">
                    {" "}
                    / {workspace.analyses_limit} analyses this month
                  </span>
                </p>
                <p className="mt-1 text-sm text-gray-500">
                  {workspace.name} &middot; {workspace.plan} plan
                </p>
              </div>
            </div>
            <div className="mt-4 h-2 w-full overflow-hidden rounded-full bg-surface-300">
              <div
                className="h-full bg-brand-gradient"
                style={{
                  width: `${Math.min(
                    100,
                    (workspace.analyses_this_month / Math.max(workspace.analyses_limit, 1)) * 100
                  )}%`,
                }}
              />
            </div>
          </>
        ) : (
          !loading && <p className="text-sm text-gray-500">No workspace found yet.</p>
        )}
      </div>

      <div>
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-white">Recent Analyses</h2>
          <Link href="/dashboard/analyses" className="text-sm text-brand-300 hover:text-brand-200">
            View all
          </Link>
        </div>

        {loading ? (
          <p className="text-sm text-gray-500">Loading analyses...</p>
        ) : analyses.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-surface-300 p-10 text-center text-sm text-gray-500">
            No analyses yet. Create your first one to get started.
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {analyses.map((a) => {
              const band = typeof a.score === "number" ? getScoreBand(a.score) : null;
              return (
                <Link
                  key={a.id}
                  href={`/dashboard/analyses/${a.id}`}
                  className="rounded-xl border border-surface-300 bg-surface-100 p-3 transition hover:border-brand-500/50"
                >
                  <div className="relative mb-2 overflow-hidden rounded-lg bg-surface-300">
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                      src={a.thumbnail_url}
                      alt={a.keyword}
                      className="aspect-video w-full object-cover"
                    />
                    {band && (
                      <span
                        className={`absolute right-1.5 top-1.5 rounded-md px-1.5 py-0.5 text-xs font-bold ${band.bgClass} ${band.textClass}`}
                      >
                        {a.score}
                      </span>
                    )}
                  </div>
                  <p className="truncate text-sm font-medium text-gray-200">{a.keyword}</p>
                  <p className="mt-1 text-xs capitalize text-gray-500">{a.status}</p>
                </Link>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
