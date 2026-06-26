"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { clearAuth, getStoredUser, isAuthenticated, storeAuth } from "@/lib/auth";
import type { AuthResponse, User } from "@/types";

interface UseAuthResult {
  user: User | null;
  loading: boolean;
  authenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, fullName: string) => Promise<void>;
  logout: () => void;
}

export function useAuth(): UseAuthResult {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [authenticated, setAuthenticated] = useState(false);

  useEffect(() => {
    setUser(getStoredUser());
    setAuthenticated(isAuthenticated());
    setLoading(false);
  }, []);

  const login = useCallback(
    async (email: string, password: string) => {
      const res = await api.post<AuthResponse>(
        "/auth/login",
        { email, password },
        false
      );
      storeAuth(res);
      setUser(res.user);
      setAuthenticated(true);
      router.push("/dashboard");
    },
    [router]
  );

  const register = useCallback(
    async (email: string, password: string, fullName: string) => {
      const res = await api.post<AuthResponse>(
        "/auth/register",
        { email, password, full_name: fullName },
        false
      );
      storeAuth(res);
      setUser(res.user);
      setAuthenticated(true);
      router.push("/dashboard");
    },
    [router]
  );

  const logout = useCallback(() => {
    clearAuth();
    setUser(null);
    setAuthenticated(false);
    router.push("/login");
  }, [router]);

  return { user, loading, authenticated, login, register, logout };
}
