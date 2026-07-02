"use client";

import { useState } from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import {
  BookOpen,
  GraduationCap,
  Lightbulb,
  List,
  Microscope,
  Scale,
  Trophy,
  Zap,
  type LucideIcon,
} from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { Alert } from "@/components/ui/Alert";
import { Badge, type BadgeVariant } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { EmptyState } from "@/components/ui/EmptyState";
import { Input } from "@/components/ui/Input";
import { PageHeader } from "@/components/ui/PageHeader";
import { Skeleton } from "@/components/ui/Skeleton";
import { useMotionVariants } from "@/lib/motion";

interface VideoIdea {
  title: string;
  hook: string;
  format: string;
  ctr_potential: string;
}

const CTR_BADGE: Record<string, BadgeVariant> = {
  "Very High": "success",
  High: "success",
  Medium: "warning",
  Low: "neutral",
};

const FORMAT_ICONS: Record<string, LucideIcon> = {
  Listicle: List,
  "How-To": GraduationCap,
  Story: BookOpen,
  Challenge: Trophy,
  Comparison: Scale,
  "Case Study": Microscope,
  Reaction: Zap,
};

export default function IdeasPage() {
  const { fadeInUp, staggerContainer } = useMotionVariants();
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
        `/keywords/${encodeURIComponent(keyword.trim())}/ideas`,
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
      <PageHeader
        title="Video Ideas"
        description="Enter a keyword to get AI-powered video ideas based on what's already working for top competitors."
      />

      <form onSubmit={handleGenerate} className="flex gap-3">
        <Input
          type="text"
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          placeholder="e.g. crypto trading for beginners"
          className="flex-1"
        />
        <Button type="submit" loading={loading} disabled={!keyword.trim()}>
          Generate Ideas
        </Button>
      </form>

      {error && <Alert variant="danger">{error}</Alert>}

      {loading && (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {[...Array(6)].map((_, i) => (
            <Card key={i} className="space-y-3">
              <Skeleton className="h-4 w-3/4" />
              <Skeleton className="h-3 w-full" />
              <Skeleton className="h-3 w-1/2" />
            </Card>
          ))}
        </div>
      )}

      {!loading && ideas.length > 0 && (
        <>
          <p className="text-sm text-gray-500">
            {ideas.length} ideas generated for &ldquo;<span className="text-white">{searched}</span>&rdquo;
          </p>
          <motion.div
            initial="hidden"
            animate="visible"
            variants={staggerContainer}
            className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3"
          >
            {ideas.map((idea, i) => {
              const Icon = FORMAT_ICONS[idea.format] ?? Lightbulb;
              return (
                <motion.div key={i} variants={fadeInUp}>
                  <Card hover className="flex h-full flex-col">
                    <div className="mb-3 flex items-center gap-2">
                      <span className="flex h-7 w-7 items-center justify-center rounded-lg bg-brand-600/15 text-brand-300">
                        <Icon className="h-3.5 w-3.5" aria-hidden="true" />
                      </span>
                      <span className="rounded-full bg-surface-300 px-2 py-0.5 text-xs text-gray-300">
                        {idea.format}
                      </span>
                      <Badge variant={CTR_BADGE[idea.ctr_potential] ?? "neutral"} className="ml-auto">
                        {idea.ctr_potential} CTR
                      </Badge>
                    </div>

                    <h3 className="mb-2 text-sm font-semibold leading-snug text-white">{idea.title}</h3>
                    <p className="flex-1 text-xs leading-relaxed text-gray-400">{idea.hook}</p>

                    <Link
                      href={`/dashboard/analyses/new?keyword=${encodeURIComponent(idea.title)}`}
                      className="mt-4 block rounded-lg border border-surface-300 py-2 text-center text-xs font-medium text-gray-300 transition-colors duration-150 hover:border-brand-500 hover:text-brand-300"
                    >
                      Use this idea →
                    </Link>
                  </Card>
                </motion.div>
              );
            })}
          </motion.div>
        </>
      )}

      {!loading && ideas.length === 0 && !error && (
        <EmptyState
          icon={<Lightbulb className="h-5 w-5" aria-hidden="true" />}
          title="Enter a keyword to generate video ideas"
          description="We study your top competitors' titles and use AI to suggest original angles that stand out."
        />
      )}
    </div>
  );
}
