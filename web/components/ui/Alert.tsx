import type { ReactNode } from "react";
import { AlertCircle, CheckCircle2, Info, TriangleAlert } from "lucide-react";
import { cn } from "@/lib/cn";

type AlertVariant = "success" | "warning" | "danger" | "info";

const CONFIG: Record<AlertVariant, { icon: typeof Info; classes: string }> = {
  success: { icon: CheckCircle2, classes: "border-success/30 bg-success/10 text-success" },
  warning: { icon: TriangleAlert, classes: "border-warning/30 bg-warning/10 text-warning" },
  danger: { icon: AlertCircle, classes: "border-danger/30 bg-danger/10 text-danger" },
  info: { icon: Info, classes: "border-info/30 bg-info/10 text-info" },
};

export function Alert({
  variant = "danger",
  children,
  className,
}: {
  variant?: AlertVariant;
  children: ReactNode;
  className?: string;
}) {
  const { icon: Icon, classes } = CONFIG[variant];
  return (
    <div
      role={variant === "danger" ? "alert" : "status"}
      className={cn(
        "flex items-start gap-2.5 rounded-lg border px-4 py-3 text-sm leading-relaxed",
        classes,
        className,
      )}
    >
      <Icon className="mt-0.5 h-4 w-4 shrink-0" aria-hidden="true" />
      <div className="flex-1">{children}</div>
    </div>
  );
}
