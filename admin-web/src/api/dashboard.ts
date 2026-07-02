import { api } from "@/lib/axios";
import type { DashboardStats } from "@/types";

export function getDashboard() {
  return api.get<DashboardStats>("/admin/dashboard").then((r) => r.data);
}
