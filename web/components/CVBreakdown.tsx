import type { CVResults, CompetitorAvg } from "@/types";

interface CVBreakdownProps {
  cvResults: CVResults;
  competitorAvg?: CompetitorAvg | null;
}

function StatRow({
  label,
  value,
  compare,
}: {
  label: string;
  value: string;
  compare?: string;
}) {
  return (
    <div className="flex items-center justify-between border-b border-surface-300 py-2.5 last:border-b-0">
      <span className="text-sm text-gray-400">{label}</span>
      <div className="text-right">
        <span className="text-sm font-medium text-white">{value}</span>
        {compare && <span className="ml-2 text-xs text-gray-500">{compare}</span>}
      </div>
    </div>
  );
}

export default function CVBreakdown({ cvResults, competitorAvg }: CVBreakdownProps) {
  return (
    <div className="rounded-xl border border-surface-300 bg-surface-100 p-5">
      <h3 className="mb-2 text-base font-semibold text-white">Computer Vision Breakdown</h3>
      <div>
        <StatRow label="Face count" value={String(cvResults.face_count)} compare={competitorAvg ? `avg ${competitorAvg.face_count}` : undefined} />
        <StatRow label="Primary emotion" value={cvResults.primary_emotion ?? "None detected"} />
        <StatRow
          label="Text detected"
          value={cvResults.text_detected ? "Yes" : "No"}
        />
        <StatRow
          label="Text density"
          value={`${cvResults.text_density_pct}%`}
          compare={competitorAvg ? `avg ${competitorAvg.text_density_pct}%` : undefined}
        />
        <StatRow label="Word count" value={String(cvResults.word_count)} />
        <StatRow
          label="Clutter score"
          value={cvResults.clutter_score.toFixed(1)}
          compare={competitorAvg ? `avg ${competitorAvg.clutter_score.toFixed(1)}` : undefined}
        />
        <StatRow label="Visual complexity" value={cvResults.visual_complexity.toFixed(1)} />
        <StatRow label="Contrast (WCAG)" value={cvResults.contrast_score.toFixed(1)} />
        <StatRow label="Brightness" value={cvResults.brightness_score.toFixed(1)} />
        <StatRow label="Saturation" value={cvResults.saturation_score.toFixed(1)} />
        <StatRow label="Object count" value={String(cvResults.object_count)} />
      </div>

      {cvResults.text_strings.length > 0 && (
        <div className="mt-4">
          <p className="mb-2 text-sm text-gray-400">Detected text</p>
          <div className="flex flex-wrap gap-2">
            {cvResults.text_strings.map((text, idx) => (
              <span
                key={`${text}-${idx}`}
                className="rounded-md bg-surface-300 px-2 py-1 text-xs text-gray-200"
              >
                {text}
              </span>
            ))}
          </div>
        </div>
      )}

      {cvResults.dominant_colors.length > 0 && (
        <div className="mt-4">
          <p className="mb-2 text-sm text-gray-400">Dominant colors</p>
          <div className="flex gap-2">
            {cvResults.dominant_colors.map((color) => (
              <div key={color.hex} className="flex flex-col items-center gap-1">
                <span
                  className="h-8 w-8 rounded-full border border-surface-400"
                  style={{ backgroundColor: color.hex }}
                  title={color.hex}
                />
                <span className="text-[10px] text-gray-500">
                  {color.percentage.toFixed(0)}%
                </span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
