"use client";

import { useRef, useState } from "react";
import { api, ApiError } from "@/lib/api";
import CompareGrid, { CompareItem } from "@/components/CompareGrid";
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
      setItems((prev) =>
        prev.map((i) => (i.id === itemId ? { ...i, status: "error" } : i))
      );
      return;
    }
    setItems((prev) =>
      prev.map((i) => (i.id === itemId ? { ...i, status: "uploading" } : i))
    );
    try {
      const formData = new FormData();
      formData.append("thumbnail", file);
      formData.append("keyword", keyword.trim());
      const created = await api.postForm<AnalysisCreateResponse>("/analyses", formData);
      const score = await pollForScore(created.id);
      setItems((prev) =>
        prev.map((i) => (i.id === itemId ? { ...i, status: "scored", score } : i))
      );
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to score one or more thumbnails.");
      setItems((prev) =>
        prev.map((i) => (i.id === itemId ? { ...i, status: "error" } : i))
      );
    }
  }

  function handleRemove(id: string) {
    setItems((prev) => prev.filter((i) => i.id !== id));
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">Compare Thumbnails</h1>
        <p className="mt-1 text-sm text-gray-400">
          Upload up to {MAX_ITEMS} thumbnail variants for the same keyword and see which one
          predicts the highest score.
        </p>
      </div>

      <div className="rounded-2xl border border-surface-300 bg-surface-100 p-6">
        <label htmlFor="keyword" className="mb-1 block text-sm font-medium text-gray-300">
          Target keyword
        </label>
        <input
          id="keyword"
          type="text"
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          placeholder="e.g. best budget gaming laptop 2026"
          className="w-full rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500"
        />

        <div className="mt-4 flex items-center gap-3">
          <button
            type="button"
            onClick={() => inputRef.current?.click()}
            disabled={items.length >= MAX_ITEMS}
            className="rounded-lg bg-brand-gradient px-4 py-2.5 text-sm font-semibold text-white shadow-glow transition hover:opacity-90 disabled:opacity-50"
          >
            + Add Thumbnails
          </button>
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
      </div>

      {error && (
        <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">{error}</p>
      )}

      <CompareGrid items={items} onRemove={handleRemove} />
    </div>
  );
}
