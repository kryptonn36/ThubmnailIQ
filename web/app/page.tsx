import Link from "next/link";
import Navbar from "@/components/Navbar";

const PILLARS = [
  {
    title: "Competitive Intelligence",
    description:
      "See exactly how your thumbnail performs against what's actually ranking for your target keyword.",
    icon: "🎯",
  },
  {
    title: "Objective Scoring",
    description:
      "Replace gut feeling with a repeatable 0-100 ThumbnailIQ Score built from six measurable sub-dimensions.",
    icon: "📈",
  },
  {
    title: "Actionable Guidance",
    description:
      "Get a ranked list of specific improvements, each with an estimated CTR impact range.",
    icon: "✅",
  },
  {
    title: "Historical Insight",
    description:
      "Track competitor channels and keywords over time to see what changes the winners make.",
    icon: "🕒",
  },
  {
    title: "Team Workflow",
    description:
      "Bring your whole team into the review process with shared workspaces and version comparisons.",
    icon: "🤝",
  },
];

export default function LandingPage() {
  return (
    <div className="min-h-screen">
      <Navbar />

      <main>
        <section className="relative overflow-hidden px-6 pt-24 pb-20 text-center">
          <div
            className="pointer-events-none absolute inset-0 -z-10 bg-brand-gradient opacity-20 blur-3xl"
            aria-hidden
          />
          <p className="mb-4 inline-block rounded-full border border-brand-500/40 bg-brand-500/10 px-4 py-1 text-sm font-medium text-brand-300">
            YouTube Thumbnail Intelligence
          </p>
          <h1 className="mx-auto max-w-3xl text-4xl font-bold leading-tight text-white sm:text-5xl">
            Know Your Thumbnail Before You Publish
          </h1>
          <p className="mx-auto mt-6 max-w-xl text-lg text-gray-400">
            ThumbnailIQ scores your thumbnail against real competitors, scores it
            objectively out of 100, and tells you exactly what to fix before you
            hit publish.
          </p>
          <div className="mt-10 flex items-center justify-center gap-4">
            <Link
              href="/register"
              className="rounded-lg bg-brand-gradient px-6 py-3 text-base font-semibold text-white shadow-glow transition hover:opacity-90"
            >
              Get Started Free
            </Link>
            <Link
              href="/login"
              className="rounded-lg border border-surface-300 px-6 py-3 text-base font-medium text-gray-200 transition hover:bg-surface-200"
            >
              Log In
            </Link>
          </div>
        </section>

        <section className="px-6 py-16">
          <div className="mx-auto max-w-6xl">
            <h2 className="mb-12 text-center text-2xl font-semibold text-white">
              Built on five strategic pillars
            </h2>
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {PILLARS.map((pillar) => (
                <div
                  key={pillar.title}
                  className="rounded-2xl border border-surface-300 bg-surface-100 p-6 transition hover:border-brand-500/50"
                >
                  <span className="mb-4 block text-3xl">{pillar.icon}</span>
                  <h3 className="mb-2 text-lg font-semibold text-white">
                    {pillar.title}
                  </h3>
                  <p className="text-sm text-gray-400">{pillar.description}</p>
                </div>
              ))}
            </div>
          </div>
        </section>

        <section className="px-6 py-20">
          <div className="mx-auto max-w-3xl rounded-3xl border border-brand-500/30 bg-brand-500/5 p-10 text-center">
            <h2 className="text-2xl font-semibold text-white">
              Stop guessing. Start scoring.
            </h2>
            <p className="mt-3 text-gray-400">
              Join creators who use data, not gut feeling, to win the click.
            </p>
            <Link
              href="/register"
              className="mt-6 inline-block rounded-lg bg-brand-gradient px-6 py-3 text-base font-semibold text-white shadow-glow transition hover:opacity-90"
            >
              Create Your Free Account
            </Link>
          </div>
        </section>
      </main>

      <footer className="border-t border-surface-300 px-6 py-8 text-center text-sm text-gray-500">
        © {new Date().getFullYear()} ThumbnailIQ. All rights reserved.
      </footer>
    </div>
  );
}
