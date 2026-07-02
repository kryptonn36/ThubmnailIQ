"use client";

import { forwardRef, type SelectHTMLAttributes } from "react";
import { ChevronDown } from "lucide-react";
import { cn } from "@/lib/cn";

export interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  error?: boolean;
}

export const Select = forwardRef<HTMLSelectElement, SelectProps>(function Select(
  { error, className, children, ...rest },
  ref,
) {
  return (
    <div className="relative">
      <select
        ref={ref}
        aria-invalid={error || undefined}
        className={cn(
          "w-full appearance-none rounded-lg border bg-surface-200 px-3 py-2.5 pr-9 text-sm text-white outline-none",
          "transition-colors duration-150",
          "focus:border-brand-500 focus:ring-2 focus:ring-brand-500/20",
          error
            ? "border-danger focus:border-danger focus:ring-danger/20"
            : "border-surface-300",
          className,
        )}
        {...rest}
      >
        {children}
      </select>
      <ChevronDown
        className="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-500"
        aria-hidden="true"
      />
    </div>
  );
});
