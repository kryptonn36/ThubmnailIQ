"use client";

import { useRef, useState } from "react";
import { useParams } from "next/navigation";
import { motion } from "framer-motion";
import { AlertTriangle, Check, Plus, Printer } from "lucide-react";
import { useAnalysis } from "@/hooks/useAnalysis";
import ScoreGauge from "@/components/ScoreGauge";
import SuggestionList from "@/components/SuggestionList";
import CompetitorGrid from "@/components/CompetitorGrid";
import CVBreakdown from "@/components/CVBreakdown";
import ImageLightbox from "@/components/ImageLightbox";
import SERPPreview from "@/components/SERPPreview";
import NicheBenchmark from "@/components/NicheBenchmark";
import VersionComparison from "@/components/VersionComparison";
import { api, ApiError } from "@/lib/api";
import { Alert } from "@/components/ui/Alert";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { Spinner } from "@/components/ui/Spinner";
import { useMotionVariants } from "@/lib/motion";

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
  const { fadeInUp, staggerContainer } = useMotionVariants();

  const addVersionInputRef = useRef<HTMLInputElement>(null);
  const [compareSubmitting, setCompareSubmitting] = useState(false);
  const [compareError, setCompareError] = useState<string | null>(null);
  const [lightbox, setLightbox] = useState<{ src: string; alt: string } | null>(null);
  const [ctr, setCtr] = useState("");
  const [ctrSaving, setCtrSaving] = useState(false);
  const [ctrMsg, setCtrMsg] = useState<string | null>(null);

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
      <div className="flex min-h-[40vh] items-center justify-center gap-2.5 text-gray-500">
        <Spinner className="h-4 w-4" />
        <span className="text-sm">Loading analysis…</span>
      </div>
    );
  }

  if (error) {
    return <Alert variant="danger">{error}</Alert>;
  }

  if (!analysis) {
    return <p className="text-sm text-gray-500">Analysis not found.</p>;
  }

  if (analysis.status === "pending" || analysis.status === "processing") {
    return (
      <div className="flex min-h-[50vh] flex-col items-center justify-center gap-4 text-center">
        <div className="h-10 w-10 animate-spin rounded-full border-4 border-surface-300 border-t-brand-500" />
        <p className="text-lg font-medium text-white">Analyzing…</p>
        <p className="max-w-sm text-sm text-gray-500">
          We&apos;re scoring your thumbnail against live competitor data. This usually takes a
          few seconds.
        </p>
      </div>
    );
  }

  if (analysis.status === "failed") {
    return (
      <Card className="border-danger/30 bg-danger/5 p-8 text-center">
        <p className="text-lg font-semibold text-danger">Analysis failed</p>
        <p className="mt-2 text-sm text-gray-400">
          Something went wrong while processing this thumbnail. Please try again.
        </p>
      </Card>
    );
  }

  return (
    <motion.div initial="hidden" animate="visible" variants={staggerContainer} className="space-y-8">
      <motion.div variants={fadeInUp} className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-white sm:text-3xl">Analysis Results</h1>
          <p className="mt-1.5 text-sm text-gray-400">Keyword: &ldquo;{analysis.keyword}&rdquo;</p>
        </div>
        <div className="flex flex-wrap gap-2 print:hidden">
          <Button variant="secondary" type="button" onClick={() => window.print()} icon={<Printer className="h-4 w-4" aria-hidden="true" />}>
            Export PDF
          </Button>
          <Button
            variant="secondary"
            type="button"
            loading={compareSubmitting}
            icon={<Plus className="h-4 w-4" aria-hidden="true" />}
            onClick={() => addVersionInputRef.current?.click()}
          >
            Add Version
          </Button>
          <input
            ref={addVersionInputRef}
            type="file"
            accept="image/png,image/jpeg"
            className="hidden"
            onChange={handleAddVersion}
            disabled={compareSubmitting}
          />
        </div>
      </motion.div>

      {compareError && (
        <motion.div variants={fadeInUp}>
          <Alert variant="danger">{compareError}</Alert>
        </motion.div>
      )}

      {analysis.relevance_score !== null && analysis.relevance_score < 70 && (
        <motion.div variants={fadeInUp}>
          <Alert variant="warning">
            <p className="font-medium">
              Low keyword relevance ({analysis.relevance_score}%) — this thumbnail may not match &ldquo;{analysis.keyword}&rdquo;, which reduced your score and rank.
            </p>
            <p className="mt-1 opacity-80">
              Try uploading a thumbnail whose visuals and text are clearly related to the keyword.
            </p>
          </Alert>
        </motion.div>
      )}

      <motion.div variants={fadeInUp} className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <Card className="p-6 lg:col-span-1">
          <button
            type="button"
            onClick={() => setLightbox({ src: analysis.thumbnail_url, alt: analysis.keyword })}
            className="mb-4 block w-full overflow-hidden rounded-xl bg-surface-300 transition-opacity duration-150 hover:opacity-90"
          >
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={analysis.thumbnail_url}
              alt={analysis.keyword}
              className="aspect-video w-full cursor-zoom-in object-cover"
            />
          </button>
          <div className="flex items-center justify-center">
            <ScoreGauge score={analysis.score ?? 0} />
          </div>
          {analysis.rank_in_competitors !== null && analysis.competitor_count > 0 && (
            <p className="mt-4 text-center text-sm text-gray-400">
              Ranked <span className="font-semibold text-white">#{analysis.rank_in_competitors}</span> of{" "}
              {analysis.competitor_count} competitors
            </p>
          )}
        </Card>

        <Card className="p-6 lg:col-span-2">
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
                    <motion.div
                      initial={{ width: 0 }}
                      animate={{ width: `${Math.max(0, Math.min(100, value))}%` }}
                      transition={{ duration: 0.5, ease: [0.16, 1, 0.3, 1] }}
                      className="h-full bg-brand-gradient"
                    />
                  </div>
                </div>
              );
            })}
          </div>
        </Card>
      </motion.div>

      {analysis.status === "complete" && (
        <motion.div variants={fadeInUp}>
          <SERPPreview
            thumbnailUrl={analysis.thumbnail_url}
            keyword={analysis.keyword}
            rank={analysis.rank_in_competitors}
            competitorCount={analysis.competitor_count}
            competitors={analysis.competitors}
          />
        </motion.div>
      )}

      {analysis.cv_results && (
        <motion.div variants={fadeInUp}>
          <CVBreakdown cvResults={analysis.cv_results} competitorAvg={analysis.competitor_avg} />
        </motion.div>
      )}

      {analysis.status === "complete" && analysis.score !== null && analysis.competitors.length > 0 && (
        <motion.div variants={fadeInUp}>
          <NicheBenchmark score={analysis.score} competitors={analysis.competitors} keyword={analysis.keyword} />
        </motion.div>
      )}

      <motion.div variants={fadeInUp}>
        <Card className="p-6">
          <h2 className="mb-4 text-base font-semibold text-white">Suggestions</h2>
          <SuggestionList suggestions={analysis.suggestions} />
        </Card>
      </motion.div>

      {/* CTR Tracker — record real-world performance after publishing */}
      {analysis.status === "complete" && (
        <motion.div variants={fadeInUp} className="print:hidden">
          <Card className="p-6">
            <h2 className="mb-1 text-base font-semibold text-white">Track Performance</h2>
            <p className="mb-4 text-xs text-gray-500">
              After publishing, enter your actual CTR to see how the ThumbnailIQ score correlated with real performance.
            </p>
            {analysis.actual_ctr !== null && analysis.actual_ctr !== undefined ? (
              <div className="mb-4 flex items-center gap-3 rounded-lg bg-success/10 px-4 py-3">
                <span className="text-2xl font-bold text-success">{analysis.actual_ctr}%</span>
                <div>
                  <p className="text-sm font-medium text-success">Actual CTR recorded</p>
                  {analysis.published_at && (
                    <p className="text-xs text-gray-500">
                      Published {new Date(analysis.published_at).toLocaleDateString()}
                    </p>
                  )}
                </div>
              </div>
            ) : null}
            <form
              onSubmit={async (e) => {
                e.preventDefault();
                setCtrSaving(true);
                setCtrMsg(null);
                try {
                  await api.patch(`/analyses/${id}/ctr`, {
                    actual_ctr: parseFloat(ctr),
                    published_at: new Date().toISOString(),
                  });
                  setCtrMsg("CTR saved");
                  setCtr("");
                } catch {
                  setCtrMsg("Failed to save CTR.");
                } finally {
                  setCtrSaving(false);
                }
              }}
              className="flex items-end gap-3"
            >
              <div className="flex-1">
                <label className="mb-1 block text-xs font-medium text-gray-400">
                  CTR % (from YouTube Studio)
                </label>
                <Input
                  type="number"
                  min="0"
                  max="100"
                  step="0.1"
                  value={ctr}
                  onChange={(e) => setCtr(e.target.value)}
                  placeholder="e.g. 4.7"
                />
              </div>
              <Button type="submit" loading={ctrSaving} disabled={!ctr}>
                Save CTR
              </Button>
            </form>
            {ctrMsg && (
              <p className="mt-2 flex items-center gap-1 text-xs text-success">
                {ctrMsg === "CTR saved" && <Check className="h-3 w-3" aria-hidden="true" />}
                {ctrMsg}
              </p>
            )}
          </Card>
        </motion.div>
      )}

      <motion.div variants={fadeInUp}>
        <Card className="p-6">
          <h2 className="mb-4 text-base font-semibold text-white">Competitors ({analysis.competitor_count})</h2>
          <CompetitorGrid
            competitors={analysis.competitors}
            userScore={analysis.score}
            userRank={analysis.rank_in_competitors}
          />
        </Card>
      </motion.div>

      {analysis.versions && analysis.versions.length > 0 && (
        <motion.div variants={fadeInUp}>
          <VersionComparison analysis={analysis} versions={analysis.versions} />
        </motion.div>
      )}

      {lightbox && <ImageLightbox src={lightbox.src} alt={lightbox.alt} onClose={() => setLightbox(null)} />}
    </motion.div>
  );
}
