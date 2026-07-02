import { motion } from "framer-motion";
import type { ImpactLevel, Suggestion } from "@/types";
import { Card } from "@/components/ui/Card";
import { useMotionVariants } from "@/lib/motion";

const IMPACT_ORDER: ImpactLevel[] = ["high", "medium", "low"];

const IMPACT_STYLES: Record<ImpactLevel, { label: string; badge: string }> = {
  high: { label: "High Impact", badge: "bg-danger/15 text-danger border-danger/30" },
  medium: { label: "Medium Impact", badge: "bg-warning/15 text-warning border-warning/30" },
  low: { label: "Low Impact", badge: "bg-info/15 text-info border-info/30" },
};

interface SuggestionListProps {
  suggestions: Suggestion[];
}

export default function SuggestionList({ suggestions }: SuggestionListProps) {
  const { staggerContainer, fadeInUp } = useMotionVariants();

  if (!suggestions || suggestions.length === 0) {
    return <p className="text-sm text-gray-400">No suggestions available yet.</p>;
  }

  return (
    <motion.div initial="hidden" animate="visible" variants={staggerContainer} className="space-y-6">
      {IMPACT_ORDER.map((level) => {
        const group = suggestions.filter((s) => s.impact_level === level);
        if (group.length === 0) return null;
        const style = IMPACT_STYLES[level];
        return (
          <motion.div key={level} variants={fadeInUp}>
            <h4 className="mb-3 text-sm font-semibold uppercase tracking-wide text-gray-400">
              {style.label}
            </h4>
            <div className="space-y-3">
              {group.map((suggestion, idx) => (
                <Card key={`${suggestion.type}-${idx}`} className="p-4">
                  <div className="mb-2 flex items-start justify-between gap-3">
                    <p className="font-medium text-white">{suggestion.headline}</p>
                    <span
                      className={`shrink-0 rounded-full border px-2 py-0.5 text-xs font-semibold ${style.badge}`}
                    >
                      {style.label}
                    </span>
                  </div>
                  <p className="text-sm leading-relaxed text-gray-400">{suggestion.explanation}</p>
                  <p className="mt-2 text-xs font-medium text-brand-300">
                    Estimated CTR impact: +{suggestion.estimated_ctr_min}% to +
                    {suggestion.estimated_ctr_max}%
                  </p>
                </Card>
              ))}
            </div>
          </motion.div>
        );
      })}
    </motion.div>
  );
}
