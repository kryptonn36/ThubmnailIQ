"use client";

import { useState } from "react";
import { useParams } from "next/navigation";
import { useAnalysis } from "@/hooks/useAnalysis";
import ScoreGauge from "@/components/ScoreGauge";
import SuggestionList from "@/components/SuggestionList";
import CompetitorGrid from "@/components/CompetitorGrid";
import CVBreakdown from "@/components/CVBreakdown";
import { api, ApiError } from "@/lib/api";

const SUB_SCORES: { key: "visibility_score" | "contrast_score" | "attention_score" | "mobile_score" | "branding_score" | "curiosity_score"; label: string }[] = [
  { key: "visibility_score", label: "Visibility" },
  { key: "contrast_score", label: "Contrast" },
  { key: "attention_score", label: "Attention" },
  { key: "mobile_score", label: "Mobile" },
  { key: "branding_score", label: "Branding" },
  { key: "curiosity_score", label: "Curiosity" },
];

export default function AnalysisDetailPage() {
  const params = useParams<{ id: string }>();
  const id = params.id;
  const { analysis, loading, error } = useAnalysis(id);

  const [compareSubmitting, setCompareSubmitting] = useState(false);
  const [compareError, setCompareError] = useState<string | null>(null);

  async function handleAddVersion(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    setCompareError(null);
    setCompareSubmitting(true);
    try {
      const formData = new FormData();
      formData.append("thumbnail", file);
      await api.postForm(`/analyses/${id}/compare`, formData);
      window.location.reload();
    } catch (err) {
      setCompareError(err instanceof ApiError ? err.message : "Failed to add version.");
    } finally {
      setCompareSubmitting(false);
    }
  }

  if (loading && !analysis) {
    return (
      <div className="flex min-h-[40vh] items-center justify-center">
        <p className="text-sm text-gray-500">Loading analysis...</p>
      </div>
    );
  }

  if (error) {
    return (
      <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">{error}</p>
    );
  }

  if (!analysis) {
    return <p className="text-sm text-gray-500">Analysis not found.</p>;
  }

  if (analysis.status === "pending" || analysis.status === "processing") {
    return (
      <div className="flex min-h-[50vh] flex-col items-center justify-center gap-4 text-center">
        <div className="h-10 w-10 animate-spin rounded-full border-4 border-surface-300 border-t-brand-500" />
        <p className="text-lg font-medium text-white">Analyzing...</p>
        <p className="max-w-sm text-sm text-gray-500">
          We&apos;re scoring your thumbnail against live competitor data. This usually takes a
          few seconds.
        </p>
      </div>
    );
  }

  if (analysis.status === "failed") {
    return (
      <div className="rounded-2xl border border-red-500/30 bg-red-500/5 p-8 text-center">
        <p className="text-lg font-semibold text-red-400">Analysis failed</p>
        <p className="mt-2 text-sm text-gray-400">
          Something went wrong while processing this thumbnail. Please try again.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white">Analysis Results</h1>
          <p className="mt-1 text-sm text-gray-400">Keyword: &ldquo;{analysis.keyword}&rdquo;</p>
        </div>
        <label className="cursor-pointer rounded-lg border border-surface-300 px-4 py-2.5 text-sm font-medium text-gray-200 transition hover:bg-surface-200">
          {compareSubmitting ? "Uploading..." : "+ Add Version"}
          <input
            type="file"
            accept="image/png,image/jpeg"
            className="hidden"
            onChange={handleAddVersion}
            disabled={compareSubmitting}
          />
        </label>
      </div>

      {compareError && (
        <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">{compareError}</p>
      )}

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <div className="rounded-2xl border border-surface-300 bg-surface-100 p-6 lg:col-span-1">
          <div className="mb-4 overflow-hidden rounded-xl bg-surface-300">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={analysis.thumbnail_url}
              alt={analysis.keyword}
              className="aspect-video w-full object-cover"
            />
          </div>
          <div className="flex items-center justify-center">
            <ScoreGauge score={analysis.score ?? 0} />
          </div>
          {analysis.rank_in_competitors !== null && analysis.competitor_count > 0 && (
            <p className="mt-4 text-center text-sm text-gray-400">
              Ranked{" "}
              <span className="font-semibold text-white">#{analysis.rank_in_competitors}</span>{" "}
              of {analysis.competitor_count} competitors
            </p>
          )}
        </div>

        <div className="rounded-2xl border border-surface-300 bg-surface-100 p-6 lg:col-span-2">
          <h2 className="mb-4 text-base font-semibold text-white">Sub-Score Breakdown</h2>
          <div className="space-y-4">
            {SUB_SCORES.map(({ key, label }) => {
              const value = analysis[key] ?? 0;
              return (
                <div key={key}>
                  <div className="mb-1 flex items-center justify-between text-sm">
                    <span className="text-gray-300">{label}</span>
                    <span className="font-medium text-gray-200">{value}</span>
                  </div>
                  <div className="h-2 w-full overflow-hidden rounded-full bg-surface-300">
                    <div
                      className="h-full bg-brand-gradient"
                      style={{ width: `${Math.max(0, Math.min(100, value))}%` }}
                    />
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>

      {analysis.cv_results && (
        <CVBreakdown cvResults={analysis.cv_results} competitorAvg={analysis.competitor_avg} />
      )}

      <div className="rounded-2xl border border-surface-300 bg-surface-100 p-6">
        <h2 className="mb-4 text-base font-semibold text-white">Suggestions</h2>
        <SuggestionList suggestions={analysis.suggestions} />
      </div>

      <div className="rounded-2xl border border-surface-300 bg-surface-100 p-6">
        <h2 className="mb-4 text-base font-semibold text-white">
          Competitors ({analysis.competitor_count})
        </h2>
        <CompetitorGrid
          competitors={analysis.competitors}
          userScore={analysis.score}
          userRank={analysis.rank_in_competitors}
        />
      </div>

      {analysis.versions && analysis.versions.length > 0 && (
        <div className="rounded-2xl border border-surface-300 bg-surface-100 p-6">
          <h2 className="mb-4 text-base font-semibold text-white">Versions</h2>
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
            {analysis.versions.map((v) => (
              <div key={v.id} className="rounded-xl border border-surface-300 bg-surface-200 p-3">
                <div className="mb-2 overflow-hidden rounded-lg bg-surface-300">
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img
                    src={v.thumbnail_url}
                    alt={`Version ${v.version_number}`}
                    className="aspect-video w-full object-cover"
                  />
                </div>
                <p className="text-xs text-gray-400">Version {v.version_number}</p>
                <p className="text-lg font-bold text-white">{v.score}</p>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
