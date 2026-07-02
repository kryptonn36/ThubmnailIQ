import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { MotionConfig } from "framer-motion";
import { ToastProvider } from "@/components/ui/Toast";
import "./globals.css";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
  display: "swap",
});

export const metadata: Metadata = {
  title: "ThumbnailIQ — YouTube Thumbnail Intelligence",
  description:
    "Know your thumbnail's performance before you publish. Competitive intelligence, objective scoring, and actionable guidance for YouTube creators.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className={inter.variable}>
      <body className="min-h-screen bg-surface font-sans text-gray-100 antialiased">
        {/* reducedMotion="user" respects the OS prefers-reduced-motion
            setting for every Framer Motion animation in the app, applied
            post-hydration so it can never cause a server/client mismatch
            (see lib/motion.ts for why that matters). */}
        <MotionConfig reducedMotion="user">
          <ToastProvider>{children}</ToastProvider>
        </MotionConfig>
      </body>
    </html>
  );
}
