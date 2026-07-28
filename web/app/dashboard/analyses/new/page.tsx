"use client";

import { useCallback, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { Sparkles, Upload } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { Alert } from "@/components/ui/Alert";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { useMotionVariants } from "@/lib/motion";
import type { AnalysisCreateResponse } from "@/types";

// Client-side image compression function
async function compressImage(
  file: File,
  options: {
    maxSizeMB?: number;
    maxWidthOrHeight?: number;
    useWebWorker?: boolean;
    initialQuality?: number;
  } = {}
): Promise<File> {
  const { maxSizeMB = 0.5, maxWidthOrHeight = 1920, useWebWorker = false, initialQuality = 0.8 } = options;

  // Return early if file is already small enough
  if (file.size <= maxSizeMB * 1024 * 1024) {
    return file;
  }

  // For simplicity and reliability, we'll use main-thread compression only
  // Web workers can cause issues in some environments, especially with Next.js
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.onerror = () => {
      console.error('Image loading error');
      reject(new Error('Failed to load image'));
    };
    img.onload = () => {
      try {
        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');

        if (!ctx) {
          reject(new Error('Could not get canvas context'));
          return;
        }

        // Calculate dimensions maintaining aspect ratio
        let width = img.width;
        let height = img.height;

        if (width > height && width > maxWidthOrHeight) {
          height = Math.round((height * maxWidthOrHeight) / width);
          width = maxWidthOrHeight;
        } else if (height > width && height > maxWidthOrHeight) {
          width = Math.round((width * maxWidthOrHeight) / height);
          height = maxWidthOrHeight;
        }

        canvas.width = width;
        canvas.height = height;

        // Draw image on canvas
        ctx.drawImage(img, 0, 0, width, height);

        // Convert to blob with quality adjustment
        canvas.toBlob(
          (blob) => {
            if (!blob) {
              reject(new Error('Failed to create blob from canvas'));
              return;
            }

            // Create File object from blob
            const compressedFile = new File([blob], file.name, {
              type: blob.type,
              lastModified: Date.now()
            });

            console.log(`Image compressed: ${Math.round(file.size / 1024)} KB → ${Math.round(compressedFile.size / 1024)} KB`);
            resolve(compressedFile);
          },
          file.type || 'image/jpeg',
          initialQuality // Use the initial quality without recursive reduction for simplicity
        );
      } catch (error) {
        reject(error);
      }
    };

    img.src = URL.createObjectURL(file);
  });
}

export default function NewAnalysisPage() {
  const router = useRouter();
  const { fadeInUp } = useMotionVariants();
  const [file, setFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [keyword, setKeyword] = useState("");
  const [dragActive, setDragActive] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleFile = useCallback((selected: File | null) => {
    if (!selected) return;
    console.log('File selected:', selected.name, `(${Math.round(selected.size / 1024)} KB)`);
    setFile(selected);
    setPreviewUrl(URL.createObjectURL(selected));
  }, []);

  function handleDrop(e: React.DragEvent<HTMLDivElement>) {
    e.preventDefault();
    setDragActive(false);
    const dropped = e.dataTransfer.files?.[0] ?? null;
    if (dropped) {
      console.log('File dropped:', dropped.name, `(${Math.round(dropped.size / 1024)} KB)`);
      handleFile(dropped);
    }
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
      console.log('Starting image compression...');
      // Compress the image before upload
      const compressedFile = await compressImage(file, {
        maxSizeMB: 0.5,          // target max size ~500KB
        maxWidthOrHeight: 1920,  // cap dimensions
        useWebWorker: false,     // disabled for reliability
        initialQuality: 0.8,
      });

      console.log(`Original file size: ${(file.size / 1024).toFixed(1)} KB`);
      console.log(`Compressed file size: ${(compressedFile.size / 1024).toFixed(1)} KB`);
      console.log(`Compression ratio: ${(file.size / compressedFile.size).toFixed(1)}x`);

      const formData = new FormData();
      formData.append("thumbnail", compressedFile);
      formData.append("keyword", keyword.trim());

      console.log('Submitting to API...');
      const res = await api.postForm<AnalysisCreateResponse>("/analyses", formData);
      console.log('Analysis created successfully:', res.id);
      router.push(`/dashboard/analyses/${res.id}`);
    } catch (err) {
      console.error('Error in handleSubmit:', err);
      setError(err instanceof ApiError ? err.message : "Failed to compress or upload image.");
      setSubmitting(false);
    }
  }

  return (
    <motion.div initial="hidden" animate="visible" variants={fadeInUp} className="mx-auto max-w-2xl space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-white sm:text-3xl">New Analysis</h1>
        <p className="mt-1.5 text-sm text-gray-400">
          Upload a thumbnail and target keyword to score it against real competitors.
        </p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-5">
        {error && <Alert variant="danger">{error}</Alert>}

        <div
          onDragOver={(e) => {
            e.preventDefault();
            setDragActive(true);
          }}
          onDragLeave={() => setDragActive(false)}
          onDrop={handleDrop}
          onClick={() => inputRef.current?.click()}
          role="button"
          tabIndex={0}
          onKeyDown={(e) => {
            if (e.key === "Enter" || e.key === " ") inputRef.current?.click();
          }}
          className={`flex min-h-[220px] cursor-pointer flex-col items-center justify-center rounded-2xl border-2 border-dashed p-6 text-center transition-colors duration-150 ${
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
            <img src={previewUrl} alt="Thumbnail preview" className="max-h-48 rounded-lg object-contain" />
          ) : (
            <>
              <div className="mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-brand-600/15 text-brand-300">
                <Upload className="h-5 w-5" aria-hidden="true" />
              </div>
              <p className="text-sm font-medium text-gray-300">Drag and drop your thumbnail here</p>
              <p className="mt-1 text-xs text-gray-500">
                or click to browse — PNG/JPG, 1280×720 recommended
              </p>
            </>
          )}
        </div>

        <div>
          <label htmlFor="keyword" className="mb-1.5 block text-sm font-medium text-gray-300">
            Target keyword
          </label>
          <Input
            id="keyword"
            type="text"
            required
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            placeholder="e.g. how to lose weight fast"
          />
        </div>

        <Button
          type="submit"
          loading={submitting}
          icon={<Sparkles className="h-4 w-4" aria-hidden="true" />}
          className="w-full"
        >
          Analyze Thumbnail
        </Button>
      </form>
    </motion.div>
  );
}
