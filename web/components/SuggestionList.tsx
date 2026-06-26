import type { ImpactLevel, Suggestion } from "@/types";

const IMPACT_ORDER: ImpactLevel[] = ["high", "medium", "low"];

const IMPACT_STYLES: Record<ImpactLevel, { label: string; badge: string }> = {
  high: { label: "High Impact", badge: "bg-red-500/15 text-red-400 border-red-500/30" },
  medium: {
    label: "Medium Impact",
    badge: "bg-amber-500/15 text-amber-400 border-amber-500/30",
  },
  low: {
    label: "Low Impact",
    badge: "bg-sky-500/15 text-sky-400 border-sky-500/30",
  },
};

interface SuggestionListProps {
  suggestions: Suggestion[];
}

export default function SuggestionList({ suggestions }: SuggestionListProps) {
  if (!suggestions || suggestions.length === 0) {
    return (
      <p className="text-sm text-gray-400">No suggestions available yet.</p>
    );
  }

  return (
    <div className="space-y-6">
      {IMPACT_ORDER.map((level) => {
        const group = suggestions.filter((s) => s.impact_level === level);
        if (group.length === 0) return null;
        const style = IMPACT_STYLES[level];
        return (
          <div key={level}>
            <h4 className="mb-3 text-sm font-semibold uppercase tracking-wide text-gray-400">
              {style.label}
            </h4>
            <div className="space-y-3">
              {group.map((suggestion, idx) => (
                <div
                  key={`${suggestion.type}-${idx}`}
                  className="rounded-xl border border-surface-300 bg-surface-100 p-4"
                >
                  <div className="mb-2 flex items-start justify-between gap-3">
                    <p className="font-medium text-white">{suggestion.headline}</p>
                    <span
                      className={`shrink-0 rounded-full border px-2 py-0.5 text-xs font-semibold ${style.badge}`}
                    >
                      {style.label}
                    </span>
                  </div>
                  <p className="text-sm text-gray-400">{suggestion.explanation}</p>
                  <p className="mt-2 text-xs font-medium text-brand-300">
                    Estimated CTR impact: +{suggestion.estimated_ctr_min}% to +
                    {suggestion.estimated_ctr_max}%
                  </p>
                </div>
              ))}
            </div>
          </div>
        );
      })}
    </div>
  );
}
