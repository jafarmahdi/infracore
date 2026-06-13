import { useState } from "react";
import { CalendarClock, FileCheck2, Search, Shield } from "lucide-react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { MetricCard } from "@/components/dashboard/metric-card";

const CONTRACTS = [
  { id: "CTR-001", number: "SA-2024-0042", vendor: "Dell Technologies",  type: "Maintenance", coverage: "Hardware support", assets: 24, start: "2024-01-01", end: "2027-01-01", value: "$48,000",  status: "active" },
  { id: "CTR-002", number: "SA-2024-0057", vendor: "Cisco Systems",      type: "SmartNet",    coverage: "Network devices",  assets: 18, start: "2024-03-01", end: "2026-09-01", value: "$32,500",  status: "expiring" },
  { id: "CTR-003", number: "SA-2024-0081", vendor: "Palo Alto Networks", type: "Support",     coverage: "Security appliances",assets: 4, start: "2024-06-01", end: "2027-06-01", value: "$18,200",  status: "active" },
  { id: "CTR-004", number: "SA-2023-0014", vendor: "HP Enterprise",      type: "Maintenance", coverage: "Server hardware",  assets: 12, start: "2023-01-01", end: "2026-07-01", value: "$22,800",  status: "expiring" },
  { id: "CTR-005", number: "SA-2025-0003", vendor: "NetApp",             type: "Premium",     coverage: "Storage systems",  assets: 3,  start: "2025-01-01", end: "2028-01-01", value: "$56,000",  status: "active" },
  { id: "CTR-006", number: "SA-2023-0099", vendor: "SolarWinds",         type: "SaaS",        coverage: "NPM platform",     assets: 1,  start: "2023-01-01", end: "2025-12-01", value: "$9,600",   status: "expired" },
  { id: "CTR-007", number: "SA-2024-0112", vendor: "Veeam Software",     type: "Production",  coverage: "Backup solution",  assets: 10, start: "2024-04-01", end: "2027-04-01", value: "$14,400",  status: "active" },
  { id: "CTR-008", number: "SA-2024-0088", vendor: "Red Hat",            type: "Premium",     coverage: "Linux servers",    assets: 60, start: "2024-06-01", end: "2027-06-01", value: "$28,800",  status: "active" },
];

const ST_VARIANT: Record<string, "success" | "warning" | "destructive"> = {
  active: "success", expiring: "warning", expired: "destructive",
};

function daysLeft(dateStr: string): number {
  return Math.floor((new Date(dateStr).getTime() - Date.now()) / 86400000);
}

export function ContractsPage() {
  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState<"all" | "active" | "expiring" | "expired">("all");

  const counts = {
    active: CONTRACTS.filter(c => c.status === "active").length,
    expiring: CONTRACTS.filter(c => c.status === "expiring").length,
    expired: CONTRACTS.filter(c => c.status === "expired").length,
  };
  const filtered = CONTRACTS.filter(c =>
    (filter === "all" || c.status === filter) &&
    (c.vendor.toLowerCase().includes(search.toLowerCase()) ||
      c.number.toLowerCase().includes(search.toLowerCase()) ||
      c.coverage.toLowerCase().includes(search.toLowerCase()))
  );

  return (
    <div className="space-y-5">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <MetricCard title="Active contracts" value={counts.active} description="in coverage" icon={FileCheck2} iconClassName="bg-emerald-500/10 text-emerald-500" />
        <MetricCard title="Expiring soon" value={counts.expiring} description="within 90 days" icon={CalendarClock} iconClassName="bg-amber-500/10 text-amber-500" trend={{ value: 1, direction: "up" }} />
        <MetricCard title="Expired" value={counts.expired} description="out of support" icon={Shield} iconClassName="bg-red-500/10 text-red-500" />
        <MetricCard title="Total value" value={230300} description="annual contract spend" icon={FileCheck2} iconClassName="bg-blue-500/10 text-blue-500" />
      </div>

      <Card>
        <CardHeader className="pb-0">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex flex-wrap gap-2">
              <div className="relative">
                <Search className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
                <Input placeholder="Vendor or contract no." className="pl-8 h-8 w-56 text-sm" value={search} onChange={e => setSearch(e.target.value)} />
              </div>
              <div className="flex gap-1 rounded-lg border bg-muted/40 p-1">
                {(["all", "active", "expiring", "expired"] as const).map(f => (
                  <button key={f} onClick={() => setFilter(f)}
                    className={`rounded-md px-3 py-1 text-xs font-medium capitalize transition-colors ${filter === f ? "bg-background shadow-sm" : "text-muted-foreground hover:text-foreground"}`}>
                    {f.charAt(0).toUpperCase() + f.slice(1)}
                  </button>
                ))}
              </div>
            </div>
            <span className="text-xs text-muted-foreground">{filtered.length} contracts</span>
          </div>
        </CardHeader>
        <CardContent className="pt-4 overflow-x-auto">
          <table className="w-full text-sm">
            <thead><tr className="border-b text-left">
              {["Contract #", "Vendor", "Type", "Coverage", "Assets", "Expires", "Days left", "Value", "Status"].map(h => (
                <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground whitespace-nowrap">{h}</th>
              ))}
            </tr></thead>
            <tbody className="divide-y">
              {filtered.map(c => {
                const days = daysLeft(c.end);
                return (
                  <tr key={c.id} className="hover:bg-muted/30">
                    <td className="py-3 pr-4 font-mono text-xs">{c.number}</td>
                    <td className="py-3 pr-4 font-medium">{c.vendor}</td>
                    <td className="py-3 pr-4"><Badge variant="secondary">{c.type}</Badge></td>
                    <td className="py-3 pr-4 text-muted-foreground">{c.coverage}</td>
                    <td className="py-3 pr-4 tabular-nums">{c.assets}</td>
                    <td className="py-3 pr-4 font-mono text-xs">{c.end}</td>
                    <td className="py-3 pr-4">
                      <span className={`text-xs tabular-nums font-medium ${days < 0 ? "text-red-500" : days < 90 ? "text-amber-500" : "text-muted-foreground"}`}>
                        {days < 0 ? "Expired" : `${days}d`}
                      </span>
                    </td>
                    <td className="py-3 pr-4 font-medium">{c.value}</td>
                    <td className="py-3"><Badge variant={ST_VARIANT[c.status]}>{c.status}</Badge></td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
