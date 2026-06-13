import { useState } from "react";
import { Boxes, HardDrive, Search, Server, Trash2 } from "lucide-react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { MetricCard } from "@/components/dashboard/metric-card";

type Category = "all" | "server" | "switch" | "router" | "firewall" | "storage" | "vm";
type Status = "all" | "active" | "decommissioned" | "maintenance" | "spare";

const ASSETS = [
  { id: "AST-0001", name: "SRV-PROD-01", category: "server", vendor: "Dell", model: "PowerEdge R750", serial: "SN-84720A", site: "DC1 – Riyadh Core", rack: "R01 / U12", status: "active", warranty: "2027-06-01" },
  { id: "AST-0002", name: "SRV-PROD-02", category: "server", vendor: "HP", model: "ProLiant DL380 G10", serial: "SN-39271B", site: "DC1 – Riyadh Core", rack: "R01 / U14", status: "active", warranty: "2026-08-15" },
  { id: "AST-0003", name: "SW-CORE-01", category: "switch", vendor: "Cisco", model: "Catalyst 9500", serial: "SN-C950011", site: "DC1 – Riyadh Core", rack: "R01 / U01", status: "active", warranty: "2028-01-01" },
  { id: "AST-0004", name: "FW-EDGE-01", category: "firewall", vendor: "Palo Alto", model: "PA-5250", serial: "SN-PA52501", site: "DC1 – Riyadh Core", rack: "R01 / U03", status: "active", warranty: "2027-03-20" },
  { id: "AST-0005", name: "SRV-DB-01", category: "server", vendor: "Dell", model: "PowerEdge R840", serial: "SN-R84001", site: "DC1 – Riyadh Core", rack: "R02 / U08", status: "active", warranty: "2025-12-01" },
  { id: "AST-0006", name: "SAN-PROD-01", category: "storage", vendor: "NetApp", model: "AFF A400", serial: "SN-AFF4001", site: "DC1 – Riyadh Core", rack: "R03 / U20", status: "active", warranty: "2028-06-01" },
  { id: "AST-0007", name: "SW-ACCESS-01", category: "switch", vendor: "Juniper", model: "EX4300", serial: "SN-JEX4301", site: "DC2 – Jeddah Edge", rack: "R04 / U01", status: "active", warranty: "2027-09-01" },
  { id: "AST-0008", name: "SRV-EDGE-01", category: "server", vendor: "Lenovo", model: "ThinkSystem SR650", serial: "SN-TSR6501", site: "DC2 – Jeddah Edge", rack: "R04 / U10", status: "maintenance", warranty: "2026-11-15" },
  { id: "AST-0009", name: "RTR-WAN-01", category: "router", vendor: "Cisco", model: "ASR 1002-X", serial: "SN-ASR10021", site: "DC2 – Jeddah Edge", rack: "R04 / U05", status: "active", warranty: "2026-07-01" },
  { id: "AST-0010", name: "SRV-OLD-01", category: "server", vendor: "Dell", model: "PowerEdge R620", serial: "SN-R62001", site: "DC3 – Dammam West", rack: "R05 / U30", status: "decommissioned", warranty: "2022-01-01" },
  { id: "AST-0011", name: "VM-WEB-CLUSTER", category: "vm", vendor: "VMware", model: "vSphere 8.0", serial: "–", site: "DC1 – Riyadh Core", rack: "Virtual", status: "active", warranty: "2027-01-01" },
  { id: "AST-0012", name: "SRV-SPARE-01", category: "server", vendor: "HP", model: "ProLiant DL360 G9", serial: "SN-DL3601", site: "DR – Backup Site", rack: "R07 / U02", status: "spare", warranty: "2025-09-01" },
];

const statusVariant: Record<string, "success" | "warning" | "destructive" | "secondary"> = {
  active: "success", maintenance: "warning", decommissioned: "destructive", spare: "secondary",
};

const CATEGORIES: { value: Category; label: string }[] = [
  { value: "all", label: "All categories" },
  { value: "server", label: "Servers" },
  { value: "switch", label: "Switches" },
  { value: "router", label: "Routers" },
  { value: "firewall", label: "Firewalls" },
  { value: "storage", label: "Storage" },
  { value: "vm", label: "Virtual machines" },
];

