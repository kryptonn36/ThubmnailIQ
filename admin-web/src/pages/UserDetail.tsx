import { useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import toast from "react-hot-toast";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { DataTable, type Column } from "@/components/table/DataTable";
import {
  useActivateUser,
  useChangeUserRole,
  useDeleteUser,
  useResetUserPassword,
  useSuspendUser,
  useUser,
  useUserUploads,
} from "@/hooks/useUsers";
import { extractErrorMessage } from "@/lib/axios";
import type { UploadSummary } from "@/types";

function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB"];
  const i = Math.min(units.length - 1, Math.floor(Math.log(bytes) / Math.log(1024)));
  return `${(bytes / 1024 ** i).toFixed(1)} ${units[i]}`;
}

export function UserDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [uploadsPage, setUploadsPage] = useState(1);
  const [tempPassword, setTempPassword] = useState<string | null>(null);
  const [role, setRole] = useState("editor");

  const { data: user, isLoading, error, refetch } = useUser(id);
  const { data: uploads, isLoading: uploadsLoading } = useUserUploads(id, { page: uploadsPage, per_page: 10 });

  const suspend = useSuspendUser();
  const activate = useActivateUser();
  const del = useDeleteUser();
  const resetPassword = useResetUserPassword();
  const changeRole = useChangeUserRole();

  if (isLoading) return <p className="text-sm text-gray-500">Loading…</p>;
  if (error || !user) {
    return (
      <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">
        {error ? extractErrorMessage(error) : "User not found."}
      </p>
    );
  }

  function handleAction(promise: Promise<unknown>, successMsg: string) {
    promise
      .then(() => {
        toast.success(successMsg);
        void refetch();
      })
      .catch((err) => toast.error(extractErrorMessage(err)));
  }

  const uploadColumns: Column<UploadSummary>[] = [
    { key: "keyword", header: "Keyword", render: (u) => u.keyword },
    { key: "status", header: "Status", render: (u) => <Badge status={u.status} /> },
    { key: "score", header: "Score", render: (u) => u.score ?? "—" },
    {
      key: "created",
      header: "Created",
      render: (u) => new Date(u.created_at).toLocaleDateString(),
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">{user.full_name || user.email}</h1>
          <p className="mt-1 text-sm text-gray-400">{user.email}</p>
        </div>
        <Badge status={user.status} />
      </div>

      <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
        <Card>
          <p className="text-xs uppercase text-gray-500">Plan</p>
          <p className="mt-1 text-lg font-semibold capitalize text-white">{user.plan || "—"}</p>
        </Card>
        <Card>
          <p className="text-xs uppercase text-gray-500">Uploads</p>
          <p className="mt-1 text-lg font-semibold text-white">{user.analyses_count}</p>
        </Card>
        <Card>
          <p className="text-xs uppercase text-gray-500">Storage used</p>
          <p className="mt-1 text-lg font-semibold text-white">{formatBytes(user.storage_used_bytes)}</p>
        </Card>
        <Card>
          <p className="text-xs uppercase text-gray-500">Joined</p>
          <p className="mt-1 text-lg font-semibold text-white">
            {new Date(user.created_at).toLocaleDateString()}
          </p>
        </Card>
      </div>

      <Card>
        <p className="mb-3 text-sm font-semibold text-white">Actions</p>
        <div className="flex flex-wrap gap-3">
          {user.status === "active" ? (
            <Button
              variant="secondary"
              loading={suspend.isPending}
              onClick={() => handleAction(suspend.mutateAsync(user.id), "User suspended")}
            >
              Suspend
            </Button>
          ) : (
            <Button
              variant="secondary"
              loading={activate.isPending}
              onClick={() => handleAction(activate.mutateAsync(user.id), "User activated")}
            >
              Activate
            </Button>
          )}
          <Button
            variant="secondary"
            loading={resetPassword.isPending}
            onClick={() =>
              resetPassword.mutate(user.id, {
                onSuccess: (res) => setTempPassword(res.temporary_password),
                onError: (err) => toast.error(extractErrorMessage(err)),
              })
            }
          >
            Reset password
          </Button>
          <Button
            variant="danger"
            loading={del.isPending}
            onClick={() => {
              if (!confirm("Delete this user? This can be undone by an engineer via the database only.")) return;
              del.mutate(user.id, {
                onSuccess: () => {
                  toast.success("User deleted");
                  navigate("/users");
                },
                onError: (err) => toast.error(extractErrorMessage(err)),
              });
            }}
          >
            Delete user
          </Button>
        </div>

        {tempPassword && (
          <div className="mt-4 rounded-lg border border-emerald-500/40 bg-emerald-500/10 p-3">
            <p className="text-sm font-medium text-emerald-300">
              Temporary password (copy now — shown only once):
            </p>
            <code className="mt-1 block break-all font-mono text-sm text-emerald-200">{tempPassword}</code>
          </div>
        )}

        {user.workspace_id && (
          <div className="mt-5 border-t border-surface-300 pt-4">
            <p className="mb-2 text-sm font-semibold text-white">Change workspace role</p>
            <div className="flex flex-wrap items-center gap-3">
              <select
                value={role}
                onChange={(e) => setRole(e.target.value)}
                className="rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500"
              >
                <option value="editor">Editor</option>
                <option value="viewer">Viewer</option>
              </select>
              <Button
                variant="secondary"
                loading={changeRole.isPending}
                onClick={() =>
                  handleAction(
                    changeRole.mutateAsync({ id: user.id, workspaceId: user.workspace_id!, role }),
                    "Role updated",
                  )
                }
              >
                Update role
              </Button>
            </div>
          </div>
        )}
      </Card>

      <div>
        <p className="mb-3 text-sm font-semibold text-white">Uploads</p>
        <DataTable
          columns={uploadColumns}
          rows={uploads?.data ?? []}
          getRowId={(u) => u.id}
          loading={uploadsLoading}
          emptyMessage="This user hasn't uploaded anything yet."
          page={uploadsPage}
          perPage={10}
          total={uploads?.total ?? 0}
          onPageChange={setUploadsPage}
          onRowClick={(u) => navigate(`/uploads/${u.id}`)}
        />
      </div>

      <Link to="/users" className="inline-block text-sm text-brand-500 hover:underline">
        ← Back to users
      </Link>
    </div>
  );
}
