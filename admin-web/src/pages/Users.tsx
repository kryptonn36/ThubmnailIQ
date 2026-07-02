import { useState } from "react";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";
import { Badge } from "@/components/ui/Badge";
import { DataTable, type Column } from "@/components/table/DataTable";
import { useActivateUser, useDeleteUser, useSuspendUser, useUsers } from "@/hooks/useUsers";
import { extractErrorMessage } from "@/lib/axios";
import type { UserSummary } from "@/types";

function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB"];
  const i = Math.min(units.length - 1, Math.floor(Math.log(bytes) / Math.log(1024)));
  return `${(bytes / 1024 ** i).toFixed(1)} ${units[i]}`;
}

export function UsersPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState("");

  const { data, isLoading, error } = useUsers({ page, per_page: 20, search, status });
  const suspend = useSuspendUser();
  const activate = useActivateUser();
  const del = useDeleteUser();

  function runBulk(action: (id: string) => Promise<unknown>, label: string, rows: UserSummary[]) {
    Promise.allSettled(rows.map((r) => action(r.id))).then((results) => {
      const failed = results.filter((r) => r.status === "rejected").length;
      if (failed > 0) {
        toast.error(`${label}: ${failed} of ${rows.length} failed`);
      } else {
        toast.success(`${label}: ${rows.length} user(s) updated`);
      }
    });
  }

  const columns: Column<UserSummary>[] = [
    {
      key: "user",
      header: "User",
      render: (u) => (
        <div>
          <p className="font-medium text-white">{u.full_name || u.email}</p>
          <p className="text-xs text-gray-500">{u.email}</p>
        </div>
      ),
      sortValue: (u) => u.email,
    },
    { key: "status", header: "Status", render: (u) => <Badge status={u.status} /> },
    { key: "plan", header: "Plan", render: (u) => <span className="capitalize">{u.plan || "—"}</span> },
    { key: "analyses", header: "Uploads", render: (u) => u.analyses_count, sortValue: (u) => u.analyses_count },
    {
      key: "storage",
      header: "Storage",
      render: (u) => formatBytes(u.storage_used_bytes),
      sortValue: (u) => u.storage_used_bytes,
    },
    {
      key: "created",
      header: "Joined",
      render: (u) => new Date(u.created_at).toLocaleDateString(),
      sortValue: (u) => u.created_at,
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">Users</h1>
        <p className="mt-1 text-sm text-gray-400">Manage customer accounts.</p>
      </div>

      <DataTable
        columns={columns}
        rows={data?.data ?? []}
        getRowId={(u) => u.id}
        loading={isLoading}
        error={error ? extractErrorMessage(error) : null}
        emptyMessage="No users match your filters."
        page={page}
        perPage={20}
        total={data?.total ?? 0}
        onPageChange={setPage}
        search={search}
        onSearchChange={(v) => {
          setSearch(v);
          setPage(1);
        }}
        searchPlaceholder="Search by name or email…"
        filters={
          <select
            value={status}
            onChange={(e) => {
              setStatus(e.target.value);
              setPage(1);
            }}
            className="rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500"
          >
            <option value="">All statuses</option>
            <option value="active">Active</option>
            <option value="suspended">Suspended</option>
          </select>
        }
        onRowClick={(u) => navigate(`/users/${u.id}`)}
        bulkActions={[
          {
            label: "Suspend",
            onClick: (rows) => runBulk((id) => suspend.mutateAsync(id), "Suspend", rows),
          },
          {
            label: "Activate",
            onClick: (rows) => runBulk((id) => activate.mutateAsync(id), "Activate", rows),
          },
          {
            label: "Delete",
            variant: "danger",
            onClick: (rows) => {
              if (!confirm(`Delete ${rows.length} user(s)? This can't be undone from here.`)) return;
              runBulk((id) => del.mutateAsync(id), "Delete", rows);
            },
          },
        ]}
      />
    </div>
  );
}
