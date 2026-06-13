import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Bell, ChevronDown, Command, Languages, LogOut, Menu, Moon, Search, Sun, UserRound } from "lucide-react";
import { useAuthStore } from "@/stores/auth-store";
import { useUIStore } from "@/stores/ui-store";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export function Topbar() {
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);
  const { user, clearSession } = useAuthStore();
  const { theme, locale, toggleTheme, setLocale, setMobileSidebarOpen } = useUIStore();
  const backendName = [user?.first_name, user?.last_name]
    .filter(Boolean)
    .join(" ");
  const displayName = user?.name || backendName || user?.username || "";
  const initials = displayName
    .split(" ")
    .filter(Boolean)
    .map((part) => part[0])
    .slice(0, 2)
    .join("");

  const logout = () => {
    clearSession();
    navigate("/login");
  };

  return (
    <header className="sticky top-0 z-30 flex h-16 items-center gap-3 border-b bg-background/90 px-4 backdrop-blur-xl md:px-6">
      <Button variant="ghost" size="icon" className="lg:hidden" onClick={() => setMobileSidebarOpen(true)}>
        <Menu className="h-5 w-5" />
      </Button>
      <div className="relative hidden w-full max-w-md md:block">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input className="h-9 bg-muted/60 pl-9 pr-16" placeholder="Search assets, IPs, alerts..." />
        <div className="absolute right-2 top-1/2 flex -translate-y-1/2 items-center gap-1 rounded border bg-background px-1.5 py-0.5 text-[10px] text-muted-foreground">
          <Command className="h-3 w-3" /> K
        </div>
      </div>
      <div className="ml-auto flex items-center gap-1">
        <Button variant="ghost" size="icon" title="Toggle language" onClick={() => setLocale(locale === "en" ? "ar" : "en")}>
          <Languages className="h-[18px] w-[18px]" />
        </Button>
        <Button variant="ghost" size="icon" title="Toggle theme" onClick={toggleTheme}>
          {theme === "light" ? <Moon className="h-[18px] w-[18px]" /> : <Sun className="h-[18px] w-[18px]" />}
        </Button>
        <Button variant="ghost" size="icon" className="relative">
          <Bell className="h-[18px] w-[18px]" />
          <span className="absolute right-2 top-2 h-2 w-2 rounded-full border-2 border-background bg-red-500" />
        </Button>
        <div className="mx-2 h-6 w-px bg-border" />
        <div className="relative">
          <button className="flex items-center gap-2 rounded-lg p-1.5 hover:bg-muted" onClick={() => setMenuOpen(!menuOpen)}>
            <div className="grid h-8 w-8 place-items-center rounded-lg bg-gradient-to-br from-blue-500 to-cyan-500 text-xs font-bold text-white">
              {initials}
            </div>
            <div className="hidden text-left md:block">
              <div className="text-xs font-semibold">{displayName}</div>
              <div className="text-[10px] capitalize text-muted-foreground">{user?.role}</div>
            </div>
            <ChevronDown className="hidden h-3.5 w-3.5 text-muted-foreground md:block" />
          </button>
          {menuOpen && (
            <div className="absolute right-0 top-12 w-48 rounded-xl border bg-card p-1.5 shadow-panel">
              <button className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm hover:bg-muted"><UserRound className="h-4 w-4" /> Profile</button>
              <button onClick={logout} className="flex w-full items-center gap-2 rounded-lg px-3 py-2 text-sm text-red-500 hover:bg-red-500/10"><LogOut className="h-4 w-4" /> Sign out</button>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
