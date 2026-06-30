import type { Competitor } from "@/types";

interface NicheBenchmarkProps {
  score: number;
  competitors: Competitor[];
  keyword: string;
}

export default function NicheBenchmark({ score, competitors, keyword }: NicheBenchmarkProps) {
  if (competitors.length === 0) return null;

  const scores = competitors.map((c) => c.score ?? 0).filter((s) => s > 0);
  if (scores.length === 0) return null;

  const avg = Math.round(scores.reduce((a, b) => a + b, 0) / scores.length);
  const max = Math.max(...scores);
  const min = Math.min(...scores);
  const beaten = scores.filter((s) => score > s).length;
  const percentile = Math.round((beaten / scores.length) * 100);

  // Build 10-bucket histogram
  const bucketSize = 10;
  const buckets = Array.from({ length: 10 }, (_, i) => ({
    range: `${i * bucketSize}–${(i + 1) * bucketSize}`,
    low: i * bucketSize,
    count: scores.filter((s) => s >= i * bucketSize && s < (i + 1) * bucketSize).length,
    hasYours: score >= i * bucketSize && score < (i + 1) * bucketSize,
  }));
  const maxBucketCount = Math.max(...buckets.map((b) => b.count), 1);

  return (
    <div className="rounded-2xl border border-surface-300 bg-surface-100 p-6">
      <h3 className="mb-1 text-base font-semibold text-white">Niche Benchmark</h3>
      <p className="mb-5 text-xs text-gray-500">
        How your score compares to {scores.length} competitors for &ldquo;{keyword}&rdquo;
      </p>

      {/* Key stats */}
      <div className="mb-6 grid grid-cols-4 gap-3 text-center">
        {[
          { label: "Your score", value: score, accent: true },
          { label: "Niche avg", value: avg, accent: false },
          { label: "Top score", value: max, accent: false },
          { label: "Percentile", value: `${percentile}%`, accent: false },
        ].map(({ label, value, accent }) => (
          <div key={label} className={`rounded-xl p-3 ${accent ? "bg-brand-600/20 ring-1 ring-brand-500/40" : "bg-surface-200"}`}>
            <p className={`text-xl font-bold ${accent ? "text-brand-300" : "text-white"}`}>{value}</p>
            <p className="mt-0.5 text-[10px] text-gray-500">{label}</p>
          </div>
        ))}
      </div>

      {/* Distribution chart */}
      <div>
        <p className="mb-2 text-[11px] font-medium text-gray-500">Score distribution</p>
        <div className="flex items-end gap-1 h-20">
          {buckets.map((b) => (
            <div key={b.range} className="flex flex-1 flex-col items-center gap-0.5">
              <div
                className={`w-full rounded-t transition ${b.hasYours ? "bg-brand-500" : "bg-surface-300"}`}
                style={{ height: `${(b.count / maxBucketCount) * 100}%`, minHeight: b.count > 0 ? "4px" : "0" }}
                title={`${b.range}: ${b.count} competitors${b.hasYours ? " (you)" : ""}`}
              />
            </div>
          ))}
        </div>
        <div className="mt-1 flex justify-between text-[9px] text-gray-600">
          <span>0</span><span>50</span><span>100</span>
        </div>
      </div>

      <p className="mt-4 rounded-lg bg-surface-200 px-3 py-2 text-xs text-gray-400">
        {score > avg
          ? `✅ Your thumbnail scores ${score - avg} points above the niche average — top ${100 - percentile}% of competitors.`
          : `⚠ Your thumbnail scores ${avg - score} points below the niche average. Focus on the suggestions below to close the gap.`}
      </p>

      {/* Range bar */}
      <div className="mt-4">
        <div className="relative h-2 overflow-visible rounded-full bg-surface-300">
          <div
            className="absolute top-0 h-2 rounded-full bg-surface-400"
            style={{ left: `${min}%`, width: `${max - min}%` }}
          />
          <div
            className="absolute -top-1 h-4 w-1 rounded-full bg-brand-400 shadow-glow"
            style={{ left: `${Math.min(score, 99)}%` }}
            title={`Your score: ${score}`}
          />
        </div>
        <div className="mt-1 flex justify-between text-[9px] text-gray-600">
          <span>Min {min}</span><span>Avg {avg}</span><span>Max {max}</span>
        </div>
      </div>
    </div>
  );
}
