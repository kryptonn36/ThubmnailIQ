import { useNavigate } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";

export function Topbar({ onMenuClick }: { onMenuClick: () => void }) {
  const { admin, logout } = useAuth();
  const navigate = useNavigate();

  function handleLogout() {
    logout();
    navigate("/login", { replace: true });
  }

  return (
    <header className="flex h-14 items-center justify-between border-b border-surface-300 bg-surface-100 px-4 md:px-6">
      <button
        type="button"
        onClick={onMenuClick}
        className="rounded-lg p-2 text-gray-300 hover:bg-surface-200 md:hidden"
        aria-label="Toggle menu"
      >
        ☰
      </button>
      <div className="hidden md:block" />
      <div className="flex items-center gap-3">
        <span className="text-sm text-gray-400">{admin?.email}</span>
        <button
          type="button"
          onClick={handleLogout}
          className="rounded-lg border border-surface-300 px-3 py-1.5 text-xs font-medium text-gray-300 hover:bg-surface-200"
        >
          Log out
        </button>
      </div>
    </header>
  );
}
