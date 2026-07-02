"use client";

import { useEffect } from "react";
import { motion } from "framer-motion";
import { AlertTriangle, RotateCcw } from "lucide-react";
import { Button } from "@/components/ui/Button";
import { useMotionVariants } from "@/lib/motion";
import { logger } from "@/lib/logger";

export default function ErrorPage({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  const { fadeInUp } = useMotionVariants();

  useEffect(() => {
    logger.error("unhandled route error", { message: error.message, digest: error.digest });
  }, [error]);

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-surface px-6 text-center">
      <motion.div initial="hidden" animate="visible" variants={fadeInUp} className="max-w-md">
        <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-danger/15">
          <AlertTriangle className="h-8 w-8 text-danger" aria-hidden="true" />
        </div>
        <p className="text-sm font-semibold uppercase tracking-widest text-danger">Error</p>
        <h1 className="mt-2 text-3xl font-bold text-white">Something went wrong</h1>
        <p className="mt-3 text-sm leading-relaxed text-gray-400">
          An unexpected error occurred. You can try again, or head back to your dashboard.
        </p>
        <div className="mt-8 flex justify-center">
          <Button
            variant="secondary"
            icon={<RotateCcw className="h-4 w-4" aria-hidden="true" />}
            onClick={() => reset()}
          >
            Try again
          </Button>
        </div>
      </motion.div>
    </div>
  );
}
