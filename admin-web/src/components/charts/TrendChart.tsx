import { Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { Card } from "@/components/ui/Card";
import type { TrendPoint } from "@/types";

export function TrendChart({ title, data }: { title: string; data: TrendPoint[] }) {
  return (
    <Card>
      <p className="mb-4 text-sm font-semibold text-white">{title}</p>
      {data.length === 0 ? (
        <p className="py-10 text-center text-sm text-gray-500">No data for this period.</p>
      ) : (
        <div className="h-56">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={data}>
              <XAxis dataKey="date" stroke="#6b7280" fontSize={11} tickLine={false} />
              <YAxis stroke="#6b7280" fontSize={11} tickLine={false} allowDecimals={false} />
              <Tooltip
                contentStyle={{ background: "#1e1e29", border: "1px solid #2a2a3a", borderRadius: 8 }}
                labelStyle={{ color: "#e5e5ec" }}
              />
              <Line type="monotone" dataKey="count" stroke="#6366f1" strokeWidth={2} dot={false} />
            </LineChart>
          </ResponsiveContainer>
        </div>
      )}
    </Card>
  );
}
