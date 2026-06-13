import { useState } from "react";
import { Building2, PlugZap, Search, Server, Thermometer } from "lucide-react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { MetricCard } from "@/components/dashboard/metric-card";

const DATA_CENTERS = [
  { id: "dc-1", name: "DC1 – Riyadh Core", location: "Riyadh, SA", tier: "Tier IV", racks: 42, racksFilled: 38, powerKW: 800, powerUsed: 612, tempC: 21, status: "online" },
  { id: "dc-2", name: "DC2 – Jeddah Edge", location: "Jeddah, SA", tier: "Tier III", racks: 18, racksFilled: 14, powerKW: 320, powerUsed: 230, tempC: 22, status: "online" },
  { id: "dc-3", name: "DC3 – Dammam West", location: "Dammam, SA", tier: "Tier III", racks: 24, racksFilled: 22, powerKW: 480, powerUsed: 451, tempC: 23, status: "warning" },
  { id: "dc-4", name: "DR – Backup Site", location: "Medina, SA", tier: "Tier II", racks: 10, racksFilled: 4, powerKW: 200, powerUsed: 88, tempC: 20, status: "online" },
];

const RACKS = [
  { id: "R01", dc: "DC1 – Riyadh Core", row: "A", units: 42, used: 38, powerA: 16, powerB: 14, status: "active" },
  { id: "R02", dc: "DC1 – Riyadh Core", row: "A", units: 42, used: 40, powerA: 18, powerB: 16, status: "active" },
  { id: "R03", dc: "DC1 – Riyadh Core", row: "B", units: 42, used: 21, powerA: 10, powerB: 9, status: "active" },
  { id: "R04", dc: "DC2 – Jeddah Edge", row: "A", units: 42, used: 36, powerA: 14, powerB: 14, status: "active" },
  { id: "R05", dc: "DC3 – Dammam West", row: "A", units: 42, used: 42, powerA: 20, powerB: 19, status: "critical" },
  { id: "R06", dc: "DC3 – Dammam West", row: "A", units: 42, used: 39, powerA: 17, powerB: 17, status: "warning" },
  { id: "R07", dc: "DR – Backup Site", row: "A", units: 42, used: 12, powerA: 6, powerB: 4, status: "active" },
];

const statusVariant: Record<string, "success" | "warning" | "destructive" | "secondary"> = {
  online: "success", active: "success", warning: "warning", critical: "destructive", offline: "destructive",
};

function UtilBar({ pct }: { pct: number }) {
  const color = pct >= 95 ? "bg-red-500" : pct >= 80 ? "bg-amber-500" : "bg-emerald-500";
  return (
    <div className="flex items-center gap-2">
      <div className="h-1.5 w-20 rounded-full bg-muted overflow-hidden">
        <div className={`h-full rounded-full ${color}`} style={{ width: `${Math.min(pct, 100)}%` }} />
      </div>
      <span className="text-xs tabular-nums">{pct.toFixed(0)}%</span>
    </div>
  );
}

export function DCIMPage() {
  const [tab, setTab] = useState<"datacenters" | "racks">("datacenters");
  const [search, setSearch] = useState("");

  const filteredDCs = DATA_CENTERS.filter(dc => dc.name.toLowerCase().includes(search.toLowerCase()));
  const filteredRacks = RACKS.filter(r => r.id.toLowerCase().includes(search.toLowerCase()) || r.dc.toLowerCase().includes(search.toLowerCase()));

  const totalRacks = DATA_CENTERS.reduce((s, d) => s + d.racks, 0);
  const usedRacks = DATA_CENTERS.reduce((s, d) => s + d.racksFilled, 0);
  const totalPower = DATA_CENTERS.reduce((s, d) => s + d.powerKW, 0);
  const usedPower = DATA_CENTERS.reduce((s, d) => s + d.powerUsed, 0);

  return (
    <div className="space-y-5">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <MetricCard title="Data centers" value={DATA_CENTERS.length} description="all regions" icon={Building2} iconClassName="bg-blue-500/10 text-blue-500" />
        <MetricCard title="Total racks" value={totalRacks} description={`${usedRacks} occupied`} icon={Server} iconClassName="bg-violet-500/10 text-violet-500" trend={{ value: 2.4, direction: "up" }} />
        <MetricCard title="Power capacity" value={totalPower} suffix=" kW" description={`${usedPower} kW in use`} icon={PlugZap} iconClassName="bg-amber-500/10 text-amber-500" />
        <MetricCard title="Avg. temp" value={21.5} suffix="°C" description="across all sites" icon={Thermometer} iconClassName="bg-emerald-500/10 text-emerald-500" />
      </div>

      <Card>
        <CardHeader className="pb-0">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex gap-1 rounded-lg border bg-muted/40 p-1 w-fit">
              {(["datacenters", "racks"] as const).map((t) => (
                <button key={t} onClick={() => setTab(t)}
                  className={`rounded-md px-4 py-1.5 text-sm font-medium capitalize transition-colors ${tab === t ? "bg-background shadow-sm" : "text-muted-foreground hover:text-foreground"}`}>
                  {t === "datacenters" ? "Data Centers" : "Racks"}
                </button>
              ))}
            </div>
            <div className="relative">
              <Search className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
              <Input placeholder="Search…" className="pl-8 h-8 w-56 text-sm" value={search} onChange={e => setSearch(e.target.value)} />
            </div>
          </div>
        </CardHeader>
        <CardContent className="pt-4">
          {tab === "datacenters" ? (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b text-left">
                  {["Name", "Location", "Tier", "Rack util.", "Power util.", "Avg temp", "Status"].map(h => (
                    <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground last:pr-0">{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y">
                {filteredDCs.map(dc => (
                  <tr key={dc.id} className="hover:bg-muted/30 transition-colors">
                    <td className="py-3 pr-4 font-medium">{dc.name}</td>
                    <td className="py-3 pr-4 text-muted-foreground">{dc.location}</td>
                    <td className="py-3 pr-4"><Badge variant="secondary">{dc.tier}</Badge></td>
                    <td className="py-3 pr-4"><UtilBar pct={(dc.racksFilled / dc.racks) * 100} /></td>
                    <td className="py-3 pr-4"><UtilBar pct={(dc.powerUsed / dc.powerKW) * 100} /></td>
                    <td className="py-3 pr-4 font-mono">{dc.tempC}°C</td>
                    <td className="py-3"><Badge variant={statusVariant[dc.status]}>{dc.status}</Badge></td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : (
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b text-left">
                  {["Rack ID", "Data center", "Row", "Unit util.", "Power A (kW)", "Power B (kW)", "Status"].map(h => (
                    <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground last:pr-0">{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y">
                {filteredRacks.map(r => (
                  <tr key={r.id} className="hover:bg-muted/30 transition-colors">
                    <td className="py-3 pr-4 font-mono font-medium">{r.id}</td>
                    <td className="py-3 pr-4 text-muted-foreground">{r.dc}</td>
                    <td className="py-3 pr-4">{r.row}</td>
                    <td className="py-3 pr-4"><UtilBar pct={(r.used / r.units) * 100} /></td>
                    <td className="py-3 pr-4 font-mono">{r.powerA}</td>
                    <td className="py-3 pr-4 font-mono">{r.powerB}</td>
                    <td className="py-3"><Badge variant={statusVariant[r.status]}>{r.status}</Badge></td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
