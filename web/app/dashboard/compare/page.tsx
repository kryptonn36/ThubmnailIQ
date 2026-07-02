"use client";

import { useRef, useState } from "react";
import { motion } from "framer-motion";
import { ImagePlus } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import CompareGrid, { CompareItem } from "@/components/CompareGrid";
import { Alert } from "@/components/ui/Alert";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { useMotionVariants } from "@/lib/motion";
import type { AnalysisCreateResponse, AnalysisFull } from "@/types";

const MAX_ITEMS = 5;
const POLL_INTERVAL_MS = 2000;
const POLL_TIMEOUT_MS = 60000;

function pollForScore(id: string): Promise<number> {
  const startedAt = Date.now();
  return new Promise((resolve, reject) => {
    const tick = async () => {
      try {
        const data = await api.get<AnalysisFull>(`/analyses/${id}`);
        if (data.status === "complete") {
          resolve(data.score ?? 0);
          return;
        }
        if (data.status === "failed") {
          reject(new Error("Scoring failed"));
          return;
        }
        if (Date.now() - startedAt > POLL_TIMEOUT_MS) {
          reject(new Error("Timed out waiting for score"));
          return;
        }
        setTimeout(tick, POLL_INTERVAL_MS);
      } catch (err) {
        reject(err);
      }
    };
    tick();
  });
}

export default function ComparePage() {
  const { fadeInUp } = useMotionVariants();
  const [keyword, setKeyword] = useState("");
  const [items, setItems] = useState<CompareItem[]>([]);
  const [error, setError] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  function handleFiles(files: FileList | null) {
    if (!files) return;
    const remaining = MAX_ITEMS - items.length;
    const selected = Array.from(files).slice(0, remaining);
    const newItems: CompareItem[] = selected.map((file) => ({
      id: `${file.name}-${Date.now()}-${Math.random().toString(36).slice(2)}`,
      fileName: file.name,
      previewUrl: URL.createObjectURL(file),
      status: "idle",
    }));
    setItems((prev) => [...prev, ...newItems]);
    selected.forEach((file, idx) => scoreFile(file, newItems[idx].id));
  }

  async function scoreFile(file: File, itemId: string) {
    if (!keyword.trim()) {
      setError("Please enter a keyword before uploading thumbnails.");
      setItems((prev) => prev.map((i) => (i.id === itemId ? { ...i, status: "error" } : i)));
      return;
    }
    setItems((prev) => prev.map((i) => (i.id === itemId ? { ...i, status: "uploading" } : i)));
    try {
      const formData = new FormData();
      formData.append("thumbnail", file);
      formData.append("keyword", keyword.trim());
      const created = await api.postForm<AnalysisCreateResponse>("/analyses", formData);
      const score = await pollForScore(created.id);
      setItems((prev) => prev.map((i) => (i.id === itemId ? { ...i, status: "scored", score } : i)));
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to score one or more thumbnails.");
      setItems((prev) => prev.map((i) => (i.id === itemId ? { ...i, status: "error" } : i)));
    }
  }

  function handleRemove(id: string) {
    setItems((prev) => prev.filter((i) => i.id !== id));
  }

  return (
    <motion.div initial="hidden" animate="visible" variants={fadeInUp} className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-white sm:text-3xl">Compare Thumbnails</h1>
        <p className="mt-1.5 text-sm text-gray-400">
          Upload up to {MAX_ITEMS} thumbnail variants for the same keyword and see which one
          predicts the highest score.
        </p>
      </div>

      <Card>
        <label htmlFor="keyword" className="mb-1.5 block text-sm font-medium text-gray-300">
          Target keyword
        </label>
        <Input
          id="keyword"
          type="text"
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          placeholder="e.g. best budget gaming laptop 2026"
        />

        <div className="mt-4 flex items-center gap-3">
          <Button
            type="button"
            onClick={() => inputRef.current?.click()}
            disabled={items.length >= MAX_ITEMS}
            icon={<ImagePlus className="h-4 w-4" aria-hidden="true" />}
          >
            Add Thumbnails
          </Button>
          <span className="text-xs text-gray-500">{items.length}/{MAX_ITEMS} added</span>
          <input
            ref={inputRef}
            type="file"
            accept="image/png,image/jpeg"
            multiple
            className="hidden"
            onChange={(e) => handleFiles(e.target.files)}
          />
        </div>
      </Card>

      {error && <Alert variant="danger">{error}</Alert>}

      <CompareGrid items={items} onRemove={handleRemove} />
    </motion.div>
  );
}
