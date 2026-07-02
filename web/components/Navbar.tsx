"use client";

import Link from "next/link";
import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/Button";

export default function Navbar() {
  return (
    <header className="sticky top-0 z-40 border-b border-surface-300 bg-surface/80 backdrop-blur-md">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
        <Link href="/" className="flex items-center gap-2.5 transition-opacity hover:opacity-80">
          <span className="flex h-8 w-8 items-center justify-center rounded-lg bg-brand-gradient text-sm font-bold text-white shadow-glow">
            IQ
          </span>
          <span className="text-lg font-semibold tracking-tight text-white">ThumbnailIQ</span>
        </Link>
        <nav className="flex items-center gap-2" aria-label="Account">
          <Link
            href="/login"
            className="rounded-lg px-4 py-2 text-sm font-medium text-gray-300 transition-colors hover:text-white"
          >
            Log In
          </Link>
          <Link href="/register">
            <Button size="md" icon={<ArrowRight className="h-4 w-4" aria-hidden="true" />} iconPosition="right">
              Get Started
            </Button>
          </Link>
        </nav>
      </div>
    </header>
  );
}
