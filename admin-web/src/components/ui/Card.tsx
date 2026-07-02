import type { ReactNode } from "react";

export function Card({ children, className = "" }: { children: ReactNode; className?: string }) {
  return (
    <div className={`rounded-2xl border border-surface-300 bg-surface-100 p-5 ${className}`}>
      {children}
    </div>
  );
}
