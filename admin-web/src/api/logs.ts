import { api } from "@/lib/axios";
import type { AuditLog, PaginatedResponse } from "@/types";

export function listAuditLogs(params: { page?: number; per_page?: number }) {
  return api.get<PaginatedResponse<AuditLog>>("/admin/logs/audit", { params }).then((r) => r.data);
}
