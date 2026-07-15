"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Menu } from "lucide-react";
import Sidebar, { MobileSidebar } from "@/components/Sidebar";
import { Spinner } from "@/components/ui/Spinner";
import { isAuthenticated } from "@/lib/auth";
import { WorkspaceProvider } from "@/hooks/useWorkspace";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const [checked, setChecked] = useState(false);
  const [mobileNavOpen, setMobileNavOpen] = useState(false);

  useEffect(() => {
    if (!isAuthenticated()) {
      router.replace("/login");
      return;
    }
    setChecked(true);
  }, [router]);

  if (!checked) {
    return (
      <div className="flex min-h-screen items-center justify-center gap-2.5 text-gray-500">
        <Spinner className="h-4 w-4" />
        <span className="text-sm">Loading…</span>
      </div>
    );
  }

  return (
    <WorkspaceProvider>
    <div className="flex min-h-screen">
      <Sidebar />
      <MobileSidebar open={mobileNavOpen} onClose={() => setMobileNavOpen(false)} />
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex h-14 items-center gap-3 border-b border-surface-300 bg-surface-50 px-4 lg:hidden">
          <button
            type="button"
            onClick={() => setMobileNavOpen(true)}
            aria-label="Open navigation menu"
            className="rounded-lg p-2 text-gray-300 transition-colors hover:bg-surface-200 hover:text-white"
          >
            <Menu className="h-5 w-5" aria-hidden="true" />
          </button>
          <Link href="/dashboard" className="flex items-center gap-2">
            <span className="flex h-7 w-7 items-center justify-center rounded-lg bg-brand-gradient text-xs font-bold text-white">
              IQ
            </span>
            <span className="text-sm font-semibold text-white">ThumbnailIQ</span>
          </Link>
        </header>
        <main className="mx-auto w-full max-w-6xl flex-1 px-4 py-6 sm:px-6 sm:py-8">{children}</main>
      </div>
    </div>
    </WorkspaceProvider>
  );
}