function warrantyBadge(date: string) {
  if (date === "–") return <span className="text-muted-foreground">–</span>;
  const daysLeft = Math.floor((new Date(date).getTime() - Date.now()) / 86400000);
  if (daysLeft < 0) return <Badge variant="destructive">Expired</Badge>;
  if (daysLeft < 180) return <Badge variant="warning">{date}</Badge>;
  return <span className="text-xs text-muted-foreground">{date}</span>;
}

export function AssetsPage() {
  const [search, setSearch] = useState("");
  const [category, setCategory] = useState<Category>("all");
  const [status, setStatus] = useState<Status>("all");

  const filtered = ASSETS.filter(a =>
    (category === "all" || a.category === category) &&
    (status === "all" || a.status === status) &&
    (a.name.toLowerCase().includes(search.toLowerCase()) ||
      a.serial.toLowerCase().includes(search.toLowerCase()) ||
      a.site.toLowerCase().includes(search.toLowerCase()))
  );

  const counts = {
    total: ASSETS.length,
    active: ASSETS.filter(a => a.status === "active").length,
    maintenance: ASSETS.filter(a => a.status === "maintenance").length,
    decommissioned: ASSETS.filter(a => a.status === "decommissioned").length,
  };

  return (
    <div className="space-y-5">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <MetricCard title="Total assets" value={counts.total} description="managed inventory" icon={Boxes} iconClassName="bg-blue-500/10 text-blue-500" trend={{ value: 3.8, direction: "up" }} />
        <MetricCard title="Active" value={counts.active} description="in production" icon={Server} iconClassName="bg-emerald-500/10 text-emerald-500" />
        <MetricCard title="In maintenance" value={counts.maintenance} description="scheduled work" icon={HardDrive} iconClassName="bg-amber-500/10 text-amber-500" />
        <MetricCard title="Decommissioned" value={counts.decommissioned} description="pending disposal" icon={Trash2} iconClassName="bg-red-500/10 text-red-500" />
      </div>

      <Card>
        <CardHeader className="pb-0">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex flex-wrap gap-2">
              <div className="relative">
                <Search className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
                <Input placeholder="Name, serial, site…" className="pl-8 h-8 w-56 text-sm" value={search} onChange={e => setSearch(e.target.value)} />
              </div>
              <select value={category} onChange={e => setCategory(e.target.value as Category)}
                className="h-8 rounded-md border bg-background px-3 text-sm text-foreground">
                {CATEGORIES.map(c => <option key={c.value} value={c.value}>{c.label}</option>)}
              </select>
              <select value={status} onChange={e => setStatus(e.target.value as Status)}
                className="h-8 rounded-md border bg-background px-3 text-sm text-foreground">
                {["all", "active", "maintenance", "decommissioned", "spare"].map(s => (
                  <option key={s} value={s}>{s === "all" ? "All statuses" : s.charAt(0).toUpperCase() + s.slice(1)}</option>
                ))}
              </select>
            </div>
            <span className="text-xs text-muted-foreground">{filtered.length} results</span>
          </div>
        </CardHeader>
        <CardContent className="pt-4 overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b text-left">
                {["Asset ID", "Name", "Category", "Vendor / Model", "Serial", "Location", "Warranty", "Status"].map(h => (
                  <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground whitespace-nowrap last:pr-0">{h}</th>
                ))}
              </tr>
            </thead>
            <tbody className="divide-y">
              {filtered.map(a => (
                <tr key={a.id} className="hover:bg-muted/30 transition-colors">
                  <td className="py-3 pr-4 font-mono text-xs">{a.id}</td>
                  <td className="py-3 pr-4 font-medium whitespace-nowrap">{a.name}</td>
                  <td className="py-3 pr-4">
                    <Badge variant="secondary">{a.category}</Badge>
                  </td>
                  <td className="py-3 pr-4 text-muted-foreground whitespace-nowrap">{a.vendor} {a.model}</td>
                  <td className="py-3 pr-4 font-mono text-xs">{a.serial}</td>
                  <td className="py-3 pr-4 text-muted-foreground whitespace-nowrap">{a.site} · {a.rack}</td>
                  <td className="py-3 pr-4">{warrantyBadge(a.warranty)}</td>
                  <td className="py-3"><Badge variant={statusVariant[a.status]}>{a.status}</Badge></td>
                </tr>
              ))}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
