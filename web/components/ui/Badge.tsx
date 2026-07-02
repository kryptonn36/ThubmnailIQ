import type { ReactNode } from "react";
import { cn } from "@/lib/cn";

export type BadgeVariant =
  | "pending"
  | "processing"
  | "complete"
  | "failed"
  | "success"
  | "warning"
  | "danger"
  | "info"
  | "neutral";

const VARIANT_CLASSES: Record<BadgeVariant, string> = {
  pending: "bg-surface-300 text-gray-300",
  processing: "bg-warning/15 text-warning",
  complete: "bg-success/15 text-success",
  success: "bg-success/15 text-success",
  failed: "bg-danger/15 text-danger",
  danger: "bg-danger/15 text-danger",
  warning: "bg-warning/15 text-warning",
  info: "bg-info/15 text-info",
  neutral: "bg-surface-300 text-gray-400",
};

export function Badge({
  variant = "neutral",
  children,
  className,
}: {
  variant?: BadgeVariant;
  children: ReactNode;
  className?: string;
}) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-xs font-medium capitalize",
        VARIANT_CLASSES[variant],
        className,
      )}
    >
      {children}
    </span>
  );
}
