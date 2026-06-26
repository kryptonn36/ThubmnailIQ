"use client";

import { useCallback, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { api, ApiError } from "@/lib/api";
import type { AnalysisCreateResponse } from "@/types";

export default function NewAnalysisPage() {
  const router = useRouter();
  const [file, setFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [keyword, setKeyword] = useState("");
  const [dragActive, setDragActive] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleFile = useCallback((selected: File | null) => {
    if (!selected) return;
    setFile(selected);
    setPreviewUrl(URL.createObjectURL(selected));
  }, []);

  function handleDrop(e: React.DragEvent<HTMLDivElement>) {
    e.preventDefault();
    setDragActive(false);
    const dropped = e.dataTransfer.files?.[0] ?? null;
    handleFile(dropped);
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!file || !keyword.trim()) {
      setError("Please provide both a thumbnail image and a target keyword.");
      return;
    }
    setError(null);
    setSubmitting(true);
    try {
      const formData = new FormData();
      formData.append("thumbnail", file);
      formData.append("keyword", keyword.trim());
      const res = await api.postForm<AnalysisCreateResponse>("/analyses", formData);
      router.push(`/dashboard/analyses/${res.id}`);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to start analysis.");
      setSubmitting(false);
    }
  }

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">New Analysis</h1>
        <p className="mt-1 text-sm text-gray-400">
          Upload a thumbnail and target keyword to score it against real competitors.
        </p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-5">
        {error && (
          <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">{error}</p>
        )}

        <div
          onDragOver={(e) => {
            e.preventDefault();
            setDragActive(true);
          }}
          onDragLeave={() => setDragActive(false)}
          onDrop={handleDrop}
          onClick={() => inputRef.current?.click()}
          className={`flex min-h-[220px] cursor-pointer flex-col items-center justify-center rounded-2xl border-2 border-dashed p-6 text-center transition ${
            dragActive
              ? "border-brand-500 bg-brand-500/10"
              : "border-surface-300 bg-surface-100 hover:border-brand-500/50"
          }`}
        >
          <input
            ref={inputRef}
            type="file"
            accept="image/png,image/jpeg"
            className="hidden"
            onChange={(e) => handleFile(e.target.files?.[0] ?? null)}
          />
          {previewUrl ? (
            // eslint-disable-next-line @next/next/no-img-element
            <img
              src={previewUrl}
              alt="Thumbnail preview"
              className="max-h-48 rounded-lg object-contain"
            />
          ) : (
            <>
              <span className="mb-2 text-3xl">📤</span>
              <p className="text-sm font-medium text-gray-300">
                Drag and drop your thumbnail here
              </p>
              <p className="mt-1 text-xs text-gray-500">
                or click to browse — PNG/JPG, 1280×720 recommended
              </p>
            </>
          )}
        </div>

        <div>
          <label htmlFor="keyword" className="mb-1 block text-sm font-medium text-gray-300">
            Target keyword
          </label>
          <input
            id="keyword"
            type="text"
            required
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            placeholder="e.g. how to lose weight fast"
            className="w-full rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500"
          />
        </div>

        <button
          type="submit"
          disabled={submitting}
          className="w-full rounded-lg bg-brand-gradient px-4 py-2.5 text-sm font-semibold text-white shadow-glow transition hover:opacity-90 disabled:opacity-60"
        >
          {submitting ? "Starting analysis..." : "Analyze Thumbnail"}
        </button>
      </form>
    </div>
  );
}
