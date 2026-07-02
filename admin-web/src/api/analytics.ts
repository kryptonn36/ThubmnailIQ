import { api } from "@/lib/axios";
import type { Analytics } from "@/types";

export function getAnalytics() {
  return api.get<Analytics>("/admin/analytics").then((r) => r.data);
}
