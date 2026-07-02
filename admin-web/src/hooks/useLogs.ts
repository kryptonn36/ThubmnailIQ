import { useQuery } from "@tanstack/react-query";
import { listAuditLogs } from "@/api/logs";

export function useAuditLogs(params: { page?: number; per_page?: number }) {
  return useQuery({
    queryKey: ["logs", "audit", params],
    queryFn: () => listAuditLogs(params),
  });
}
