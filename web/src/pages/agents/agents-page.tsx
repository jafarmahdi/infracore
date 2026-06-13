import { useState } from "react";
import { Bot, Radio, Search, Wifi, WifiOff } from "lucide-react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { MetricCard } from "@/components/dashboard/metric-card";

type AgentStatus = "online" | "offline" | "degraded" | "updating";

const AGENTS = [
  { id: "AGT-001", name: "agent-dc1-01",  version: "2.4.1", ip: "10.10.0.50", site: "DC1 – Riyadh Core",  status: "online" as AgentStatus,   heartbeat: "5s ago",   capabilities: ["snmp", "icmp", "ssh", "metrics"] },
  { id: "AGT-002", name: "agent-dc1-02",  version: "2.4.1", ip: "10.10.0.51", site: "DC1 – Riyadh Core",  status: "online" as AgentStatus,   heartbeat: "8s ago",   capabilities: ["snmp", "icmp", "metrics"] },
  { id: "AGT-003", name: "agent-dc1-03",  version: "2.3.9", ip: "10.10.0.52", site: "DC1 – Riyadh Core",  status: "online" as AgentStatus,   heartbeat: "12s ago",  capabilities: ["icmp", "ssh"] },
  { id: "AGT-004", name: "agent-dc2-01",  version: "2.4.1", ip: "10.20.0.50", site: "DC2 – Jeddah Edge",  status: "online" as AgentStatus,   heartbeat: "7s ago",   capabilities: ["snmp", "icmp", "metrics"] },
  { id: "AGT-005", name: "agent-dc2-02",  version: "2.4.0", ip: "10.20.0.51", site: "DC2 – Jeddah Edge",  status: "degraded" as AgentStatus, heartbeat: "2m ago",   capabilities: ["icmp"] },
  { id: "AGT-006", name: "agent-dc3-01",  version: "2.4.1", ip: "10.30.0.50", site: "DC3 – Dammam West",  status: "online" as AgentStatus,   heartbeat: "4s ago",   capabilities: ["snmp", "icmp", "ssh", "metrics"] },
  { id: "AGT-007", name: "agent-dc3-02",  version: "2.3.8", ip: "10.30.0.51", site: "DC3 – Dammam West",  status: "offline" as AgentStatus,  heartbeat: "4h ago",   capabilities: ["icmp", "metrics"] },
  { id: "AGT-008", name: "agent-dr-01",   version: "2.4.1", ip: "10.40.0.50", site: "DR – Backup Site",   status: "updating" as AgentStatus, heartbeat: "30s ago",  capabilities: ["snmp", "icmp"] },
];

const ST_VARIANT: Record<AgentStatus, "success" | "warning" | "destructive" | "secondary"> = {
  online: "success", degraded: "warning", offline: "destructive", updating: "secondary",
};

const CAPABILITY_COLORS: Record<string, string> = {
  snmp: "bg-blue-500/10 text-blue-600 dark:text-blue-400",
  icmp: "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400",
  ssh: "bg-violet-500/10 text-violet-600 dark:text-violet-400",
  metrics: "bg-amber-500/10 text-amber-600 dark:text-amber-400",
};

export function AgentsPage() {
  const [search, setSearch] = useState("");
  const [stFilter, setStFilter] = useState<"all" | AgentStatus>("all");

  const counts = {
    online: AGENTS.filter(a => a.status === "online").length,
    offline: AGENTS.filter(a => a.status === "offline").length,
    degraded: AGENTS.filter(a => a.status === "degraded").length,
    outdated: AGENTS.filter(a => a.version !== "2.4.1").length,
  };

  const filtered = AGENTS.filter(a =>
    (stFilter === "all" || a.status === stFilter) &&
    (a.name.toLowerCase().includes(search.toLowerCase()) || a.ip.includes(search) || a.site.toLowerCase().includes(search.toLowerCase()))
  );

  return (
    <div className="space-y-5">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <MetricCard title="Online agents" value={counts.online} description={`of ${AGENTS.length} total`} icon={Wifi} iconClassName="bg-emerald-500/10 text-emerald-500" />
        <MetricCard title="Offline" value={counts.offline} description="no heartbeat" icon={WifiOff} iconClassName="bg-red-500/10 text-red-500" />
        <MetricCard title="Degraded" value={counts.degraded} description="partial function" icon={Radio} iconClassName="bg-amber-500/10 text-amber-500" />
        <MetricCard title="Outdated" value={counts.outdated} description="not on v2.4.1" icon={Bot} iconClassName="bg-violet-500/10 text-violet-500" />
      </div>

      <Card>
        <CardHeader className="pb-0">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex flex-wrap gap-2">
              <div className="relative">
                <Search className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
                <Input placeholder="Name, IP, site…" className="pl-8 h-8 w-52 text-sm" value={search} onChange={e => setSearch(e.target.value)} />
              </div>
              <select value={stFilter} onChange={e => setStFilter(e.target.value as "all" | AgentStatus)}
                className="h-8 rounded-md border bg-background px-3 text-sm text-foreground">
                {["all", "online", "degraded", "offline", "updating"].map(s => (
                  <option key={s} value={s}>{s === "all" ? "All statuses" : s.charAt(0).toUpperCase() + s.slice(1)}</option>
                ))}
              </select>
            </div>
            <span className="text-xs text-muted-foreground">{filtered.length} agents</span>
          </div>
        </CardHeader>
        <CardContent className="pt-4 overflow-x-auto">
          <table className="w-full text-sm">
            <thead><tr className="border-b text-left">
              {["Agent ID", "Name", "Version", "IP", "Site", "Status", "Last heartbeat", "Capabilities"].map(h => (
                <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground whitespace-nowrap">{h}</th>
              ))}
            </tr></thead>
            <tbody className="divide-y">
              {filtered.map(a => (
                <tr key={a.id} className="hover:bg-muted/30">
                  <td className="py-3 pr-4 font-mono text-xs">{a.id}</td>
                  <td className="py-3 pr-4 font-medium">{a.name}</td>
                  <td className="py-3 pr-4 font-mono text-xs">
                    <span className={a.version !== "2.4.1" ? "text-amber-500" : ""}>{a.version}</span>
                  </td>
                  <td className="py-3 pr-4 font-mono text-xs">{a.ip}</td>
                  <td className="py-3 pr-4 text-muted-foreground">{a.site}</td>
                  <td className="py-3 pr-4"><Badge variant={ST_VARIANT[a.status]}>{a.status}</Badge></td>
                  <td className="py-3 pr-4 text-muted-foreground">{a.heartbeat}</td>
                  <td className="py-3">
                    <div className="flex flex-wrap gap-1">
                      {a.capabilities.map(c => (
                        <span key={c} className={`rounded px-1.5 py-0.5 text-[10px] font-medium ${CAPABILITY_COLORS[c] ?? "bg-muted text-muted-foreground"}`}>{c}</span>
                      ))}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
