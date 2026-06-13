import { useState } from "react";
import { AlertTriangle, CheckCheck, Clock, Search, ShieldAlert, XCircle } from "lucide-react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { MetricCard } from "@/components/dashboard/metric-card";

type Severity = "critical" | "high" | "warning" | "info";
type AlertStatus = "firing" | "acknowledged" | "resolved";

const ALERTS = [
  { id: "ALT-001", title: "Core switch packet loss above 5%",       source: "SW-CORE-01",   severity: "critical" as Severity, status: "firing" as AlertStatus,       duration: "2h 14m", assignee: "ops-team",    fired: "2026-06-13T10:51:00Z" },
  { id: "ALT-002", title: "Storage volume SAN-PROD-02 at 91%",     source: "SAN-PROD-02",  severity: "high" as Severity,     status: "firing" as AlertStatus,       duration: "45m",    assignee: "–",           fired: "2026-06-13T12:20:00Z" },
  { id: "ALT-003", title: "FW-EDGE-01 high CPU (88%)",              source: "FW-EDGE-01",   severity: "warning" as Severity,  status: "acknowledged" as AlertStatus, duration: "1h 02m", assignee: "jafar.m",     fired: "2026-06-13T12:03:00Z" },
  { id: "ALT-004", title: "SRV-DR-01 unreachable (ICMP timeout)",  source: "SRV-DR-01",    severity: "critical" as Severity, status: "firing" as AlertStatus,       duration: "4h 01m", assignee: "–",           fired: "2026-06-13T09:04:00Z" },
  { id: "ALT-005", title: "SRV-DR-02 unreachable (ICMP timeout)",  source: "SRV-DR-02",    severity: "critical" as Severity, status: "firing" as AlertStatus,       duration: "4h 00m", assignee: "–",           fired: "2026-06-13T09:05:00Z" },
  { id: "ALT-006", title: "Certificate expiry in 14 days",          source: "FW-EDGE-01",   severity: "warning" as Severity,  status: "acknowledged" as AlertStatus, duration: "12h",    assignee: "sec-team",    fired: "2026-06-13T01:00:00Z" },
  { id: "ALT-007", title: "Backup job failed: DC3 nightly",         source: "BACKUP-SVC",   severity: "high" as Severity,     status: "firing" as AlertStatus,       duration: "6h",     assignee: "–",           fired: "2026-06-13T07:00:00Z" },
  { id: "ALT-008", title: "Agent version out of compliance",        source: "BR-BGD-FW01",  severity: "info" as Severity,     status: "firing" as AlertStatus,       duration: "25m",    assignee: "–",           fired: "2026-06-13T12:40:00Z" },
  { id: "ALT-009", title: "Power supply redundancy restored",       source: "DC1-R14-SRV08",severity: "info" as Severity,     status: "resolved" as AlertStatus,     duration: "30m",    assignee: "auto-resolve",fired: "2026-06-13T11:30:00Z" },
  { id: "ALT-010", title: "BGP session flap RTR-WAN-01 → ISP1",   source: "RTR-WAN-01",   severity: "high" as Severity,     status: "resolved" as AlertStatus,     duration: "8m",     assignee: "auto-resolve",fired: "2026-06-13T08:10:00Z" },
];

const SEV_VARIANT: Record<Severity, "destructive" | "warning" | "secondary"> = {
  critical: "destructive", high: "destructive", warning: "warning", info: "secondary",
};
const ST_VARIANT: Record<AlertStatus, "destructive" | "warning" | "success" | "secondary"> = {
  firing: "destructive", acknowledged: "warning", resolved: "success",
};

