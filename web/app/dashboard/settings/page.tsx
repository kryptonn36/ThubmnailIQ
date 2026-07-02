"use client";

import { LogOut } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { PageHeader } from "@/components/ui/PageHeader";

export default function SettingsPage() {
  const { user, logout } = useAuth();

  return (
    <div className="max-w-xl space-y-6">
      <PageHeader title="Settings" description="Manage your account." />

      <Card>
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
      </Card>

      <Button variant="danger" onClick={logout} icon={<LogOut className="h-4 w-4" aria-hidden="true" />}>
        Log Out
      </Button>
    </div>
  );
}
