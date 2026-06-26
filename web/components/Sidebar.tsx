"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const NAV_ITEMS = [
  { href: "/dashboard", label: "Overview", icon: "📊" },
  { href: "/dashboard/analyses", label: "Analyses", icon: "🖼️" },
  { href: "/dashboard/analyses/new", label: "New Analysis", icon: "➕" },
  { href: "/dashboard/compare", label: "Compare", icon: "⚖️" },
  { href: "/dashboard/tracking", label: "Tracking", icon: "📡" },
  { href: "/dashboard/database", label: "Viral Database", icon: "🔥" },
  { href: "/dashboard/billing", label: "Billing", icon: "💳" },
  { href: "/dashboard/settings", label: "Settings", icon: "⚙️" },
];

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="hidden w-60 shrink-0 border-r border-surface-300 bg-surface-50 px-4 py-6 lg:block">
      <Link href="/dashboard" className="mb-8 flex items-center gap-2 px-2">
        <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-brand-gradient text-sm font-bold text-white">
          IQ
        </span>
        <span className="text-lg font-semibold text-white">ThumbnailIQ</span>
      </Link>
      <nav className="space-y-1">
        {NAV_ITEMS.map((item) => {
          const active =
            item.href === "/dashboard"
              ? pathname === "/dashboard"
              : pathname?.startsWith(item.href);
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition ${
                active
                  ? "bg-brand-600/20 text-brand-300"
                  : "text-gray-400 hover:bg-surface-200 hover:text-gray-200"
              }`}
            >
              <span aria-hidden>{item.icon}</span>
              {item.label}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
}