export function AlertsPage() {
  const [search, setSearch] = useState("");
  const [sevFilter, setSevFilter] = useState<"all" | Severity>("all");
  const [stFilter, setStFilter] = useState<"all" | AlertStatus>("all");

  const counts = {
    critical: ALERTS.filter(a => a.severity === "critical" && a.status === "firing").length,
    high: ALERTS.filter(a => a.severity === "high" && a.status === "firing").length,
    warning: ALERTS.filter(a => a.severity === "warning" && a.status === "firing").length,
    resolved: ALERTS.filter(a => a.status === "resolved").length,
  };

  const filtered = ALERTS.filter(a =>
    (sevFilter === "all" || a.severity === sevFilter) &&
    (stFilter === "all" || a.status === stFilter) &&
    (a.title.toLowerCase().includes(search.toLowerCase()) || a.source.toLowerCase().includes(search.toLowerCase()))
  );

  return (
    <div className="space-y-5">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <MetricCard title="Critical" value={counts.critical} description="firing now" icon={XCircle} iconClassName="bg-red-500/10 text-red-500" trend={{ value: 1, direction: "up" }} />
        <MetricCard title="High" value={counts.high} description="firing now" icon={ShieldAlert} iconClassName="bg-orange-500/10 text-orange-500" />
        <MetricCard title="Warning" value={counts.warning} description="active" icon={AlertTriangle} iconClassName="bg-amber-500/10 text-amber-500" />
        <MetricCard title="Resolved today" value={counts.resolved} description="auto or manual" icon={CheckCheck} iconClassName="bg-emerald-500/10 text-emerald-500" trend={{ value: 12, direction: "down" }} />
      </div>

      <Card>
        <CardHeader className="pb-0">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex flex-wrap gap-2">
              <div className="relative">
                <Search className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
                <Input placeholder="Title or source…" className="pl-8 h-8 w-52 text-sm" value={search} onChange={e => setSearch(e.target.value)} />
              </div>
              <select value={sevFilter} onChange={e => setSevFilter(e.target.value as "all" | Severity)}
                className="h-8 rounded-md border bg-background px-3 text-sm text-foreground">
                {["all", "critical", "high", "warning", "info"].map(s => (
                  <option key={s} value={s}>{s === "all" ? "All severities" : s.charAt(0).toUpperCase() + s.slice(1)}</option>
                ))}
              </select>
              <select value={stFilter} onChange={e => setStFilter(e.target.value as "all" | AlertStatus)}
                className="h-8 rounded-md border bg-background px-3 text-sm text-foreground">
                {["all", "firing", "acknowledged", "resolved"].map(s => (
                  <option key={s} value={s}>{s === "all" ? "All statuses" : s.charAt(0).toUpperCase() + s.slice(1)}</option>
                ))}
              </select>
            </div>
            <span className="text-xs text-muted-foreground">{filtered.length} alerts</span>
          </div>
        </CardHeader>
        <CardContent className="pt-4 overflow-x-auto">
          <table className="w-full text-sm">
            <thead><tr className="border-b text-left">
              {["ID", "Title", "Source", "Severity", "Status", "Duration", "Assignee"].map(h => (
                <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground whitespace-nowrap">{h}</th>
              ))}
            </tr></thead>
            <tbody className="divide-y">
              {filtered.map(a => (
                <tr key={a.id} className="hover:bg-muted/30">
                  <td className="py-3 pr-4 font-mono text-xs">{a.id}</td>
                  <td className="py-3 pr-4 max-w-xs">
                    <span className="line-clamp-1">{a.title}</span>
                  </td>
                  <td className="py-3 pr-4 font-mono text-xs text-muted-foreground">{a.source}</td>
                  <td className="py-3 pr-4"><Badge variant={SEV_VARIANT[a.severity]}>{a.severity}</Badge></td>
                  <td className="py-3 pr-4"><Badge variant={ST_VARIANT[a.status]}>{a.status}</Badge></td>
                  <td className="py-3 pr-4 text-muted-foreground flex items-center gap-1 whitespace-nowrap">
                    <Clock className="h-3 w-3" />{a.duration}
                  </td>
                  <td className="py-3 text-muted-foreground">{a.assignee}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
