import { useState } from "react";
import { Globe, Network, Search } from "lucide-react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { MetricCard } from "@/components/dashboard/metric-card";

const PREFIXES = [
  { prefix: "10.0.0.0/8",     vrf: "Global",    site: "All sites",          ips: 16777216, used: 8432,    role: "Container",  status: "active" },
  { prefix: "10.10.0.0/16",   vrf: "Global",    site: "DC1 – Riyadh Core",  ips: 65536,   used: 48120,   role: "Production", status: "active" },
  { prefix: "10.20.0.0/16",   vrf: "Global",    site: "DC2 – Jeddah Edge",  ips: 65536,   used: 12844,   role: "Production", status: "active" },
  { prefix: "10.30.0.0/16",   vrf: "Global",    site: "DC3 – Dammam West",  ips: 65536,   used: 31000,   role: "Production", status: "active" },
  { prefix: "10.100.0.0/24",  vrf: "MGMT",      site: "DC1 – Riyadh Core",  ips: 254,     used: 210,     role: "Management", status: "active" },
  { prefix: "172.16.0.0/12",  vrf: "DMZ",       site: "DC1 – Riyadh Core",  ips: 1048576, used: 2100,    role: "DMZ",        status: "active" },
  { prefix: "192.168.1.0/24", vrf: "MGMT",      site: "DR – Backup Site",   ips: 254,     used: 44,      role: "Management", status: "active" },
  { prefix: "10.50.0.0/20",   vrf: "VPN-Users", site: "Global",             ips: 4094,    used: 3800,    role: "VPN",        status: "warning" },
];

const IP_ADDRESSES = [
  { address: "10.10.0.1",   family: 4, vrf: "Global",    assigned: "SW-CORE-01",   type: "Interface",  status: "active" },
  { address: "10.10.0.2",   family: 4, vrf: "Global",    assigned: "FW-EDGE-01",   type: "Interface",  status: "active" },
  { address: "10.10.1.10",  family: 4, vrf: "Global",    assigned: "SRV-PROD-01",  type: "Primary",    status: "active" },
  { address: "10.10.1.11",  family: 4, vrf: "Global",    assigned: "SRV-PROD-02",  type: "Primary",    status: "active" },
  { address: "10.10.1.20",  family: 4, vrf: "Global",    assigned: "SRV-DB-01",    type: "Primary",    status: "active" },
  { address: "10.100.0.5",  family: 4, vrf: "MGMT",      assigned: "RTR-WAN-01",   type: "Management", status: "active" },
  { address: "10.100.0.50", family: 4, vrf: "MGMT",      assigned: "–",            type: "–",          status: "available" },
  { address: "172.16.0.1",  family: 4, vrf: "DMZ",       assigned: "FW-EDGE-01",   type: "Interface",  status: "active" },
  { address: "10.50.0.100", family: 4, vrf: "VPN-Users", assigned: "jafar.mahdi",  type: "Dynamic",    status: "dhcp" },
];

const VLANS = [
  { id: 10, name: "SERVERS",    site: "DC1 – Riyadh Core",  role: "Production", prefix: "10.10.1.0/24",  status: "active" },
  { id: 20, name: "MANAGEMENT", site: "All",                 role: "Management", prefix: "10.100.0.0/24", status: "active" },
  { id: 30, name: "DMZ",        site: "DC1 – Riyadh Core",  role: "DMZ",        prefix: "172.16.0.0/24", status: "active" },
  { id: 40, name: "STORAGE",    site: "DC1 – Riyadh Core",  role: "Storage",    prefix: "10.10.2.0/24",  status: "active" },
  { id: 50, name: "BACKUP",     site: "DR – Backup Site",   role: "Backup",     prefix: "10.30.1.0/24",  status: "active" },
  { id: 99, name: "NATIVE",     site: "All",                 role: "Native",     prefix: "–",             status: "active" },
  { id: 200, name: "VPN-POOL",  site: "DC1 – Riyadh Core",  role: "VPN",        prefix: "10.50.0.0/20",  status: "warning" },
];

const statusVariant: Record<string, "success" | "warning" | "destructive" | "secondary"> = {
  active: "success", warning: "warning", available: "secondary", dhcp: "secondary", deprecated: "destructive",
};

function UtilBar({ used, total }: { used: number; total: number }) {
  const pct = Math.min((used / total) * 100, 100);
  const color = pct >= 90 ? "bg-red-500" : pct >= 75 ? "bg-amber-500" : "bg-emerald-500";
  return (
    <div className="flex items-center gap-2">
      <div className="h-1.5 w-24 rounded-full bg-muted overflow-hidden">
        <div className={`h-full rounded-full ${color}`} style={{ width: `${pct}%` }} />
      </div>
      <span className="text-xs tabular-nums text-muted-foreground">{pct.toFixed(0)}%</span>
    </div>
  );
}

type Tab = "prefixes" | "addresses" | "vlans";

