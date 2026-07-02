"use client";

import Link from "next/link";
import { motion } from "framer-motion";
import { Compass, Home } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { useMotionVariants } from "@/lib/motion";

export default function NotFound() {
  const { fadeInUp } = useMotionVariants();

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-surface px-6 text-center">
      <motion.div initial="hidden" animate="visible" variants={fadeInUp} className="max-w-md">
        <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-brand-gradient shadow-glow">
          <Compass className="h-8 w-8 text-white" aria-hidden="true" />
        </div>
        <p className="text-sm font-semibold uppercase tracking-widest text-brand-400">404</p>
        <h1 className="mt-2 text-3xl font-bold text-white">This page doesn&apos;t exist</h1>
        <p className="mt-3 text-sm leading-relaxed text-gray-400">
          The page you&apos;re looking for may have been moved, renamed, or never existed.
          Let&apos;s get you back on track.
        </p>
        <div className="mt-8 flex justify-center">
          <Link href="/">
            <Button icon={<Home className="h-4 w-4" aria-hidden="true" />}>Back to home</Button>
          </Link>
        </div>
      </motion.div>
    </div>
  );
}
