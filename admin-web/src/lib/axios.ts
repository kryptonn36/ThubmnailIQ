import axios, { AxiosError, type InternalAxiosRequestConfig } from "axios";
import { clearTokens, getAccessToken, getRefreshToken, setTokens } from "./tokenStorage";
import type { ApiErrorBody, AuthResponse } from "@/types";

export const API_BASE_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080/api/v1";

export const api = axios.create({ baseURL: API_BASE_URL });

api.interceptors.request.use((config) => {
  const token = getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Single-flight refresh: if several requests 401 at once, only the first
// triggers a refresh call; the rest wait on the same in-flight promise
// instead of each racing to refresh (which would revoke each other's tokens
// — the backend rotates the refresh token on every use).
let refreshPromise: Promise<string> | null = null;

async function refreshAccessToken(): Promise<string> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) {
    throw new Error("no refresh token");
  }
  const res = await axios.post<AuthResponse>(`${API_BASE_URL}/admin/auth/refresh`, {
    refresh_token: refreshToken,
  });
  setTokens(res.data.access_token, res.data.refresh_token);
  return res.data.access_token;
}

interface RetriableConfig extends InternalAxiosRequestConfig {
  _retried?: boolean;
  _transientRetryCount?: number;
}

// Transient-failure retry: network errors and 5xx are usually momentary
// (a cold-starting service, a dropped connection) rather than a real
// client error, so a couple of short, backed-off retries smooth those over
// instead of surfacing a scary error on the first blip. Limited to GET —
// retrying a POST/PATCH/DELETE blind risks double-applying a mutation the
// server actually received but whose response was lost in transit.
const TRANSIENT_RETRY_DELAYS_MS = [300, 900];

function isTransientError(error: AxiosError): boolean {
  if (!error.response) return true; // network error / timeout, no response at all
  return error.response.status >= 500;
}

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiErrorBody>) => {
    const original = error.config as RetriableConfig | undefined;
    const isAuthEndpoint = original?.url?.includes("/admin/auth/");

    if (error.response?.status === 401 && original && !original._retried && !isAuthEndpoint) {
      original._retried = true;
      try {
        refreshPromise ??= refreshAccessToken();
        const newAccessToken = await refreshPromise;
        refreshPromise = null;
        original.headers.Authorization = `Bearer ${newAccessToken}`;
        return api(original);
      } catch {
        refreshPromise = null;
        clearTokens();
        if (typeof window !== "undefined") {
          window.location.href = "/login";
        }
      }
      return Promise.reject(error);
    }

    if (
      original &&
      original.method?.toLowerCase() === "get" &&
      isTransientError(error) &&
      (original._transientRetryCount ?? 0) < TRANSIENT_RETRY_DELAYS_MS.length
    ) {
      const attempt = original._transientRetryCount ?? 0;
      original._transientRetryCount = attempt + 1;
      await new Promise((resolve) => setTimeout(resolve, TRANSIENT_RETRY_DELAYS_MS[attempt]));
      return api(original);
    }

    return Promise.reject(error);
  },
);

export function extractErrorMessage(error: unknown): string {
  if (axios.isAxiosError<ApiErrorBody>(error)) {
    return error.response?.data?.error ?? error.response?.data?.message ?? error.message;
  }
  if (error instanceof Error) return error.message;
  return "Something went wrong.";
}
