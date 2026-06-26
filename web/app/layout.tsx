import type { Metadata } from "next";
import "./globals.css";

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
    <html lang="en">
      <body className="min-h-screen bg-surface text-gray-100 antialiased">
        {children}
      </body>
    </html>
  );
}
