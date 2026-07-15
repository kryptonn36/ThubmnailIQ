import Link from "next/link";

// Shared site footer for public (marketing/legal) pages. Server-safe: no
// client hooks, so pages that use it can stay server components and export
// their own metadata.
export default function Footer() {
  return (
    <footer className="border-t border-surface-300 px-6 py-8">
      <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-3 text-sm text-gray-500 sm:flex-row">
        <p>© {new Date().getFullYear()} ThumbnailIQ. All rights reserved.</p>
        <nav aria-label="Legal">
          <Link
            href="/privacy"
            className="rounded transition-colors hover:text-gray-300 focus-visible:outline focus-visible:outline-2 focus-visible:outline-brand-400"
          >
            Privacy Policy
          </Link>
        </nav>
      </div>
    </footer>
  );
}
