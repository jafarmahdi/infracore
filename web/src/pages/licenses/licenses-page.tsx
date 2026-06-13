import { useState } from "react";
import { CalendarClock, Key, Search, Users } from "lucide-react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { MetricCard } from "@/components/dashboard/metric-card";

const LICENSES = [
  { id: "LIC-001", product: "VMware vSphere Enterprise+", vendor: "VMware",     type: "Perpetual",    seats: 48,  used: 44, expires: "2027-01-15", support: "2027-01-15", status: "active" },
  { id: "LIC-002", product: "Microsoft Windows Server",   vendor: "Microsoft",  type: "Subscription", seats: 120, used: 98, expires: "2026-09-30", support: "2026-09-30", status: "expiring" },
  { id: "LIC-003", product: "Red Hat Enterprise Linux",   vendor: "Red Hat",    type: "Subscription", seats: 60,  used: 52, expires: "2027-06-01", support: "2027-06-01", status: "active" },
  { id: "LIC-004", product: "Cisco DNA Center",           vendor: "Cisco",      type: "Subscription", seats: 1,   used: 1,  expires: "2026-07-01", support: "2026-07-01", status: "expiring" },
  { id: "LIC-005", product: "Palo Alto Panorama",         vendor: "Palo Alto",  type: "Perpetual",    seats: 1,   used: 1,  expires: "2028-03-01", support: "2028-03-01", status: "active" },
  { id: "LIC-006", product: "SolarWinds NPM",             vendor: "SolarWinds", type: "Perpetual",    seats: 1,   used: 1,  expires: "2025-12-01", support: "2025-12-01", status: "expired" },
  { id: "LIC-007", product: "Veeam Backup Enterprise",    vendor: "Veeam",      type: "Perpetual",    seats: 10,  used: 7,  expires: "2027-04-15", support: "2027-04-15", status: "active" },
  { id: "LIC-008", product: "Microsoft 365 Business",     vendor: "Microsoft",  type: "Subscription", seats: 200, used: 186,expires: "2026-08-01", support: "2026-08-01", status: "expiring" },
  { id: "LIC-009", product: "Zabbix Enterprise",          vendor: "Zabbix LLC", type: "Subscription", seats: 500, used: 312,expires: "2027-01-01", support: "2027-01-01", status: "active" },
  { id: "LIC-010", product: "GitLab Ultimate",            vendor: "GitLab",     type: "Subscription", seats: 25,  used: 22, expires: "2026-06-20", support: "2026-06-20", status: "expiring" },
];

const ST_VARIANT: Record<string, "success" | "warning" | "destructive"> = {
  active: "success", expiring: "warning", expired: "destructive",
};

function SeatBar({ used, seats }: { used: number; seats: number }) {
  const pct = Math.min((used / seats) * 100, 100);
  const color = pct >= 95 ? "bg-red-500" : pct >= 80 ? "bg-amber-500" : "bg-emerald-500";
  return (
    <div className="flex items-center gap-2">
      <div className="h-1.5 w-16 rounded-full bg-muted overflow-hidden">
        <div className={`h-full rounded-full ${color}`} style={{ width: `${pct}%` }} />
      </div>
      <span className="text-xs tabular-nums">{used}/{seats}</span>
    </div>
  );
}

export function LicensesPage() {
  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState<"all" | "active" | "expiring" | "expired">("all");

  const counts = {
    total: LICENSES.length,
    active: LICENSES.filter(l => l.status === "active").length,
    expiring: LICENSES.filter(l => l.status === "expiring").length,
    expired: LICENSES.filter(l => l.status === "expired").length,
  };
  const totalSeats = LICENSES.reduce((s, l) => s + l.seats, 0);
  const usedSeats = LICENSES.reduce((s, l) => s + l.used, 0);

  const filtered = LICENSES.filter(l =>
    (filter === "all" || l.status === filter) &&
    (l.product.toLowerCase().includes(search.toLowerCase()) || l.vendor.toLowerCase().includes(search.toLowerCase()))
  );

  return (
    <div className="space-y-5">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <MetricCard title="Total licenses" value={counts.total} description="tracked products" icon={Key} iconClassName="bg-blue-500/10 text-blue-500" />
        <MetricCard title="Expiring soon" value={counts.expiring} description="within 90 days" icon={CalendarClock} iconClassName="bg-amber-500/10 text-amber-500" trend={{ value: 2, direction: "up" }} />
        <MetricCard title="Expired" value={counts.expired} description="action required" icon={Key} iconClassName="bg-red-500/10 text-red-500" />
        <MetricCard title="Seat utilization" value={Math.round((usedSeats / totalSeats) * 100)} suffix="%" description={`${usedSeats} of ${totalSeats} seats used`} icon={Users} iconClassName="bg-violet-500/10 text-violet-500" />
      </div>

      <Card>
        <CardHeader className="pb-0">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex flex-wrap gap-2">
              <div className="relative">
                <Search className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
                <Input placeholder="Product or vendor…" className="pl-8 h-8 w-52 text-sm" value={search} onChange={e => setSearch(e.target.value)} />
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
            <span className="text-xs text-muted-foreground">{filtered.length} licenses</span>
          </div>
        </CardHeader>
        <CardContent className="pt-4 overflow-x-auto">
          <table className="w-full text-sm">
            <thead><tr className="border-b text-left">
              {["ID", "Product", "Vendor", "Type", "Seat usage", "Expires", "Status"].map(h => (
                <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground whitespace-nowrap">{h}</th>
              ))}
            </tr></thead>
            <tbody className="divide-y">
              {filtered.map(l => (
                <tr key={l.id} className="hover:bg-muted/30">
                  <td className="py-3 pr-4 font-mono text-xs">{l.id}</td>
                  <td className="py-3 pr-4 font-medium">{l.product}</td>
                  <td className="py-3 pr-4 text-muted-foreground">{l.vendor}</td>
                  <td className="py-3 pr-4"><Badge variant="secondary">{l.type}</Badge></td>
                  <td className="py-3 pr-4"><SeatBar used={l.used} seats={l.seats} /></td>
                  <td className="py-3 pr-4 font-mono text-xs">{l.expires}</td>
                  <td className="py-3"><Badge variant={ST_VARIANT[l.status]}>{l.status}</Badge></td>
                </tr>
              ))}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
