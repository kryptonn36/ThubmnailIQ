"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api, ApiError } from "@/lib/api";
import { getScoreBand } from "@/components/ScoreGauge";
import type { AnalysisSummary, AnalysisStatus, PaginatedResponse } from "@/types";

const STATUS_FILTERS: { label: string; value: AnalysisStatus | "" }[] = [
  { label: "All", value: "" },
  { label: "Pending", value: "pending" },
  { label: "Processing", value: "processing" },
  { label: "Complete", value: "complete" },
  { label: "Failed", value: "failed" },
];

const STATUS_BADGE: Record<AnalysisStatus, string> = {
  pending: "bg-gray-500/15 text-gray-300",
  processing: "bg-amber-500/15 text-amber-400",
  complete: "bg-emerald-500/15 text-emerald-400",
  failed: "bg-red-500/15 text-red-400",
};

const PER_PAGE = 12;

export default function AnalysesPage() {
  const [analyses, setAnalyses] = useState<AnalysisSummary[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [status, setStatus] = useState<AnalysisStatus | "">("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    async function load() {
      setLoading(true);
      try {
        const res = await api.get<PaginatedResponse<AnalysisSummary>>("/analyses", {
          page,
          per_page: PER_PAGE,
          status: status || undefined,
        });
        if (cancelled) return;
        setAnalyses(res.data);
        setTotal(res.total);
      } catch (err) {
        if (cancelled) return;
        setError(err instanceof ApiError ? err.message : "Failed to load analyses");
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    load();
    return () => {
      cancelled = true;
    };
  }, [page, status]);

  const totalPages = Math.max(1, Math.ceil(total / PER_PAGE));

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white">All Analyses</h1>
          <p className="mt-1 text-sm text-gray-400">Browse and filter your thumbnail analyses.</p>
        </div>
        <Link
          href="/dashboard/analyses/new"
          className="rounded-lg bg-brand-gradient px-4 py-2.5 text-sm font-semibold text-white shadow-glow transition hover:opacity-90"
        >
          + New Analysis
        </Link>
      </div>

      <div className="flex flex-wrap gap-2">
        {STATUS_FILTERS.map((f) => (
          <button
            key={f.value}
            onClick={() => {
              setStatus(f.value);
              setPage(1);
            }}
            className={`rounded-full px-3 py-1.5 text-sm font-medium transition ${
              status === f.value
                ? "bg-brand-600 text-white"
                : "bg-surface-200 text-gray-400 hover:text-gray-200"
            }`}
          >
            {f.label}
          </button>
        ))}
      </div>

      {error && (
        <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">{error}</p>
      )}

      {loading ? (
        <p className="text-sm text-gray-500">Loading...</p>
      ) : analyses.length === 0 ? (
        <div className="rounded-2xl border border-dashed border-surface-300 p-10 text-center text-sm text-gray-500">
          No analyses found.
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
                <span
                  className={`mt-1 inline-block rounded-full px-2 py-0.5 text-xs font-medium capitalize ${STATUS_BADGE[a.status]}`}
                >
                  {a.status}
                </span>
              </Link>
            );
          })}
        </div>
      )}

      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-3">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page <= 1}
            className="rounded-lg border border-surface-300 px-3 py-1.5 text-sm text-gray-300 disabled:opacity-40"
          >
            Previous
          </button>
          <span className="text-sm text-gray-500">
            Page {page} of {totalPages}
          </span>
          <button
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page >= totalPages}
            className="rounded-lg border border-surface-300 px-3 py-1.5 text-sm text-gray-300 disabled:opacity-40"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}
