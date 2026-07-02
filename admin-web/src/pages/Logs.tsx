import { useState } from "react";
import { DataTable, type Column } from "@/components/table/DataTable";
import { useAuditLogs } from "@/hooks/useLogs";
import { extractErrorMessage } from "@/lib/axios";
import type { AuditLog } from "@/types";

export function LogsPage() {
  const [page, setPage] = useState(1);
  const { data, isLoading, error } = useAuditLogs({ page, per_page: 25 });

  const columns: Column<AuditLog>[] = [
    {
      key: "created_at",
      header: "When",
      render: (l) => new Date(l.created_at).toLocaleString(),
      sortValue: (l) => l.created_at,
    },
    { key: "action", header: "Action", render: (l) => <span className="font-mono text-xs">{l.action}</span> },
    { key: "target_type", header: "Target type", render: (l) => l.target_type },
    { key: "target_id", header: "Target ID", render: (l) => <span className="font-mono text-xs text-gray-400">{l.target_id}</span> },
    {
      key: "metadata",
      header: "Metadata",
      render: (l) =>
        l.metadata ? (
          <code className="text-xs text-gray-400">{JSON.stringify(l.metadata)}</code>
        ) : (
          <span className="text-gray-600">—</span>
        ),
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">Audit Logs</h1>
        <p className="mt-1 text-sm text-gray-400">A record of every admin action taken in this panel.</p>
      </div>

      <DataTable
        columns={columns}
        rows={data?.data ?? []}
        getRowId={(l) => l.id}
        loading={isLoading}
        error={error ? extractErrorMessage(error) : null}
        emptyMessage="No admin actions recorded yet."
        page={page}
        perPage={25}
        total={data?.total ?? 0}
        onPageChange={setPage}
      />
    </div>
  );
}
