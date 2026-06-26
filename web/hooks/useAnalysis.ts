"use client";

import { useEffect, useRef, useState } from "react";
import { api } from "@/lib/api";
import type { AnalysisFull } from "@/types";

interface UseAnalysisResult {
  analysis: AnalysisFull | null;
  loading: boolean;
  error: string | null;
}

const POLL_INTERVAL_MS = 2000;
const ACTIVE_STATUSES = new Set(["pending", "processing"]);

export function useAnalysis(id: string): UseAnalysisResult {
  const [analysis, setAnalysis] = useState<AnalysisFull | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function fetchAnalysis() {
      try {
        const data = await api.get<AnalysisFull>(`/analyses/${id}`);
        if (cancelled) return;
        setAnalysis(data);
        setError(null);
        if (ACTIVE_STATUSES.has(data.status)) {
          timerRef.current = setTimeout(fetchAnalysis, POLL_INTERVAL_MS);
        }
      } catch (err) {
        if (cancelled) return;
        setError(err instanceof Error ? err.message : "Failed to load analysis");
      } finally {
        if (!cancelled) setLoading(false);
      }
    }

    fetchAnalysis();

    return () => {
      cancelled = true;
      if (timerRef.current) clearTimeout(timerRef.current);
    };
  }, [id]);

  return { analysis, loading, error };
}
