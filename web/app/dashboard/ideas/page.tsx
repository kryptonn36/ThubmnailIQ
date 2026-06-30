"use client";

import { useState } from "react";
import { api, ApiError } from "@/lib/api";
import Link from "next/link";

interface VideoIdea {
  title: string;
  hook: string;
  format: string;
  ctr_potential: string;
}

const CTR_COLORS: Record<string, string> = {
  "Very High": "text-green-400 bg-green-400/10",
  High: "text-emerald-400 bg-emerald-400/10",
  Medium: "text-yellow-400 bg-yellow-400/10",
  Low: "text-gray-400 bg-gray-400/10",
};

const FORMAT_ICONS: Record<string, string> = {
  Listicle: "📋",
  "How-To": "🎓",
  Story: "📖",
  Challenge: "🏆",
  Comparison: "⚖️",
  "Case Study": "🔬",
  Reaction: "😲",
};

export default function IdeasPage() {
  const [keyword, setKeyword] = useState("");
  const [ideas, setIdeas] = useState<VideoIdea[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [searched, setSearched] = useState("");

  async function handleGenerate(e: React.FormEvent) {
    e.preventDefault();
    if (!keyword.trim()) return;
    setLoading(true);
    setError(null);
    setIdeas([]);
    try {
      const data = await api.get<{ ideas: VideoIdea[] }>(
        `/keywords/${encodeURIComponent(keyword.trim())}/ideas`
      );
      setIdeas(data.ideas);
      setSearched(keyword.trim());
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to generate ideas.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-white">Video Ideas</h1>
        <p className="mt-1 text-sm text-gray-400">
          Enter a keyword to get AI-powered video ideas based on what&apos;s already working for top competitors.
        </p>
      </div>

      <form onSubmit={handleGenerate} className="flex gap-3">
        <input
          type="text"
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          placeholder="e.g. crypto trading for beginners"
          className="flex-1 rounded-xl border border-surface-300 bg-surface-100 px-4 py-3 text-sm text-white outline-none focus:border-brand-500 placeholder:text-gray-600"
        />
        <button
          type="submit"
          disabled={loading || !keyword.trim()}
          className="rounded-xl bg-brand-gradient px-6 py-3 text-sm font-semibold text-white shadow-glow transition hover:opacity-90 disabled:opacity-60"
        >
          {loading ? "Generating..." : "Generate Ideas"}
        </button>
      </form>

      {error && (
        <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">{error}</p>
      )}

      {loading && (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="animate-pulse rounded-2xl border border-surface-300 bg-surface-100 p-5 space-y-3">
              <div className="h-4 w-3/4 rounded bg-surface-300" />
              <div className="h-3 w-full rounded bg-surface-300" />
              <div className="h-3 w-1/2 rounded bg-surface-300" />
            </div>
          ))}
        </div>
      )}

      {!loading && ideas.length > 0 && (
        <>
          <p className="text-sm text-gray-500">
            {ideas.length} ideas generated for &ldquo;<span className="text-white">{searched}</span>&rdquo;
          </p>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {ideas.map((idea, i) => (
              <div
                key={i}
                className="flex flex-col rounded-2xl border border-surface-300 bg-surface-100 p-5 transition hover:border-brand-500/50"
              >
                <div className="mb-3 flex items-center gap-2">
                  <span className="text-lg">{FORMAT_ICONS[idea.format] ?? "💡"}</span>
                  <span className="rounded-full bg-surface-300 px-2 py-0.5 text-xs text-gray-300">
                    {idea.format}
                  </span>
                  <span className={`ml-auto rounded-full px-2 py-0.5 text-xs font-medium ${CTR_COLORS[idea.ctr_potential] ?? "text-gray-400 bg-gray-400/10"}`}>
                    {idea.ctr_potential} CTR
                  </span>
                </div>

                <h3 className="mb-2 text-sm font-semibold leading-snug text-white">{idea.title}</h3>
                <p className="flex-1 text-xs leading-relaxed text-gray-400">{idea.hook}</p>

                <Link
                  href={`/dashboard/analyses/new?keyword=${encodeURIComponent(idea.title)}`}
                  className="mt-4 block rounded-lg border border-surface-300 py-2 text-center text-xs font-medium text-gray-300 transition hover:border-brand-500 hover:text-brand-300"
                >
                  Use this idea →
                </Link>
              </div>
            ))}
          </div>
        </>
      )}

      {!loading && ideas.length === 0 && !error && (
        <div className="rounded-2xl border border-dashed border-surface-300 p-14 text-center">
          <p className="text-4xl mb-3">💡</p>
          <p className="font-medium text-gray-300">Enter a keyword to generate video ideas</p>
          <p className="mt-1 text-sm text-gray-500">
            We study your top competitors&apos; titles and use AI to suggest original angles that stand out.
          </p>
        </div>
      )}
    </div>
  );
}
