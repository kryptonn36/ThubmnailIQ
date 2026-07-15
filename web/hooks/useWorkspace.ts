"use client";

import {
  createContext,
  createElement,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { api } from "@/lib/api";
import { getActiveWorkspaceId, storeActiveWorkspaceId } from "@/lib/workspace";
import type { Workspace } from "@/types";

interface WorkspaceContextValue {
  /** The currently selected workspace (persisted across sessions). */
  workspace: Workspace | null;
  /** Every workspace the user belongs to, for the switcher. */
  workspaces: Workspace[];
  loading: boolean;
  switchWorkspace: (id: string) => void;
  createWorkspace: (name: string) => Promise<Workspace>;
  /** Re-fetch the list, e.g. after a rename or membership change. */
  refresh: () => Promise<void>;
}

const WorkspaceContext = createContext<WorkspaceContextValue | null>(null);

export function WorkspaceProvider({ children }: { children: React.ReactNode }) {
  const [workspaces, setWorkspaces] = useState<Workspace[]>([]);
  const [currentId, setCurrentId] = useState<string | null>(getActiveWorkspaceId);
  const [loading, setLoading] = useState(true);

  const refresh = useCallback(async () => {
    const list = await api.get<Workspace[]>("/workspaces");
    setWorkspaces(list);
  }, []);

  useEffect(() => {
    let cancelled = false;
    api
      .get<Workspace[]>("/workspaces")
      .then((list) => {
        if (!cancelled) setWorkspaces(list);
      })
      .catch(() => {})
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  // Fall back to the first workspace when nothing is selected yet or the
  // stored selection no longer exists (e.g. the user was removed from it).
  const workspace = useMemo(() => {
    if (workspaces.length === 0) return null;
    return workspaces.find((w) => w.id === currentId) ?? workspaces[0];
  }, [workspaces, currentId]);

  // Keep the persisted id in sync with what is actually active, so the API
  // layer's X-Workspace-ID header always matches what the UI shows.
  useEffect(() => {
    if (workspace && workspace.id !== currentId) {
      setCurrentId(workspace.id);
      storeActiveWorkspaceId(workspace.id);
    }
  }, [workspace, currentId]);

  const switchWorkspace = useCallback((id: string) => {
    setCurrentId(id);
    storeActiveWorkspaceId(id);
  }, []);

  const createWorkspace = useCallback(
    async (name: string) => {
      const created = await api.post<Workspace>("/workspaces", { name });
      await refresh();
      switchWorkspace(created.id);
      return created;
    },
    [refresh, switchWorkspace],
  );

  const value = useMemo(
    () => ({ workspace, workspaces, loading, switchWorkspace, createWorkspace, refresh }),
    [workspace, workspaces, loading, switchWorkspace, createWorkspace, refresh],
  );

  return createElement(WorkspaceContext.Provider, { value }, children);
}

export function useWorkspace(): WorkspaceContextValue {
  const ctx = useContext(WorkspaceContext);
  if (!ctx) {
    throw new Error("useWorkspace must be used within a WorkspaceProvider");
  }
  return ctx;
}