export function IPAMPage() {
  const [tab, setTab] = useState<Tab>("prefixes");
  const [search, setSearch] = useState("");

  const totalIPs = 16777216 + 65536 * 3;
  const usedIPs = PREFIXES.filter(p => p.role !== "Container").reduce((s, p) => s + p.used, 0);

  const filteredPrefixes = PREFIXES.filter(p => p.prefix.includes(search) || p.vrf.toLowerCase().includes(search.toLowerCase()) || p.site.toLowerCase().includes(search.toLowerCase()));
  const filteredIPs = IP_ADDRESSES.filter(ip => ip.address.includes(search) || ip.assigned.toLowerCase().includes(search.toLowerCase()));
  const filteredVLANs = VLANS.filter(v => v.name.toLowerCase().includes(search.toLowerCase()) || String(v.id).includes(search));

  return (
    <div className="space-y-5">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <MetricCard title="IP prefixes" value={PREFIXES.length} description="across all VRFs" icon={Network} iconClassName="bg-blue-500/10 text-blue-500" />
        <MetricCard title="Allocated IPs" value={usedIPs} description="from managed space" icon={Globe} iconClassName="bg-violet-500/10 text-violet-500" trend={{ value: 5.2, direction: "up" }} />
        <MetricCard title="VLANs" value={VLANS.length} description="configured" icon={Network} iconClassName="bg-emerald-500/10 text-emerald-500" />
        <MetricCard title="Utilization" value={Math.round((usedIPs / totalIPs) * 100)} suffix="%" description="address space used" icon={Network} iconClassName="bg-amber-500/10 text-amber-500" />
      </div>

      <Card>
        <CardHeader className="pb-0">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div className="flex gap-1 rounded-lg border bg-muted/40 p-1 w-fit">
              {(["prefixes", "addresses", "vlans"] as Tab[]).map(t => (
                <button key={t} onClick={() => setTab(t)}
                  className={`rounded-md px-4 py-1.5 text-sm font-medium capitalize transition-colors ${tab === t ? "bg-background shadow-sm" : "text-muted-foreground hover:text-foreground"}`}>
                  {t === "addresses" ? "IP Addresses" : t.charAt(0).toUpperCase() + t.slice(1)}
                </button>
              ))}
            </div>
            <div className="relative">
              <Search className="absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
              <Input placeholder="Search…" className="pl-8 h-8 w-52 text-sm" value={search} onChange={e => setSearch(e.target.value)} />
            </div>
          </div>
        </CardHeader>
        <CardContent className="pt-4 overflow-x-auto">
          {tab === "prefixes" && (
            <table className="w-full text-sm">
              <thead><tr className="border-b text-left">
                {["Prefix", "VRF", "Site", "Role", "Utilization", "Status"].map(h => (
                  <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground">{h}</th>
                ))}
              </tr></thead>
              <tbody className="divide-y">
                {filteredPrefixes.map(p => (
                  <tr key={p.prefix} className="hover:bg-muted/30">
                    <td className="py-3 pr-4 font-mono font-medium">{p.prefix}</td>
                    <td className="py-3 pr-4 text-muted-foreground">{p.vrf}</td>
                    <td className="py-3 pr-4 text-muted-foreground">{p.site}</td>
                    <td className="py-3 pr-4"><Badge variant="secondary">{p.role}</Badge></td>
                    <td className="py-3 pr-4"><UtilBar used={p.used} total={p.ips} /></td>
                    <td className="py-3"><Badge variant={statusVariant[p.status]}>{p.status}</Badge></td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
          {tab === "addresses" && (
            <table className="w-full text-sm">
              <thead><tr className="border-b text-left">
                {["IP Address", "Family", "VRF", "Assigned to", "Type", "Status"].map(h => (
                  <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground">{h}</th>
                ))}
              </tr></thead>
              <tbody className="divide-y">
                {filteredIPs.map(ip => (
                  <tr key={ip.address} className="hover:bg-muted/30">
                    <td className="py-3 pr-4 font-mono font-medium">{ip.address}</td>
                    <td className="py-3 pr-4"><Badge variant="secondary">IPv{ip.family}</Badge></td>
                    <td className="py-3 pr-4 text-muted-foreground">{ip.vrf}</td>
                    <td className="py-3 pr-4">{ip.assigned}</td>
                    <td className="py-3 pr-4 text-muted-foreground">{ip.type}</td>
                    <td className="py-3"><Badge variant={statusVariant[ip.status]}>{ip.status}</Badge></td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
          {tab === "vlans" && (
            <table className="w-full text-sm">
              <thead><tr className="border-b text-left">
                {["VLAN ID", "Name", "Site", "Role", "Prefix", "Status"].map(h => (
                  <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground">{h}</th>
                ))}
              </tr></thead>
              <tbody className="divide-y">
                {filteredVLANs.map(v => (
                  <tr key={v.id} className="hover:bg-muted/30">
                    <td className="py-3 pr-4 font-mono font-medium">{v.id}</td>
                    <td className="py-3 pr-4 font-medium">{v.name}</td>
                    <td className="py-3 pr-4 text-muted-foreground">{v.site}</td>
                    <td className="py-3 pr-4"><Badge variant="secondary">{v.role}</Badge></td>
                    <td className="py-3 pr-4 font-mono text-xs">{v.prefix}</td>
                    <td className="py-3"><Badge variant={statusVariant[v.status]}>{v.status}</Badge></td>
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
