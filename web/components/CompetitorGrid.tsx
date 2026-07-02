/* eslint-disable @next/next/no-img-element */
import { motion } from "framer-motion";
import type { Competitor } from "@/types";
import { useMotionVariants } from "@/lib/motion";

interface CompetitorGridProps {
  competitors: Competitor[];
  userScore?: number | null;
  userRank?: number | null;
}

function formatViews(views: number): string {
  if (views >= 1_000_000) return `${(views / 1_000_000).toFixed(1)}M views`;
  if (views >= 1_000) return `${(views / 1_000).toFixed(1)}K views`;
  return `${views} views`;
}

export default function CompetitorGrid({ competitors, userScore, userRank }: CompetitorGridProps) {
  const { staggerContainer, fadeInUp } = useMotionVariants();

  if (!competitors || competitors.length === 0) {
    return <p className="text-sm text-gray-400">No competitor data available.</p>;
  }

  return (
    <motion.div
      initial="hidden"
      animate="visible"
      variants={staggerContainer}
      className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4"
    >
      {competitors.map((c) => {
        const isUserRank = userRank !== null && userRank !== undefined && c.rank_position === userRank;
        return (
          <motion.div
            key={c.video_id}
            variants={fadeInUp}
            className={`rounded-xl border bg-surface-100 p-3 shadow-card transition-all duration-200 ${
              isUserRank ? "border-brand-500 shadow-glow" : "border-surface-300 hover:border-surface-400"
            }`}
          >
            <div className="relative mb-2 overflow-hidden rounded-lg bg-surface-300">
              <img src={c.thumbnail_url} alt={c.video_title} className="aspect-video w-full object-cover" />
              <span className="absolute left-1.5 top-1.5 rounded-md bg-black/70 px-1.5 py-0.5 text-xs font-semibold text-white">
                #{c.rank_position}
              </span>
              <span className="absolute right-1.5 top-1.5 rounded-md bg-brand-600/90 px-1.5 py-0.5 text-xs font-bold text-white">
                {c.score}
              </span>
            </div>
            <p className="line-clamp-2 text-sm font-medium text-gray-200" title={c.video_title}>
              {c.video_title}
            </p>
            <p className="mt-1 text-xs text-gray-500">{c.channel_name}</p>
            <p className="text-xs text-gray-500">{formatViews(c.view_count)}</p>
            {isUserRank && (
              <p className="mt-1 text-xs font-semibold text-brand-300">Your thumbnail scored {userScore}</p>
            )}
          </motion.div>
        );
      })}
    </motion.div>
  );
}
