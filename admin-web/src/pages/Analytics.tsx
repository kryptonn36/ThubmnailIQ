import { TrendChart } from "@/components/charts/TrendChart";
import { Card } from "@/components/ui/Card";
import { StatCard } from "@/components/ui/StatCard";
import { useAnalytics } from "@/hooks/useAnalytics";
import { extractErrorMessage } from "@/lib/axios";

function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB"];
  const i = Math.min(units.length - 1, Math.floor(Math.log(bytes) / Math.log(1024)));
  return `${(bytes / 1024 ** i).toFixed(1)} ${units[i]}`;
}

export function AnalyticsPage() {
  const { data, isLoading, error } = useAnalytics();

  if (isLoading) return <p className="text-sm text-gray-500">Loading analytics…</p>;
  if (error || !data) {
    return (
      <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">
        {error ? extractErrorMessage(error) : "Failed to load analytics."}
      </p>
    );
  }

  const fileTypeEntries = Object.entries(data.file_type_breakdown);

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-white">Analytics</h1>
        <p className="mt-1 text-sm text-gray-400">Product usage and growth over time.</p>
      </div>

      <div className="grid grid-cols-2 gap-4 sm:grid-cols-3">
        <StatCard label="Storage used" value={formatBytes(data.storage_used_bytes)} />
        <StatCard label="API requests this month" value={data.api_usage.total_requests_this_month} />
        <StatCard label="Active API keys" value={data.api_usage.active_keys} />
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <TrendChart title="Daily signups (30 days)" data={data.daily_user_signups} />
        <TrendChart title="Monthly signups (12 months)" data={data.monthly_user_signups} />
        <TrendChart title="Daily uploads (30 days)" data={data.daily_upload_trend} />
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card>
          <p className="mb-3 text-sm font-semibold text-white">Top active users</p>
          {data.top_active_users.length === 0 ? (
            <p className="py-6 text-center text-sm text-gray-500">No activity yet.</p>
          ) : (
            <div className="space-y-2">
              {data.top_active_users.map((u, i) => (
                <div key={u.id} className="flex items-center justify-between rounded-lg px-2 py-2">
                  <div className="flex items-center gap-3">
                    <span className="w-5 text-xs text-gray-500">#{i + 1}</span>
                    <div>
                      <p className="text-sm font-medium text-gray-200">{u.full_name || u.email}</p>
                      <p className="text-xs text-gray-500">{u.email}</p>
                    </div>
                  </div>
                  <span className="text-sm font-semibold text-white">{u.analyses_count} uploads</span>
                </div>
              ))}
            </div>
          )}
        </Card>

        <Card>
          <p className="mb-3 text-sm font-semibold text-white">Upload file types</p>
          {fileTypeEntries.length === 0 ? (
            <p className="py-6 text-center text-sm text-gray-500">No uploads yet.</p>
          ) : (
            <div className="space-y-2">
              {fileTypeEntries.map(([ext, count]) => (
                <div key={ext} className="flex items-center justify-between text-sm">
                  <span className="uppercase text-gray-300">{ext || "unknown"}</span>
                  <span className="font-semibold text-white">{count}</span>
                </div>
              ))}
            </div>
          )}
        </Card>
      </div>
    </div>
  );
}
