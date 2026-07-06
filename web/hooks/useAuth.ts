"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { clearAuth, getStoredUser, isAuthenticated, storeAuth } from "@/lib/auth";
import type { AuthResponse, RegisterResponse, User } from "@/types";

interface UseAuthResult {
  user: User | null;
  loading: boolean;
  authenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, fullName: string) => Promise<{ email: string }>;
  verifyEmail: (email: string, code: string) => Promise<void>;
  resendVerification: (email: string) => Promise<void>;
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

  // Registration no longer logs in — it returns the email to verify. The caller
  // routes the user to the verification screen.
  const register = useCallback(
    async (email: string, password: string, fullName: string) => {
      const res = await api.post<RegisterResponse>(
        "/auth/register",
        { email, password, full_name: fullName },
        false
      );
      return { email: res.email };
    },
    []
  );

  // Submitting the correct OTP verifies the email AND logs the user in.
  const verifyEmail = useCallback(
    async (email: string, code: string) => {
      const res = await api.post<AuthResponse>(
        "/auth/verify-email",
        { email, code },
        false
      );
      storeAuth(res);
      setUser(res.user);
      setAuthenticated(true);
      router.push("/dashboard");
    },
    [router]
  );

  const resendVerification = useCallback(async (email: string) => {
    await api.post("/auth/resend-verification", { email }, false);
  }, []);

  const logout = useCallback(() => {
    clearAuth();
    setUser(null);
    setAuthenticated(false);
    router.push("/login");
  }, [router]);

  return { user, loading, authenticated, login, register, verifyEmail, resendVerification, logout };
}
