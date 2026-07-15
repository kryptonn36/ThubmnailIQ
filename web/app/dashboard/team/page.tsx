"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Check, LogOut, Pencil, UserMinus, UserPlus, Users, X } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { useWorkspace } from "@/hooks/useWorkspace";
import { getStoredUser } from "@/lib/auth";
import { Alert } from "@/components/ui/Alert";
import { Badge, type BadgeVariant } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { EmptyState } from "@/components/ui/EmptyState";
import { Input } from "@/components/ui/Input";
import { PageHeader } from "@/components/ui/PageHeader";
import { Select } from "@/components/ui/Select";
import { Skeleton } from "@/components/ui/Skeleton";
import { useMotionVariants } from "@/lib/motion";

interface Member {
  id: string;
  user_id: string;
  full_name: string;
  email: string;
  role: string;
  joined_at: string;
}

const ROLE_BADGE: Record<string, BadgeVariant> = {
  owner: "info",
  admin: "warning",
  editor: "success",
  viewer: "neutral",
};

export default function TeamPage() {
  const { workspace, refresh } = useWorkspace();
  const { fadeInUp, staggerContainer } = useMotionVariants();
  const me = getStoredUser();

  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  const [email, setEmail] = useState("");
  const [role, setRole] = useState<"admin" | "editor" | "viewer">("editor");
  const [inviting, setInviting] = useState(false);

  const [editingName, setEditingName] = useState(false);
  const [nameDraft, setNameDraft] = useState("");
  const [savingName, setSavingName] = useState(false);

  const [removingId, setRemovingId] = useState<string | null>(null);

  // What the viewer may do here comes straight from their role in the
  // enriched workspace payload.
  const canManage = workspace?.role === "owner" || workspace?.role === "admin";
  const isOwner = workspace?.role === "owner";

  async function loadMembers() {
    if (!workspace?.id) return;
    setLoading(true);
    try {
      const data = await api.get<Member[]>(`/workspaces/${workspace.id}/members`);
      setMembers(data);
      setError(null);
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
      await Promise.all([loadMembers(), refresh()]);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to send invite.");
    } finally {
      setInviting(false);
    }
  }

  async function handleRename(e: React.FormEvent) {
    e.preventDefault();
    if (!workspace?.id || !nameDraft.trim()) return;
    setSavingName(true);
    setError(null);
    try {
      await api.patch(`/workspaces/${workspace.id}`, { name: nameDraft.trim() });
      setEditingName(false);
      await refresh();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to rename workspace.");
    } finally {
      setSavingName(false);
    }
  }

  async function handleRemove(member: Member) {
    if (!workspace?.id) return;
    const leaving = member.user_id === me?.id;
    const confirmed = window.confirm(
      leaving
        ? `Leave "${workspace.name}"? You'll lose access to this workspace.`
        : `Remove ${member.full_name || member.email} from "${workspace.name}"?`,
    );
    if (!confirmed) return;
    setRemovingId(member.user_id);
    setError(null);
    setSuccessMsg(null);
    try {
      await api.del(`/workspaces/${workspace.id}/members/${member.user_id}`);
      // Leaving drops this workspace from the list; the provider falls back
      // to the first remaining workspace automatically.
      await refresh();
      if (!leaving) {
        setSuccessMsg(`${member.full_name || member.email} was removed.`);
        await loadMembers();
      }
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to remove member.");
    } finally {
      setRemovingId(null);
    }
  }

  const ownedByYou = workspace && me ? workspace.owner_id === me.id : false;

  return (
    <div className="max-w-2xl space-y-8">
      <PageHeader title="Team" description="Manage who has access to this workspace." />

      {/* Workspace identity: always make it obvious WHOSE workspace this is. */}
      {workspace && (
        <Card className="flex items-center gap-4">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-brand-gradient text-lg font-bold text-white">
            {workspace.name.charAt(0).toUpperCase()}
          </div>
          <div className="min-w-0 flex-1">
            {editingName ? (
              <form onSubmit={handleRename} className="flex items-center gap-2">
                <Input
                  autoFocus
                  value={nameDraft}
                  onChange={(e) => setNameDraft(e.target.value)}
                  className="h-8 flex-1 text-sm"
                />
                <Button type="submit" loading={savingName} disabled={!nameDraft.trim()} className="h-8 px-2">
                  <Check className="h-4 w-4" aria-hidden="true" />
                </Button>
                <button
                  type="button"
                  onClick={() => setEditingName(false)}
                  aria-label="Cancel rename"
                  className="rounded-lg p-1.5 text-gray-400 transition-colors hover:bg-surface-200 hover:text-gray-200"
                >
                  <X className="h-4 w-4" aria-hidden="true" />
                </button>
              </form>
            ) : (
              <div className="flex items-center gap-2">
                <h2 className="truncate text-base font-semibold text-white">{workspace.name}</h2>
                {canManage && (
                  <button
                    type="button"
                    onClick={() => {
                      setNameDraft(workspace.name);
                      setEditingName(true);
                    }}
                    aria-label="Rename workspace"
                    className="rounded-lg p-1 text-gray-500 transition-colors hover:bg-surface-200 hover:text-gray-200"
                  >
                    <Pencil className="h-3.5 w-3.5" aria-hidden="true" />
                  </button>
                )}
              </div>
            )}
            <p className="truncate text-xs text-gray-500">
              {ownedByYou
                ? "Owned by you"
                : `Owned by ${workspace.owner_name || workspace.owner_email}`}
              {" · "}
              {workspace.member_count} {workspace.member_count === 1 ? "member" : "members"}
            </p>
          </div>
          <Badge variant={ROLE_BADGE[workspace.role] ?? "neutral"}>your role: {workspace.role}</Badge>
        </Card>
      )}

      {canManage && (
        <Card as="form" onSubmit={handleInvite} className="space-y-4">
          <h2 className="text-base font-semibold text-white">Invite a team member</h2>

          <div className="flex gap-3">
            <Input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="colleague@example.com"
              className="flex-1"
            />
            <Select
              value={role}
              onChange={(e) => setRole(e.target.value as "admin" | "editor" | "viewer")}
              className="w-32"
            >
              {isOwner && <option value="admin">Admin</option>}
              <option value="editor">Editor</option>
              <option value="viewer">Viewer</option>
            </Select>
          </div>

          <div className="text-xs text-gray-500">
            <span className="font-medium text-gray-400">Admin</span> — can manage members and settings.&nbsp;
            <span className="font-medium text-gray-400">Editor</span> — can create and view analyses.&nbsp;
            <span className="font-medium text-gray-400">Viewer</span> — can only view results.
          </div>

          <Button type="submit" loading={inviting} disabled={!email.trim()} icon={<UserPlus className="h-4 w-4" aria-hidden="true" />}>
            Send Invite
          </Button>
        </Card>
      )}

      {error && <Alert variant="danger">{error}</Alert>}
      {successMsg && <Alert variant="success">{successMsg}</Alert>}

      <div>
        <h2 className="mb-4 text-base font-semibold text-white">
          Members {members.length > 0 && <span className="ml-1 text-gray-500">({members.length})</span>}
        </h2>

        {loading ? (
          <div className="space-y-3">
            {[1, 2].map((i) => (
              <Skeleton key={i} className="h-16 rounded-xl" />
            ))}
          </div>
        ) : members.length === 0 ? (
          <EmptyState
            icon={<Users className="h-5 w-5" aria-hidden="true" />}
            title="Only you are here."
            description="Invite someone to collaborate."
          />
        ) : (
          <motion.div initial="hidden" animate="visible" variants={staggerContainer} className="space-y-2">
            {members.map((m) => {
              const isSelf = m.user_id === me?.id;
              const isRowOwner = m.role === "owner";
              // The owner can never be removed; others can be removed by
              // owner/admin, and anyone (but the owner) can remove themselves.
              const canRemove = !isRowOwner && (isSelf || canManage);
              return (
                <motion.div
                  key={m.id}
                  variants={fadeInUp}
                  className="flex items-center gap-4 rounded-xl border border-surface-300 bg-surface-100 px-4 py-3"
                >
                  <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-brand-gradient text-sm font-bold text-white">
                    {(m.full_name || m.email).charAt(0).toUpperCase()}
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="truncate text-sm font-medium text-white">
                      {m.full_name || m.email}
                      {isSelf && <span className="ml-1.5 text-xs font-normal text-gray-500">(you)</span>}
                    </p>
                    <p className="truncate text-xs text-gray-500">{m.email}</p>
                  </div>
                  <Badge variant={ROLE_BADGE[m.role] ?? "neutral"}>{m.role}</Badge>
                  {canRemove && (
                    <button
                      type="button"
                      onClick={() => handleRemove(m)}
                      disabled={removingId === m.user_id}
                      aria-label={isSelf ? "Leave workspace" : `Remove ${m.full_name || m.email}`}
                      title={isSelf ? "Leave workspace" : "Remove member"}
                      className="rounded-lg p-2 text-gray-500 transition-colors hover:bg-danger/15 hover:text-danger disabled:opacity-50"
                    >
                      {isSelf ? (
                        <LogOut className="h-4 w-4" aria-hidden="true" />
                      ) : (
                        <UserMinus className="h-4 w-4" aria-hidden="true" />
                      )}
                    </button>
                  )}
                </motion.div>
              );
            })}
          </motion.div>
        )}
      </div>
    </div>
  );
}
