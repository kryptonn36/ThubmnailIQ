"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Check, Video } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { useWorkspace } from "@/hooks/useWorkspace";
import { Alert } from "@/components/ui/Alert";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { PageHeader } from "@/components/ui/PageHeader";
import { useMotionVariants } from "@/lib/motion";

const FONTS = ["Inter", "Roboto", "Poppins", "Montserrat", "Lato", "Open Sans", "Nunito", "Raleway"];

const PRESET_PALETTES = [
  { name: "ThumbnailIQ", primary: "#6366F1", secondary: "#8B5CF6" },
  { name: "YouTube Red", primary: "#FF0000", secondary: "#CC0000" },
  { name: "Ocean", primary: "#0EA5E9", secondary: "#0284C7" },
  { name: "Forest", primary: "#16A34A", secondary: "#15803D" },
  { name: "Sunset", primary: "#F97316", secondary: "#EA580C" },
  { name: "Rose", primary: "#F43F5E", secondary: "#E11D48" },
];

export default function BrandKitPage() {
  const { workspace } = useWorkspace();
  const { fadeInUp } = useMotionVariants();
  const [primary, setPrimary] = useState("#6366F1");
  const [secondary, setSecondary] = useState("#8B5CF6");
  const [font, setFont] = useState("Inter");
  const [saving, setSaving] = useState(false);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (workspace) {
      setPrimary(workspace.brand_primary_color || "#6366F1");
      setSecondary(workspace.brand_secondary_color || "#8B5CF6");
      setFont(workspace.brand_font || "Inter");
    }
  }, [workspace]);

  async function handleSave(e: React.FormEvent) {
    e.preventDefault();
    if (!workspace) return;
    setSaving(true);
    setError(null);
    setSaved(false);
    try {
      await api.patch(`/workspaces/${workspace.id}/brand`, {
        primary_color: primary,
        secondary_color: secondary,
        font,
      });
      setSaved(true);
      setTimeout(() => setSaved(false), 3000);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to save brand kit.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <motion.div initial="hidden" animate="visible" variants={fadeInUp} className="max-w-2xl space-y-8">
      <PageHeader
        title="Brand Kit"
        description="Set your brand colors and font so ThumbnailIQ can give design suggestions that fit your channel identity."
      />

      {/* Live preview */}
      <Card className="overflow-hidden !p-0">
        <div className="p-4 text-xs font-medium uppercase tracking-wide text-gray-500">Preview</div>
        <div
          className="flex items-center gap-4 p-6"
          style={{ background: `linear-gradient(135deg, ${primary} 0%, ${secondary} 100%)` }}
        >
          <div className="flex h-16 w-16 items-center justify-center rounded-xl bg-white/20 text-white">
            <Video className="h-7 w-7" aria-hidden="true" />
          </div>
          <div>
            <p className="text-lg font-bold text-white" style={{ fontFamily: font }}>
              Your Channel Name
            </p>
            <p className="text-sm text-white/80" style={{ fontFamily: font }}>
              How your brand colours look on a thumbnail
            </p>
          </div>
        </div>
        <div className="flex gap-3 px-6 py-4">
          <span
            className="rounded-full px-3 py-1 text-xs font-semibold text-white"
            style={{ background: primary }}
          >
            Primary
          </span>
          <span
            className="rounded-full px-3 py-1 text-xs font-semibold text-white"
            style={{ background: secondary }}
          >
            Secondary
          </span>
          <span
            className="rounded-full border border-surface-300 px-3 py-1 text-xs font-medium text-gray-300"
            style={{ fontFamily: font }}
          >
            {font} typeface
          </span>
        </div>
      </Card>

      <Card as="form" onSubmit={handleSave} className="space-y-6">
        {/* Quick presets */}
        <div>
          <p className="mb-3 text-sm font-medium text-gray-300">Quick palettes</p>
          <div className="flex flex-wrap gap-2">
            {PRESET_PALETTES.map((p) => (
              <button
                key={p.name}
                type="button"
                onClick={() => {
                  setPrimary(p.primary);
                  setSecondary(p.secondary);
                }}
                className="flex items-center gap-2 rounded-full border border-surface-300 px-3 py-1.5 text-xs text-gray-300 transition-colors duration-150 hover:border-white/30"
              >
                <span className="flex gap-0.5">
                  <span className="h-3 w-3 rounded-full" style={{ background: p.primary }} />
                  <span className="h-3 w-3 rounded-full" style={{ background: p.secondary }} />
                </span>
                {p.name}
              </button>
            ))}
          </div>
        </div>

        {/* Color pickers */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="mb-2 block text-sm font-medium text-gray-300">Primary colour</label>
            <div className="flex items-center gap-3 rounded-lg border border-surface-300 bg-surface-200 px-3 py-2">
              <input
                type="color"
                value={primary}
                onChange={(e) => setPrimary(e.target.value)}
                className="h-7 w-7 cursor-pointer rounded border-0 bg-transparent"
              />
              <span className="font-mono text-sm text-gray-300">{primary.toUpperCase()}</span>
            </div>
          </div>
          <div>
            <label className="mb-2 block text-sm font-medium text-gray-300">Secondary colour</label>
            <div className="flex items-center gap-3 rounded-lg border border-surface-300 bg-surface-200 px-3 py-2">
              <input
                type="color"
                value={secondary}
                onChange={(e) => setSecondary(e.target.value)}
                className="h-7 w-7 cursor-pointer rounded border-0 bg-transparent"
              />
              <span className="font-mono text-sm text-gray-300">{secondary.toUpperCase()}</span>
            </div>
          </div>
        </div>

        {/* Font picker */}
        <div>
          <label className="mb-2 block text-sm font-medium text-gray-300">Brand font</label>
          <div className="grid grid-cols-2 gap-2 sm:grid-cols-4">
            {FONTS.map((f) => (
              <button
                key={f}
                type="button"
                onClick={() => setFont(f)}
                className={`rounded-lg border py-2 text-sm transition-colors duration-150 ${
                  font === f
                    ? "border-brand-500 bg-brand-600/20 text-brand-300"
                    : "border-surface-300 bg-surface-200 text-gray-400 hover:text-gray-200"
                }`}
                style={{ fontFamily: f }}
              >
                {f}
              </button>
            ))}
          </div>
        </div>

        {error && <Alert variant="danger">{error}</Alert>}
        {saved && (
          <Alert variant="success">
            <span className="flex items-center gap-1.5">
              <Check className="h-4 w-4" aria-hidden="true" /> Brand kit saved
            </span>
          </Alert>
        )}

        <Button type="submit" loading={saving} disabled={!workspace}>
          Save Brand Kit
        </Button>
      </Card>
    </motion.div>
  );
}
