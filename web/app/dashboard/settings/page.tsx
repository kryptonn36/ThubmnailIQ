"use client";

import { useAuth } from "@/hooks/useAuth";

export default function SettingsPage() {
  const { user, logout } = useAuth();

  return (
    <div className="max-w-xl space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">Settings</h1>
        <p className="mt-1 text-sm text-gray-400">Manage your account.</p>
      </div>

      <div className="rounded-2xl border border-surface-300 bg-surface-100 p-6">
        <h2 className="mb-4 text-base font-semibold text-white">Profile</h2>
        <div className="flex items-center gap-4">
          <div className="flex h-14 w-14 items-center justify-center rounded-full bg-brand-gradient text-lg font-bold text-white">
            {user?.full_name?.charAt(0)?.toUpperCase() ?? "?"}
          </div>
          <div>
            <p className="font-medium text-white">{user?.full_name ?? "Unknown user"}</p>
            <p className="text-sm text-gray-400">{user?.email ?? "—"}</p>
          </div>
        </div>
      </div>

      <button
        onClick={logout}
        className="rounded-lg border border-red-500/30 bg-red-500/10 px-4 py-2.5 text-sm font-semibold text-red-400 transition hover:bg-red-500/20"
      >
        Log Out
      </button>
    </div>
  );
}
