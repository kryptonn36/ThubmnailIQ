interface ScoreBand {
  label: string;
  textClass: string;
  ringClass: string;
  bgClass: string;
}

export function getScoreBand(score: number): ScoreBand {
  if (score >= 80) {
    return {
      label: "Optimized",
      textClass: "text-fuchsia-400",
      ringClass: "stroke-fuchsia-400",
      bgClass: "bg-fuchsia-500/10",
    };
  }
  if (score >= 60) {
    return {
      label: "Above Average",
      textClass: "text-emerald-400",
      ringClass: "stroke-emerald-400",
      bgClass: "bg-emerald-500/10",
    };
  }
  if (score >= 40) {
    return {
      label: "Room to Improve",
      textClass: "text-amber-400",
      ringClass: "stroke-amber-400",
      bgClass: "bg-amber-500/10",
    };
  }
  return {
    label: "Needs Major Work",
    textClass: "text-red-400",
    ringClass: "stroke-red-400",
    bgClass: "bg-red-500/10",
  };
}

interface ScoreGaugeProps {
  score: number;
  size?: number;
}

export default function ScoreGauge({ score, size = 200 }: ScoreGaugeProps) {
  const band = getScoreBand(score);
  const radius = (size - 20) / 2;
  const circumference = 2 * Math.PI * radius;
  const clamped = Math.max(0, Math.min(100, score));
  const offset = circumference - (clamped / 100) * circumference;

  return (
    <div className="flex flex-col items-center gap-3">
      <div className="relative" style={{ width: size, height: size }}>
        <svg width={size} height={size} className="-rotate-90">
          <circle
            cx={size / 2}
            cy={size / 2}
            r={radius}
            strokeWidth={14}
            fill="none"
            className="stroke-surface-300"
          />
          <circle
            cx={size / 2}
            cy={size / 2}
            r={radius}
            strokeWidth={14}
            fill="none"
            strokeLinecap="round"
            strokeDasharray={circumference}
            strokeDashoffset={offset}
            className={`${band.ringClass} transition-all duration-700 ease-out`}
          />
        </svg>
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <span className="text-4xl font-bold text-white">{Math.round(clamped)}</span>
          <span className="text-xs uppercase tracking-wide text-gray-400">/ 100</span>
        </div>
      </div>
      <span
        className={`rounded-full px-3 py-1 text-sm font-semibold ${band.bgClass} ${band.textClass}`}
      >
        {band.label}
      </span>
    </div>
  );
}
