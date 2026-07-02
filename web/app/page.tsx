"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import { ArrowRight, CheckCircle2, Handshake, History, Target, TrendingUp } from "lucide-react";
import Navbar from "@/components/Navbar";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { useMotionVariants } from "@/lib/motion";

const PILLARS = [
  {
    title: "Competitive Intelligence",
    description:
      "See exactly how your thumbnail performs against what's actually ranking for your target keyword.",
    icon: Target,
  },
  {
    title: "Objective Scoring",
    description:
      "Replace gut feeling with a repeatable 0-100 ThumbnailIQ Score built from six measurable sub-dimensions.",
    icon: TrendingUp,
  },
  {
    title: "Actionable Guidance",
    description:
      "Get a ranked list of specific improvements, each with an estimated CTR impact range.",
    icon: CheckCircle2,
  },
  {
    title: "Historical Insight",
    description:
      "Track competitor channels and keywords over time to see what changes the winners make.",
    icon: History,
  },
  {
    title: "Team Workflow",
    description:
      "Bring your whole team into the review process with shared workspaces and version comparisons.",
    icon: Handshake,
  },
];

export default function LandingPage() {
  const { fadeInUp, staggerContainer } = useMotionVariants();

  return (
    <div className="min-h-screen">
      <Navbar />

      <main>
        <section className="relative overflow-hidden px-6 pb-20 pt-24 text-center">
          <div
            className="pointer-events-none absolute inset-0 -z-10 bg-brand-gradient opacity-20 blur-3xl"
            aria-hidden
          />
          <motion.div initial="hidden" animate="visible" variants={staggerContainer}>
            <motion.p
              variants={fadeInUp}
              className="mb-4 inline-block rounded-full border border-brand-500/40 bg-brand-500/10 px-4 py-1 text-sm font-medium text-brand-300"
            >
              YouTube Thumbnail Intelligence
            </motion.p>
            <motion.h1
              variants={fadeInUp}
              className="mx-auto max-w-3xl text-4xl font-bold leading-[1.1] tracking-tight text-white sm:text-5xl md:text-6xl"
            >
              Know Your Thumbnail Before You Publish
            </motion.h1>
            <motion.p
              variants={fadeInUp}
              className="mx-auto mt-6 max-w-xl text-lg leading-relaxed text-gray-400"
            >
              ThumbnailIQ scores your thumbnail against real competitors, scores it
              objectively out of 100, and tells you exactly what to fix before you
              hit publish.
            </motion.p>
            <motion.div variants={fadeInUp} className="mt-10 flex items-center justify-center gap-4">
              <Link href="/register">
                <Button
                  size="lg"
                  icon={<ArrowRight className="h-4 w-4" aria-hidden="true" />}
                  iconPosition="right"
                >
                  Get Started Free
                </Button>
              </Link>
              <Link href="/login">
                <Button variant="secondary" size="lg">
                  Log In
                </Button>
              </Link>
            </motion.div>
          </motion.div>
        </section>

        <section className="px-6 py-16">
          <div className="mx-auto max-w-6xl">
            <h2 className="mb-12 text-center text-2xl font-semibold tracking-tight text-white sm:text-3xl">
              Built on five strategic pillars
            </h2>
            <motion.div
              initial="hidden"
              whileInView="visible"
              viewport={{ once: true, margin: "-80px" }}
              variants={staggerContainer}
              className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3"
            >
              {PILLARS.map((pillar) => {
                const Icon = pillar.icon;
                return (
                  <motion.div key={pillar.title} variants={fadeInUp}>
                    <Card hover className="h-full">
                      <div className="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-brand-600/15 text-brand-300">
                        <Icon className="h-5 w-5" aria-hidden="true" />
                      </div>
                      <h3 className="mb-2 text-lg font-semibold text-white">{pillar.title}</h3>
                      <p className="text-sm leading-relaxed text-gray-400">{pillar.description}</p>
                    </Card>
                  </motion.div>
                );
              })}
            </motion.div>
          </div>
        </section>

        <motion.section
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-80px" }}
          variants={fadeInUp}
          className="px-6 py-20"
        >
          <Card accent className="mx-auto max-w-3xl border-brand-500/30 bg-brand-500/5 p-10 text-center">
            <h2 className="text-2xl font-semibold tracking-tight text-white sm:text-3xl">
              Stop guessing. Start scoring.
            </h2>
            <p className="mt-3 text-gray-400">
              Join creators who use data, not gut feeling, to win the click.
            </p>
            <Link href="/register" className="mt-6 inline-block">
              <Button size="lg" icon={<ArrowRight className="h-4 w-4" aria-hidden="true" />} iconPosition="right">
                Create Your Free Account
              </Button>
            </Link>
          </Card>
        </motion.section>
      </main>

      <footer className="border-t border-surface-300 px-6 py-8 text-center text-sm text-gray-500">
        © {new Date().getFullYear()} ThumbnailIQ. All rights reserved.
      </footer>
    </div>
  );
}
