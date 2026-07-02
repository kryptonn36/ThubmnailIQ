import { createContext, useContext, useMemo, useState, type ReactNode } from "react";
import * as authApi from "@/api/auth";
import { clearTokens, getAccessToken, setTokens } from "@/lib/tokenStorage";
import type { Admin } from "@/types";

const ADMIN_KEY = "admin_current_admin";

function loadStoredAdmin(): Admin | null {
  const raw = localStorage.getItem(ADMIN_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as Admin;
  } catch {
    return null;
  }
}

interface AuthContextValue {
  admin: Admin | null;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [admin, setAdmin] = useState<Admin | null>(loadStoredAdmin);

  const value = useMemo<AuthContextValue>(
    () => ({
      admin,
      isAuthenticated: Boolean(admin && getAccessToken()),
      async login(email, password) {
        const res = await authApi.login(email, password);
        setTokens(res.access_token, res.refresh_token);
        localStorage.setItem(ADMIN_KEY, JSON.stringify(res.admin));
        setAdmin(res.admin);
      },
      logout() {
        clearTokens();
        localStorage.removeItem(ADMIN_KEY);
        setAdmin(null);
      },
    }),
    [admin],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
