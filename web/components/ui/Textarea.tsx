"use client";

import { forwardRef, type TextareaHTMLAttributes } from "react";
import { cn } from "@/lib/cn";

export interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  error?: boolean;
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(function Textarea(
  { error, className, ...rest },
  ref,
) {
  return (
    <textarea
      ref={ref}
      aria-invalid={error || undefined}
      className={cn(
        "w-full rounded-lg border bg-surface-200 px-3 py-2.5 text-sm text-white outline-none placeholder:text-gray-600",
        "transition-colors duration-150",
        "focus:border-brand-500 focus:ring-2 focus:ring-brand-500/20",
        error
          ? "border-danger focus:border-danger focus:ring-danger/20"
          : "border-surface-300",
        className,
      )}
      {...rest}
    />
  );
});
