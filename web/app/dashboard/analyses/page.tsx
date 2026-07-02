"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import { Plus } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { getScoreBand } from "@/components/ScoreGauge";
import { Alert } from "@/components/ui/Alert";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { EmptyState } from "@/components/ui/EmptyState";
import { PageHeader } from "@/components/ui/PageHeader";
import { Skeleton } from "@/components/ui/Skeleton";
import { useMotionVariants } from "@/lib/motion";
import type { AnalysisSummary, AnalysisStatus, PaginatedResponse } from "@/types";

const STATUS_FILTERS: { label: string; value: AnalysisStatus | "" }[] = [
  { label: "All", value: "" },
  { label: "Pending", value: "pending" },
  { label: "Processing", value: "processing" },
  { label: "Complete", value: "complete" },
  { label: "Failed", value: "failed" },
];

const PER_PAGE = 12;

export default function AnalysesPage() {
  const [analyses, setAnalyses] = useState<AnalysisSummary[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [status, setStatus] = useState<AnalysisStatus | "">("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { fadeInUp, staggerContainer } = useMotionVariants();

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
      <PageHeader
        title="All Analyses"
        description="Browse and filter your thumbnail analyses."
        actions={
          <Link href="/dashboard/analyses/new">
            <Button icon={<Plus className="h-4 w-4" aria-hidden="true" />}>New Analysis</Button>
          </Link>
        }
      />

      <div className="flex flex-wrap gap-2">
        {STATUS_FILTERS.map((f) => (
          <button
            key={f.value}
            onClick={() => {
              setStatus(f.value);
              setPage(1);
            }}
            className={`rounded-full px-3 py-1.5 text-sm font-medium transition-colors duration-150 ${
              status === f.value
                ? "bg-brand-600 text-white"
                : "bg-surface-200 text-gray-400 hover:text-gray-200"
            }`}
          >
            {f.label}
          </button>
        ))}
      </div>

      {error && <Alert variant="danger">{error}</Alert>}

      {loading ? (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {[...Array(8)].map((_, i) => (
            <div key={i} className="rounded-xl border border-surface-300 bg-surface-100 p-3">
              <Skeleton className="mb-2 aspect-video w-full" />
              <Skeleton className="mb-1 h-3 w-3/4" />
              <Skeleton className="h-4 w-16 rounded-full" />
            </div>
          ))}
        </div>
      ) : analyses.length === 0 ? (
        <EmptyState title="No analyses found." description="Try a different filter or create a new analysis." />
      ) : (
        <motion.div
          initial="hidden"
          animate="visible"
          variants={staggerContainer}
          className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4"
        >
          {analyses.map((a) => {
            const band = typeof a.score === "number" ? getScoreBand(a.score) : null;
            return (
              <motion.div key={a.id} variants={fadeInUp}>
                <Link
                  href={`/dashboard/analyses/${a.id}`}
                  className="block rounded-xl border border-surface-300 bg-surface-100 p-3 transition-colors duration-150 hover:border-brand-500/50"
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
                  <Badge variant={a.status} className="mt-1">
                    {a.status}
                  </Badge>
                </Link>
              </motion.div>
            );
          })}
        </motion.div>
      )}

      {totalPages > 1 && (
        <div className="flex items-center justify-center gap-3">
          <Button
            variant="secondary"
            size="sm"
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page <= 1}
          >
            Previous
          </Button>
          <span className="text-sm text-gray-500">
            Page {page} of {totalPages}
          </span>
          <Button
            variant="secondary"
            size="sm"
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page >= totalPages}
          >
            Next
          </Button>
        </div>
      )}
    </div>
  );
}
