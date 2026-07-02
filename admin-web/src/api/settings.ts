import { api } from "@/lib/axios";
import type { Settings } from "@/types";

export function getSettings() {
  return api.get<Settings>("/admin/settings").then((r) => r.data);
}

export function updateSettings(settings: Settings) {
  return api.patch<Settings>("/admin/settings", settings).then((r) => r.data);
}
