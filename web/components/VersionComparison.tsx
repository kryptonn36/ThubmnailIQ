"use client";

import { useState } from "react";
import type { AnalysisFull, ThumbnailVersion, CVResults } from "@/types";
import ImageLightbox from "./ImageLightbox";

// ---------- scoring helpers (mirrors the Go scoring engine) ----------

interface SubScores {
  visibility: number;
  contrast: number;
  attention: number;
  mobile: number;
  branding: number;
  curiosity: number;
}

function subScoresFromCV(cv: CVResults): SubScores {
  const contrastNorm = Math.min((cv.colors.contrast_score / 21) * 100, 100);
  const contrast = Math.min(
    Math.round(contrastNorm + (cv.colors.contrast_score >= 4.5 ? 10 : 0)),
    100
  );

  // Visibility: saturation 30% + contrast 40% + saturation-distance 30%
  const colorDist = Math.min(Math.abs(cv.colors.saturation_score - 50) * 1.4, 100);
  const visibility = Math.round(
    cv.colors.saturation_score * 0.3 + contrastNorm * 0.4 + colorDist * 0.3
  );

  // Attention: face 45% + eye-contact 30% + large-text 25%
  const facePresent = cv.face.face_count > 0 ? 1 : 0;
  const eyeContact = cv.face.has_eye_contact ? 1 : 0;
  const largeText = (cv.ocr.avg_text_height_pct ?? 0) >= 8 ? 1 : 0;
  const attention = Math.round((facePresent * 0.45 + eyeContact * 0.3 + largeText * 0.25) * 100);

  // Mobile: 100 − text-density-penalty − clutter-penalty − small-text-penalty
  const textDensityPenalty = Math.max(0, (cv.ocr.text_density_pct - 30) * 2);
  const clutterPenalty = cv.clutter.clutter_score * 0.5;
  const avgH = cv.ocr.avg_text_height_pct ?? 0;
  const smallTextPenalty =
    cv.ocr.text_detected && avgH > 0 && avgH < 8 ? 20 : 0;
  const mobile = Math.min(
    Math.max(0, Math.round(100 - textDensityPenalty - clutterPenalty - smallTextPenalty)),
    100
  );

  return { visibility, contrast, attention, mobile, branding: 50, curiosity: 50 };
}

// ---------- types ----------

interface VersionData {
  label: string;
  thumbnailUrl: string;
  totalScore: number;
  sub: SubScores;
  cv: CVResults | null;
  isOriginal?: boolean;
}

// ---------- constants ----------

const SUB_SCORE_META: Array<{
  key: keyof SubScores;
  label: string;
  icon: string;
  tip: (winner: string) => string;
}> = [
  {
    key: "visibility",
    label: "Visibility",
    icon: "👁",
    tip: (w) =>
      `${w} stands out the most in crowded search results due to stronger colour contrast and saturation.`,
  },
  {
    key: "contrast",
    label: "Contrast",
    icon: "🌓",
    tip: (w) =>
      `${w} achieves the highest WCAG contrast ratio — text and elements are easiest to read at a glance.`,
  },
  {
    key: "attention",
    label: "Attention",
    icon: "🎯",
    tip: (w) =>
      `${w} has the strongest emotional hook — face presence and large text grab attention fastest.`,
  },
  {
    key: "mobile",
    label: "Mobile",
    icon: "📱",
    tip: (w) =>
      `${w} has the cleanest mobile layout — least clutter and optimal text size for small screens.`,
  },
  {
    key: "branding",
    label: "Branding",
    icon: "🏷",
    tip: () => `Branding is consistent across all versions (requires channel history to differentiate).`,
  },
  {
    key: "curiosity",
    label: "Curiosity",
    icon: "🧠",
    tip: () => `Curiosity scoring requires Gemini API for live analysis.`,
  },
];

