/* eslint-disable @next/next/no-img-element */
import { motion } from "framer-motion";
import { Loader2, Trophy, X } from "lucide-react";
import { getScoreBand } from "@/components/ScoreGauge";
import { useMotionVariants } from "@/lib/motion";

export interface CompareItem {
  id: string;
  fileName: string;
  previewUrl: string;
  score?: number | null;
  status: "idle" | "uploading" | "scored" | "error";
}

interface CompareGridProps {
  items: CompareItem[];
  onRemove?: (id: string) => void;
}

export default function CompareGrid({ items, onRemove }: CompareGridProps) {
  const { staggerContainer, fadeInUp } = useMotionVariants();

  if (items.length === 0) {
    return <p className="text-sm text-gray-400">Add up to 5 thumbnails to compare their predicted scores.</p>;
  }

  const scored = items.filter((i) => i.status === "scored" && typeof i.score === "number");
  const winnerId =
    scored.length > 0
      ? scored.reduce((best, item) => ((item.score ?? 0) > (best.score ?? 0) ? item : best)).id
      : null;

  return (
    <motion.div
      initial="hidden"
      animate="visible"
      variants={staggerContainer}
      className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3"
    >
      {items.map((item) => {
        const isWinner = item.id === winnerId;
        const band = typeof item.score === "number" ? getScoreBand(item.score) : null;
        return (
          <motion.div
            key={item.id}
            variants={fadeInUp}
            className={`relative rounded-xl border bg-surface-100 p-4 shadow-card transition-all duration-200 ${
              isWinner ? "border-brand-500 shadow-glow" : "border-surface-300"
            }`}
          >
            {isWinner && (
              <span className="absolute -top-3 left-4 flex items-center gap-1 rounded-full bg-brand-gradient px-3 py-1 text-xs font-bold text-white">
                <Trophy className="h-3 w-3" aria-hidden="true" />
                Predicted Winner
              </span>
            )}
            {onRemove && (
              <button
                onClick={() => onRemove(item.id)}
                aria-label={`Remove ${item.fileName}`}
                className="absolute right-3 top-3 flex h-6 w-6 items-center justify-center rounded-full bg-black/50 text-gray-300 transition-colors hover:text-white"
              >
                <X className="h-3.5 w-3.5" aria-hidden="true" />
              </button>
            )}
            <div className="mb-3 overflow-hidden rounded-lg bg-surface-300">
              <img src={item.previewUrl} alt={item.fileName} className="aspect-video w-full object-cover" />
            </div>
            <p className="truncate text-sm text-gray-300">{item.fileName}</p>
            <div className="mt-2 flex items-center justify-between">
              {item.status === "uploading" && (
                <span className="flex items-center gap-1.5 text-xs text-gray-500">
                  <Loader2 className="h-3 w-3 animate-spin" aria-hidden="true" />
                  Scoring…
                </span>
              )}
              {item.status === "error" && <span className="text-xs text-danger">Failed to score</span>}
              {item.status === "scored" && band && (
                <>
                  <span className={`text-2xl font-bold tracking-tight ${band.textClass}`}>{item.score}</span>
                  <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${band.bgClass} ${band.textClass}`}>
                    {band.label}
                  </span>
                </>
              )}
            </div>
          </motion.div>
        );
      })}
    </motion.div>
  );
}
