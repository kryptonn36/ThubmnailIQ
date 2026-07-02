"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { UserPlus, Users } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { useWorkspace } from "@/hooks/useWorkspace";
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
  editor: "success",
  viewer: "neutral",
};

export default function TeamPage() {
  const { workspace } = useWorkspace();
  const { fadeInUp, staggerContainer } = useMotionVariants();
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
      <PageHeader title="Team" description="Manage who has access to your workspace." />

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
            onChange={(e) => setRole(e.target.value as "editor" | "viewer")}
            className="w-32"
          >
            <option value="editor">Editor</option>
            <option value="viewer">Viewer</option>
          </Select>
        </div>

        <div className="text-xs text-gray-500">
          <span className="font-medium text-gray-400">Editor</span> — can create and view analyses.&nbsp;
          <span className="font-medium text-gray-400">Viewer</span> — can only view results.
        </div>

        {error && <Alert variant="danger">{error}</Alert>}
        {successMsg && <Alert variant="success">{successMsg}</Alert>}

        <Button type="submit" loading={inviting} disabled={!email.trim()} icon={<UserPlus className="h-4 w-4" aria-hidden="true" />}>
          Send Invite
        </Button>
      </Card>

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
            {members.map((m) => (
              <motion.div
                key={m.id}
                variants={fadeInUp}
                className="flex items-center gap-4 rounded-xl border border-surface-300 bg-surface-100 px-4 py-3"
              >
                <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-brand-gradient text-sm font-bold text-white">
                  {(m.full_name || m.email).charAt(0).toUpperCase()}
                </div>
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-medium text-white">{m.full_name || m.email}</p>
                  <p className="truncate text-xs text-gray-500">{m.email}</p>
                </div>
                <Badge variant={ROLE_BADGE[m.role] ?? "neutral"}>{m.role}</Badge>
              </motion.div>
            ))}
          </motion.div>
        )}
      </div>
    </div>
  );
}
