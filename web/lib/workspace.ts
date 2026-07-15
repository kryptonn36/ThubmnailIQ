// The active workspace selection, persisted so it survives reloads and is
// readable by the API layer without importing React context. Lives in its own
// module (not hooks/useWorkspace) to avoid an import cycle: the workspace
// provider imports lib/api, and lib/api needs this value for every request.

const CURRENT_WORKSPACE_KEY = "tiq_workspace_id";

export function getActiveWorkspaceId(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(CURRENT_WORKSPACE_KEY);
}

export function storeActiveWorkspaceId(id: string): void {
  if (typeof window === "undefined") return;
  window.localStorage.setItem(CURRENT_WORKSPACE_KEY, id);
}
