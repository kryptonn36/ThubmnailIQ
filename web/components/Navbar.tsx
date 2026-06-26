"use client";

import Link from "next/link";

export default function Navbar() {
  return (
    <header className="border-b border-surface-300 bg-surface/80 backdrop-blur">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <Link href="/" className="flex items-center gap-2">
          <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-brand-gradient text-sm font-bold text-white">
            IQ
          </span>
          <span className="text-lg font-semibold text-white">ThumbnailIQ</span>
        </Link>
        <nav className="flex items-center gap-3">
          <Link
            href="/login"
            className="rounded-lg px-4 py-2 text-sm font-medium text-gray-300 hover:text-white"
          >
            Log In
          </Link>
          <Link
            href="/register"
            className="rounded-lg bg-brand-gradient px-4 py-2 text-sm font-semibold text-white shadow-glow transition hover:opacity-90"
          >
            Get Started
          </Link>
        </nav>
      </div>
    </header>
  );
}
