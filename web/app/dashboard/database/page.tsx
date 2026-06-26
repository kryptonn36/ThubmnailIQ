"use client";

import { useEffect, useState } from "react";
import { api, ApiError } from "@/lib/api";
import type { ViralThumbnail } from "@/types";

const NICHES = [
  "All",
  "Gaming",
  "Finance",
  "Fitness",
  "Beauty",
  "Tech",
  "Education",
  "Entertainment",
];

function formatViews(views: number): string {
  if (views >= 1_000_000) return `${(views / 1_000_000).toFixed(1)}M views`;
  if (views >= 1_000) return `${(views / 1_000).toFixed(1)}K views`;
  return `${views} views`;
}

export default function ViralDatabasePage() {
  const [niche, setNiche] = useState("All");
  const [minScore, setMinScore] = useState(0);
  const [hasFace, setHasFace] = useState<"" | "true" | "false">("");
  const [results, setResults] = useState<ViralThumbnail[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    async function load() {
      setLoading(true);
      try {
        const data = await api.get<{ data: ViralThumbnail[] }>("/viral-db", {
          niche: niche === "All" ? undefined : niche,
          min_score: minScore > 0 ? minScore : undefined,
          has_face: hasFace || undefined,
        });
        if (cancelled) return;
        setResults(data.data);
      } catch (err) {
        if (cancelled) return;
        setError(err instanceof ApiError ? err.message : "Failed to load viral database.");
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    load();
    return () => {
      cancelled = true;
    };
  }, [niche, minScore, hasFace]);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">Viral Thumbnail Database</h1>
        <p className="mt-1 text-sm text-gray-400">
          Browse high-performing thumbnails collected across all analyses.
        </p>
      </div>

      <div className="flex flex-wrap items-end gap-4 rounded-2xl border border-surface-300 bg-surface-100 p-5">
        <div>
          <label htmlFor="niche" className="mb-1 block text-xs font-medium text-gray-400">
            Niche
          </label>
          <select
            id="niche"
            value={niche}
            onChange={(e) => setNiche(e.target.value)}
            className="rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500"
          >
            {NICHES.map((n) => (
              <option key={n} value={n}>
                {n}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label htmlFor="minScore" className="mb-1 block text-xs font-medium text-gray-400">
            Min score
          </label>
          <input
            id="minScore"
            type="number"
            min={0}
            max={100}
            value={minScore}
            onChange={(e) => setMinScore(Number(e.target.value))}
            className="w-24 rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500"
          />
        </div>

        <div>
          <label htmlFor="hasFace" className="mb-1 block text-xs font-medium text-gray-400">
            Has face
          </label>
          <select
            id="hasFace"
            value={hasFace}
            onChange={(e) => setHasFace(e.target.value as "" | "true" | "false")}
            className="rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500"
          >
            <option value="">Any</option>
            <option value="true">Yes</option>
            <option value="false">No</option>
          </select>
        </div>
      </div>

      {error && (
        <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">{error}</p>
      )}

      {loading ? (
        <p className="text-sm text-gray-500">Loading...</p>
      ) : results.length === 0 ? (
        <div className="rounded-2xl border border-dashed border-surface-300 p-10 text-center text-sm text-gray-500">
          No thumbnails match your filters.
        </div>
      ) : (
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
          {results.map((item) => (
            <div
              key={item.id}
              className="rounded-xl border border-surface-300 bg-surface-100 p-3"
            >
              <div className="relative mb-2 overflow-hidden rounded-lg bg-surface-300">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={item.thumbnail_url}
                  alt={item.video_title}
                  className="aspect-video w-full object-cover"
                />
                <span className="absolute right-1.5 top-1.5 rounded-md bg-brand-600/90 px-1.5 py-0.5 text-xs font-bold text-white">
                  {item.score}
                </span>
              </div>
              <p className="line-clamp-2 text-sm font-medium text-gray-200">{item.video_title}</p>
              <p className="mt-1 text-xs text-gray-500">{item.channel_name}</p>
              <p className="text-xs text-gray-500">
                {item.niche} &middot; {formatViews(item.view_count)}
              </p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
