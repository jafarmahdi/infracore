import { useState } from "react";
import { Cable, Link2, Search } from "lucide-react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { MetricCard } from "@/components/dashboard/metric-card";

const CABLES = [
  { id: "CBL-001", label: "P-SW-CORE-E1",    fromDev: "SW-CORE-01",   fromPort: "Gi1/0/1",  toDev: "FW-EDGE-01",   toPort: "eth0",       type: "CAT6A",    color: "blue",   lengthM: 3,   status: "connected" },
  { id: "CBL-002", label: "P-SW-CORE-E2",    fromDev: "SW-CORE-01",   fromPort: "Gi1/0/2",  toDev: "SRV-PROD-01",  toPort: "eno1",       type: "CAT6A",    color: "blue",   lengthM: 5,   status: "connected" },
  { id: "CBL-003", label: "P-SW-CORE-E3",    fromDev: "SW-CORE-01",   fromPort: "Gi1/0/3",  toDev: "SRV-PROD-02",  toPort: "eno1",       type: "CAT6A",    color: "blue",   lengthM: 5,   status: "connected" },
  { id: "CBL-004", label: "P-SW-CORE-DB",    fromDev: "SW-CORE-01",   fromPort: "Gi1/0/4",  toDev: "SRV-DB-01",    toPort: "eno1",       type: "CAT6A",    color: "yellow", lengthM: 8,   status: "connected" },
  { id: "CBL-005", label: "F-SRV01-SAN",     fromDev: "SRV-PROD-01",  fromPort: "HBA0",     toDev: "SAN-PROD-01",  toPort: "FC-01",      type: "Fibre 32G",color: "orange", lengthM: 10,  status: "connected" },
  { id: "CBL-006", label: "F-SRV02-SAN",     fromDev: "SRV-PROD-02",  fromPort: "HBA0",     toDev: "SAN-PROD-01",  toPort: "FC-02",      type: "Fibre 32G",color: "orange", lengthM: 10,  status: "connected" },
  { id: "CBL-007", label: "P-UPLINK-WAN",    fromDev: "RTR-WAN-01",   fromPort: "Gi0/0/0",  toDev: "ISP-HANDOFF",  toPort: "Port1",      type: "CAT6",     color: "red",    lengthM: 15,  status: "connected" },
  { id: "CBL-008", label: "P-R01-PATCH-01",  fromDev: "R01-PP-01",    fromPort: "P01",      toDev: "SW-CORE-01",   toPort: "Gi2/0/1",    type: "CAT6A",    color: "green",  lengthM: 2,   status: "connected" },
  { id: "CBL-009", label: "F-DR-STORAGE",    fromDev: "SRV-DR-01",    fromPort: "HBA0",     toDev: "SAN-DR-01",    toPort: "FC-01",      type: "Fibre 16G",color: "orange", lengthM: 6,   status: "disconnected" },
  { id: "CBL-010", label: "P-SPARE-UNUSED",  fromDev: "SW-ACCESS-01", fromPort: "Gi0/12",   toDev: "–",            toPort: "–",          type: "CAT6A",    color: "grey",   lengthM: 1,   status: "planned" },
];

const TYPE_COLORS: Record<string, string> = {
  "CAT6A": "bg-blue-500/10 text-blue-600",
  "CAT6": "bg-blue-400/10 text-blue-500",
  "Fibre 32G": "bg-orange-500/10 text-orange-600",
  "Fibre 16G": "bg-orange-400/10 text-orange-500",
};

const ST_VARIANT: Record<string, "success" | "warning" | "destructive" | "secondary"> = {
  connected: "success", disconnected: "destructive", planned: "secondary",
};

export function CablingPage() {
  const [search, setSearch] = useState("");
  const [typeFilter, setTypeFilter] = useState("all");

  const counts = {
    total: CABLES.length,
    connected: CABLES.filter(c => c.status === "connected").length,
    copper: CABLES.filter(c => c.type.startsWith("CAT")).length,
    fibre: CABLES.filter(c => c.type.startsWith("Fibre")).length,
  };

  const filtered = CABLES.filter(c =>
    (typeFilter === "all" || c.type.startsWith(typeFilter)) &&
    (c.label.toLowerCase().includes(search.toLowerCase()) ||
      c.fromDev.toLowerCase().includes(search.toLowerCase()) ||
      c.toDev.toLowerCase().includes(search.toLowerCase()))
  );

  return (
    <div className="space-y-5">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <MetricCard title="Total cables" value={counts.total} description="documented" icon={Cable} iconClassName="bg-blue-500/10 text-blue-500" />
        <MetricCard title="Connected" value={counts.connected} description="active links" icon={Link2} iconClassName="bg-emerald-500/10 text-emerald-500" />
        <MetricCard title="Copper runs" value={counts.copper} description="CAT6/CAT6A" icon={Cable} iconClassName="bg-violet-500/10 text-violet-500" />
        <MetricCard title="Fibre runs" value={counts.fibre} description="single / multi-mode" icon={Cable} iconClassName="bg-orange-500/10 text-orange-500" />
      </div>

      <Card>
        <CardHeader className="pb-0">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex flex-wrap gap-2">
              <div className="relative">
                <Search className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
                <Input placeholder="Label or device…" className="pl-8 h-8 w-52 text-sm" value={search} onChange={e => setSearch(e.target.value)} />
              </div>
              <select value={typeFilter} onChange={e => setTypeFilter(e.target.value)}
                className="h-8 rounded-md border bg-background px-3 text-sm text-foreground">
                <option value="all">All types</option>
                <option value="CAT">Copper (CAT)</option>
                <option value="Fibre">Fibre</option>
              </select>
            </div>
            <span className="text-xs text-muted-foreground">{filtered.length} cables</span>
          </div>
        </CardHeader>
        <CardContent className="pt-4 overflow-x-auto">
          <table className="w-full text-sm">
            <thead><tr className="border-b text-left">
              {["ID", "Label", "From (device / port)", "To (device / port)", "Type", "Color", "Length", "Status"].map(h => (
                <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground whitespace-nowrap">{h}</th>
              ))}
            </tr></thead>
            <tbody className="divide-y">
              {filtered.map(c => (
                <tr key={c.id} className="hover:bg-muted/30">
                  <td className="py-3 pr-4 font-mono text-xs">{c.id}</td>
                  <td className="py-3 pr-4 font-medium">{c.label}</td>
                  <td className="py-3 pr-4">
                    <span className="font-medium">{c.fromDev}</span>
                    <span className="text-muted-foreground"> · {c.fromPort}</span>
                  </td>
                  <td className="py-3 pr-4">
                    <span className="font-medium">{c.toDev}</span>
                    {c.toPort !== "–" && <span className="text-muted-foreground"> · {c.toPort}</span>}
                  </td>
                  <td className="py-3 pr-4">
                    <span className={`rounded px-1.5 py-0.5 text-[10px] font-medium ${TYPE_COLORS[c.type] ?? "bg-muted text-muted-foreground"}`}>{c.type}</span>
                  </td>
                  <td className="py-3 pr-4 capitalize text-muted-foreground">{c.color}</td>
                  <td className="py-3 pr-4 font-mono text-xs">{c.lengthM}m</td>
                  <td className="py-3"><Badge variant={ST_VARIANT[c.status]}>{c.status}</Badge></td>
                </tr>
              ))}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
