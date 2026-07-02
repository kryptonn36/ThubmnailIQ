import { useMemo, useState, type ReactNode } from "react";

export interface Column<T> {
  key: string;
  header: string;
  render: (row: T) => ReactNode;
  // sortValue is used for the client-side sort of the *current page* only —
  // the backend list endpoints always order by created_at desc server-side,
  // so sorting here re-orders whatever page is currently loaded rather than
  // the full dataset. That's an intentional, documented limitation rather
  // than adding dynamic ORDER BY to every sqlc query.
  sortValue?: (row: T) => string | number;
}

export interface BulkAction<T> {
  label: string;
  variant?: "primary" | "danger";
  onClick: (rows: T[]) => void;
}

interface DataTableProps<T> {
  columns: Column<T>[];
  rows: T[];
  getRowId: (row: T) => string;
  loading?: boolean;
  error?: string | null;
  emptyMessage?: string;
  page: number;
  perPage: number;
  total: number;
  onPageChange: (page: number) => void;
  search?: string;
  onSearchChange?: (value: string) => void;
  searchPlaceholder?: string;
  filters?: ReactNode;
  bulkActions?: BulkAction<T>[];
  onRowClick?: (row: T) => void;
}

export function DataTable<T>({
  columns,
  rows,
  getRowId,
  loading,
  error,
  emptyMessage = "Nothing here yet.",
  page,
  perPage,
  total,
  onPageChange,
  search,
  onSearchChange,
  searchPlaceholder = "Search…",
  filters,
  bulkActions,
  onRowClick,
}: DataTableProps<T>) {
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [sort, setSort] = useState<{ key: string; dir: "asc" | "desc" } | null>(null);

  const sortedRows = useMemo(() => {
    if (!sort) return rows;
    const column = columns.find((c) => c.key === sort.key);
    if (!column?.sortValue) return rows;
    const copy = [...rows];
    copy.sort((a, b) => {
      const av = column.sortValue!(a);
      const bv = column.sortValue!(b);
      if (av < bv) return sort.dir === "asc" ? -1 : 1;
      if (av > bv) return sort.dir === "asc" ? 1 : -1;
      return 0;
    });
    return copy;
  }, [rows, sort, columns]);

  const totalPages = Math.max(1, Math.ceil(total / perPage));
  const allSelected = rows.length > 0 && rows.every((r) => selected.has(getRowId(r)));

  function toggleAll() {
    if (allSelected) {
      setSelected(new Set());
    } else {
      setSelected(new Set(rows.map(getRowId)));
    }
  }

  function toggleRow(id: string) {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  }

  function handleSort(key: string) {
    setSort((prev) => {
      if (!prev || prev.key !== key) return { key, dir: "asc" };
      if (prev.dir === "asc") return { key, dir: "desc" };
      return null;
    });
  }

  const selectedRows = rows.filter((r) => selected.has(getRowId(r)));

  return (
    <div className="space-y-3">
      <div className="flex flex-wrap items-center gap-3">
        {onSearchChange && (
          <input
            type="text"
            value={search ?? ""}
            onChange={(e) => onSearchChange(e.target.value)}
            placeholder={searchPlaceholder}
            className="min-w-[200px] flex-1 rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500 sm:flex-none"
          />
        )}
        {filters}
      </div>

      {bulkActions && selectedRows.length > 0 && (
        <div className="flex items-center gap-3 rounded-lg border border-brand-500/30 bg-brand-500/10 px-4 py-2 text-sm text-brand-300">
          <span>{selectedRows.length} selected</span>
          {bulkActions.map((action) => (
            <button
              key={action.label}
              type="button"
              onClick={() => {
                action.onClick(selectedRows);
                setSelected(new Set());
              }}
              className={`rounded-md px-3 py-1 text-xs font-semibold ${
                action.variant === "danger"
                  ? "bg-red-500/20 text-red-300 hover:bg-red-500/30"
                  : "bg-brand-600/30 text-brand-200 hover:bg-brand-600/40"
              }`}
            >
              {action.label}
            </button>
          ))}
        </div>
      )}

      <div className="overflow-x-auto rounded-2xl border border-surface-300 bg-surface-100">
        <table className="w-full text-left text-sm">
          <thead>
            <tr className="border-b border-surface-300 text-xs uppercase tracking-wide text-gray-500">
              {bulkActions && (
                <th className="w-10 px-4 py-3">
                  <input type="checkbox" checked={allSelected} onChange={toggleAll} />
                </th>
              )}
              {columns.map((col) => (
                <th
                  key={col.key}
                  className={`px-4 py-3 font-medium ${col.sortValue ? "cursor-pointer select-none hover:text-gray-300" : ""}`}
                  onClick={col.sortValue ? () => handleSort(col.key) : undefined}
                >
                  {col.header}
                  {sort?.key === col.key ? (sort.dir === "asc" ? " ↑" : " ↓") : ""}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {loading ? (
              Array.from({ length: 5 }).map((_, i) => (
                <tr key={i} className="border-b border-surface-300/50">
                  <td colSpan={columns.length + (bulkActions ? 1 : 0)} className="px-4 py-4">
                    <div className="h-4 w-full animate-pulse rounded bg-surface-300" />
                  </td>
                </tr>
              ))
            ) : error ? (
              <tr>
                <td
                  colSpan={columns.length + (bulkActions ? 1 : 0)}
                  className="px-4 py-10 text-center text-sm text-red-400"
                >
                  {error}
                </td>
              </tr>
            ) : sortedRows.length === 0 ? (
              <tr>
                <td
                  colSpan={columns.length + (bulkActions ? 1 : 0)}
                  className="px-4 py-10 text-center text-sm text-gray-500"
                >
                  {emptyMessage}
                </td>
              </tr>
            ) : (
              sortedRows.map((row) => {
                const id = getRowId(row);
                return (
                  <tr
                    key={id}
                    onClick={() => onRowClick?.(row)}
                    className={`border-b border-surface-300/50 last:border-0 ${onRowClick ? "cursor-pointer hover:bg-surface-200/50" : ""}`}
                  >
                    {bulkActions && (
                      <td className="px-4 py-3" onClick={(e) => e.stopPropagation()}>
                        <input
                          type="checkbox"
                          checked={selected.has(id)}
                          onChange={() => toggleRow(id)}
                        />
                      </td>
                    )}
                    {columns.map((col) => (
                      <td key={col.key} className="px-4 py-3 text-gray-200">
                        {col.render(row)}
                      </td>
                    ))}
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>

      {!loading && !error && total > 0 && (
        <div className="flex items-center justify-between text-sm text-gray-400">
          <span>
            Page {page} of {totalPages} &middot; {total} total
          </span>
          <div className="flex gap-2">
            <button
              type="button"
              disabled={page <= 1}
              onClick={() => onPageChange(page - 1)}
              className="rounded-lg border border-surface-300 px-3 py-1.5 text-xs font-medium disabled:opacity-40"
            >
              Previous
            </button>
            <button
              type="button"
              disabled={page >= totalPages}
              onClick={() => onPageChange(page + 1)}
              className="rounded-lg border border-surface-300 px-3 py-1.5 text-xs font-medium disabled:opacity-40"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