function scoreColor(score: number) {
  if (score >= 70) return "text-emerald-400";
  if (score >= 50) return "text-amber-400";
  return "text-red-400";
}

function barColor(score: number) {
  if (score >= 70) return "bg-emerald-500";
  if (score >= 50) return "bg-amber-500";
  return "bg-red-500";
}

// ---------- component ----------

interface Props {
  analysis: AnalysisFull;
  versions: ThumbnailVersion[];
}

export default function VersionComparison({ analysis, versions }: Props) {
  const [lightbox, setLightbox] = useState<{ src: string; alt: string } | null>(null);

  if (versions.length === 0) return null;

  // Build a unified row for original + each version
  const originalSub: SubScores = {
    visibility: analysis.visibility_score ?? 0,
    contrast:   analysis.contrast_score   ?? 0,
    attention:  analysis.attention_score  ?? 0,
    mobile:     analysis.mobile_score     ?? 0,
    branding:   analysis.branding_score   ?? 0,
    curiosity:  analysis.curiosity_score  ?? 0,
  };

  const all: VersionData[] = [
    {
      label: "Original",
      thumbnailUrl: analysis.thumbnail_url,
      totalScore: analysis.score ?? 0,
      sub: originalSub,
      cv: analysis.cv_results,
      isOriginal: true,
    },
    ...versions.map((v, i) => ({
      label: `Version ${i + 1}`,
      thumbnailUrl: v.thumbnail_url,
      totalScore: v.score ?? 0,
      sub: v.cv_results ? subScoresFromCV(v.cv_results) : originalSub,
      cv: v.cv_results ?? null,
    })),
  ];

  const bestOverallIdx = all.reduce(
    (best, cur, i) => (cur.totalScore > all[best].totalScore ? i : best),
    0
  );

  // For each sub-score key, find which version index wins
  const winners: Record<keyof SubScores, number> = {
    visibility: 0, contrast: 0, attention: 0, mobile: 0, branding: 0, curiosity: 0,
  };
  (Object.keys(winners) as (keyof SubScores)[]).forEach((key) => {
    winners[key] = all.reduce(
      (best, cur, i) => (cur.sub[key] > all[best].sub[key] ? i : best),
      0
    );
  });

  // "Best of each" unique insights (skip branding/curiosity as they're fixed)
  const actionableKeys: (keyof SubScores)[] = ["visibility", "contrast", "attention", "mobile"];
  const insights = actionableKeys
    .map((key) => {
      const winnerIdx = winners[key];
      const meta = SUB_SCORE_META.find((m) => m.key === key)!;
      const winnerLabel = all[winnerIdx].label;
      return { icon: meta.icon, label: meta.label, winner: winnerLabel, tip: meta.tip(winnerLabel) };
    })
    .filter(({ winner }) => winner !== "Original" || all.some((v) => !v.isOriginal));

  return (
    <>
      {lightbox && (
        <ImageLightbox
          src={lightbox.src}
          alt={lightbox.alt}
          onClose={() => setLightbox(null)}
        />
      )}

      <div className="rounded-2xl border border-surface-300 bg-surface-100 p-6 space-y-6">
        <div>
          <h2 className="text-base font-semibold text-white">Version Comparison</h2>
          <p className="text-xs text-gray-500 mt-0.5">
            Side-by-side breakdown of every sub-score across all versions
          </p>
        </div>

        {/* Thumbnail cards row */}
        <div className="grid gap-3" style={{ gridTemplateColumns: `repeat(${all.length}, minmax(0, 1fr))` }}>
          {all.map((v, i) => (
            <div
              key={i}
              className={`flex flex-col rounded-xl border p-3 ${
                i === bestOverallIdx
                  ? "border-brand-500 bg-brand-600/10"
                  : "border-surface-300 bg-surface-200"
              }`}
            >
              {i === bestOverallIdx && (
                <span className="mb-2 self-start rounded-full bg-brand-500 px-2 py-0.5 text-[10px] font-bold text-white">
                  BEST OVERALL
                </span>
              )}
              <button
                type="button"
                onClick={() => setLightbox({ src: v.thumbnailUrl, alt: v.label })}
                className="overflow-hidden rounded-lg bg-surface-300 transition hover:opacity-90 mb-2"
              >
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={v.thumbnailUrl}
                  alt={v.label}
                  className="aspect-video w-full cursor-zoom-in object-cover"
                />
              </button>
              <p className="text-xs font-medium text-gray-400">{v.label}</p>
              <p className={`text-2xl font-bold ${scoreColor(v.totalScore)}`}>
                {v.totalScore}
              </p>
              <p className="text-[10px] text-gray-600">overall score</p>
            </div>
          ))}
        </div>

        {/* Sub-score comparison table */}
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-surface-300">
                <th className="py-2 text-left text-xs font-medium text-gray-500 w-28">Metric</th>
                {all.map((v, i) => (
                  <th key={i} className="py-2 text-center text-xs font-medium text-gray-400 min-w-[80px]">
                    {v.label}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {SUB_SCORE_META.map(({ key, label, icon }) => {
                const winnerIdx = winners[key];
                return (
                  <tr key={key} className="border-b border-surface-300/50">
                    <td className="py-3 pr-4">
                      <span className="mr-1">{icon}</span>
                      <span className="text-xs text-gray-400">{label}</span>
                    </td>
                    {all.map((v, i) => {
                      const val = v.sub[key];
                      const isWinner = i === winnerIdx;
                      return (
                        <td key={i} className="py-3 px-2 text-center">
                          <div className={`inline-flex flex-col items-center gap-1 rounded-lg px-2 py-1 ${isWinner ? "bg-emerald-500/15 ring-1 ring-emerald-500/40" : ""}`}>
                            <span className={`text-sm font-bold ${isWinner ? "text-emerald-400" : "text-white"}`}>
                              {val}
                              {isWinner && <span className="ml-1 text-[9px]">▲</span>}
                            </span>
                            <div className="w-12 h-1.5 rounded-full bg-surface-300 overflow-hidden">
                              <div
                                className={`h-full rounded-full ${barColor(val)}`}
                                style={{ width: `${val}%` }}
                              />
                            </div>
                          </div>
                        </td>
                      );
                    })}
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>

        {/* Best-of-each insights */}
        <div>
          <h3 className="mb-3 text-sm font-semibold text-white">
            💡 Best elements from each version
          </h3>
          <div className="space-y-2">
            {insights.map(({ icon, label, winner, tip }) => (
              <div
                key={label}
                className="flex items-start gap-3 rounded-lg border border-surface-300 bg-surface-200 px-4 py-3"
              >
                <span className="mt-0.5 text-base">{icon}</span>
                <div className="min-w-0">
                  <p className="text-sm text-white">
                    <span className="font-semibold text-brand-300">{winner}</span>
                    {" "}wins on <span className="font-medium">{label}</span>
                  </p>
                  <p className="mt-0.5 text-xs text-gray-500">{tip}</p>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Final recommendation */}
        <div className="rounded-lg border border-brand-500/30 bg-brand-500/10 px-4 py-3">
          <p className="text-sm font-semibold text-brand-300">
            🏆 Recommendation: Publish {all[bestOverallIdx].label}
          </p>
          <p className="mt-1 text-xs text-gray-400">
            It scores the highest overall ({all[bestOverallIdx].totalScore}/100).
            To create a perfect thumbnail, take the face/emotion composition from{" "}
            <strong className="text-gray-300">{all[winners.attention].label}</strong>,
            the colour palette from{" "}
            <strong className="text-gray-300">{all[winners.visibility].label}</strong>,
            and the mobile-friendly layout from{" "}
            <strong className="text-gray-300">{all[winners.mobile].label}</strong>.
          </p>
        </div>
      </div>
    </>
  );
}
