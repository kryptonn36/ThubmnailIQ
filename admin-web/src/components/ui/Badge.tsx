const VARIANTS: Record<string, string> = {
  active: "bg-emerald-500/20 text-emerald-300",
  complete: "bg-emerald-500/20 text-emerald-300",
  suspended: "bg-red-500/20 text-red-300",
  failed: "bg-red-500/20 text-red-300",
  pending: "bg-amber-500/20 text-amber-300",
  processing: "bg-amber-500/20 text-amber-300",
  deleted: "bg-red-500/20 text-red-300",
  default: "bg-surface-300 text-gray-300",
};

export function Badge({ status }: { status: string }) {
  const classes = VARIANTS[status.toLowerCase()] ?? VARIANTS.default;
  return (
    <span className={`inline-block rounded-full px-2.5 py-1 text-xs font-medium capitalize ${classes}`}>
      {status}
    </span>
  );
}
