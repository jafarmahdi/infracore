import { Outlet, useLocation } from "react-router-dom";
import { ChevronRight } from "lucide-react";
import { Sidebar } from "./sidebar";
import { Topbar } from "./topbar";
import { useUIStore } from "@/stores/ui-store";
import { cn } from "@/lib/utils";

const titles: Record<string, string> = {
  dashboard: "Operations overview", dcim: "Data center infrastructure", assets: "Asset management",
  ipam: "IP address management", monitoring: "Monitoring", alerts: "Alerts & escalation",
  agents: "Agent management", licenses: "Software licenses", contracts: "Contracts",
  reports: "Reports & analytics", cabling: "Cabling", security: "Security", settings: "Settings",
};

export function AppLayout() {
  const { pathname } = useLocation();
  const { sidebarCollapsed, theme, locale } = useUIStore();
  const section = pathname.split("/")[1] || "dashboard";

  return (
    <div className={cn(theme === "dark" && "dark")} dir={locale === "ar" ? "rtl" : "ltr"}>
      <div className="min-h-screen bg-background text-foreground">
        <Sidebar />
        <div className={cn("transition-[margin] duration-200 lg:ml-[260px]", sidebarCollapsed && "lg:ml-[72px]")}>
          <Topbar />
          <main className="px-4 py-5 md:px-6 lg:px-8">
            <div className="mb-5 flex items-center gap-1.5 text-xs text-muted-foreground">
              <span>InfraCore</span><ChevronRight className="h-3 w-3" /><span className="text-foreground">{titles[section]}</span>
            </div>
            <Outlet />
          </main>
        </div>
      </div>
    </div>
  );
}
