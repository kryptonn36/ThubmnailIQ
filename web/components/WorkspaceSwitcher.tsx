"use client";

import { useEffect, useRef, useState } from "react";
import { Check, ChevronsUpDown, Loader2, Plus } from "lucide-react";
import { useWorkspace } from "@/hooks/useWorkspace";
import { getStoredUser } from "@/lib/auth";
import { cn } from "@/lib/cn";
import { Badge, type BadgeVariant } from "@/components/ui/Badge";
import type { Workspace, WorkspaceRole } from "@/types";

const ROLE_BADGE: Record<WorkspaceRole, BadgeVariant> = {
  owner: "info",
  admin: "warning",
  editor: "success",
  viewer: "neutral",
};

// ownerLabel answers "whose workspace is this?" from the viewer's
// perspective: their own workspaces say so explicitly, anyone else's are
// attributed to the owner by name.
function ownerLabel(ws: Workspace, viewerId: string | undefined): string {
  if (viewerId && ws.owner_id === viewerId) return "Your workspace";
  return `Owned by ${ws.owner_name || ws.owner_email}`;
}

function WorkspaceAvatar({ name, className }: { name: string; className?: string }) {
  return (
    <span
      className={cn(
        "flex shrink-0 items-center justify-center rounded-lg bg-brand-gradient font-bold text-white",
        className,
      )}
    >
      {(name || "?").charAt(0).toUpperCase()}
    </span>
  );
}

// A Render-style workspace switcher: shows the active workspace and who owns
// it, and expands into the full list of workspaces the user belongs to, each
// attributed to its owner, plus inline workspace creation.
export default function WorkspaceSwitcher() {
  const { workspace, workspaces, loading, switchWorkspace, createWorkspace } = useWorkspace();
  const [open, setOpen] = useState(false);
  const [creating, setCreating] = useState(false);
  const [newName, setNewName] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [createError, setCreateError] = useState<string | null>(null);
  const rootRef = useRef<HTMLDivElement>(null);
  const viewerId = getStoredUser()?.id;

  // Close on outside click / Escape, standard menu behavior.
  useEffect(() => {
    if (!open) return;
    function onPointerDown(e: MouseEvent) {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) setOpen(false);
    }
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === "Escape") setOpen(false);
    }
    document.addEventListener("mousedown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("mousedown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [open]);

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    const name = newName.trim();
    if (!name) return;
    setSubmitting(true);
    setCreateError(null);
    try {
      await createWorkspace(name);
      setNewName("");
      setCreating(false);
      setOpen(false);
    } catch {
      setCreateError("Could not create workspace.");
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) {
    return (
      <div className="mb-6 flex h-14 items-center gap-3 rounded-xl border border-surface-300 bg-surface-100 px-3">
        <Loader2 className="h-4 w-4 animate-spin text-gray-500" aria-hidden="true" />
        <span className="text-sm text-gray-500">Loading workspace…</span>
      </div>
    );
  }

  if (!workspace) return null;

  return (
    <div ref={rootRef} className="relative mb-6">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        aria-expanded={open}
        aria-haspopup="listbox"
        className="flex w-full items-center gap-3 rounded-xl border border-surface-300 bg-surface-100 px-3 py-2.5 text-left transition-colors hover:border-surface-400 hover:bg-surface-200"
      >
        <WorkspaceAvatar name={workspace.name} className="h-8 w-8 text-sm" />
        <span className="min-w-0 flex-1">
          <span className="block truncate text-sm font-semibold text-white">{workspace.name}</span>
          <span className="block truncate text-xs text-gray-500">
            {ownerLabel(workspace, viewerId)}
          </span>
        </span>
        <ChevronsUpDown className="h-4 w-4 shrink-0 text-gray-500" aria-hidden="true" />
      </button>

      {open && (
        <div
          role="listbox"
          aria-label="Switch workspace"
          className="absolute left-0 right-0 top-full z-50 mt-2 overflow-hidden rounded-xl border border-surface-300 bg-surface-100 shadow-xl"
        >
          <p className="px-3 pb-1 pt-3 text-[10px] font-semibold uppercase tracking-widest text-gray-600">
            Workspaces
          </p>
          <div className="max-h-64 overflow-y-auto p-1">
            {workspaces.map((ws) => {
              const active = ws.id === workspace.id;
              return (
                <button
                  key={ws.id}
                  type="button"
                  role="option"
                  aria-selected={active}
                  onClick={() => {
                    switchWorkspace(ws.id);
                    setOpen(false);
                  }}
                  className={cn(
                    "flex w-full items-center gap-3 rounded-lg px-2 py-2 text-left transition-colors",
                    active ? "bg-brand-600/15" : "hover:bg-surface-200",
                  )}
                >
                  <WorkspaceAvatar name={ws.name} className="h-7 w-7 text-xs" />
                  <span className="min-w-0 flex-1">
                    <span className="block truncate text-sm font-medium text-white">{ws.name}</span>
                    <span className="block truncate text-xs text-gray-500">
                      {ownerLabel(ws, viewerId)} · {ws.member_count}{" "}
                      {ws.member_count === 1 ? "member" : "members"}
                    </span>
                  </span>
                  <Badge variant={ROLE_BADGE[ws.role] ?? "neutral"}>{ws.role}</Badge>
                  {active && <Check className="h-4 w-4 shrink-0 text-brand-300" aria-hidden="true" />}
                </button>
              );
            })}
          </div>

          <div className="border-t border-surface-300 p-1">
            {creating ? (
              <form onSubmit={handleCreate} className="space-y-2 p-2">
                <input
                  autoFocus
                  value={newName}
                  onChange={(e) => setNewName(e.target.value)}
                  placeholder="Workspace name"
                  className="w-full rounded-lg border border-surface-300 bg-surface-50 px-3 py-1.5 text-sm text-white placeholder:text-gray-600 focus:border-brand-500 focus:outline-none"
                />
                {createError && <p className="text-xs text-danger">{createError}</p>}
                <div className="flex gap-2">
                  <button
                    type="submit"
                    disabled={!newName.trim() || submitting}
                    className="flex-1 rounded-lg bg-brand-600 px-3 py-1.5 text-xs font-semibold text-white transition-colors hover:bg-brand-500 disabled:opacity-50"
                  >
                    {submitting ? "Creating…" : "Create"}
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setCreating(false);
                      setCreateError(null);
                    }}
                    className="rounded-lg px-3 py-1.5 text-xs font-medium text-gray-400 transition-colors hover:bg-surface-200 hover:text-gray-200"
                  >
                    Cancel
                  </button>
                </div>
              </form>
            ) : (
              <button
                type="button"
                onClick={() => setCreating(true)}
                className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium text-gray-400 transition-colors hover:bg-surface-200 hover:text-gray-200"
              >
                <Plus className="h-4 w-4" aria-hidden="true" />
                Create workspace
              </button>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
