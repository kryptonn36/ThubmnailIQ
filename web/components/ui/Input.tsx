"use client";

import { forwardRef, type InputHTMLAttributes } from "react";
import { cn } from "@/lib/cn";

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  error?: boolean;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(function Input(
  { error, className, ...rest },
  ref,
) {
  return (
    <input
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
