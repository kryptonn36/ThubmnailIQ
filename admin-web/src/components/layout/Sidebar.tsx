import { NavLink } from "react-router-dom";

const NAV_ITEMS = [
  { to: "/", label: "Dashboard", icon: "📊", end: true },
  { to: "/users", label: "Users", icon: "👥" },
  { to: "/uploads", label: "Uploads", icon: "🖼️" },
  { to: "/analytics", label: "Analytics", icon: "📈" },
  { to: "/settings", label: "Settings", icon: "⚙️" },
  { to: "/logs", label: "Logs", icon: "📜" },
  { to: "/profile", label: "Profile", icon: "👤" },
];

function SidebarContent({ onNavigate }: { onNavigate?: () => void }) {
  return (
    <>
      <div className="flex h-14 items-center gap-2 border-b border-surface-300 px-5">
        <span className="flex h-7 w-7 items-center justify-center rounded-lg bg-brand-600 text-xs font-bold text-white">
          IQ
        </span>
        <span className="text-sm font-semibold text-white">Admin Panel</span>
      </div>
      <nav className="flex-1 space-y-1 p-3">
        {NAV_ITEMS.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.end}
            onClick={onNavigate}
            className={({ isActive }) =>
              `flex items-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition ${
                isActive
                  ? "bg-brand-600/20 text-brand-500"
                  : "text-gray-400 hover:bg-surface-200 hover:text-gray-200"
              }`
            }
          >
            <span>{item.icon}</span>
            {item.label}
          </NavLink>
        ))}
      </nav>
    </>
  );
}

export function Sidebar() {
  return (
    <aside className="hidden w-56 shrink-0 flex-col border-r border-surface-300 bg-surface-100 md:flex">
      <SidebarContent />
    </aside>
  );
}

// MobileSidebar is a slide-in overlay used on small screens, toggled by the
// hamburger button in Topbar — the desktop <Sidebar> above stays fixed and
// never shows this overlay (it's hidden below the md breakpoint).
export function MobileSidebar({ open, onClose }: { open: boolean; onClose: () => void }) {
  if (!open) return null;
  return (
    <div className="fixed inset-0 z-40 md:hidden">
      <div className="absolute inset-0 bg-black/60" onClick={onClose} aria-hidden />
      <aside className="relative flex h-full w-64 flex-col border-r border-surface-300 bg-surface-100">
        <SidebarContent onNavigate={onClose} />
      </aside>
    </div>
  );
}
