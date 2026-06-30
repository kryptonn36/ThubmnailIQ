"use client";

import { useEffect, useState } from "react";
import { api, ApiError } from "@/lib/api";
import { useWorkspace } from "@/hooks/useWorkspace";

interface Member {
  id: string;
  user_id: string;
  full_name: string;
  email: string;
  role: string;
  joined_at: string;
}

const ROLE_COLORS: Record<string, string> = {
  owner: "bg-brand-600/20 text-brand-300",
  editor: "bg-emerald-500/20 text-emerald-300",
  viewer: "bg-surface-300 text-gray-400",
};

export default function TeamPage() {
  const { workspace } = useWorkspace();
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [email, setEmail] = useState("");
  const [role, setRole] = useState<"editor" | "viewer">("editor");
  const [inviting, setInviting] = useState(false);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  async function loadMembers() {
    if (!workspace?.id) return;
    setLoading(true);
    try {
      const data = await api.get<Member[]>(`/workspaces/${workspace.id}/members`);
      setMembers(data);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to load members.");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadMembers();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [workspace?.id]);

  async function handleInvite(e: React.FormEvent) {
    e.preventDefault();
    if (!workspace?.id || !email.trim()) return;
    setInviting(true);
    setError(null);
    setSuccessMsg(null);
    try {
      await api.post(`/workspaces/${workspace.id}/members`, {
        email: email.trim(),
        role,
      });
      setEmail("");
      setSuccessMsg(`Invite sent to ${email.trim()}`);
      await loadMembers();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to send invite.");
    } finally {
      setInviting(false);
    }
  }

  return (
    <div className="max-w-2xl space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-white">Team</h1>
        <p className="mt-1 text-sm text-gray-400">
          Manage who has access to your workspace.
        </p>
      </div>

      {/* Invite form */}
      <form
        onSubmit={handleInvite}
        className="space-y-4 rounded-2xl border border-surface-300 bg-surface-100 p-6"
      >
        <h2 className="text-base font-semibold text-white">Invite a team member</h2>

        <div className="flex gap-3">
          <input
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="colleague@example.com"
            className="flex-1 rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500 placeholder:text-gray-600"
          />
          <select
            value={role}
            onChange={(e) => setRole(e.target.value as "editor" | "viewer")}
            className="rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500"
          >
            <option value="editor">Editor</option>
            <option value="viewer">Viewer</option>
          </select>
        </div>

        <div className="text-xs text-gray-500">
          <span className="text-gray-400 font-medium">Editor</span> — can create and view analyses.&nbsp;
          <span className="text-gray-400 font-medium">Viewer</span> — can only view results.
        </div>

        {error && <p className="text-sm text-red-400">{error}</p>}
        {successMsg && <p className="text-sm text-emerald-400">{successMsg}</p>}

        <button
          type="submit"
          disabled={inviting || !email.trim()}
          className="rounded-lg bg-brand-gradient px-5 py-2.5 text-sm font-semibold text-white shadow-glow transition hover:opacity-90 disabled:opacity-60"
        >
          {inviting ? "Sending invite..." : "Send Invite"}
        </button>
      </form>

      {/* Member list */}
      <div>
        <h2 className="mb-4 text-base font-semibold text-white">
          Members {members.length > 0 && <span className="ml-1 text-gray-500">({members.length})</span>}
        </h2>

        {loading ? (
          <div className="space-y-3">
            {[1,2].map(i => (
              <div key={i} className="animate-pulse h-16 rounded-xl border border-surface-300 bg-surface-100" />
            ))}
          </div>
        ) : members.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-surface-300 p-10 text-center text-sm text-gray-500">
            Only you are here. Invite someone to collaborate.
          </div>
        ) : (
          <div className="space-y-2">
            {members.map((m) => (
              <div
                key={m.id}
                className="flex items-center gap-4 rounded-xl border border-surface-300 bg-surface-100 px-4 py-3"
              >
                <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-brand-gradient text-sm font-bold text-white">
                  {(m.full_name || m.email).charAt(0).toUpperCase()}
                </div>
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-medium text-white">{m.full_name || m.email}</p>
                  <p className="truncate text-xs text-gray-500">{m.email}</p>
                </div>
                <span className={`rounded-full px-2.5 py-1 text-xs font-medium capitalize ${ROLE_COLORS[m.role] ?? "bg-surface-300 text-gray-400"}`}>
                  {m.role}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
