import { Link } from "react-router-dom";
import { Card } from "@/components/ui/Card";
import { useProfile } from "@/hooks/useProfile";
import { extractErrorMessage } from "@/lib/axios";

export function ProfilePage() {
  const { data, isLoading, error } = useProfile();

  if (isLoading) return <p className="text-sm text-gray-500">Loading profile…</p>;
  if (error || !data) {
    return (
      <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">
        {error ? extractErrorMessage(error) : "Failed to load profile."}
      </p>
    );
  }

  return (
    <div className="max-w-xl space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">Profile</h1>
        <p className="mt-1 text-sm text-gray-400">Your admin account details.</p>
      </div>

      <Card>
        <div className="flex items-center gap-4">
          <div className="flex h-14 w-14 items-center justify-center rounded-full bg-brand-600 text-lg font-bold text-white">
            {(data.full_name || data.email).charAt(0).toUpperCase()}
          </div>
          <div>
            <p className="font-medium text-white">{data.full_name}</p>
            <p className="text-sm text-gray-400">{data.email}</p>
          </div>
        </div>

        <dl className="mt-6 grid grid-cols-2 gap-4 text-sm">
          <div>
            <dt className="text-gray-500">Role</dt>
            <dd className="mt-1 font-medium capitalize text-white">{data.role}</dd>
          </div>
          <div>
            <dt className="text-gray-500">Last login</dt>
            <dd className="mt-1 font-medium text-white">
              {data.last_login_at ? new Date(data.last_login_at).toLocaleString() : "—"}
            </dd>
          </div>
          <div>
            <dt className="text-gray-500">Account created</dt>
            <dd className="mt-1 font-medium text-white">{new Date(data.created_at).toLocaleDateString()}</dd>
          </div>
        </dl>

        <Link
          to="/profile/change-password"
          className="mt-6 inline-block rounded-lg border border-surface-300 px-4 py-2 text-sm font-medium text-gray-200 hover:bg-surface-200"
        >
          Change password
        </Link>
      </Card>
    </div>
  );
}
