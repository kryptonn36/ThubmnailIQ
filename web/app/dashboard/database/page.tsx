"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Search } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { Alert } from "@/components/ui/Alert";
import { Button } from "@/components/ui/Button";
import { EmptyState } from "@/components/ui/EmptyState";
import { PageHeader } from "@/components/ui/PageHeader";
import { Select } from "@/components/ui/Select";
import { Skeleton } from "@/components/ui/Skeleton";
import { useMotionVariants } from "@/lib/motion";
import type { ViralThumbnail } from "@/types";

const NICHES = ["All", "Gaming", "Finance", "Fitness", "Beauty", "Tech", "Education", "Entertainment"];

function formatViews(views: number): string {
  if (views >= 1_000_000) return `${(views / 1_000_000).toFixed(1)}M views`;
  if (views >= 1_000) return `${(views / 1_000).toFixed(1)}K views`;
  return `${views} views`;
}

export default function ViralDatabasePage() {
  const { fadeInUp, staggerContainer } = useMotionVariants();
  const [keywordInput, setKeywordInput] = useState("");
  const [keyword, setKeyword] = useState("");
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
          keyword: keyword.trim() || undefined,
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
  }, [keyword, niche, minScore, hasFace]);

  return (
    <div className="space-y-6">
      <PageHeader
        title="Viral Thumbnail Database"
        description="Browse high-performing thumbnails collected across all analyses."
      />

      <form
        onSubmit={(e) => {
          e.preventDefault();
          setKeyword(keywordInput.trim());
        }}
        className="flex flex-wrap items-end gap-4 rounded-2xl border border-surface-300 bg-surface-100 p-5 shadow-card"
      >
        <div>
          <label htmlFor="keyword" className="mb-1 block text-xs font-medium text-gray-400">
            Keyword
          </label>
          <input
            id="keyword"
            type="search"
            value={keywordInput}
            onChange={(e) => setKeywordInput(e.target.value)}
            placeholder="Search YouTube"
            className="w-56 rounded-lg border border-surface-300 bg-surface-200 px-3 py-2.5 text-sm text-white outline-none transition-colors placeholder:text-gray-600 focus:border-brand-500 focus:ring-2 focus:ring-brand-500/20"
          />
        </div>

        <Button type="submit" icon={<Search className="h-4 w-4" aria-hidden="true" />}>
          Search
        </Button>

        <div>
          <label htmlFor="niche" className="mb-1 block text-xs font-medium text-gray-400">
            Niche
          </label>
          <Select id="niche" value={niche} onChange={(e) => setNiche(e.target.value)} className="w-40">
            {NICHES.map((n) => (
              <option key={n} value={n}>
                {n}
              </option>
            ))}
          </Select>
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
            className="w-24 rounded-lg border border-surface-300 bg-surface-200 px-3 py-2.5 text-sm text-white outline-none transition-colors focus:border-brand-500 focus:ring-2 focus:ring-brand-500/20"
          />
        </div>

        <div>
          <label htmlFor="hasFace" className="mb-1 block text-xs font-medium text-gray-400">
            Has face
          </label>
          <Select
            id="hasFace"
            value={hasFace}
            onChange={(e) => setHasFace(e.target.value as "" | "true" | "false")}
            className="w-28"
          >
            <option value="">Any</option>
            <option value="true">Yes</option>
            <option value="false">No</option>
          </Select>
        </div>
      </form>

      {error && <Alert variant="danger">{error}</Alert>}

      {loading ? (
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
          {[...Array(8)].map((_, i) => (
            <div key={i} className="rounded-xl border border-surface-300 bg-surface-100 p-3">
              <Skeleton className="mb-2 aspect-video w-full" />
              <Skeleton className="mb-1 h-3 w-full" />
              <Skeleton className="h-3 w-1/2" />
            </div>
          ))}
        </div>
      ) : results.length === 0 ? (
        <EmptyState title="No thumbnails match your filters." />
      ) : (
        <motion.div
          initial="hidden"
          animate="visible"
          variants={staggerContainer}
          className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4"
        >
          {results.map((item) => (
            <motion.div
              key={item.id}
              variants={fadeInUp}
              className="rounded-xl border border-surface-300 bg-surface-100 p-3 transition-colors duration-150 hover:border-brand-500/40"
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
            </motion.div>
          ))}
        </motion.div>
      )}
    </div>
  );
}
