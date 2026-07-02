import { Link } from "react-router-dom";
import { TrendChart } from "@/components/charts/TrendChart";
import { Badge } from "@/components/ui/Badge";
import { Card } from "@/components/ui/Card";
import { StatCard } from "@/components/ui/StatCard";
import { useAnalytics } from "@/hooks/useAnalytics";
import { useDashboard } from "@/hooks/useDashboard";
import { extractErrorMessage } from "@/lib/axios";
import type { SystemHealth } from "@/types";

function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB"];
  const i = Math.min(units.length - 1, Math.floor(Math.log(bytes) / Math.log(1024)));
  return `${(bytes / 1024 ** i).toFixed(1)} ${units[i]}`;
}

const HEALTH_LABELS: Record<keyof SystemHealth, string> = {
  database: "Database",
  redis: "Redis",
  cv_service: "CV Service",
};

function HealthPill({ label, healthy }: { label: string; healthy: boolean }) {
  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium ${
        healthy ? "bg-emerald-500/15 text-emerald-300" : "bg-red-500/15 text-red-300"
      }`}
    >
      <span className={`h-1.5 w-1.5 rounded-full ${healthy ? "bg-emerald-400" : "bg-red-400"}`} />
      {label}
    </span>
  );
}

export function DashboardPage() {
  const { data, isLoading, error } = useDashboard();
  const { data: analytics } = useAnalytics();

  if (isLoading) {
    return (
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
        {Array.from({ length: 6 }).map((_, i) => (
          <div key={i} className="h-24 animate-pulse rounded-2xl border border-surface-300 bg-surface-100" />
        ))}
      </div>
    );
  }

  if (error || !data) {
    return (
      <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">
        {error ? extractErrorMessage(error) : "Failed to load dashboard."}
      </p>
    );
  }

  return (
    <div className="space-y-8">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white">Dashboard</h1>
          <p className="mt-1 text-sm text-gray-400">Platform-wide overview.</p>
        </div>
        <div className="flex flex-wrap gap-2">
          {(Object.keys(HEALTH_LABELS) as (keyof SystemHealth)[]).map((key) => (
            <HealthPill key={key} label={HEALTH_LABELS[key]} healthy={data.system_health[key]} />
          ))}
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
        <StatCard label="Total users" value={data.total_users} />
        <StatCard label="Active users" value={data.active_users} />
        <StatCard label="Total uploads" value={data.total_uploads} />
        <StatCard label="Storage used" value={formatBytes(data.storage_used_bytes)} />
        <StatCard label="Uploads today" value={data.daily_uploads} />
        <StatCard label="Uploads this month" value={data.monthly_uploads} />
      </div>

      {analytics && (
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <TrendChart title="User growth (30 days)" data={analytics.daily_user_signups} />
          <TrendChart title="Upload growth (30 days)" data={analytics.daily_upload_trend} />
        </div>
      )}

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card>
          <div className="mb-3 flex items-center justify-between">
            <p className="text-sm font-semibold text-white">Recent users</p>
            <Link to="/users" className="text-xs text-brand-500 hover:underline">
              View all →
            </Link>
          </div>
          {data.recent_users.length === 0 ? (
            <p className="py-6 text-center text-sm text-gray-500">No users yet.</p>
          ) : (
            <div className="space-y-2">
              {data.recent_users.map((u) => (
                <Link
                  key={u.id}
                  to={`/users/${u.id}`}
                  className="flex items-center justify-between rounded-lg px-2 py-2 hover:bg-surface-200"
                >
                  <div className="min-w-0">
                    <p className="truncate text-sm font-medium text-gray-200">{u.full_name || u.email}</p>
                    <p className="truncate text-xs text-gray-500">{u.email}</p>
                  </div>
                  <Badge status={u.status} />
                </Link>
              ))}
            </div>
          )}
        </Card>

        <Card>
          <div className="mb-3 flex items-center justify-between">
            <p className="text-sm font-semibold text-white">Recent uploads</p>
            <Link to="/uploads" className="text-xs text-brand-500 hover:underline">
              View all →
            </Link>
          </div>
          {data.recent_uploads.length === 0 ? (
            <p className="py-6 text-center text-sm text-gray-500">No uploads yet.</p>
          ) : (
            <div className="space-y-2">
              {data.recent_uploads.map((u) => (
                <Link
                  key={u.id}
                  to={`/uploads/${u.id}`}
                  className="flex items-center justify-between rounded-lg px-2 py-2 hover:bg-surface-200"
                >
                  <div className="min-w-0">
                    <p className="truncate text-sm font-medium text-gray-200">{u.keyword}</p>
                    <p className="text-xs text-gray-500">{u.score !== null ? `Score ${u.score}` : "—"}</p>
                  </div>
                  <Badge status={u.status} />
                </Link>
              ))}
            </div>
          )}
        </Card>
      </div>
    </div>
  );
}
