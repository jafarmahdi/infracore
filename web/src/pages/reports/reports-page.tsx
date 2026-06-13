import { Download, FileBarChart, FileText, RefreshCw, TrendingUp } from "lucide-react";
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, BarChart, Bar, Legend } from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { MetricCard } from "@/components/dashboard/metric-card";

const ASSET_GROWTH = [
  { month: "Jan", assets: 2680, vms: 124 }, { month: "Feb", assets: 2710, vms: 131 },
  { month: "Mar", assets: 2730, vms: 138 }, { month: "Apr", assets: 2765, vms: 145 },
  { month: "May", assets: 2800, vms: 154 }, { month: "Jun", assets: 2847, vms: 163 },
];

const ALERT_TREND = [
  { week: "W1", critical: 3, high: 7, warning: 14 }, { week: "W2", critical: 5, high: 9, warning: 11 },
  { week: "W3", critical: 2, high: 6, warning: 16 }, { week: "W4", critical: 4, high: 8, warning: 12 },
];

const AVAILABILITY = [
  { site: "DC1 – Riyadh Core", sla: 99.98, incidents: 1, downtime: "8m" },
  { site: "DC2 – Jeddah Edge",  sla: 99.80, incidents: 3, downtime: "1h 04m" },
  { site: "DC3 – Dammam West",  sla: 97.20, incidents: 8, downtime: "8h 12m" },
  { site: "DR – Backup Site",   sla: 99.50, incidents: 2, downtime: "25m" },
];

const SAVED_REPORTS = [
  { name: "Monthly asset inventory",      type: "Asset",       schedule: "Monthly",  lastRun: "2026-06-01", format: "XLSX" },
  { name: "Weekly availability report",   type: "Monitoring",  schedule: "Weekly",   lastRun: "2026-06-09", format: "PDF" },
  { name: "License compliance summary",   type: "License",     schedule: "Monthly",  lastRun: "2026-06-01", format: "PDF" },
  { name: "Contract expiry forecast",     type: "Contract",    schedule: "Quarterly",lastRun: "2026-04-01", format: "XLSX" },
  { name: "Security audit trail",         type: "Security",    schedule: "Weekly",   lastRun: "2026-06-09", format: "CSV" },
];

function slaColor(sla: number) {
  if (sla >= 99.9) return "text-emerald-600 dark:text-emerald-400";
  if (sla >= 99) return "text-amber-500";
  return "text-red-500";
}

export function ReportsPage() {
  return (
    <div className="space-y-5">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <MetricCard title="Reports scheduled" value={SAVED_REPORTS.length} description="automated delivery" icon={FileBarChart} iconClassName="bg-blue-500/10 text-blue-500" />
        <MetricCard title="Asset growth" value={3.8} suffix="%" description="month-over-month" icon={TrendingUp} iconClassName="bg-emerald-500/10 text-emerald-500" trend={{ value: 3.8, direction: "up" }} />
        <MetricCard title="Avg. platform SLA" value={99.12} suffix="%" description="30-day rolling" icon={FileText} iconClassName="bg-violet-500/10 text-violet-500" />
        <MetricCard title="Total incidents" value={14} description="this month" icon={FileBarChart} iconClassName="bg-amber-500/10 text-amber-500" trend={{ value: 8, direction: "down" }} />
      </div>

      <div className="grid gap-4 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-semibold">Asset growth (6 months)</CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <ResponsiveContainer width="100%" height={220}>
              <AreaChart data={ASSET_GROWTH} margin={{ top: 4, right: 4, left: -20, bottom: 0 }}>
                <defs>
                  <linearGradient id="assets" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.15} />
                    <stop offset="95%" stopColor="#3b82f6" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
                <XAxis dataKey="month" tick={{ fontSize: 11 }} className="fill-muted-foreground" />
                <YAxis tick={{ fontSize: 11 }} className="fill-muted-foreground" domain={[2600, 2900]} />
                <Tooltip contentStyle={{ fontSize: 12 }} />
                <Area type="monotone" dataKey="assets" stroke="#3b82f6" strokeWidth={2} fill="url(#assets)" name="Assets" />
                <Area type="monotone" dataKey="vms" stroke="#8b5cf6" strokeWidth={2} fill="transparent" name="VMs" />
              </AreaChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-sm font-semibold">Alert volume by severity (4 weeks)</CardTitle>
          </CardHeader>
          <CardContent className="pt-0">
            <ResponsiveContainer width="100%" height={220}>
              <BarChart data={ALERT_TREND} margin={{ top: 4, right: 4, left: -20, bottom: 0 }}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
                <XAxis dataKey="week" tick={{ fontSize: 11 }} className="fill-muted-foreground" />
                <YAxis tick={{ fontSize: 11 }} className="fill-muted-foreground" />
                <Tooltip contentStyle={{ fontSize: 12 }} />
                <Legend iconType="circle" iconSize={8} wrapperStyle={{ fontSize: 11 }} />
                <Bar dataKey="critical" fill="#ef4444" name="Critical" radius={[3, 3, 0, 0]} />
                <Bar dataKey="high" fill="#f97316" name="High" radius={[3, 3, 0, 0]} />
                <Bar dataKey="warning" fill="#f59e0b" name="Warning" radius={[3, 3, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-semibold">Site availability – 30 day</CardTitle>
        </CardHeader>
        <CardContent className="pt-0 overflow-x-auto">
          <table className="w-full text-sm">
            <thead><tr className="border-b text-left">
              {["Site", "SLA", "Incidents", "Total downtime"].map(h => (
                <th key={h} className="pb-2.5 pr-8 text-xs font-medium text-muted-foreground">{h}</th>
              ))}
            </tr></thead>
            <tbody className="divide-y">
              {AVAILABILITY.map(r => (
                <tr key={r.site} className="hover:bg-muted/30">
                  <td className="py-3 pr-8 font-medium">{r.site}</td>
                  <td className={`py-3 pr-8 font-mono font-semibold ${slaColor(r.sla)}`}>{r.sla.toFixed(2)}%</td>
                  <td className="py-3 pr-8 tabular-nums">{r.incidents}</td>
                  <td className="py-3 font-mono text-xs">{r.downtime}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="text-sm font-semibold">Scheduled reports</CardTitle>
            <Button size="sm"><RefreshCw className="h-3.5 w-3.5" /> Run all</Button>
          </div>
        </CardHeader>
        <CardContent className="pt-0 overflow-x-auto">
          <table className="w-full text-sm">
            <thead><tr className="border-b text-left">
              {["Report name", "Type", "Schedule", "Last run", "Format", ""].map(h => (
                <th key={h} className="pb-2.5 pr-4 text-xs font-medium text-muted-foreground">{h}</th>
              ))}
            </tr></thead>
            <tbody className="divide-y">
              {SAVED_REPORTS.map(r => (
                <tr key={r.name} className="hover:bg-muted/30">
                  <td className="py-3 pr-4 font-medium">{r.name}</td>
                  <td className="py-3 pr-4"><Badge variant="secondary">{r.type}</Badge></td>
                  <td className="py-3 pr-4 text-muted-foreground">{r.schedule}</td>
                  <td className="py-3 pr-4 font-mono text-xs">{r.lastRun}</td>
                  <td className="py-3 pr-4"><Badge variant="secondary">{r.format}</Badge></td>
                  <td className="py-3 text-right">
                    <Button variant="ghost" size="sm"><Download className="h-3.5 w-3.5" /></Button>
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
