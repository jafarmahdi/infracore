import { useQuery } from "@tanstack/react-query";
import { Activity, AlertTriangle, Boxes, CalendarClock, RefreshCw, Server, ShieldCheck } from "lucide-react";
import { getDashboard } from "@/services/api/dashboard";
import { MetricCard } from "@/components/dashboard/metric-card";
import { OperationsChart } from "@/components/dashboard/operations-chart";
import { TrafficChart } from "@/components/dashboard/traffic-chart";
import { RecentEvents } from "@/components/dashboard/recent-events";
import { HealthPanel } from "@/components/dashboard/health-panel";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

function DashboardLoading() {
  return (
    <div className="space-y-5">
      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-5">{Array.from({ length: 5 }).map((_, i) => <Skeleton key={i} className="h-[142px]" />)}</div>
      <div className="grid gap-4 lg:grid-cols-3"><Skeleton className="h-[380px] lg:col-span-2" /><Skeleton className="h-[380px]" /></div>
    </div>
  );
}

export function DashboardPage() {
  const query = useQuery({ queryKey: ["dashboard"], queryFn: getDashboard, refetchInterval: 60_000 });

  if (query.isLoading) return <DashboardLoading />;
  if (query.isError) {
    return (
      <Card><CardContent className="flex min-h-[420px] flex-col items-center justify-center text-center">
        <div className="mb-4 rounded-xl bg-red-500/10 p-3 text-red-500"><AlertTriangle className="h-6 w-6" /></div>
        <h2 className="font-semibold">Dashboard data is unavailable</h2>
        <p className="mt-2 text-sm text-muted-foreground">The operations API did not respond. Your session is still active.</p>
        <Button className="mt-5" variant="outline" onClick={() => query.refetch()}><RefreshCw className="h-4 w-4" /> Try again</Button>
      </CardContent></Card>
    );
  }
  if (!query.data) return <div className="py-20 text-center text-sm text-muted-foreground">No dashboard data available.</div>;

  const { summary, utilization, traffic, events, deviceHealth } = query.data;
  return (
    <div className="space-y-5">
      <div className="flex flex-col justify-between gap-3 sm:flex-row sm:items-end">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-2xl font-semibold tracking-tight">Operations overview</h1>
            <span className="flex items-center gap-1.5 rounded-full bg-emerald-500/10 px-2 py-1 text-[10px] font-semibold text-emerald-600 dark:text-emerald-400">
              <span className="h-1.5 w-1.5 animate-pulse rounded-full bg-emerald-500" /> LIVE
            </span>
          </div>
          <p className="mt-1 text-sm text-muted-foreground">Infrastructure health and performance across all environments.</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={() => query.refetch()}><RefreshCw className="h-3.5 w-3.5" /> Refresh</Button>
          <Button size="sm"><Activity className="h-3.5 w-3.5" /> Open NOC view</Button>
        </div>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-5">
        <MetricCard title="Total assets" value={summary.totalAssets} description="vs. last month" icon={Boxes} iconClassName="bg-blue-500/10 text-blue-500" trend={summary.trends.totalAssets} />
        <MetricCard title="Online devices" value={summary.onlineDevices} description={`${summary.offlineDevices} offline`} icon={Server} iconClassName="bg-emerald-500/10 text-emerald-500" trend={summary.trends.onlineDevices} />
        <MetricCard title="Active alerts" value={summary.activeAlerts} description="since yesterday" icon={AlertTriangle} iconClassName="bg-red-500/10 text-red-500" trend={summary.trends.activeAlerts} />
        <MetricCard title="License expiry" value={summary.expiringLicenses} description="within 30 days" icon={CalendarClock} iconClassName="bg-amber-500/10 text-amber-500" trend={summary.trends.expiringLicenses} />
        <MetricCard title="Platform SLA" value={summary.sla} suffix="%" description="30-day availability" icon={ShieldCheck} iconClassName="bg-violet-500/10 text-violet-500" />
      </div>

      <div className="grid gap-4 lg:grid-cols-3">
        <OperationsChart data={utilization} />
        <HealthPanel data={deviceHealth} />
      </div>
      <div className="grid gap-4 lg:grid-cols-3">
        <RecentEvents events={events} />
        <TrafficChart data={traffic} />
      </div>
    </div>
  );
}
