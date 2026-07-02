"use client";

import { useState } from "react";
import { Monitor, Search, Smartphone } from "lucide-react";
import type { Competitor } from "@/types";
import { Card } from "@/components/ui/Card";

interface SERPPreviewProps {
  thumbnailUrl: string;
  keyword: string;
  rank: number | null;
  competitorCount: number;
  competitors: Competitor[];
}

function formatViews(n: number) {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`;
  if (n >= 1_000) return `${(n / 1_000).toFixed(0)}K`;
  return String(n);
}

function VideoCard({
  thumbnailUrl,
  title,
  channel,
  views,
  isYours,
  rank,
  small,
}: {
  thumbnailUrl: string;
  title: string;
  channel: string;
  views: number;
  isYours?: boolean;
  rank?: number;
  small?: boolean;
}) {
  return (
    <div
      className={`group relative flex ${small ? "gap-2" : "flex-col gap-2"} ${isYours ? "rounded-xl p-1 ring-2 ring-brand-500" : ""}`}
    >
      {isYours && (
        <div className="absolute -top-3 left-1/2 z-10 -translate-x-1/2 whitespace-nowrap rounded-full bg-brand-500 px-2 py-0.5 text-[10px] font-bold text-white shadow">
          YOUR THUMBNAIL · #{rank}
        </div>
      )}
      <div
        className={`relative overflow-hidden rounded-lg bg-surface-300 ${small ? "h-16 w-28 shrink-0" : "aspect-video w-full"}`}
      >
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src={thumbnailUrl}
          alt={title}
          className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
        />
        {isYours && <div className="absolute inset-0 rounded-lg ring-2 ring-inset ring-brand-500" />}
      </div>
      <div className={small ? "min-w-0 flex-1" : ""}>
        <p className={`line-clamp-2 font-medium leading-snug text-white ${small ? "text-xs" : "text-sm"}`}>
          {title}
        </p>
        <p className={`mt-0.5 text-gray-500 ${small ? "text-[10px]" : "text-xs"}`}>
          {channel} · {views > 0 ? `${formatViews(views)} views` : ""}
        </p>
      </div>
    </div>
  );
}

export default function SERPPreview({ thumbnailUrl, keyword, rank, competitors }: SERPPreviewProps) {
  const [view, setView] = useState<"desktop" | "mobile">("desktop");

  // Build a list of result slots: insert user at their rank position
  const userSlot = rank ?? competitors.length + 1;
  const slots: Array<{ isYours: boolean; comp?: Competitor; rank: number }> = [];

  let compIdx = 0;
  for (let pos = 1; pos <= Math.min(12, (competitors.length || 0) + 1); pos++) {
    if (pos === userSlot) {
      slots.push({ isYours: true, rank: pos });
    } else {
      if (compIdx < competitors.length) {
        slots.push({ isYours: false, comp: competitors[compIdx], rank: pos });
        compIdx++;
      }
    }
  }

  return (
    <Card className="p-6">
      <div className="mb-5 flex flex-wrap items-center justify-between gap-3">
        <div>
          <h3 className="text-base font-semibold text-white">SERP Preview</h3>
          <p className="text-xs text-gray-500">How your thumbnail looks in YouTube search results</p>
        </div>
        <div className="flex overflow-hidden rounded-lg border border-surface-300">
          {(["desktop", "mobile"] as const).map((v) => {
            const Icon = v === "desktop" ? Monitor : Smartphone;
            return (
              <button
                key={v}
                type="button"
                onClick={() => setView(v)}
                className={`flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium capitalize transition-colors duration-150 ${
                  view === v ? "bg-brand-600 text-white" : "text-gray-400 hover:text-gray-200"
                }`}
              >
                <Icon className="h-3.5 w-3.5" aria-hidden="true" />
                {v}
              </button>
            );
          })}
        </div>
      </div>

      {/* Browser chrome mockup */}
      <div className="overflow-hidden rounded-xl border border-surface-200 bg-[#0f0f0f]">
        {/* Address bar */}
        <div className="flex items-center gap-2 border-b border-surface-300 bg-surface-200 px-3 py-2">
          <div className="flex gap-1">
            {["#ef4444", "#f59e0b", "#22c55e"].map((c) => (
              <span key={c} className="h-2.5 w-2.5 rounded-full" style={{ backgroundColor: c }} />
            ))}
          </div>
          <div className="flex-1 rounded-full bg-surface-300 px-3 py-1 text-center text-[10px] text-gray-400">
            youtube.com/results?search_query={encodeURIComponent(keyword)}
          </div>
        </div>

        {/* YouTube header mockup */}
        <div className="flex items-center gap-4 border-b border-surface-300/50 px-4 py-2">
          <span className="text-base font-bold text-white">▶ YouTube</span>
          <div
            className={`flex items-center rounded-full border border-surface-300 bg-surface-200 ${view === "desktop" ? "w-72" : "w-40"} px-3 py-1`}
          >
            <span className="flex-1 text-xs text-gray-300">{keyword}</span>
            <Search className="h-3.5 w-3.5 text-gray-400" aria-hidden="true" />
          </div>
        </div>

        {/* Results */}
        <div className={`p-4 ${view === "mobile" ? "mx-auto max-w-sm" : ""}`}>
          <p className="mb-3 text-[10px] text-gray-500">
            About {((competitors.length || 0) * 1247 + 842).toLocaleString()} results
          </p>

          {view === "desktop" ? (
            <div className="grid grid-cols-3 gap-x-4 gap-y-6">
              {slots.map((slot, i) => (
                <VideoCard
                  key={i}
                  thumbnailUrl={slot.isYours ? thumbnailUrl : (slot.comp?.thumbnail_url ?? thumbnailUrl)}
                  title={slot.isYours ? "Your Video Title (How it appears to viewers)" : (slot.comp?.video_title ?? "Competitor Video")}
                  channel={slot.isYours ? "Your Channel" : (slot.comp?.channel_name ?? "Competitor")}
                  views={slot.isYours ? 0 : (slot.comp?.view_count ?? 0)}
                  isYours={slot.isYours}
                  rank={slot.rank}
                />
              ))}
            </div>
          ) : (
            <div className="space-y-3">
              {slots.map((slot, i) => (
                <VideoCard
                  key={i}
                  thumbnailUrl={slot.isYours ? thumbnailUrl : (slot.comp?.thumbnail_url ?? thumbnailUrl)}
                  title={slot.isYours ? "Your Video Title (How it appears to viewers)" : (slot.comp?.video_title ?? "Competitor Video")}
                  channel={slot.isYours ? "Your Channel" : (slot.comp?.channel_name ?? "Competitor")}
                  views={slot.isYours ? 0 : (slot.comp?.view_count ?? 0)}
                  isYours={slot.isYours}
                  rank={slot.rank}
                  small
                />
              ))}
            </div>
          )}
        </div>
      </div>

      {rank !== null && (
        <p className="mt-3 text-center text-xs text-gray-500">
          Your thumbnail is predicted to rank <span className="font-semibold text-white">#{rank}</span> for &ldquo;{keyword}&rdquo;
        </p>
      )}
    </Card>
  );
}
