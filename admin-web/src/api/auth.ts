import { api } from "@/lib/axios";
import type { AuthResponse } from "@/types";

export function login(email: string, password: string) {
  return api.post<AuthResponse>("/admin/auth/login", { email, password }).then((r) => r.data);
}
