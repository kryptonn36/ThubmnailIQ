"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import type { Workspace } from "@/types";

interface UseWorkspaceResult {
  workspace: Workspace | null;
  loading: boolean;
}

export function useWorkspace(): UseWorkspaceResult {
  const [workspace, setWorkspace] = useState<Workspace | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;
    api
      .get<Workspace[]>("/workspaces")
      .then((list) => {
        if (!cancelled && list.length > 0) setWorkspace(list[0]);
      })
      .catch(() => {})
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  return { workspace, loading };
}
