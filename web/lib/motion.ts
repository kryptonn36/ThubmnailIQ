"use client";

import type { Variants } from "framer-motion";

// Reduced-motion handling is done globally via <MotionConfig reducedMotion="user">
// in app/layout.tsx — NOT by branching these variants on useReducedMotion().
// Branching here would make the server (which can't see the OS preference)
// and the client (which can) render different initial `hidden` styles on
// the very first paint, producing a React hydration mismatch. MotionConfig
// instead lets every component render identical markup and only changes
// how Framer Motion *animates* between states, which is a client-only,
// post-hydration concern.

export const fadeInUp: Variants = {
  hidden: { opacity: 0, y: 16 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.4, ease: [0.16, 1, 0.3, 1] } },
};

export const fadeIn: Variants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1, transition: { duration: 0.3 } },
};

export const scaleIn: Variants = {
  hidden: { opacity: 0, scale: 0.96 },
  visible: { opacity: 1, scale: 1, transition: { duration: 0.25, ease: [0.16, 1, 0.3, 1] } },
};

export const staggerContainer: Variants = {
  hidden: {},
  visible: { transition: { staggerChildren: 0.06, delayChildren: 0.05 } },
};

/**
 * Kept as a hook (rather than exporting the constants directly) so every
 * call site already using `const { fadeInUp } = useMotionVariants()`
 * continues to work unchanged.
 */
export function useMotionVariants() {
  return { fadeInUp, fadeIn, scaleIn, staggerContainer };
}
