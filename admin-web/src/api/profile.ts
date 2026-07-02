import { api } from "@/lib/axios";
import type { Admin } from "@/types";

export function getProfile() {
  return api.get<Admin>("/admin/profile").then((r) => r.data);
}

export function changePassword(currentPassword: string, newPassword: string) {
  return api
    .patch("/admin/profile/password", { current_password: currentPassword, new_password: newPassword })
    .then((r) => r.data);
}
