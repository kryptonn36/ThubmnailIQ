import { type ElementType, type HTMLAttributes, type ReactNode } from "react";
import { cn } from "@/lib/cn";

interface CardProps extends HTMLAttributes<HTMLElement> {
  /** Lifts and glows slightly on hover — use for clickable/interactive cards. */
  hover?: boolean;
  /** Adds a thin brand-gradient accent line along the top edge. */
  accent?: boolean;
  /** Render as a different element — e.g. "form" for a card-wrapped form. */
  as?: ElementType;
  children: ReactNode;
}

export function Card({
  hover = false,
  accent = false,
  as: Component = "div",
  className,
  children,
  ...rest
}: CardProps) {
  return (
    <Component
      className={cn(
        "rounded-2xl border border-surface-300 bg-surface-100 p-5 shadow-card transition-all duration-200 ease-out",
        hover && "hover:-translate-y-0.5 hover:border-surface-400 hover:shadow-card-hover",
        accent &&
          "relative overflow-hidden before:absolute before:inset-x-0 before:top-0 before:h-0.5 before:bg-brand-gradient",
        className,
      )}
      {...rest}
    >
      {children}
    </Component>
  );
}
