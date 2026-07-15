"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { motion } from "framer-motion";
import {
  CreditCard,
  Flame,
  GalleryHorizontalEnd,
  LayoutDashboard,
  Lightbulb,
  Palette,
  Plus,
  Radar,
  Scale,
  Settings,
  Users,
  type LucideIcon,
} from "lucide-react";
import { cn } from "@/lib/cn";
import { useMotionVariants } from "@/lib/motion";
import WorkspaceSwitcher from "@/components/WorkspaceSwitcher";

interface NavItem {
  href: string;
  label: string;
  icon: LucideIcon;
}

interface NavSection {
  title: string;
  items: NavItem[];
}

const NAV_SECTIONS: NavSection[] = [
  {
    title: "Analyze",
    items: [
      { href: "/dashboard", label: "Overview", icon: LayoutDashboard },
      { href: "/dashboard/analyses/new", label: "New Analysis", icon: Plus },
      { href: "/dashboard/analyses", label: "All Analyses", icon: GalleryHorizontalEnd },
      { href: "/dashboard/compare", label: "Compare Versions", icon: Scale },
    ],
  },
  {
    title: "Research",
    items: [
      { href: "/dashboard/ideas", label: "Video Ideas", icon: Lightbulb },
      { href: "/dashboard/tracking", label: "Tracking", icon: Radar },
      { href: "/dashboard/database", label: "Viral Database", icon: Flame },
    ],
  },
  {
    title: "Workspace",
    items: [
      { href: "/dashboard/team", label: "Team", icon: Users },
      { href: "/dashboard/brand", label: "Brand Kit", icon: Palette },
      { href: "/dashboard/billing", label: "Billing", icon: CreditCard },
      { href: "/dashboard/settings", label: "Settings", icon: Settings },
    ],
  },
];

const ALL_ITEMS = NAV_SECTIONS.flatMap((s) => s.items);

// Several hrefs can be textual prefixes of the current path at once (e.g.
// "/dashboard/analyses" is a prefix of "/dashboard/analyses/new") — only the
// single most specific (longest) match should ever be highlighted, so
// "New Analysis" and "All Analyses" don't both light up together.
function getActiveHref(pathname: string | null): string | null {
  if (!pathname) return null;
  let best: string | null = null;
  for (const item of ALL_ITEMS) {
    const matches =
      item.href === "/dashboard"
        ? pathname === "/dashboard"
        : pathname === item.href || pathname.startsWith(`${item.href}/`);
    if (matches && (!best || item.href.length > best.length)) {
      best = item.href;
    }
  }
  return best;
}

function SidebarLogo() {
  return (
    <Link href="/dashboard" className="mb-8 flex items-center gap-2 px-2">
      <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-brand-gradient text-sm font-bold text-white shadow-glow">
        IQ
      </span>
      <span className="text-lg font-semibold tracking-tight text-white">ThumbnailIQ</span>
    </Link>
  );
}

function SidebarNav({ pathname, onNavigate }: { pathname: string | null; onNavigate?: () => void }) {
  const activeHref = getActiveHref(pathname);
  return (
    <nav className="space-y-5" aria-label="Main navigation">
      {NAV_SECTIONS.map((section) => (
        <div key={section.title}>
          <p className="mb-1 px-3 text-[10px] font-semibold uppercase tracking-widest text-gray-600">
            {section.title}
          </p>
          <div className="space-y-0.5">
            {section.items.map((item) => {
              const active = item.href === activeHref;
              const Icon = item.icon;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  onClick={onNavigate}
                  aria-current={active ? "page" : undefined}
                  className={cn(
                    "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors duration-150",
                    active
                      ? "bg-brand-600/20 text-brand-300"
                      : "text-gray-400 hover:bg-surface-200 hover:text-gray-200",
                  )}
                >
                  <Icon className="h-4 w-4 shrink-0" aria-hidden="true" />
                  {item.label}
                </Link>
              );
            })}
          </div>
        </div>
      ))}
    </nav>
  );
}

export default function Sidebar() {
  const pathname = usePathname();
  const { fadeIn } = useMotionVariants();

  return (
    <motion.aside
      initial="hidden"
      animate="visible"
      variants={fadeIn}
      className="hidden w-60 shrink-0 border-r border-surface-300 bg-surface-50 px-4 py-6 lg:block"
    >
      <SidebarLogo />
      <WorkspaceSwitcher />
      <SidebarNav pathname={pathname} />
    </motion.aside>
  );
}

// Slide-in mobile drawer — the desktop <Sidebar> above is `hidden` below
// the `lg` breakpoint with no mobile equivalent today, so this is new
// functionality (mobile nav access), not just a restyle.
export function MobileSidebar({ open, onClose }: { open: boolean; onClose: () => void }) {
  const pathname = usePathname();

  return (
    <motion.div
      initial={false}
      animate={open ? "visible" : "hidden"}
      className={cn("fixed inset-0 z-40 lg:hidden", !open && "pointer-events-none")}
    >
      <motion.div
        variants={{ hidden: { opacity: 0 }, visible: { opacity: 1 } }}
        transition={{ duration: 0.2 }}
        onClick={onClose}
        className="absolute inset-0 bg-black/60"
        aria-hidden="true"
      />
      <motion.aside
        variants={{ hidden: { x: "-100%" }, visible: { x: 0 } }}
        transition={{ duration: 0.25, ease: [0.16, 1, 0.3, 1] }}
        className="relative flex h-full w-72 flex-col overflow-y-auto border-r border-surface-300 bg-surface-50 px-4 py-6"
        role="dialog"
        aria-modal="true"
        aria-label="Navigation menu"
      >
        <SidebarLogo />
        <WorkspaceSwitcher />
        <SidebarNav pathname={pathname} onNavigate={onClose} />
      </motion.aside>
    </motion.div>
  );
}
