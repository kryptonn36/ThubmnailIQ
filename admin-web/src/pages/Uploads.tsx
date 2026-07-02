import { useState } from "react";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";
import { Badge } from "@/components/ui/Badge";
import { DataTable, type Column } from "@/components/table/DataTable";
import { useDeleteUpload, useRestoreUpload, useUploads } from "@/hooks/useUploads";
import { extractErrorMessage } from "@/lib/axios";
import type { UploadSummary } from "@/types";

function formatBytes(bytes: number | null): string {
  if (!bytes) return "—";
  const units = ["B", "KB", "MB", "GB"];
  const i = Math.min(units.length - 1, Math.floor(Math.log(bytes) / Math.log(1024)));
  return `${(bytes / 1024 ** i).toFixed(1)} ${units[i]}`;
}

export function UploadsPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState("");
  const [includeDeleted, setIncludeDeleted] = useState(false);

  const { data, isLoading, error } = useUploads({
    page,
    per_page: 20,
    search,
    status,
    include_deleted: includeDeleted,
  });
  const del = useDeleteUpload();
  const restore = useRestoreUpload();

  function runBulk(action: (id: string) => Promise<unknown>, label: string, rows: UploadSummary[]) {
    Promise.allSettled(rows.map((r) => action(r.id))).then((results) => {
      const failed = results.filter((r) => r.status === "rejected").length;
      if (failed > 0) toast.error(`${label}: ${failed} of ${rows.length} failed`);
      else toast.success(`${label}: ${rows.length} upload(s) updated`);
    });
  }

  const columns: Column<UploadSummary>[] = [
    {
      key: "thumbnail",
      header: "",
      render: (u) => (
        <img src={u.thumbnail_url} alt={u.keyword} className="h-10 w-16 rounded object-cover" />
      ),
    },
    { key: "keyword", header: "Keyword", render: (u) => u.keyword, sortValue: (u) => u.keyword },
    { key: "status", header: "Status", render: (u) => <Badge status={u.status} /> },
    { key: "score", header: "Score", render: (u) => u.score ?? "—", sortValue: (u) => u.score ?? -1 },
    {
      key: "size",
      header: "Size",
      render: (u) => formatBytes(u.file_size_bytes),
      sortValue: (u) => u.file_size_bytes ?? 0,
    },
    {
      key: "created",
      header: "Created",
      render: (u) => new Date(u.created_at).toLocaleDateString(),
      sortValue: (u) => u.created_at,
    },
    {
      key: "download",
      header: "",
      render: (u) => (
        <a
          href={u.thumbnail_url}
          target="_blank"
          rel="noreferrer"
          onClick={(e) => e.stopPropagation()}
          className="text-xs font-medium text-brand-500 hover:underline"
        >
          Download
        </a>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">Uploads</h1>
        <p className="mt-1 text-sm text-gray-400">Moderate thumbnail analyses across all workspaces.</p>
      </div>

      <DataTable
        columns={columns}
        rows={data?.data ?? []}
        getRowId={(u) => u.id}
        loading={isLoading}
        error={error ? extractErrorMessage(error) : null}
        emptyMessage="No uploads match your filters."
        page={page}
        perPage={20}
        total={data?.total ?? 0}
        onPageChange={setPage}
        search={search}
        onSearchChange={(v) => {
          setSearch(v);
          setPage(1);
        }}
        searchPlaceholder="Search by keyword…"
        filters={
          <>
            <select
              value={status}
              onChange={(e) => {
                setStatus(e.target.value);
                setPage(1);
              }}
              className="rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500"
            >
              <option value="">All statuses</option>
              <option value="pending">Pending</option>
              <option value="processing">Processing</option>
              <option value="complete">Complete</option>
              <option value="failed">Failed</option>
            </select>
            <label className="flex items-center gap-2 text-sm text-gray-300">
              <input
                type="checkbox"
                checked={includeDeleted}
                onChange={(e) => {
                  setIncludeDeleted(e.target.checked);
                  setPage(1);
                }}
              />
              Show deleted
            </label>
          </>
        }
        onRowClick={(u) => navigate(`/uploads/${u.id}`)}
        bulkActions={[
          {
            label: "Restore",
            onClick: (rows) => runBulk((id) => restore.mutateAsync(id), "Restore", rows),
          },
          {
            label: "Delete",
            variant: "danger",
            onClick: (rows) => {
              if (!confirm(`Delete ${rows.length} upload(s)?`)) return;
              runBulk((id) => del.mutateAsync(id), "Delete", rows);
            },
          },
        ]}
      />
    </div>
  );
}
