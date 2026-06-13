import { NavLink } from "react-router-dom";
import { PanelLeftClose, PanelLeftOpen, X } from "lucide-react";
import { navigation } from "@/config/navigation";
import { useAuthStore } from "@/stores/auth-store";
import { useUIStore } from "@/stores/ui-store";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Logo } from "@/components/common/logo";

export function Sidebar() {
  const user = useAuthStore((state) => state.user);
  const { sidebarCollapsed, mobileSidebarOpen, toggleSidebar, setMobileSidebarOpen } = useUIStore();
  const items = navigation.filter(
    (item) => !item.roles || (user?.role && item.roles.includes(user.role)),
  );

  const content = (
    <>
      <div className="flex h-16 items-center justify-between border-b px-4">
        <Logo compact={sidebarCollapsed} />
        <Button variant="ghost" size="icon" className="lg:hidden" onClick={() => setMobileSidebarOpen(false)}>
          <X className="h-5 w-5" />
        </Button>
      </div>
      <nav className="flex-1 space-y-1 overflow-y-auto px-3 py-4">
        {!sidebarCollapsed && <div className="mb-2 px-3 text-[10px] font-semibold uppercase tracking-[0.18em] text-muted-foreground">Workspace</div>}
        {items.map(({ label, path, icon: Icon, badge }) => (
          <NavLink
            key={path}
            to={path}
            onClick={() => setMobileSidebarOpen(false)}
            className={({ isActive }) =>
              cn(
                "group flex h-10 items-center gap-3 rounded-lg px-3 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted hover:text-foreground",
                isActive && "bg-blue-600/10 text-blue-600 dark:text-blue-400",
                sidebarCollapsed && "justify-center px-0",
              )
            }
            title={sidebarCollapsed ? label : undefined}
          >
            <Icon className="h-[18px] w-[18px] shrink-0" />
            {!sidebarCollapsed && <span className="flex-1">{label}</span>}
            {!sidebarCollapsed && badge && <span className="rounded-full bg-red-500 px-2 py-0.5 text-[10px] text-white">{badge}</span>}
          </NavLink>
        ))}
      </nav>
      <div className="border-t p-3">
        <button
          onClick={toggleSidebar}
          className="hidden h-9 w-full items-center justify-center gap-2 rounded-lg text-xs font-medium text-muted-foreground hover:bg-muted lg:flex"
        >
          {sidebarCollapsed ? <PanelLeftOpen className="h-4 w-4" /> : <><PanelLeftClose className="h-4 w-4" /> Collapse sidebar</>}
        </button>
      </div>
    </>
  );

  return (
    <>
      {mobileSidebarOpen && <div className="fixed inset-0 z-40 bg-slate-950/60 backdrop-blur-sm lg:hidden" onClick={() => setMobileSidebarOpen(false)} />}
      <aside className={cn(
        "fixed inset-y-0 left-0 z-50 flex w-[260px] flex-col border-r bg-card transition-transform duration-200 lg:translate-x-0",
        mobileSidebarOpen ? "translate-x-0" : "-translate-x-full",
        sidebarCollapsed && "lg:w-[72px]",
      )}>
        {content}
      </aside>
    </>
  );
}
