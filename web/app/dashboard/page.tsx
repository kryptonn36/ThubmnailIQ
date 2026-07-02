"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import {
  Clock,
  Flame,
  Image as ImageIcon,
  Lightbulb,
  Plus,
  Radar,
  type LucideIcon,
} from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { getScoreBand } from "@/components/ScoreGauge";
import { Alert } from "@/components/ui/Alert";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { EmptyState } from "@/components/ui/EmptyState";
import { PageHeader } from "@/components/ui/PageHeader";
import { Skeleton } from "@/components/ui/Skeleton";
import { useMotionVariants } from "@/lib/motion";
import type { AnalysisSummary, PaginatedResponse, Workspace } from "@/types";

const QUICK_ACTIONS: {
  href: string;
  icon: LucideIcon;
  title: string;
  desc: string;
  cta: string;
  accent: "brand" | "emerald" | "amber" | "rose";
}[] = [
  {
    href: "/dashboard/analyses/new",
    icon: ImageIcon,
    title: "Analyze a Thumbnail",
    desc: "Upload and score a new thumbnail against live competitors",
    cta: "Start analysis",
    accent: "brand",
  },
  {
    href: "/dashboard/ideas",
    icon: Lightbulb,
    title: "Generate Video Ideas",
    desc: "AI-powered ideas based on what's working in your niche",
    cta: "Get ideas",
    accent: "emerald",
  },
  {
    href: "/dashboard/tracking",
    icon: Radar,
    title: "Track Competitors",
    desc: "Monitor how rival thumbnails change over time",
    cta: "Set up tracking",
    accent: "amber",
  },
  {
    href: "/dashboard/database",
    icon: Flame,
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

const ACCENT_ICON_CLASSES: Record<string, string> = {
  brand: "bg-brand-600/15 text-brand-300",
  emerald: "bg-emerald-500/15 text-emerald-400",
  amber: "bg-amber-500/15 text-amber-400",
  rose: "bg-rose-500/15 text-rose-400",
};

export default function DashboardPage() {
  const [workspace, setWorkspace] = useState<Workspace | null>(null);
  const [analyses, setAnalyses] = useState<AnalysisSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { fadeInUp, staggerContainer } = useMotionVariants();

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
    <motion.div initial="hidden" animate="visible" variants={staggerContainer} className="space-y-8">
      <motion.div variants={fadeInUp}>
        <PageHeader
          title={workspace ? workspace.name : "Overview"}
          description={
            workspace
              ? `${workspace.plan} plan · ${workspace.analyses_this_month} of ${workspace.analyses_limit} analyses used this month`
              : "Your ThumbnailIQ workspace"
          }
          actions={
            <Link href="/dashboard/analyses/new">
              <Button icon={<Plus className="h-4 w-4" aria-hidden="true" />}>New Analysis</Button>
            </Link>
          }
        />
      </motion.div>

      {error && (
        <motion.div variants={fadeInUp}>
          <Alert variant="danger">{error}</Alert>
        </motion.div>
      )}

      {/* Stats row */}
      <motion.div variants={fadeInUp} className="grid grid-cols-2 gap-4 sm:grid-cols-4">
        {[
          {
            label: "Analyses this month",
            value: workspace?.analyses_this_month ?? "—",
            sub: `of ${workspace?.analyses_limit ?? "—"} limit`,
          },
          { label: "Average score", value: avgScore !== null ? avgScore : "—", sub: "across completed analyses" },
          { label: "Best score", value: bestScore !== null ? bestScore : "—", sub: "highest scoring thumbnail" },
          {
            label: "Current plan",
            value: workspace ? workspace.plan.charAt(0).toUpperCase() + workspace.plan.slice(1) : "—",
            sub: (
              <Link href="/dashboard/billing" className="text-brand-300 hover:underline">
                Upgrade plan
              </Link>
            ),
          },
        ].map((s, i) => (
          <Card key={i}>
            <p className="text-xs font-medium uppercase tracking-wide text-gray-500">{s.label}</p>
            <p className="mt-2 text-3xl font-bold tracking-tight text-white">{s.value}</p>
            <p className="mt-1 text-xs text-gray-500">{s.sub}</p>
          </Card>
        ))}
      </motion.div>

      {/* Usage bar */}
      {workspace && (
        <motion.div variants={fadeInUp}>
          <Card>
            <div className="mb-2 flex items-center justify-between text-sm">
              <span className="text-gray-300">Monthly usage</span>
              <span className={`font-medium ${usagePct > 80 ? "text-warning" : "text-white"}`}>
                {workspace.analyses_this_month} / {workspace.analyses_limit}
              </span>
            </div>
            <div className="h-2 overflow-hidden rounded-full bg-surface-300">
              <motion.div
                initial={{ width: 0 }}
                animate={{ width: `${usagePct}%` }}
                transition={{ duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
                className={`h-full ${usagePct > 80 ? "bg-warning" : "bg-brand-gradient"}`}
              />
            </div>
            {usagePct > 80 && (
              <p className="mt-2 text-xs text-warning">
                Running low —{" "}
                <Link href="/dashboard/billing" className="underline">
                  upgrade your plan
                </Link>
              </p>
            )}
          </Card>
        </motion.div>
      )}

      {/* Quick actions */}
      <motion.div variants={fadeInUp}>
        <h2 className="mb-4 text-base font-semibold text-white">Tools</h2>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {QUICK_ACTIONS.map((action) => {
            const Icon = action.icon;
            return (
              <Link
                key={action.href}
                href={action.href}
                className={`group flex flex-col rounded-2xl border bg-surface-100 p-5 shadow-card transition-all duration-200 hover:-translate-y-0.5 hover:shadow-card-hover ${ACCENT_CLASSES[action.accent]}`}
              >
                <div
                  className={`mb-3 flex h-10 w-10 items-center justify-center rounded-xl ${ACCENT_ICON_CLASSES[action.accent]}`}
                >
                  <Icon className="h-5 w-5" aria-hidden="true" />
                </div>
                <p className="text-sm font-semibold text-white">{action.title}</p>
                <p className="mt-1 flex-1 text-xs leading-relaxed text-gray-500">{action.desc}</p>
                <p className="mt-4 text-xs font-medium text-brand-300 transition-colors group-hover:text-brand-200">
                  {action.cta} →
                </p>
              </Link>
            );
          })}
        </div>
      </motion.div>

      {/* Recent analyses */}
      <motion.div variants={fadeInUp}>
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-base font-semibold text-white">Recent Analyses</h2>
          <Link href="/dashboard/analyses" className="text-sm text-brand-300 transition-colors hover:text-brand-200">
            View all →
          </Link>
        </div>

        {loading ? (
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="rounded-xl border border-surface-300 bg-surface-100 p-3">
                <Skeleton className="mb-2 aspect-video w-full" />
                <Skeleton className="mb-1 h-3 w-3/4" />
                <Skeleton className="h-3 w-1/2" />
              </div>
            ))}
          </div>
        ) : analyses.length === 0 ? (
          <EmptyState
            icon={<ImageIcon className="h-5 w-5" aria-hidden="true" />}
            title="No analyses yet"
            description="Upload a thumbnail to see how it compares to your niche."
            action={
              <Link href="/dashboard/analyses/new">
                <Button>Create First Analysis</Button>
              </Link>
            }
          />
        ) : (
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
            {analyses.map((a) => {
              const band = typeof a.score === "number" ? getScoreBand(a.score) : null;
              return (
                <Link
                  key={a.id}
                  href={`/dashboard/analyses/${a.id}`}
                  className="group rounded-xl border border-surface-300 bg-surface-100 p-3 transition-colors duration-150 hover:border-brand-500/50"
                >
                  <div className="relative mb-2 overflow-hidden rounded-lg bg-surface-300">
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                      src={a.thumbnail_url}
                      alt={a.keyword}
                      className="aspect-video w-full object-cover transition-transform duration-300 group-hover:scale-105"
                    />
                    {band && (
                      <span
                        className={`absolute right-1.5 top-1.5 rounded-md px-1.5 py-0.5 text-xs font-bold ${band.bgClass} ${band.textClass}`}
                      >
                        {a.score}
                      </span>
                    )}
                    {(a.status === "pending" || a.status === "processing") && (
                      <span className="absolute right-1.5 top-1.5 flex items-center gap-1 rounded-md bg-warning/90 px-1.5 py-0.5 text-xs font-bold text-black">
                        <Clock className="h-3 w-3" aria-hidden="true" />
                      </span>
                    )}
                  </div>
                  <p className="truncate text-sm font-medium text-gray-200">{a.keyword}</p>
                  <Badge variant={a.status} className="mt-1">
                    {a.status}
                  </Badge>
                </Link>
              );
            })}
          </div>
        )}
      </motion.div>
    </motion.div>
  );
}
